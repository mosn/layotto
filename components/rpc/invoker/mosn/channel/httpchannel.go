package channel

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/valyala/fasthttp"
	"mosn.io/layotto/components/rpc"
	_ "mosn.io/mosn/pkg/stream/http"
)

func init() {
	RegistChannel("http", newHttpChannel)
}

type httpChannel struct {
	pool *connPool
}

func newHttpChannel(config ChannelConfig) (rpc.Channel, error) {
	return &httpChannel{
		pool: newConnPool(
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
			nil,
			nil,
		),
	}, nil
}

func (h *httpChannel) Do(req *rpc.RPCRequest) (*rpc.RPCResponse, error) {
	timeout := time.Duration(req.Timeout) * time.Millisecond
	ctx, cancel := context.WithTimeout(req.Ctx, timeout)
	defer cancel()

	conn, err := h.pool.Get(ctx)
	if err != nil {
		return nil, err
	}

	deadline, _ := ctx.Deadline()
	if err = conn.SetDeadline(deadline); err != nil {
		h.pool.Put(conn, true)
		return nil, err
	}

	httpReq := h.constructReq(req)
	defer fasthttp.ReleaseRequest(httpReq)

	if _, err = httpReq.WriteTo(conn); err != nil {
		h.pool.Put(conn, true)
		return nil, err
	}

	httpResp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(httpResp)
	if err = httpResp.Read(bufio.NewReader(conn)); err != nil {
		h.pool.Put(conn, true)
		return nil, err
	}
	body := httpResp.Body()
	h.pool.Put(conn, false)
	if httpResp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("http response code %d, body: %s", httpResp.StatusCode(), string(body))
	}

	data := make([]byte, len(body))
	copy(data, body)
	rpcResp := &rpc.RPCResponse{
		ContentType: string(httpResp.Header.ContentType()),
		Data:        data,
		Header:      map[string][]string{},
	}
	httpResp.Header.VisitAll(func(key, value []byte) {
		rpcResp.Header[string(key)] = []string{string(value)}
	})
	return rpcResp, nil
}

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
