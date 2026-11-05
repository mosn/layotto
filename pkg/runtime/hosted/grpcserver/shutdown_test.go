/*
 * Copyright 2026 Layotto Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package grpcserver

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "google.golang.org/grpc/health/grpc_health_v1"

	"mosn.io/layotto/pkg/runtime/hosted"
)

func testContext(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	t.Cleanup(cancel)
	return ctx
}
func await(t *testing.T, ctx context.Context, ch <-chan struct{}) {
	t.Helper()
	select {
	case <-ch:
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}
}

func TestClosePendingHandshakeWithinBudget(t *testing.T) {
	for _, drain := range []bool{false, true} {
		name := "direct close"
		if drain {
			name = "drain then close"
		}
		t.Run(name, func(t *testing.T) {
			m := New(Config{Address: "127.0.0.1:0", Options: []grpc.ServerOption{grpc.ConnectionTimeout(2 * time.Second)}})
			ctx := testContext(t)
			if err := m.Init(ctx); err != nil {
				t.Fatal(err)
			}
			if err := m.Start(ctx); err != nil {
				t.Fatal(err)
			}
			conn, err := net.DialTimeout("tcp", m.Address(), time.Second)
			if err != nil {
				t.Fatal(err)
			}
			defer conn.Close()
			_ = conn.SetReadDeadline(time.Now().Add(time.Second))
			// Receiving SETTINGS proves gRPC has accepted the socket and entered the
			// HTTP/2 handshake. Never send the client preface.
			if _, err = io.ReadFull(conn, make([]byte, 9)); err != nil {
				t.Fatal(err)
			}
			result := make(chan error, 1)
			go func() {
				if drain {
					d, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
					err := m.Drain(d)
					cancel()
					if !errors.Is(err, context.DeadlineExceeded) {
						result <- fmt.Errorf("Drain = %v, want deadline exceeded", err)
						return
					}
				}
				c, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				result <- m.Close(c)
			}()
			select {
			case err = <-result:
				if err != nil {
					t.Fatal(err)
				}
			case <-time.After(time.Second):
				_ = conn.Close()
				<-result
				t.Fatal("shutdown waited for handshake timeout")
			}
		})
	}
}

type cleanupHealth struct {
	pb.UnimplementedHealthServer
	entered, cancelled, permit, finished chan struct{}
	once                                 sync.Once
}

func newCleanupHealth() *cleanupHealth {
	return &cleanupHealth{entered: make(chan struct{}), cancelled: make(chan struct{}), permit: make(chan struct{}), finished: make(chan struct{})}
}
func (s *cleanupHealth) unblock() { s.once.Do(func() { close(s.permit) }) }
func (s *cleanupHealth) cleanup(ctx context.Context) error {
	close(s.entered)
	<-ctx.Done()
	close(s.cancelled)
	<-s.permit
	close(s.finished)
	return ctx.Err()
}
func (s *cleanupHealth) Check(ctx context.Context, _ *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	return nil, s.cleanup(ctx)
}
func (s *cleanupHealth) Watch(_ *pb.HealthCheckRequest, stream pb.Health_WatchServer) error {
	return s.cleanup(stream.Context())
}

func startCall(t *testing.T, ctx context.Context, m *Module, s *cleanupHealth, stream bool) {
	t.Helper()
	c, err := grpc.DialContext(ctx, m.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { s.unblock(); _ = c.Close() })
	go func() {
		if stream {
			watch, err := pb.NewHealthClient(c).Watch(ctx, &pb.HealthCheckRequest{})
			if err == nil {
				_, _ = watch.Recv()
			}
			return
		}
		_, _ = pb.NewHealthClient(c).Check(ctx, &pb.HealthCheckRequest{})
	}()
	await(t, ctx, s.entered)
}

func TestCloseWaitsForUnaryAndStreamCleanup(t *testing.T) {
	for _, stream := range []bool{false, true} {
		name := "unary"
		if stream {
			name = "stream"
		}
		t.Run(name, func(t *testing.T) {
			s := newCleanupHealth()
			released := make(chan struct{})
			m := New(Config{Address: "127.0.0.1:0", Register: func(_ context.Context, g *grpc.Server) error { pb.RegisterHealthServer(g, s); return nil }, Release: func(context.Context) error { close(released); return nil }})
			ctx := testContext(t)
			if err := m.Init(ctx); err != nil {
				t.Fatal(err)
			}
			if err := m.Start(ctx); err != nil {
				t.Fatal(err)
			}
			startCall(t, ctx, m, s, stream)
			d, cancel := context.WithTimeout(ctx, 20*time.Millisecond)
			defer cancel()
			_ = m.Drain(d)
			await(t, ctx, s.cancelled)
			result := make(chan error, 1)
			go func() { result <- m.Close(ctx) }()
			select {
			case <-released:
				t.Fatal("released during handler cleanup")
			case <-time.After(20 * time.Millisecond):
			}
			s.unblock()
			select {
			case err := <-result:
				if err != nil {
					t.Fatal(err)
				}
			case <-ctx.Done():
				t.Fatal(ctx.Err())
			}
			await(t, ctx, s.finished)
			await(t, ctx, released)
		})
	}
}

type dependency struct{ closed int32 }

func (*dependency) Init(context.Context) error  { return nil }
func (*dependency) Start(context.Context) error { return nil }
func (*dependency) Status(context.Context) (hosted.Health, error) {
	return hosted.Health{Ready: true}, nil
}
func (*dependency) Drain(context.Context) error   { return nil }
func (d *dependency) Close(context.Context) error { atomic.AddInt32(&d.closed, 1); return nil }

func TestCloseTimeoutRetainsReleaseAndRuntimeDependencies(t *testing.T) {
	s := newCleanupHealth()
	var released int32
	m := New(Config{Address: "127.0.0.1:0", Register: func(_ context.Context, g *grpc.Server) error { pb.RegisterHealthServer(g, s); return nil }, Release: func(context.Context) error { atomic.AddInt32(&released, 1); return nil }})
	dep := &dependency{}
	r, err := hosted.New([]hosted.Registration{
		{Name: "dependency", New: func(context.Context) (hosted.Module, error) { return dep, nil }},
		{Name: "api", DependsOn: []string{"dependency"}, New: func(context.Context) (hosted.Module, error) { return m, nil }},
	}, hosted.Options{DrainTimeout: 20 * time.Millisecond, CloseTimeout: 80 * time.Millisecond})
	if err != nil {
		t.Fatal(err)
	}
	ctx := testContext(t)
	if err = r.Start(ctx); err != nil {
		t.Fatal(err)
	}
	startCall(t, ctx, m, s, false)
	err = r.Shutdown(ctx)
	if !errors.Is(err, hosted.ErrResourcesInUse) || !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("shutdown=%v", err)
	}
	if atomic.LoadInt32(&released) != 0 || atomic.LoadInt32(&dep.closed) != 0 {
		t.Fatal("released resources still used by handler")
	}
	status := r.Status(ctx)
	if status.State != hosted.Failed || status.Modules[0].CloseAttempted {
		t.Fatalf("snapshot=%+v", status)
	}
	if err2 := r.Shutdown(ctx); err2.Error() != err.Error() {
		t.Fatal("shutdown result changed")
	}
	s.unblock()
	await(t, ctx, m.calls.stop())
	if atomic.LoadInt32(&released) != 0 {
		t.Fatal("unexpected asynchronous release after failed close")
	}
}

func TestCloseTracksMiddlewareCleanup(t *testing.T) {
	entered, permit := make(chan struct{}), make(chan struct{})
	var once sync.Once
	t.Cleanup(func() { once.Do(func() { close(permit) }) })
	m := New(Config{Address: "127.0.0.1:0", Options: []grpc.ServerOption{grpc.ChainUnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, next grpc.UnaryHandler) (interface{}, error) {
		reply, err := next(ctx, req)
		close(entered)
		<-permit
		return reply, err
	})}, Register: func(_ context.Context, g *grpc.Server) error {
		pb.RegisterHealthServer(g, &pb.UnimplementedHealthServer{})
		return nil
	}})
	ctx := testContext(t)
	if err := m.Init(ctx); err != nil {
		t.Fatal(err)
	}
	if err := m.Start(ctx); err != nil {
		t.Fatal(err)
	}
	c, err := grpc.DialContext(ctx, m.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	go func() { _, _ = pb.NewHealthClient(c).Check(ctx, &pb.HealthCheckRequest{}) }()
	await(t, ctx, entered)
	closing, cancel := context.WithTimeout(ctx, 20*time.Millisecond)
	defer cancel()
	if err = m.Close(closing); !errors.Is(err, hosted.ErrResourcesInUse) {
		t.Fatalf("close=%v", err)
	}
	once.Do(func() { close(permit) })
	await(t, ctx, m.calls.stop())
}

func TestRejectsSingleInterceptorOptions(t *testing.T) {
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(func(context.Context, interface{}, *grpc.UnaryServerInfo, grpc.UnaryHandler) (interface{}, error) {
			return nil, nil
		}),
		grpc.StreamInterceptor(func(interface{}, grpc.ServerStream, *grpc.StreamServerInfo, grpc.StreamHandler) error { return nil }),
	}
	for _, opt := range opts {
		m := New(Config{Address: "127.0.0.1:0", Options: []grpc.ServerOption{opt}})
		if err := m.Init(testContext(t)); err == nil || !strings.Contains(err.Error(), "ChainUnaryInterceptor/ChainStreamInterceptor") {
			t.Fatalf("init=%v", err)
		}
	}
}

func TestNoNewHandlersAfterStop(t *testing.T) {
	var c calls
	await(t, testContext(t), c.stop())
	called := false
	_, err := c.unary(context.Background(), nil, nil, func(context.Context, interface{}) (interface{}, error) { called = true; return nil, nil })
	if err == nil || called {
		t.Fatal("accepted late handler")
	}
}
