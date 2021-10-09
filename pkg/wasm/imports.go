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

package wasm

import (
	"context"
	"encoding/json"

	"github.com/golang/protobuf/proto"
	"mosn.io/layotto/pkg/grpc"
	"mosn.io/layotto/spec/proto/runtime/v1"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
	"mosn.io/mosn/pkg/wasm/abi/proxywasm010"
	"mosn.io/proxy-wasm-go-host/common"
	"mosn.io/proxy-wasm-go-host/proxywasm"
)

// LayottoHandler implement proxywasm.ImportsHandler
type LayottoHandler struct {
	proxywasm010.DefaultImportsHandler

	IoBuffer common.IoBuffer
}

var _ proxywasm.ImportsHandler = &LayottoHandler{}

var Layotto grpc.API

func (d *LayottoHandler) CallForeignFunction(funcName string, param string) (string, proxywasm.WasmResult) {

	switch funcName {
	case "SayHello":
		isJson := false
		req := &runtimev1pb.SayHelloRequest{}
		err := proto.Unmarshal([]byte(param), req)
		if err != nil {
			jsonReq := &helloRequest{}
			err = json.Unmarshal([]byte(param), jsonReq)
			if err != nil {
				return "", proxywasm.WasmResultBadArgument
			}
			req.ServiceName = jsonReq.ServiceName
			req.Name = jsonReq.Name
			isJson = true
		}
		resp, err := Layotto.SayHello(context.Background(), req)
		if err != nil {
			return "", proxywasm.WasmResultInternalFailure
		}
		if isJson {
			return resp.Hello, proxywasm.WasmResultOk
		}

		b, err := proto.Marshal(resp)
		if err != nil {
			return "", proxywasm.WasmResultSerializationFailure
		}
		return string(b), proxywasm.WasmResultOk
	case "State":
		isJson := false
		req := &runtimev1pb.GetStateRequest{}
		if err := proto.Unmarshal([]byte(param), req); err != nil {
			jsonReq := &getStateRequest{}
			err = json.Unmarshal([]byte(param), jsonReq)
			if err != nil {
				return "", proxywasm.WasmResultBadArgument
			}
			req.Key = jsonReq.Key
			req.StoreName = jsonReq.StoreName
			req.Metadata = jsonReq.Metadata
			req.Consistency = jsonReq.Consistency
			isJson = true
		}
		resp, err := Layotto.GetState(context.Background(), req)
		if err != nil {
			return "", proxywasm.WasmResultInternalFailure
		}
		if isJson {
			return string(resp.Data), proxywasm.WasmResultOk
		}

		b, err := proto.Marshal(resp)
		if err != nil {
			return "", proxywasm.WasmResultSerializationFailure
		}
		return string(b), proxywasm.WasmResultOk
	}
	return "", proxywasm.WasmResultOk
}

type helloRequest struct {
	ServiceName string `json:"service_name"`
	Name        string `json:"name"`
}

func (d *LayottoHandler) GetFuncCallData() common.IoBuffer {
	if d.IoBuffer == nil {
		d.IoBuffer = common.NewIoBufferBytes(make([]byte, 0))
	}
	return d.IoBuffer
}

type getStateRequest struct {
	// Required. The name of state store.
	StoreName string `json:"store_name,omitempty"`
	// Required. The key of the desired state
	Key string `json:"key,omitempty"`
	// (optional) read consistency mode
	Consistency runtime.StateOptions_StateConsistency `json:"consistency,omitempty"`
	// (optional) The metadata which will be sent to state store components.
	Metadata map[string]string `json:"metadata,omitempty"`
}
