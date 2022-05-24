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
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"mosn.io/api"
	"mosn.io/mosn/pkg/protocol/xprotocol/bolt"
	"mosn.io/pkg/buffer"

	"mosn.io/layotto/components/rpc"
)

var (
	proto = "bolt"
)

type testserver struct {
	api.XProtocol
}

func (ts *testserver) accept(conn net.Conn, listener string) error {
	go ts.readLoop(conn)
	return nil
}

func (ts *testserver) handleRequest(frame api.XFrame) ([]byte, error) {
	data := frame.GetData()
	if data != nil {
		data := string(data.Bytes())
		switch data {
		case "deformity":
			return []byte("deformity"), nil
		case "timeout":
			time.Sleep(time.Second)
		case "close":
			return nil, errors.New("trigger close")
		default:
			if strings.Contains(data, "echo") {
				resp := bolt.NewRpcResponse(uint32(frame.GetRequestId()), bolt.ResponseStatusSuccess, nil, buffer.NewIoBufferBytes([]byte(data)))
				buf, _ := ts.XProtocol.Encode(context.TODO(), resp)
				return buf.Bytes(), nil
			}
		}
	}
	resp := bolt.NewRpcResponse(uint32(frame.GetRequestId()), bolt.ResponseStatusSuccess, nil, buffer.NewIoBufferBytes([]byte("ok")))
	buf, _ := ts.XProtocol.Encode(context.TODO(), resp)
	return buf.Bytes(), nil
}

func (ts *testserver) readLoop(conn net.Conn) {
	data := buffer.GetIoBuffer(1024)

	defer conn.Close()

	for {
		p := make([]byte, 64)
		n, err := conn.Read(p)
		if err != nil {
			fmt.Println("readloop return")
			return
		}

		data.Write(p[:n])

		for {
			packet, err := ts.XProtocol.Decode(context.TODO(), data)
			if err != nil {
				return
			}
			if packet == nil {
				break
			}

			go func() {
				frame := packet.(*bolt.Request)
				bytes, err := ts.handleRequest(frame)
				if err != nil {
					conn.Close()
					return
				}
				if _, err = conn.Write(bytes); err != nil {
					conn.Close()
					return
				}
			}()
		}
	}
}

func startTestServer() {
	ts := &testserver{
		XProtocol: (&bolt.XCodec{}).NewXProtocol(context.TODO()), // always bolt here.
	}
	acceptFunc = ts.accept
}

func TestChannel(t *testing.T) {
	startTestServer()

	config := ChannelConfig{Size: 1, Protocol: proto, Ext: map[string]interface{}{"class": "xxx"}}
	channel, err := newXChannel(config)
	assert.Nil(t, err)

	req := &rpc.RPCRequest{Ctx: context.TODO(), Id: "foo", Method: "bar", Data: []byte("hello world"), Timeout: 1000}
	resp, err := channel.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, "ok", string(resp.Data))
}

func TestInvalidProtocal(t *testing.T) {
	config := ChannelConfig{Size: 1, Protocol: "dubbogo", Ext: map[string]interface{}{"class": "xxx"}}
	_, err := newXChannel(config)
	assert.NotNil(t, err)

	assert.Equal(t, err.Error(), "protocol dubbogo not found")
}

func TestChannelTimeout(t *testing.T) {
	startTestServer()

	config := ChannelConfig{Size: 1, Protocol: proto, Ext: map[string]interface{}{"class": "xxx"}}
	channel, err := newXChannel(config)
	assert.Nil(t, err)

	req := &rpc.RPCRequest{Ctx: context.TODO(), Id: "foo", Method: "bar", Data: []byte("timeout"), Timeout: 990}
	_, err = channel.Do(req)
	t.Log(err)
	assert.True(t, strings.Contains(err.Error(), ErrTimeout.Error()))
}

func TestMemConnClosed(t *testing.T) {
	startTestServer()

	config := ChannelConfig{Size: 1, Protocol: proto, Ext: map[string]interface{}{"class": "xxx"}}
	channel, err := newXChannel(config)
	assert.Nil(t, err)

	req := &rpc.RPCRequest{Ctx: context.TODO(), Id: "foo", Method: "bar", Data: []byte("close"), Timeout: 1000}
	_, err = channel.Do(req)
	assert.True(t, strings.Contains(err.Error(), "EOF"))
}

func TestReturnInvalidPacket(t *testing.T) {
	startTestServer()

	config := ChannelConfig{Size: 1, Protocol: proto, Ext: map[string]interface{}{"class": "xxx"}}
	channel, err := newXChannel(config)
	assert.Nil(t, err)

	req := &rpc.RPCRequest{Ctx: context.TODO(), Id: "foo", Method: "bar", Data: []byte("deformity"), Timeout: 1000}
	_, err = channel.Do(req)
	assert.True(t, strings.Contains(err.Error(), ErrTimeout.Error()))
}

