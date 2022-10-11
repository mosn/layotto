/*
 * Copyright 2021 Layotto Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package transport_protocol

import (
	"context"
	"errors"
	"reflect"
	"unsafe"

	"mosn.io/api"
	"mosn.io/mosn/pkg/protocol/xprotocol/bolt"
	"mosn.io/mosn/pkg/protocol/xprotocol/boltv2"
	"mosn.io/pkg/buffer"
	"mosn.io/pkg/header"

	"mosn.io/layotto/components/pkg/common"
	"mosn.io/layotto/components/rpc"
)

// init mosn bolt or boltv2 protocol
func init() {
	RegistProtocol("bolt", newBoltProtocol())
	RegistProtocol("boltv2", newBoltV2Protocol())
}

type boltCommon struct {
	className string
	fromFrame
}

// Init is init boltCommon info
func (b *boltCommon) Init(conf map[string]interface{}) error {
	if len(conf) == 0 {
		return errors.New("missing bolt classname")
	}
	class, ok := conf["class"]
	if !ok {
		return errors.New("bolt need class")
	}
	classStr, ok := class.(string)
	if !ok {
		return errors.New("bolt class not string")
	}
	b.className = classStr
	return nil
}

// FromFrame is boltProtocol transform
func (b *boltCommon) FromFrame(resp api.XRespFrame) (*rpc.RPCResponse, error) {
	respCode := uint16(resp.GetStatusCode())
	rpcResp, err := b.fromFrame.FromFrame(resp)
	rpcResp.Success = false
	if err != nil {
		return nil, common.Errorf(common.InternalCode, "bolt error code %d, transform RPCResponse fail with err: %+v", respCode, err)
	}
	switch respCode {
	case bolt.ResponseStatusSuccess:
		rpcResp.Success = true
	case bolt.ResponseStatusServerDeserialException:
		rpcResp.Error = common.Errorf(common.InternalCode, "bolt error code %d, ServerDeserializeException", respCode)
	case bolt.ResponseStatusServerSerialException:
		rpcResp.Error = common.Errorf(common.InternalCode, "bolt error code %d, ServerSerializeException", respCode)
	case bolt.ResponseStatusCodecException:
		rpcResp.Error = common.Errorf(common.InternalCode, "bolt error code %d, CodecException", respCode)
	default:
		rpcResp.Error = common.Errorf(common.UnavailebleCode, "bolt error code %d", respCode)
	}
	return rpcResp, nil
}

// newBoltProtocol is create boltProtocol
func newBoltProtocol() TransportProtocol {
	return &boltProtocol{XProtocol: (&bolt.XCodec{}).NewXProtocol(context.TODO()), boltCommon: boltCommon{}}
}

// boltProtocol is one of TransportProtocol
type boltProtocol struct {
	boltCommon
	api.XProtocol
}

// ToFrame is boltProtocol transform
func (b *boltProtocol) ToFrame(req *rpc.RPCRequest) api.XFrame {
	buf := buffer.NewIoBufferBytes(req.Data)
	headerrLen := len(req.Header)
	boltreq := bolt.NewRpcRequest(0, nil, buf)
	boltreq.Class = b.className
	boltreq.Timeout = req.Timeout
	boltreq.BytesHeader = header.BytesHeader{
		Kvs:     make([]header.BytesKV, headerrLen),
		Changed: true,
	}

	i := 0
	req.Header.Range(func(key string, value string) bool {
		kv := &boltreq.BytesHeader.Kvs[i]
		kv.Key = s2b(key)
		kv.Value = s2b(value)
		i++
		return true
	})
	return boltreq
}

// newBoltV2Protocol is create boltV2Protocol
func newBoltV2Protocol() TransportProtocol {
	return &boltv2Protocol{XProtocol: (&bolt.XCodec{}).NewXProtocol(context.TODO()), boltCommon: boltCommon{}}
}

// boltv2Protocol is one of TransportProtocol
type boltv2Protocol struct {
	boltCommon
	api.XProtocol
}

// ToFrame is boltv2Protocol transform
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
	boltv2Req.Class = b.className
	boltv2Req.Timeout = req.Timeout

	req.Header.Range(func(key string, value string) bool {
		boltv2Req.BytesHeader.Set(key, value)
		return true
	})
	return boltv2Req
}

// s2b is convert string to byte slice
func s2b(s string) []byte {
	ps := (*reflect.StringHeader)(unsafe.Pointer(&s))
	b := reflect.SliceHeader{
		Data: ps.Data,
		Len:  ps.Len,
		Cap:  ps.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&b))
}
