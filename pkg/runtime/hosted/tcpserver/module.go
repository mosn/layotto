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

// Package tcpserver hosts native TCP protocol handlers without an API server.
// A product adapter supplies the protocol handler and borrows its engine from
// an explicit runtime dependency. This package owns listeners and connections.
package tcpserver

import (
	"context"
	"errors"
	"net"
	"sync"

	"mosn.io/layotto/pkg/runtime/hosted"
)

type Handler func(context.Context, net.Conn)

type Module struct {
	address                                string
	handler                                Handler
	mu                                     sync.RWMutex
	listener                               net.Listener
	connections                            map[net.Conn]struct{}
	acceptDone                             chan struct{}
	sessionsDone                           chan struct{}
	wg                                     sync.WaitGroup
	cancel                                 context.CancelFunc
	initialized, running, draining, closed bool
	serveErr                               error
}

// New takes a handler that must return when its connection is closed or its
// context is cancelled. Request/transaction boundaries belong to the product.
// Drain waits for sessions to finish; Close cancels and closes remaining ones.
func New(address string, handler Handler) *Module {
	return &Module{address: address, handler: handler, connections: make(map[net.Conn]struct{})}
}

func (m *Module) Init(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if m.address == "" || m.handler == nil {
		return errors.New("TCP address and handler are required")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closed {
		return hosted.ErrStopped
	}
	m.initialized = true
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
	if m.acceptDone != nil {
		return errors.New("TCP module cannot restart after Accept exits")
	}
	if !m.initialized {
		return errors.New("TCP module not initialized")
	}
	var lc net.ListenConfig
	ln, err := lc.Listen(ctx, "tcp", m.address)
	if err != nil {
		return err
	}
	serviceCtx, cancel := context.WithCancel(context.Background())
	m.listener, m.cancel, m.running = ln, cancel, true
	m.acceptDone, m.sessionsDone = make(chan struct{}), make(chan struct{})
	go m.accept(serviceCtx, ln)
	// WaitGroup.Add only runs in accept; wait until accept has finished before
	// waiting for sessions so there is no Add/Wait race on the empty group.
	go func() { <-m.acceptDone; m.wg.Wait(); close(m.sessionsDone) }()
	return nil
}

func (m *Module) accept(ctx context.Context, ln net.Listener) {
	defer close(m.acceptDone)
	for {
		conn, err := ln.Accept()
		if err != nil {
			m.mu.Lock()
			m.running = false
			if !m.draining && !m.closed {
				m.serveErr = err
			}
			m.mu.Unlock()
			return
		}
		m.mu.Lock()
		if m.draining || m.closed {
			m.mu.Unlock()
			_ = conn.Close()
			continue
		}
		m.connections[conn] = struct{}{}
		m.wg.Add(1)
		m.mu.Unlock()
		go func() {
			defer m.wg.Done()
			defer func() { _ = conn.Close(); m.mu.Lock(); delete(m.connections, conn); m.mu.Unlock() }()
			m.handler(ctx, conn)
		}()
	}
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
	ln, done := m.listener, m.sessionsDone
	m.mu.Unlock()
	if ln != nil {
		_ = ln.Close()
	}
	if done == nil {
		return nil
	}
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (m *Module) Close(ctx context.Context) error {
	m.mu.Lock()
	m.closed, m.running = true, false
	ln, cancel, done := m.listener, m.cancel, m.sessionsDone
	connections := make([]net.Conn, 0, len(m.connections))
	for conn := range m.connections {
		connections = append(connections, conn)
	}
	m.mu.Unlock()
	if ln != nil {
		_ = ln.Close()
	}
	if cancel != nil {
		cancel()
	}
	for _, conn := range connections {
		_ = conn.Close()
	}
	if done == nil {
		return nil
	}
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return hosted.Errors{hosted.ErrResourcesInUse, ctx.Err()}
	}
}

var _ hosted.Module = (*Module)(nil)
