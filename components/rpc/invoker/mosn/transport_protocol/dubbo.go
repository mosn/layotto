package transport_protocol

import (
	"fmt"

	"mosn.io/api"
	"mosn.io/layotto/components/rpc"
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

func (d *dubboProtocol) Init(map[string]interface{}) error {
	return nil
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

	return d.fromFrame.FromFrame(resp)
}
