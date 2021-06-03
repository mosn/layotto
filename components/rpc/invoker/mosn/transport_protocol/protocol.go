package transport_protocol

import (
	"github.com/layotto/layotto/components/rpc"
	"mosn.io/api"
)

var protocolRegistry = map[string]TransportProtocol{}

// transport protocol support by mosn(bolt/boltv2...)
type TransportProtocol interface {
	api.Encoder
	api.Decoder
	ToFrame(*rpc.RPCRequest) api.XFrame
	FromFrame(api.XRespFrame) (*rpc.RPCResponse, error)
}

func GetProtocol(protocol string) TransportProtocol {
	return protocolRegistry[protocol]
}

func RegistProtocol(protocol string, proto TransportProtocol) {
	protocolRegistry[protocol] = proto
}

type fromFrame struct{}

func (f *fromFrame) FromFrame(resp api.XRespFrame) (*rpc.RPCResponse, error) {
	rpcResp := &rpc.RPCResponse{}
	resp.GetHeader().Range(func(Key, Value string) bool {
		if rpcResp.Header == nil {
			rpcResp.Header = make(map[string][]string)
		}
		rpcResp.Header[Key] = []string{Value}
		return true
	})

	if data := resp.GetData(); data != nil {
		rpcResp.Data = data.Bytes()
	}
	return rpcResp, nil
}
