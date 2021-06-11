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
	"mosn.io/layotto/components/rpc"
	"mosn.io/mosn/pkg/protocol/xprotocol"
	"mosn.io/mosn/pkg/protocol/xprotocol/bolt"
	"mosn.io/mosn/pkg/types"
	"mosn.io/pkg/buffer"
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

func (s *testserver) handleRequest(frame api.XFrame) ([]byte, error) {
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
				buf, _ := s.XProtocol.Encode(context.TODO(), resp)
				return buf.Bytes(), nil
			}
		}
	}
	resp := bolt.NewRpcResponse(uint32(frame.GetRequestId()), bolt.ResponseStatusSuccess, nil, buffer.NewIoBufferBytes([]byte("ok")))
	buf, _ := s.XProtocol.Encode(context.TODO(), resp)
	return buf.Bytes(), nil
}

func (s *testserver) readLoop(conn net.Conn) {
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
			packet, err := s.XProtocol.Decode(context.TODO(), data)
			if err != nil {
				return
			}
			if packet == nil {
				break
			}

			go func() {
				frame := packet.(*bolt.Request)
				bytes, err := s.handleRequest(frame)
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
		XProtocol: xprotocol.GetProtocol(types.ProtocolName(proto)),
	}
	acceptFunc = ts.accept
}

func TestChannel(t *testing.T) {
	startTestServer()

	config := ChannelConfig{Size: 1, Protocol: proto, Ext: map[string]interface{}{"class": "xxx"}}
	channel, err := newXChannel(config)
	assert.Nil(t, err)

	req := &rpc.RPCRequest{Id: "foo", Method: "bar", Data: []byte("hello world"), Timeout: 1000}
	resp, err := channel.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, "ok", string(resp.Data))
}

func TestChannelTimeout(t *testing.T) {
	startTestServer()

	config := ChannelConfig{Size: 1, Protocol: proto, Ext: map[string]interface{}{"class": "xxx"}}
	channel, err := newXChannel(config)
	assert.Nil(t, err)

	req := &rpc.RPCRequest{Id: "foo", Method: "bar", Data: []byte("timeout"), Timeout: 500}
	_, err = channel.Do(req)
	t.Log(err)
	assert.Equal(t, ErrTimeout, err)
}

func TestMemConnClosed(t *testing.T) {
	startTestServer()

	config := ChannelConfig{Size: 1, Protocol: proto, Ext: map[string]interface{}{"class": "xxx"}}
	channel, err := newXChannel(config)
	assert.Nil(t, err)

	req := &rpc.RPCRequest{Id: "foo", Method: "bar", Data: []byte("close"), Timeout: 1000}
	_, err = channel.Do(req)
	assert.Equal(t, err, ErrConnClosed)
}

func TestReturnInvalidPacket(t *testing.T) {
	startTestServer()

	config := ChannelConfig{Size: 1, Protocol: proto, Ext: map[string]interface{}{"class": "xxx"}}
	channel, err := newXChannel(config)
	assert.Nil(t, err)

	req := &rpc.RPCRequest{Id: "foo", Method: "bar", Data: []byte("deformity"), Timeout: 1000}
	_, err = channel.Do(req)
	assert.Equal(t, err, ErrTimeout)
}

func TestRenewConn(t *testing.T) {
	startTestServer()

	config := ChannelConfig{Size: 1, Protocol: proto, Ext: map[string]interface{}{"class": "xxx"}}
	channel, err := newXChannel(config)
	assert.Nil(t, err)

	req := &rpc.RPCRequest{Id: "foo", Method: "bar", Data: []byte("close"), Timeout: 1000}
	_, err = channel.Do(req)
	assert.Equal(t, err, ErrConnClosed)

	var errcount int32
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			req := &rpc.RPCRequest{Id: "foo", Method: "bar", Data: []byte("hello world"), Timeout: 1000}
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
	assert.Equal(t, int32(1), atomic.LoadInt32(&errcount))
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
				req := &rpc.RPCRequest{Id: "foo", Method: "bar", Data: []byte("close"), Timeout: 1000}
				channel.Do(req)
			}(i)
		}
	}()

	wg.Add(size)
	go func() {
		for i := 0; i < size; i++ {
			go func(i int) {
				defer wg.Done()
				req := &rpc.RPCRequest{Id: "foo", Method: "bar", Data: []byte("timeout"), Timeout: 1000}
				channel.Do(req)
			}(i)
		}
	}()

	wg.Add(size)
	go func() {
		for i := 0; i < size; i++ {
			go func(i int) {
				defer wg.Done()
				req := &rpc.RPCRequest{Id: "foo", Method: "bar", Data: []byte("hello world"), Timeout: 1000}
				channel.Do(req)
			}(i)
		}
	}()

	wg.Wait()
	req := &rpc.RPCRequest{Id: "foo", Method: "bar", Data: []byte("hello world"), Timeout: 1000}
	_, err = channel.Do(req)
	//assert.Nil(t, err)

	size = 100
	wg.Add(size)
	for i := 0; i < size; i++ {
		go func(i int) {
			defer wg.Done()

			data := "echo" + strconv.Itoa(i)
			req := &rpc.RPCRequest{Id: "foo", Method: "bar", Data: []byte(data), Timeout: 1000}
			resp, err := channel.Do(req)
			assert.Nil(t, err)
			assert.Equal(t, data, string(resp.Data))
		}(i)
	}
	wg.Wait()
}
