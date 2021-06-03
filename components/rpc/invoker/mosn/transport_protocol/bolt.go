package transport_protocol

import (
	"fmt"

	"github.com/layotto/components/rpc"
	"mosn.io/api"
	"mosn.io/mosn/pkg/protocol/xprotocol"
	"mosn.io/mosn/pkg/protocol/xprotocol/bolt"
	"mosn.io/mosn/pkg/protocol/xprotocol/boltv2"
	"mosn.io/pkg/buffer"
)

var sofaRequestClassName = "com.alipay.sofa.rpc.core.request.SofaRequest"

func init() {
	RegistProtocol("bolt", newBoltProtocol())
	RegistProtocol("boltv2", newBoltV2Protocol())
}

type boltCommon struct {
	fromFrame
}

func (b *boltCommon) FromFrame(resp api.XRespFrame) (*rpc.RPCResponse, error) {
	if resp.GetStatusCode() != uint32(bolt.ResponseStatusSuccess) {
		return nil, fmt.Errorf("bolt error code %d", resp.GetStatusCode())
	}

	return b.fromFrame.FromFrame(resp)
}

func newBoltProtocol() TransportProtocol {
	return &boltProtocol{XProtocol: xprotocol.GetProtocol(bolt.ProtocolName), boltCommon: boltCommon{}}
}

type boltProtocol struct {
	boltCommon
	api.XProtocol
}

func (b *boltProtocol) ToFrame(req *rpc.RPCRequest) api.XFrame {
	buf := buffer.NewIoBufferBytes(req.Data)
	boltreq := bolt.NewRpcRequest(0, nil, buf)
	boltreq.Class = sofaRequestClassName
	boltreq.Timeout = req.Timeout

	req.Header.Range(func(key string, value string) bool {
		boltreq.Header.Set(key, value)
		return true
	})
	return boltreq
}

func newBoltV2Protocol() TransportProtocol {
	return &boltv2Protocol{XProtocol: xprotocol.GetProtocol(boltv2.ProtocolName), boltCommon: boltCommon{}}
}

type boltv2Protocol struct {
	boltCommon
	api.XProtocol
}

func (b *boltv2Protocol) ToFrame(req *rpc.RPCRequest) api.XFrame {
	boltv2Req := &boltv2.Request{
		RequestHeader: boltv2.RequestHeader{
			Version1: boltv2.ProtocolVersion,
			RequestHeader: bolt.RequestHeader{
				Protocol: boltv2.ProtocolCode,
				CmdType:  bolt.CmdTypeRequest,
				CmdCode:  bolt.CmdCodeRpcRequest,
				Version:  boltv2.ProtocolVersion,
				Codec:    bolt.Hessian2Serialize,
				Timeout:  req.Timeout,
			},
		},
	}

	buf := buffer.NewIoBufferBytes(req.Data)
	boltv2Req.SetData(buf)
	boltv2Req.Class = sofaRequestClassName
	boltv2Req.Timeout = req.Timeout

	req.Header.Range(func(key string, value string) bool {
		boltv2Req.Header.Set(key, value)
		return true
	})
	return boltv2Req
}
