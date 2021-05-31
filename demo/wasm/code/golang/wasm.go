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

package main

import (
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func main() {
	proxywasm.SetNewRootContext(newContext)
}

type rootContext struct {
	// You'd better embed the default root context
	// so that you don't need to reimplement all the methods by yourself.
	proxywasm.DefaultRootContext
}

func newContext(uint32) proxywasm.RootContext { return &rootContext{} }

// Override DefaultRootContext.
func (*rootContext) NewHttpContext(contextID uint32) proxywasm.HttpContext {
	return &myHttpContext{contextID: contextID}
}

type myHttpContext struct {
	// you must embed the default context so that you need not to re-implement all the methods by yourself
	proxywasm.DefaultHttpContext
	contextID uint32
}

// override
func (ctx *myHttpContext) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {
	proxywasm.LogInfo("golang wasm receive a http request")

	hs, err := proxywasm.GetHttpRequestHeaders()
	var name string
	for _, h := range hs {
		if h[0] == "Name" {
			name = h[1]
		}
	}

	result, err := proxywasm.CallForeignFunction("SayHello", []byte(`{"service_name":"helloworld","name":"`+name+`"}`))
	if err != nil {
		proxywasm.LogErrorf("call foreign func failed: %v", err)
	}
	proxywasm.SetHttpResponseBody(result)

	return types.ActionContinue
}

//export malloc
func proxyOnMemoryAllocate(size uint) *byte {
	buf := make([]byte, size)
	return &buf[0]
}
