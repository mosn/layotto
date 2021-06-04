/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
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
	"github.com/layotto/layotto/pkg/grpc"
	runtimev1pb "github.com/layotto/layotto/spec/proto/runtime/v1"
	"mosn.io/mosn/pkg/wasm/abi/proxywasm010"
	"mosn.io/proxy-wasm-go-host/proxywasm"
)

type LayottoHandler struct {
	proxywasm010.DefaultImportsHandler
}

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
	}
	return "", proxywasm.WasmResultOk
}

type helloRequest struct {
	ServiceName string `json:"service_name"`
	Name        string `json:"name"`
}
