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
	"net"
	"sync"
)

// Track connections as soon as Accept returns, before gRPC's handshake can
// block. Closing the listener alone does not close those pending handshakes.
type trackedListener struct {
	net.Listener
	mu          sync.Mutex
	closed      bool
	connections map[*trackedConn]struct{}
}

func newTrackedListener(ln net.Listener) *trackedListener {
	return &trackedListener{Listener: ln, connections: make(map[*trackedConn]struct{})}
}
func (l *trackedListener) Accept() (net.Conn, error) {
	conn, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}
	l.mu.Lock()
	if l.closed {
		l.mu.Unlock()
		_ = conn.Close()
		return nil, net.ErrClosed
	}
	c := &trackedConn{Conn: conn, owner: l}
	l.connections[c] = struct{}{}
	l.mu.Unlock()
	return c, nil
}
func (l *trackedListener) Close() error {
	l.mu.Lock()
	l.closed = true
	l.mu.Unlock()
	return l.Listener.Close()
}
func (l *trackedListener) closeConnections() {
	_ = l.Close()
	l.mu.Lock()
	conns := make([]*trackedConn, 0, len(l.connections))
	for c := range l.connections {
		conns = append(conns, c)
	}
	l.mu.Unlock()
	for _, c := range conns {
		_ = c.Close()
	}
}

type trackedConn struct {
	net.Conn
	owner *trackedListener
	once  sync.Once
	err   error
}

func (c *trackedConn) Close() error {
	c.once.Do(func() {
		c.err = c.Conn.Close()
		c.owner.mu.Lock()
		delete(c.owner.connections, c)
		c.owner.mu.Unlock()
	})
	return c.err
}
