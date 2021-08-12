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
	"mosn.io/api"
	"mosn.io/layotto/components/rpc"
	"mosn.io/mosn/pkg/protocol/xprotocol/bolt"
)

var protocolRegistry = map[string]TransportProtocol{}

// transport protocol support by mosn(bolt/boltv2...)
type TransportProtocol interface {
	Init(map[string]interface{}) error
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
	if boltResp, ok := resp.(*bolt.Response); ok {
		rpcResp.Header = make(map[string][]string, len(boltResp.Header.Kvs))
	}
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
