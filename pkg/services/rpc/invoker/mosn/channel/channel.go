package channel

import (
	"errors"
	"fmt"
	"net"

	"github.com/layotto/layotto/pkg/services/rpc"
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
