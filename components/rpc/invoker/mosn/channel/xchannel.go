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
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"mosn.io/api"
	"mosn.io/layotto/components/rpc"
	"mosn.io/layotto/components/rpc/invoker/mosn/transport_protocol"
)

func init() {
	RegistChannel("bolt", newXChannel)
	RegistChannel("boltv2", newXChannel)
	RegistChannel("dubbo", newXChannel)
}

func newXChannel(config ChannelConfig) (rpc.Channel, error) {
	proto := transport_protocol.GetProtocol(config.Protocol)
	if proto == nil {
		return nil, fmt.Errorf("protocol %s not found", proto)
	}
	if err := proto.Init(config.Ext); err != nil {
		return nil, err
	}

	m := &xChannel{proto: proto}
	m.pool = newConnPool(
		config.Size,
		func() (net.Conn, error) {
			local, remote := net.Pipe()
			localTcpConn := &fakeTcpConn{c: local}
			remoteTcpConn := &fakeTcpConn{c: remote}
			if err := acceptFunc(remoteTcpConn, config.Listener); err != nil {
				return nil, status.Error(codes.Internal, err.Error())
			}
			return localTcpConn, nil
		},
		func() interface{} {
			return &xstate{respChans: map[uint64]chan api.XRespFrame{}}
		},
		m.onData,
	)

	return m, nil
}

type xstate struct {
	reqid     uint64
	mu        sync.Mutex
	respChans map[uint64]chan api.XRespFrame
}

type xChannel struct {
	proto transport_protocol.TransportProtocol
	pool  *connPool
}

func (m *xChannel) Do(req *rpc.RPCRequest) (*rpc.RPCResponse, error) {
	timeout := time.Duration(req.Timeout) * time.Millisecond
	ctx, cancel := context.WithTimeout(req.Ctx, timeout)
	defer cancel()

	conn, err := m.pool.Get(ctx)
	if err != nil {
		return nil, err
	}

	xstate := conn.state.(*xstate)

	// encode request
	frame := m.proto.ToFrame(req)
	id := atomic.AddUint64(&xstate.reqid, 1)
	frame.SetRequestId(id)
	buf, encErr := m.proto.Encode(req.Ctx, frame)
	if encErr != nil {
		m.pool.Put(conn, false)
		return nil, status.Error(codes.InvalidArgument, encErr.Error())
	}

	respChan := make(chan api.XRespFrame, 1)
	// set timeout
	deadline, _ := ctx.Deadline()
	if err := conn.SetWriteDeadline(deadline); err != nil {
		m.pool.Put(conn, true)
		return nil, status.Error(codes.Unavailable, err.Error())
	}
	// register response channel
	xstate.mu.Lock()
	xstate.respChans[id] = respChan
	xstate.mu.Unlock()

	// write packet
	if _, err := conn.Write(buf.Bytes()); err != nil {
		m.removeRespChan(xstate, id)
		m.pool.Put(conn, true)
		return nil, status.Error(codes.Unavailable, err.Error())
	}
	m.pool.Put(conn, false)

	select {
	case resp := <-respChan:
		return m.proto.FromFrame(resp)
	case <-ctx.Done():
		m.removeRespChan(xstate, id)
		return nil, status.Error(codes.DeadlineExceeded, ErrTimeout.Error())
	case <-conn.closeChan:
		m.removeRespChan(xstate, id)
		return nil, status.Error(codes.Unavailable, ErrConnClosed.Error())
	}
}

func (m *xChannel) removeRespChan(xstate *xstate, id uint64) {
	xstate.mu.Lock()
	delete(xstate.respChans, id)
	xstate.mu.Unlock()
}

func (m *xChannel) onData(conn *wrapConn) error {
	for {
		iframe, err := m.proto.Decode(context.TODO(), conn.buf)
		if err != nil {
			return err
		}

		if iframe == nil {
			break
		}

		frame, ok := iframe.(api.XRespFrame)
		if !ok {
			return errors.New("[runtime][rpc]xchannel type not XRespFrame")
		}

		reqID := frame.GetRequestId()
		xstate := conn.state.(*xstate)
		xstate.mu.Lock()
		notifyChan, ok := xstate.respChans[reqID]
		if ok {
			delete(xstate.respChans, reqID)
		}
		xstate.mu.Unlock()
		if ok {
			notifyChan <- frame
		}
	}
	return nil
}
