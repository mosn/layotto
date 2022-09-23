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

package main

import (
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func main() {
	proxywasm.SetVMContext(&vmContext{})
}

type vmContext struct {
	// Embed the default VM context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultVMContext
}

// Override types.DefaultVMContext.
func (*vmContext) NewPluginContext(contextID uint32) types.PluginContext {
	return &pluginContext{}
}

type pluginContext struct {
	// Embed the default plugin context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultPluginContext
}

// Override types.DefaultPluginContext.
func (*pluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	return &httpHeaders{contextID: contextID}
}

type httpHeaders struct {
	// Embed the default http context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultHttpContext
	contextID uint32
}

// Override types.DefaultHttpContext.
func (ctx *httpHeaders) OnHttpRequestBody(bodySize int, endOfStream bool) types.Action {
	//1. get request body
	body, err := proxywasm.GetHttpRequestBody(0, bodySize)
	if err != nil {
		proxywasm.LogErrorf("GetHttpRequestBody failed: %v", err)
		return types.ActionPause
	}
	bookName := string(body)

	//2. get request state from redis by specific key through ABI
	inventories, err := proxywasm.GetState("state_demo", bookName)
	if err != nil {
		proxywasm.LogErrorf("GetState failed: %v", err)
		return types.ActionPause
	}

	//3. return result
	proxywasm.AppendHttpResponseBody([]byte(inventories))
	return types.ActionContinue
}

// Override types.DefaultHttpContext.
func (ctx *httpHeaders) OnHttpStreamDone() {
	proxywasm.LogInfof("%d finished", ctx.contextID)
}

const ID = "id_2"

// DO NOT MODIFY THE FOLLOWING FUNCTIONS!
//
//export proxy_get_id
func GetID() {
	_ = ID[len(ID)-1]
	proxywasm.SetCallData([]byte(ID))
}
