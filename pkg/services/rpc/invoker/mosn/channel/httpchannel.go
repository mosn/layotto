package channel

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/layotto/layotto/pkg/services/rpc"
	"github.com/valyala/fasthttp"
	_ "mosn.io/mosn/pkg/stream/http"
)

func init() {
	RegistChannel("http", newHttpChannel)
}

type httpChannel struct {
	listenerName string
	conns        chan net.Conn
}

func newHttpChannel(config ChannelConfig) (rpc.Channel, error) {
	return &httpChannel{
		listenerName: config.Listener,
		conns:        make(chan net.Conn, 2),
	}, nil
}

func (h *httpChannel) Do(req *rpc.RPCRequest) (*rpc.RPCResponse, error) {
	conn, err := h.getConn()
	if err != nil {
		return nil, err
	}

	if err = conn.SetDeadline(time.Now().Add(time.Duration(req.Timeout) * time.Millisecond)); err != nil {
		conn.Close()
		return nil, err
	}

	httpReq := h.constructReq(req)
	defer fasthttp.ReleaseRequest(httpReq)

	if _, err = httpReq.WriteTo(conn); err != nil {
		conn.Close()
		return nil, err
	}

	httpResp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(httpResp)
	if err = httpResp.Read(bufio.NewReader(conn)); err != nil {
		conn.Close()
		return nil, err
	}
	h.putConn(conn)

	if httpResp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("http response code %d", httpResp.StatusCode())
	}

	body := httpResp.Body()
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

func (h *httpChannel) getConn() (net.Conn, error) {
	select {
	case conn := <-h.conns:
		return conn, nil
	default:
	}

	local, remote := net.Pipe()
	localTcpConn := &fakeTcpConn{c: local}
	remoteTcpConn := &fakeTcpConn{c: remote}
	if err := acceptFunc(remoteTcpConn, h.listenerName); err != nil {
		return nil, err
	}
	return localTcpConn, nil
}

func (h *httpChannel) putConn(conn net.Conn) {
	select {
	case h.conns <- conn:
	default:
		conn.Close()
	}
}

func (h *httpChannel) constructReq(req *rpc.RPCRequest) *fasthttp.Request {
	httpReq := fasthttp.AcquireRequest()
	httpReq.Header.Set("host", "localhost")
	httpReq.Header.Set("id", req.Id)
	httpReq.SetBody(req.Data)
	httpReq.SetRequestURI(req.Method)
	method := http.MethodGet
	if verb := req.Header.Get("verb"); verb != "" {
		method = verb
	}
	httpReq.Header.SetMethod(method)
	if query := req.Header.Get("query_string"); query != "" {
		httpReq.URI().SetQueryString(query)
	}
	return httpReq
}
