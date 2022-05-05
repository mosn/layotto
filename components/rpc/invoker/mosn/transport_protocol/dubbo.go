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

	"mosn.io/api"
	"mosn.io/mosn/pkg/protocol/xprotocol/dubbo"
	"mosn.io/pkg/buffer"

	common "mosn.io/layotto/components/pkg/common"
	"mosn.io/layotto/components/rpc"
)

// init dubbo protocol
func init() {
	RegistProtocol("dubbo", newDubboProtocol())
}

// newDubboProtocol is create dubbo TransportProtocol
func newDubboProtocol() TransportProtocol {
	return &dubboProtocol{XProtocol: (&dubbo.XCodec{}).NewXProtocol(context.TODO())}
}

type dubboProtocol struct {
	fromFrame
	api.XProtocol
}

func (d *dubboProtocol) Init(map[string]interface{}) error {
	return nil
}

// ToFrame is dubboProtocol transform
func (d *dubboProtocol) ToFrame(req *rpc.RPCRequest) api.XFrame {
	dubboReq := dubbo.NewRpcRequest(nil, buffer.NewIoBufferBytes(req.Data))
	req.Header.Range(func(key string, value string) bool {
		dubboReq.Header.Set(key, value)
		return true
	})
	return dubboReq
}

// FromFrame is dubboProtocol transform
func (d *dubboProtocol) FromFrame(resp api.XRespFrame) (*rpc.RPCResponse, error) {
	if resp.GetStatusCode() != dubbo.RespStatusOK {
		return nil, common.Errorf(common.UnavailebleCode, "dubbo error code %d", resp.GetStatusCode())
	}

	return d.fromFrame.FromFrame(resp)
}
