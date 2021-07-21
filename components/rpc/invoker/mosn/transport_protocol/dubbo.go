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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		return nil, status.Errorf(codes.Unavailable, "dubbo error code %d", resp.GetStatusCode())
	}

	return d.fromFrame.FromFrame(resp)
}
