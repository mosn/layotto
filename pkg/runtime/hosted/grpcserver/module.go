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

// Package grpcserver adapts a gRPC API server to the hosted runtime. The
// registration callback can register Layotto APIs using their existing Init
// and Register methods; neither the runtime nor this package imports MOSN.
package grpcserver

import (
	"context"
	"errors"
	"net"
	"sync"

	"google.golang.org/grpc"
	"mosn.io/layotto/pkg/runtime/hosted"
)

type Config struct {
	Address string
	Options []grpc.ServerOption
	// Register prepares and registers APIs on this owned server. It must honor
	// ctx. Release owns any API-specific resources (e.g. callback connections),
	// including those allocated before Register fails; never shared components.
	Register func(context.Context, *grpc.Server) error
	Release  func(context.Context) error
}

type Module struct {
	config                    Config
	mu                        sync.RWMutex
	server                    *grpc.Server
	listener                  net.Listener
	done                      chan struct{}
	drainDone                 chan struct{}
	running, draining, closed bool
	serveErr                  error
	closeErr                  error
}

func New(config Config) *Module {
	config.Options = append([]grpc.ServerOption(nil), config.Options...)
	return &Module{config: config}
}

func (m *Module) Init(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if m.config.Address == "" {
		return errors.New("gRPC listen address is required")
	}
	m.mu.Lock()
	if m.closed {
		m.mu.Unlock()
		return hosted.ErrStopped
	}
	if m.server != nil {
		m.mu.Unlock()
		return errors.New("gRPC module already initialized")
	}
	s := grpc.NewServer(m.config.Options...)
	m.server = s
	m.mu.Unlock()
	if m.config.Register != nil {
		return m.config.Register(ctx, s)
	}
	return nil
}

func (m *Module) Start(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closed || m.draining {
		return hosted.ErrStopped
	}
	if m.running {
		return nil
	}
	if m.done != nil {
		return errors.New("gRPC module cannot restart after Serve exits")
	}
	if m.server == nil {
		return errors.New("gRPC module not initialized")
	}
	var lc net.ListenConfig
	ln, err := lc.Listen(ctx, "tcp", m.config.Address)
	if err != nil {
		return err
	}
	m.listener, m.done, m.running = ln, make(chan struct{}), true
	s, done := m.server, m.done
	go func() {
		err := s.Serve(ln)
		m.mu.Lock()
		m.running = false
		if !m.draining && !m.closed {
			if err == nil {
				err = errors.New("gRPC Serve exited unexpectedly")
			}
			m.serveErr = err
		}
		m.mu.Unlock()
		close(done)
	}()
	return nil
}

func (m *Module) Address() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.listener == nil {
		return ""
	}
	return m.listener.Addr().String()
}

func (m *Module) Status(ctx context.Context) (hosted.Health, error) {
	if err := ctx.Err(); err != nil {
		return hosted.Health{}, err
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.serveErr != nil {
		return hosted.Health{Reason: m.serveErr.Error()}, nil
	}
	return hosted.Health{Ready: m.running && !m.draining && !m.closed}, nil
}

func (m *Module) Drain(ctx context.Context) error {
	m.mu.Lock()
	m.draining = true
	s := m.server
	if s == nil || m.closed {
		m.mu.Unlock()
		return nil
	}
	if m.drainDone == nil {
		m.drainDone = make(chan struct{})
		done := m.drainDone
		go func() { s.GracefulStop(); close(done) }()
	}
	done := m.drainDone
	m.mu.Unlock()
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		s.Stop()
		return ctx.Err()
	}
}

func (m *Module) Close(ctx context.Context) error {
	m.mu.Lock()
	if m.closed {
		err := m.closeErr
		m.mu.Unlock()
		return err
	}
	m.closed, m.running = true, false
	s, ln, done, drained := m.server, m.listener, m.done, m.drainDone
	m.mu.Unlock()
	if ln != nil {
		_ = ln.Close()
	}
	if s != nil {
		s.Stop()
	}
	var failures hosted.Errors
	for _, ch := range []chan struct{}{done, drained} {
		if ch != nil {
			select {
			case <-ch:
			case <-ctx.Done():
				failures = append(failures, ctx.Err())
			}
		}
	}
	if m.config.Release != nil {
		if err := m.config.Release(ctx); err != nil {
			failures = append(failures, err)
		}
	}
	var err error
	if len(failures) != 0 {
		err = failures
	}
	m.mu.Lock()
	m.closeErr = err
	m.mu.Unlock()
	return err
}

var _ hosted.Module = (*Module)(nil)
