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

package tcpserver

import (
	"context"
	"errors"
	"io"
	"net"
	"testing"
	"time"

	"mosn.io/layotto/pkg/runtime/hosted"
)

func TestDrainClosesListenerAndWaitsForSession(t *testing.T) {
	entered := make(chan struct{})
	m := New("127.0.0.1:0", func(ctx context.Context, c net.Conn) { close(entered); _, _ = io.Copy(c, c) })
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := m.Init(ctx); err != nil {
		t.Fatal(err)
	}
	if err := m.Start(ctx); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = m.Close(ctx) })
	c, err := net.DialTimeout("tcp", m.Address(), time.Second)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	_ = c.SetDeadline(time.Now().Add(time.Second))
	select {
	case <-entered:
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}
	drained := make(chan error, 1)
	go func() { drained <- m.Drain(ctx) }()
	for {
		h, _ := m.Status(ctx)
		if !h.Ready {
			break
		}
		select {
		case <-ctx.Done():
			t.Fatal(ctx.Err())
		case <-time.After(time.Millisecond):
		}
	}
	select {
	case err := <-drained:
		t.Fatalf("drain did not wait: %v", err)
	default:
	}
	if _, err := c.Write([]byte("ok")); err != nil {
		t.Fatal(err)
	}
	b := make([]byte, 2)
	if _, err := io.ReadFull(c, b); err != nil || string(b) != "ok" {
		t.Fatalf("inflight session: %q %v", b, err)
	}
	_ = c.Close()
	if err := <-drained; err != nil {
		t.Fatal(err)
	}
	if conn, err := net.DialTimeout("tcp", m.Address(), time.Second); err == nil {
		conn.Close()
		t.Fatal("listener still accepting")
	}
}

func TestRuntimeForcesLongConnectionClosedAfterTimeout(t *testing.T) {
	entered, finished := make(chan struct{}), make(chan struct{})
	m := New("127.0.0.1:0", func(ctx context.Context, c net.Conn) { close(entered); _, _ = io.Copy(c, c); close(finished) })
	r, err := hosted.New([]hosted.Registration{{Name: "native", New: func(context.Context) (hosted.Module, error) { return m, nil }}}, hosted.Options{DrainTimeout: 20 * time.Millisecond, CloseTimeout: time.Second})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := r.Start(ctx); err != nil {
		t.Fatal(err)
	}
	c, err := net.DialTimeout("tcp", m.Address(), time.Second)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	select {
	case <-entered:
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}
	if err := r.Shutdown(ctx); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("shutdown=%v", err)
	}
	select {
	case <-finished:
	case <-ctx.Done():
		t.Fatal("handler leaked")
	}
	_ = c.SetReadDeadline(time.Now().Add(time.Second))
	if _, err := c.Read(make([]byte, 1)); err == nil {
		t.Fatal("connection not closed")
	}
}

func TestUnexpectedAcceptFailureChangesHealth(t *testing.T) {
	m := New("127.0.0.1:0", func(context.Context, net.Conn) {})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := m.Init(ctx); err != nil {
		t.Fatal(err)
	}
	if err := m.Start(ctx); err != nil {
		t.Fatal(err)
	}
	_ = m.listener.Close()
	select {
	case <-m.acceptDone:
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}
	h, err := m.Status(ctx)
	if err != nil || h.Ready || h.Reason == "" {
		t.Fatalf("health=%+v err=%v", h, err)
	}
	if err := m.Close(ctx); err != nil {
		t.Fatal(err)
	}
}
