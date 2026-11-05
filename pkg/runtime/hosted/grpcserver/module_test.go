// Copyright 2021 Layotto Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package grpcserver

import (
	"context"
	"errors"
	"net"
	"sync/atomic"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	pb "google.golang.org/grpc/health/grpc_health_v1"

	"mosn.io/layotto/pkg/runtime/hosted"
)

func TestRealGRPCRequestAndRelease(t *testing.T) {
	var released int32
	m := New(Config{Address: "127.0.0.1:0", Register: func(ctx context.Context, s *grpc.Server) error {
		pb.RegisterHealthServer(s, health.NewServer())
		return nil
	}, Release: func(context.Context) error { atomic.AddInt32(&released, 1); return nil }})
	r, err := hosted.New([]hosted.Registration{{Name: "api", New: func(context.Context) (hosted.Module, error) { return m, nil }}}, hosted.Options{})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := r.Start(ctx); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = r.Shutdown(ctx) })
	c, err := grpc.DialContext(ctx, m.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	resp, err := pb.NewHealthClient(c).Check(ctx, &pb.HealthCheckRequest{})
	if err != nil || resp.Status != pb.HealthCheckResponse_SERVING {
		t.Fatalf("response=%v err=%v", resp, err)
	}
	if err := r.Shutdown(ctx); err != nil {
		t.Fatal(err)
	}
	if err := r.Shutdown(ctx); err != nil {
		t.Fatal(err)
	}
	if atomic.LoadInt32(&released) != 1 {
		t.Fatal("release count")
	}
}

type blockedHealth struct {
	pb.UnimplementedHealthServer
	entered, finished chan struct{}
}

func (b *blockedHealth) Check(ctx context.Context, _ *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	close(b.entered)
	<-ctx.Done()
	close(b.finished)
	return nil, ctx.Err()
}

func TestDrainDeadlineForcesGRPCStop(t *testing.T) {
	b := &blockedHealth{entered: make(chan struct{}), finished: make(chan struct{})}
	m := New(Config{Address: "127.0.0.1:0", Register: func(ctx context.Context, s *grpc.Server) error { pb.RegisterHealthServer(s, b); return nil }})
	r, err := hosted.New([]hosted.Registration{{Name: "api", New: func(context.Context) (hosted.Module, error) { return m, nil }}}, hosted.Options{DrainTimeout: 20 * time.Millisecond, CloseTimeout: time.Second})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := r.Start(ctx); err != nil {
		t.Fatal(err)
	}
	c, err := grpc.DialContext(ctx, m.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	result := make(chan error, 1)
	go func() { _, err := pb.NewHealthClient(c).Check(ctx, &pb.HealthCheckRequest{}); result <- err }()
	select {
	case <-b.entered:
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}
	if err := r.Shutdown(ctx); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatal(err)
	}
	select {
	case <-b.finished:
	case <-ctx.Done():
		t.Fatal("RPC handler not cancelled")
	}
	if err := <-result; err == nil {
		t.Fatal("RPC unexpectedly completed")
	}
}

func TestInitAndBindFailuresReleaseResources(t *testing.T) {
	for _, mode := range []string{"register", "bind"} {
		t.Run(mode, func(t *testing.T) {
			ln, err := net.Listen("tcp", "127.0.0.1:0")
			if err != nil {
				t.Fatal(err)
			}
			defer ln.Close()
			released := 0
			m := New(Config{Address: ln.Addr().String(), Register: func(context.Context, *grpc.Server) error {
				if mode == "register" {
					return errors.New("API init failed")
				}
				return nil
			}, Release: func(context.Context) error { released++; return nil }})
			r, err := hosted.New([]hosted.Registration{{Name: "api", New: func(context.Context) (hosted.Module, error) { return m, nil }}}, hosted.Options{})
			if err != nil {
				t.Fatal(err)
			}
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			if err := r.Start(ctx); err == nil {
				t.Fatal("startup succeeded")
			}
			if released != 1 {
				t.Fatalf("release calls=%d", released)
			}
		})
	}
}
