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
	"net"
	"net/http"
	"time"

	"mosn.io/pkg/buffer"

	"github.com/valyala/fasthttp"
	// bridge to mosn
	_ "mosn.io/mosn/pkg/stream/http"

	"mosn.io/layotto/components/pkg/common"
	"mosn.io/layotto/components/rpc"
)

// init is regist http channel
func init() {
	RegistChannel("http", newHttpChannel)
}

// hstate is a pipe for readloop goroutine to communicate with request goroutine
type hstate struct {
	// request goroutine will read data from it
	reader net.Conn
	// readloop goroutine will write data to it
	writer net.Conn
}

func (h *hstate) onData(b buffer.IoBuffer) error {
	data := b.Bytes()
	if _, err := h.writer.Write(data); err != nil {
		return err
	}
	b.Drain(len(data))
	return nil
}

func (h *hstate) close() {
	h.reader.Close()
	h.writer.Close()
}

// httpChannel is Channel implement
type httpChannel struct {
	pool *connPool
}

// newHttpChannel is used to create rpc.Channel according to ChannelConfig
func newHttpChannel(config ChannelConfig) (rpc.Channel, error) {
	hc := &httpChannel{}
	hc.pool = newConnPool(
		config.Size,
		// dialFunc
		func() (net.Conn, error) {
			_, _, err := net.SplitHostPort(config.Listener)
			if err == nil {
				return net.Dial("tcp", config.Listener)
			}
			local, remote := net.Pipe()
			localTcpConn := &fakeTcpConn{c: local}
			remoteTcpConn := &fakeTcpConn{c: remote}
			if err := acceptFunc(remoteTcpConn, config.Listener); err != nil {
				return nil, err
			}
			// the goroutine model is:
			// request goroutine --->  localTcpConn ---> 	mosn
			//		^											|
			//		|											|
			//		|											|
			// 		hstate(net.Pipe) <-- readloop goroutine <---
			return localTcpConn, nil
		},
		// stateFunc
		func() interface{} {
			// hstate is a pipe for readloop goroutine to communicate with request goroutine
			s := &hstate{}
			s.reader, s.writer = net.Pipe()
			return s
		},
		hc.onData,
		hc.cleanup,
	)
	return hc, nil
}

// Do is used to handle RPCRequest and return RPCResponse
func (h *httpChannel) Do(req *rpc.RPCRequest) (*rpc.RPCResponse, error) {
	// 1. context.WithTimeout
	timeout := time.Duration(req.Timeout) * time.Millisecond
	ctx, cancel := context.WithTimeout(req.Ctx, timeout)
	defer cancel()

	// 2. get a fake connection with mosn
	// The pool will start a readloop gorountine,
	// which aims to read data from mosn and then write data to the hstate.writer
	conn, err := h.pool.Get(ctx)
	if err != nil {
		return nil, err
	}
	// 3. set deadline before write data to this connection
	hstate := conn.state.(*hstate)
	deadline, _ := ctx.Deadline()
	if err = conn.SetWriteDeadline(deadline); err != nil {
		hstate.close()
		h.pool.Put(conn, true)
		return nil, common.Error(common.UnavailebleCode, err.Error())
	}
	// 4. write data to this fake connection
	httpReq := h.constructReq(req)
	defer fasthttp.ReleaseRequest(httpReq)

	if _, err = httpReq.WriteTo(conn); err != nil {
		hstate.close()
		h.pool.Put(conn, true)
		return nil, common.Error(common.UnavailebleCode, err.Error())
	}

	// 5. read response data and parse it into fasthttp.Response
	httpResp := &fasthttp.Response{}
	hstate.reader.SetReadDeadline(deadline)

	if err = httpResp.Read(bufio.NewReader(hstate.reader)); err != nil {
		hstate.close()
		h.pool.Put(conn, true)
		return nil, common.Error(common.UnavailebleCode, err.Error())
	}
	h.pool.Put(conn, false)

	// 6. convert result to rpc.RPCResponse,which is the response of rpc invoker
	body := httpResp.Body()
	if httpResp.StatusCode() != http.StatusOK {
		return nil, common.Errorf(common.UnavailebleCode, "http response code %d, body: %s", httpResp.StatusCode(), string(body))
	}

	rpcResp := &rpc.RPCResponse{
		ContentType: string(httpResp.Header.ContentType()),
		Data:        body,
		Header:      map[string][]string{},
	}
	httpResp.Header.VisitAll(func(key, value []byte) {
		rpcResp.Header[string(key)] = []string{string(value)}
	})
	return rpcResp, nil
}

// constructReq is handle rpc.RPCRequest to fasthttp.Request
func (h *httpChannel) constructReq(req *rpc.RPCRequest) *fasthttp.Request {
	httpReq := fasthttp.AcquireRequest()
	httpReq.SetBody(req.Data)
	httpReq.SetRequestURI(req.Method)
	method := http.MethodGet
	if verb := req.Header.Get("verb"); verb != "" {
		method = verb
		delete(req.Header, "verb")
	}
	httpReq.Header.SetMethod(method)
	if query := req.Header.Get("query_string"); query != "" {
		httpReq.URI().SetQueryString(query)
		delete(req.Header, "query_string")
	}
	req.Header.Range(func(key string, value string) bool {
		httpReq.Header.Set(key, value)
		return true
	})
	httpReq.SetHost("localhost")
	httpReq.Header.Set("id", req.Id)
	return httpReq
}

func (h *httpChannel) onData(conn *wrapConn) error {
	hstate := conn.state.(*hstate)
	return hstate.onData(conn.buf)
}

func (h *httpChannel) cleanup(conn *wrapConn, err error) {
	hstate := conn.state.(*hstate)
	hstate.close()
}
