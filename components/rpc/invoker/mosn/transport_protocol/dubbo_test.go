package transport_protocol

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"mosn.io/mosn/pkg/protocol/xprotocol/dubbo"

	"mosn.io/layotto/components/rpc"
)

func TestDubbo(t *testing.T) {
	d := dubboProtocol{}

	req := &rpc.RPCRequest{
		Data: []byte("hello"),
	}

	frame := d.ToFrame(req)
	dubboFrame := frame.(*dubbo.Frame)
	isEvent := (dubboFrame.Flag & (1 << 5)) != 0
	assert.False(t, isEvent)
	isTwoWay := (dubboFrame.Flag & (1 << 6)) != 0
	assert.True(t, isTwoWay)
	directionBool := (dubboFrame.Flag & (1 << 7))
	assert.NotEqual(t, directionBool, 0)
}
