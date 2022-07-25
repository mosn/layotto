/*
 * Copyright 2021 Layotto Authors
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

package channel

import (
	"errors"
	"fmt"
	"net"
	"time"

	"mosn.io/mosn/pkg/server"

	common "mosn.io/layotto/components/pkg/common"
	"mosn.io/layotto/components/rpc"
)

var (
	ErrTimeout = errors.New("request timeout")

	acceptFunc = func(conn net.Conn, listener string) error {
		srv := server.GetServer()
		lis := srv.Handler().FindListenerByName(listener)
		if lis == nil {
			return common.Error(common.InternalCode, "[rpc]invalid listener name")
		}
		lis.GetListenerCallbacks().OnAccept(conn, false, nil, nil, nil, nil)
		return nil
	}

	registry = map[string]func(config ChannelConfig) (rpc.Channel, error){}
)

// ChannelConfig is Channel config
type ChannelConfig struct {
	Protocol string                 `json:"protocol"`
	Listener string                 `json:"listener"`
	Size     int                    `json:"size"`
	Ext      map[string]interface{} `json:"ext"`
}

// GetChannel creates a rpc.Channel according to config.Protocol
func GetChannel(config ChannelConfig) (rpc.Channel, error) {
	c, ok := registry[config.Protocol]
	if !ok {
		return nil, fmt.Errorf("channel %s not found", config.Protocol)
	}
	return c(config)
}

// RegistChannel is set protocol
func RegistChannel(proto string, f func(config ChannelConfig) (rpc.Channel, error)) {
	registry[proto] = f
}

// fakeTcpConn simulates tcp connection. It implements net.Conn
type fakeTcpConn struct {
	c net.Conn
}

// read data
func (t *fakeTcpConn) Read(b []byte) (n int, err error) {
	return t.c.Read(b)
}

// write data
func (t *fakeTcpConn) Write(b []byte) (n int, err error) {
	return t.c.Write(b)
}

// Close is closed connected
func (t *fakeTcpConn) Close() error {
	return t.c.Close()
}

// LocalAddr is get local net address
func (t *fakeTcpConn) LocalAddr() net.Addr {
	return &net.TCPAddr{}
}

// RemoteAddr is get remote address
func (t *fakeTcpConn) RemoteAddr() net.Addr {
	return &net.TCPAddr{}
}

// SetDeadline is set deadline
func (t *fakeTcpConn) SetDeadline(time time.Time) error {
	return t.c.SetDeadline(time)
}

// SetReadDeadline is set read deadline
func (t *fakeTcpConn) SetReadDeadline(time time.Time) error {
	return t.c.SetReadDeadline(time)
}

// SetWriteDeadline is set write deadline
func (t fakeTcpConn) SetWriteDeadline(time time.Time) error {
	return t.c.SetWriteDeadline(time)
}
