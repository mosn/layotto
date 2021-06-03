package transport_protocol

import (
	"fmt"

	"github.com/layotto/layotto/components/rpc"
	"mosn.io/api"
	"mosn.io/mosn/pkg/protocol/xprotocol"
	"mosn.io/mosn/pkg/protocol/xprotocol/dubbo"
	"mosn.io/pkg/buffer"
)

func init() {
	RegistProtocol("dubbo", newDubboProtocol())
}

func newDubboProtocol() TransportProtocol {
	return &dubboProtocol{XProtocol: xprotocol.GetProtocol(dubbo.ProtocolName)}
}

type dubboProtocol struct {
	fromFrame
	api.XProtocol
}

func (d *dubboProtocol) ToFrame(req *rpc.RPCRequest) api.XFrame {
	dubboReq := dubbo.NewRpcRequest(nil, buffer.NewIoBufferBytes(req.Data))
	req.Header.Range(func(key string, value string) bool {
		dubboReq.Header.Set(key, value)
		return true
	})
	return dubboReq
}

func (d *dubboProtocol) FromFrame(resp api.XRespFrame) (*rpc.RPCResponse, error) {
	if resp.GetStatusCode() != dubbo.RespStatusOK {
		return nil, fmt.Errorf("dubbo error code %d", resp.GetStatusCode())
	}

	return d.FromFrame(resp)
}