func TestRenewConn(t *testing.T) {
	startTestServer()

	config := ChannelConfig{Size: 1, Protocol: proto, Ext: map[string]interface{}{"class": "xxx"}}
	channel, err := newXChannel(config)
	assert.Nil(t, err)

	req := &rpc.RPCRequest{Ctx: context.TODO(), Id: "foo", Method: "bar", Data: []byte("close"), Timeout: 1000}
	_, err = channel.Do(req)
	assert.True(t, strings.Contains(err.Error(), "EOF"))

	var errcount int32
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			req := &rpc.RPCRequest{Ctx: context.TODO(), Id: "foo", Method: "bar", Data: []byte("hello world"), Timeout: 1000}
			resp, err := channel.Do(req)
			if err != nil {
				assert.Equal(t, "io: read/write on closed pipe", err.Error())
				atomic.AddInt32(&errcount, 1)
			}
			if err == nil {
				assert.Equal(t, "ok", string(resp.Data))
			}
		}(i)
	}
	wg.Wait()
	assert.Equal(t, int32(0), atomic.LoadInt32(&errcount))
}

func TestConncurrent(t *testing.T) {
	startTestServer()

	config := ChannelConfig{Size: 1, Protocol: proto, Ext: map[string]interface{}{"class": "xxx"}}
	channel, err := newXChannel(config)
	assert.Nil(t, err)

	var wg sync.WaitGroup
	size := 10

	wg.Add(size)
	go func() {
		for i := 0; i < size; i++ {
			go func(i int) {
				defer wg.Done()
				req := &rpc.RPCRequest{Ctx: context.TODO(), Id: "foo", Method: "bar", Data: []byte("close"), Timeout: 1000}
				channel.Do(req)
			}(i)
		}
	}()

	wg.Add(size)
	go func() {
		for i := 0; i < size; i++ {
			go func(i int) {
				defer wg.Done()
				req := &rpc.RPCRequest{Ctx: context.TODO(), Id: "foo", Method: "bar", Data: []byte("timeout"), Timeout: 1000}
				channel.Do(req)
			}(i)
		}
	}()

	wg.Add(size)
	go func() {
		for i := 0; i < size; i++ {
			go func(i int) {
				defer wg.Done()
				req := &rpc.RPCRequest{Ctx: context.TODO(), Id: "foo", Method: "bar", Data: []byte("hello world"), Timeout: 1000}
				channel.Do(req)
			}(i)
		}
	}()

	wg.Wait()
	req := &rpc.RPCRequest{Ctx: context.TODO(), Id: "foo", Method: "bar", Data: []byte("hello world"), Timeout: 1000}
	channel.Do(req)
	//_, err = channel.Do(req)
	//assert.Nil(t, err)

	size = 100
	wg.Add(size)
	for i := 0; i < size; i++ {
		go func(i int) {
			defer wg.Done()

			data := "echo" + strconv.Itoa(i)
			req := &rpc.RPCRequest{Ctx: context.TODO(), Id: "foo", Method: "bar", Data: []byte(data), Timeout: 1000}
			resp, err := channel.Do(req)
			assert.Nil(t, err)
			assert.Equal(t, data, string(resp.Data))
		}(i)
	}
	wg.Wait()
}

func TestReadloopError(t *testing.T) {
	startTestServer()

	config := ChannelConfig{Size: 1, Protocol: proto, Ext: map[string]interface{}{"class": "xxx"}}
	channel, err := newXChannel(config)
	assert.Nil(t, err)

	xchannel := channel.(*xChannel)
	conn, _ := xchannel.pool.Get(context.TODO())
	xstate := conn.state.(*xstate)
	xchannel.pool.Put(conn, false)

	go func() {
		time.AfterFunc(100*time.Millisecond, func() {
			req := &rpc.RPCRequest{Ctx: context.TODO(), Id: "foo", Method: "bar", Data: []byte("close"), Timeout: 500}
			_, err = channel.Do(req)
			assert.True(t, strings.Contains(err.Error(), "EOF"))
		})
	}()

	var wg sync.WaitGroup
	total := 10
	wg.Add(total)
	for i := 0; i < total; i++ {
		go func() {
			defer wg.Done()
			req := &rpc.RPCRequest{Ctx: context.TODO(), Id: "foo", Method: "bar", Data: []byte("timeout"), Timeout: 500}
			_, err = channel.Do(req)
			t.Log(err)
			assert.True(t, strings.Contains(err.Error(), "EOF"))
		}()
	}
	wg.Wait()

	xstate.mu.Lock()
	callsSize := len(xstate.calls)
	xstate.mu.Unlock()
	assert.Equal(t, 0, callsSize)

	req := &rpc.RPCRequest{Ctx: context.TODO(), Id: "foo", Method: "bar", Data: []byte("hello world"), Timeout: 500}
	resp, err := channel.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, string(resp.Data), "ok")
}
