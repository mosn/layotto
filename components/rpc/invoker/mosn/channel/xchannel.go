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

	"mosn.io/api"
	common "mosn.io/layotto/components/pkg/common"
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
				return nil, err
			}
			return localTcpConn, nil
		},
		func() interface{} {
			return &xstate{calls: map[uint32]chan call{}}
		},
		m.onData,
		m.cleanup,
	)

	return m, nil
}

type xstate struct {
	reqid uint32
	mu    sync.Mutex
	calls map[uint32]chan call
}

type call struct {
	resp api.XRespFrame
	err  error
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
	id := atomic.AddUint32(&xstate.reqid, 1)
	frame.SetRequestId(uint64(id))
	buf, encErr := m.proto.Encode(req.Ctx, frame)
	if encErr != nil {
		m.pool.Put(conn, false)
		return nil, common.Error(common.InternalCode, encErr.Error())
	}

	callChan := make(chan call, 1)
	// set timeout
	deadline, _ := ctx.Deadline()
	if err := conn.SetWriteDeadline(deadline); err != nil {
		m.pool.Put(conn, true)
		return nil, common.Error(common.UnavailebleCode, err.Error())
	}
	// register response channel
	xstate.mu.Lock()
	xstate.calls[id] = callChan
	xstate.mu.Unlock()

	// write packet
	if _, err := conn.Write(buf.Bytes()); err != nil {
		m.removeCall(xstate, id)
		m.pool.Put(conn, true)
		return nil, common.Error(common.UnavailebleCode, err.Error())
	}
	m.pool.Put(conn, false)

	select {
	case res := <-callChan:
		if res.err != nil {
			return nil, common.Error(common.UnavailebleCode, res.err.Error())
		}
		return m.proto.FromFrame(res.resp)
	case <-ctx.Done():
		m.removeCall(xstate, id)
		return nil, common.Error(common.TimeoutCode, ErrTimeout.Error())
	}
}

func (m *xChannel) removeCall(xstate *xstate, id uint32) {
	xstate.mu.Lock()
	delete(xstate.calls, id)
	xstate.mu.Unlock()
}

func (m *xChannel) onData(conn *wrapConn) error {
	xstate := conn.state.(*xstate)
	for {
		var iframe interface{}
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
		reqID32 := uint32(reqID)
		xstate.mu.Lock()
		notifyChan, ok := xstate.calls[reqID32]
		if ok {
			delete(xstate.calls, reqID32)
		}
		xstate.mu.Unlock()
		if ok {
			notifyChan <- call{resp: frame}
		}
	}
	return nil
}

func (m *xChannel) cleanup(c *wrapConn, err error) {
	xstate := c.state.(*xstate)
	// cleanup pending calls
	xstate.mu.Lock()
	for id, notifyChan := range xstate.calls {
		notifyChan <- call{err: err}
		delete(xstate.calls, id)
	}
	xstate.mu.Unlock()
}
