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
	"bufio"
	"context"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"

	"mosn.io/layotto/components/rpc"
)

type testhttpServer struct {
}

func (ts *testhttpServer) accept(conn net.Conn, listener string) error {
	go ts.readLoop(conn)
	return nil
}

func (ts *testhttpServer) readLoop(conn net.Conn) {
	defer conn.Close()

	for {
		req := fasthttp.AcquireRequest()
		if err := req.Read(bufio.NewReader(conn)); err != nil {
			break
		}

		content := string(req.Body())
		switch content {
		case "close":
			return
		case "timeout":
			time.Sleep(2 * time.Second)
		default:
		}

		resp := fasthttp.AcquireResponse()
		resp.SetBody(req.Body())

		if _, err := resp.WriteTo(conn); err != nil {
			log.Println("test server err:", err.Error())
			break
		}
	}
}

func startTestHttpServer() {
	ts := &testhttpServer{}
	acceptFunc = ts.accept
}

func TestHttpChannel(t *testing.T) {
	startTestHttpServer()

	channel, err := newHttpChannel(ChannelConfig{Size: 1})
	assert.Nil(t, err)

	req := &rpc.RPCRequest{Ctx: context.TODO(), Id: "foo", Method: "bar", Data: []byte("hello"), Timeout: 1000}
	resp, err := channel.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, "hello", string(resp.Data))
}

func TestRenewHttpConn(t *testing.T) {
	startTestHttpServer()

	channel, err := newHttpChannel(ChannelConfig{Size: 1})
	assert.Nil(t, err)

	req := &rpc.RPCRequest{Ctx: context.TODO(), Id: "foo", Method: "bar", Data: []byte("close"), Timeout: 1000}
	_, err = channel.Do(req)
	assert.Error(t, err)

	req = &rpc.RPCRequest{Ctx: context.TODO(), Id: "foo", Method: "bar", Data: []byte("hello"), Timeout: 1000}
	resp, err := channel.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, "hello", string(resp.Data))
}

func TestManyRequests(t *testing.T) {
	startTestHttpServer()

	channel, err := newHttpChannel(ChannelConfig{Size: 1})
	assert.Nil(t, err)

	for i := 0; i < 100; i++ {
		req := &rpc.RPCRequest{Ctx: context.TODO(), Id: "foo", Method: "bar", Data: []byte("hello"), Timeout: 1000}
		_, err = channel.Do(req)
		assert.Nil(t, err)
	}
}

func TestResponseTimeout(t *testing.T) {
	startTestHttpServer()

	channel, err := newHttpChannel(ChannelConfig{Size: 1})
	assert.Nil(t, err)

	req := &rpc.RPCRequest{Ctx: context.TODO(), Id: "foo", Method: "bar", Data: []byte("timeout"), Timeout: 1000}
	_, err = channel.Do(req)
	assert.Error(t, err)

	for i := 0; i < 100; i++ {
		req = &rpc.RPCRequest{Ctx: context.TODO(), Id: "foo", Method: "bar", Data: []byte("hello"), Timeout: 1000}
		_, err = channel.Do(req)
		assert.Nil(t, err)
	}
}

func TestConcurrent(t *testing.T) {
	startTestHttpServer()

	channel, err := newHttpChannel(ChannelConfig{Size: 1})
	assert.Nil(t, err)

	var wg sync.WaitGroup
	for i := 0; i < 30; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			req := &rpc.RPCRequest{Ctx: context.TODO(), Id: "foo", Method: "bar", Data: []byte("hello" + strconv.Itoa(i)), Timeout: 1000}
			resp, err := channel.Do(req)
			assert.Nil(t, err)
			assert.Equal(t, "hello"+strconv.Itoa(i), string(resp.Data))
		}(i)
	}
	wg.Wait()
}

func TestConstructReq(t *testing.T) {
	channel, err := newHttpChannel(ChannelConfig{Size: 1})
	assert.Nil(t, err)

	httpChannel := channel.(*httpChannel)

	reqs := []*rpc.RPCRequest{
		{Id: "foo", Method: "bar", Data: []byte("hello world")},
		{Id: "foo", Method: "/bar", Data: []byte("hello world")},
	}

	for _, req := range reqs {
		fastHttpReq := httpChannel.constructReq(req)
		sb := &strings.Builder{}
		fastHttpReq.WriteTo(sb)
		assert.Equal(t, "GET /bar HTTP/1.1\r\nHost: localhost\r\nContent-Length: 11\r\nId: foo\r\n\r\nhello world", sb.String())
	}
}
