package channel

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/layotto/components/rpc"
	"mosn.io/mosn/pkg/server"
)

var (
	ErrTimeout    = errors.New("request timeout")
	ErrConnClosed = errors.New("connection closed by mosn")

	acceptFunc = func(conn net.Conn, listener string) error {
		srv := server.GetServer()
		lis := srv.Handler().FindListenerByName(listener)
		if lis == nil {
			return errors.New("[rpc]invalid listener name")
		}
		lis.GetListenerCallbacks().OnAccept(conn, false, nil, nil, nil)
		return nil
	}

	registry = map[string]func(config ChannelConfig) (rpc.Channel, error){}
)

type ChannelConfig struct {
	Protocol string
	Listener string
}

func GetChannel(config ChannelConfig) (rpc.Channel, error) {
	c, ok := registry[config.Protocol]
	if !ok {
		return nil, fmt.Errorf("channel %s not found", config.Protocol)
	}
	return c(config)
}

func RegistChannel(proto string, f func(config ChannelConfig) (rpc.Channel, error)) {
	registry[proto] = f
}

type fakeTcpConn struct {
	c net.Conn
}

func (t *fakeTcpConn) Read(b []byte) (n int, err error) {
	return t.c.Read(b)
}

func (t *fakeTcpConn) Write(b []byte) (n int, err error) {
	return t.c.Write(b)
}

func (t *fakeTcpConn) Close() error {
	return t.c.Close()
}

func (t *fakeTcpConn) LocalAddr() net.Addr {
	return &net.TCPAddr{}
}

func (t *fakeTcpConn) RemoteAddr() net.Addr {
	return &net.TCPAddr{}
}

func (t *fakeTcpConn) SetDeadline(time time.Time) error {
	return t.c.SetDeadline(time)
}

func (t *fakeTcpConn) SetReadDeadline(time time.Time) error {
	return t.c.SetReadDeadline(time)
}

func (t fakeTcpConn) SetWriteDeadline(time time.Time) error {
	return t.c.SetWriteDeadline(time)
}
