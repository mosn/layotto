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
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"

	"mosn.io/mosn/pkg/wasm/abi"
	"mosn.io/pkg/variable"
	proxywasm "mosn.io/proxy-wasm-go-host/proxywasm/v1"

	"mosn.io/api"
	"mosn.io/mosn/pkg/log"
	"mosn.io/mosn/pkg/types"
	"mosn.io/mosn/pkg/wasm/abi/proxywasm010"
	"mosn.io/pkg/buffer"
	"mosn.io/proxy-wasm-go-host/proxywasm/common"
)

type Filter struct {
	LayottoHandler

	ctx     context.Context
	factory *FilterConfigFactory

	router  *Router
	plugins map[string]*WasmPlugin

	contextID  int32
	pluginUsed *WasmPlugin
	instance   types.WasmInstance
	abi        types.ABI
	exports    Exports

	receiverFilterHandler api.StreamReceiverFilterHandler
	senderFilterHandler   api.StreamSenderFilterHandler

	destroyOnce sync.Once

	requestBuffer  api.IoBuffer
	responseBuffer api.IoBuffer
}

type WasmPlugin struct {
	pluginName string
	plugin     types.WasmPlugin

	// useless for now
	rootContextID     int32
	config            *filterConfigItem
	vmConfigBytes     buffer.IoBuffer
	pluginConfigBytes buffer.IoBuffer
}

// Get the VmConfig of WasmPlugin
func (p *WasmPlugin) GetVmConfig() common.IoBuffer {
	if p.vmConfigBytes != nil {
		return p.vmConfigBytes
	}

	vmConfig := p.plugin.GetVmConfig()

	typeOf := reflect.TypeOf(vmConfig)
	valueOf := reflect.ValueOf(&vmConfig).Elem()
	if typeOf.Kind() != reflect.Struct || typeOf.NumField() == 0 {
		return nil
	}

	m := make(map[string]string)
	for i := 0; i < typeOf.NumField(); i++ {
		m[typeOf.Field(i).Name] = fmt.Sprintf("%v", valueOf.Field(i).Interface())
	}

	b := common.EncodeMap(m)
	if b == nil {
		return nil
	}

	p.vmConfigBytes = buffer.NewIoBufferBytes(b)
	return p.vmConfigBytes
}

// Get the plugin config of WasmPlugin
func (p *WasmPlugin) GetPluginConfig() common.IoBuffer {
	if p.pluginConfigBytes != nil {
		return p.pluginConfigBytes
	}

	b := common.EncodeMap(p.config.UserData)
	if b == nil {
		return nil
	}

	p.pluginConfigBytes = buffer.NewIoBufferBytes(b)
	return p.pluginConfigBytes
}

var contextIDGenerator int32

// new context's ID
func newContextID(rootContextID int32) int32 {
	for {
		id := atomic.AddInt32(&contextIDGenerator, 1)
		if id != rootContextID {
			return id
		}
	}
}

// NewFilter create the filter for a request
func NewFilter(ctx context.Context, factory *FilterConfigFactory) *Filter {
	filter := &Filter{
		ctx:     ctx,
		factory: factory,

		contextID:      newContextID(factory.RootContextID),
		router:         factory.router,
		plugins:        factory.plugins,
		requestBuffer:  buffer.NewIoBuffer(100),
		responseBuffer: buffer.NewIoBuffer(100),
	}

	return filter
}

func (f *Filter) releaseUsedInstance() error {
	if f.pluginUsed == nil || f.instance == nil {
		return nil
	}
	plugin := f.pluginUsed
	f.instance.Lock(f.abi)

	_, err := f.exports.ProxyOnDone(f.contextID)
	if err != nil {
		log.DefaultLogger.Errorf("[proxywasm][filter] releaseUsedInstance fail to call ProxyOnDone, err: %v", err)
		return err
	}

	err = f.exports.ProxyOnDelete(f.contextID)
	if err != nil {
		log.DefaultLogger.Errorf("[proxywasm][filter] releaseUsedInstance fail to call ProxyOnDelete, err: %v", err)
		return err
	}

	f.instance.Unlock()
	plugin.plugin.ReleaseInstance(f.instance)

	f.instance = nil
	f.pluginUsed = nil
	f.exports = nil

	return nil
}

// Destruction of filters
func (f *Filter) OnDestroy() {
	f.destroyOnce.Do(func() {
		_ = f.releaseUsedInstance()
	})
}

// Set ReceiveFilterHandler of filter
func (f *Filter) SetReceiveFilterHandler(handler api.StreamReceiverFilterHandler) {
	f.receiverFilterHandler = handler
}

// Set SenderFilterHandler of filter
func (f *Filter) SetSenderFilterHandler(handler api.StreamSenderFilterHandler) {
	f.senderFilterHandler = handler
}

// Calculate the size of headerMap
func headerMapSize(headers api.HeaderMap) int {
	size := 0

	if headers != nil {
		headers.Range(func(key, value string) bool {
			size++
			return true
		})
	}

	return size
}

// Reset the filter when receiving then return StreamFilter status
func (f *Filter) OnReceive(ctx context.Context, headers api.HeaderMap, buf buffer.IoBuffer, trailers api.HeaderMap) api.StreamFilterStatus {
	id, ok := headers.Get("id")
	if !ok {
		log.DefaultLogger.Errorf("[proxywasm][filter] OnReceive call ProxyOnRequestHeaders no id in headers")
		return api.StreamFilterStop
	}

	wasmPlugin, err := f.router.GetRandomPluginByID(id)
	if err != nil {
		log.DefaultLogger.Errorf("[proxywasm][filter] OnReceive call ProxyOnRequestHeaders id, err: %v", err)
		return api.StreamFilterStop
	}
	f.pluginUsed = wasmPlugin

	plugin := wasmPlugin.plugin
	instance := plugin.GetInstance()
	f.instance = instance
	f.LayottoHandler.Instance = instance

	pluginABI := abi.GetABI(instance, AbiV2)
	if pluginABI == nil {
		log.DefaultLogger.Errorf("[proxywasm][filter] OnReceive fail to get instance abi")
		plugin.ReleaseInstance(instance)
		return api.StreamFilterStop
	}
	pluginABI.SetABIImports(f)
	exports := pluginABI.GetABIExports().(Exports)
	f.exports = exports

	instance.Lock(pluginABI)
	defer instance.Unlock()

	err = exports.ProxyOnContextCreate(f.contextID, wasmPlugin.rootContextID)
	if err != nil {
		log.DefaultLogger.Errorf("[proxywasm][filter] NewFilter fail to create context id: %v, rootContextID: %v, err: %v",
			f.contextID, wasmPlugin.rootContextID, err)
		return api.StreamFilterStop
	}

	endOfStream := 1
	if (buf != nil && buf.Len() > 0) || trailers != nil {
		endOfStream = 0
	}

	action, err := exports.ProxyOnRequestHeaders(f.contextID, int32(headerMapSize(headers)), int32(endOfStream))
	if err != nil || action != proxywasm.ActionContinue {
		log.DefaultLogger.Errorf("[proxywasm][filter] OnReceive call ProxyOnRequestHeaders err: %v", err)
		return api.StreamFilterStop
	}

	endOfStream = 1
	if trailers != nil {
		endOfStream = 0
	}

	if buf == nil {
		arg, _ := variable.GetString(ctx, types.VarHttpRequestArg)
		f.requestBuffer = buffer.NewIoBufferString(arg)
	} else {
		f.requestBuffer = buf
	}

	if f.requestBuffer != nil && f.requestBuffer.Len() > 0 {
		action, err = exports.ProxyOnRequestBody(f.contextID, int32(f.requestBuffer.Len()), int32(endOfStream))
		if err != nil || action != proxywasm.ActionContinue {
			log.DefaultLogger.Errorf("[proxywasm][filter] OnReceive call ProxyOnRequestBody err: %v", err)
			return api.StreamFilterStop
		}
	}

	if trailers != nil {
		action, err = exports.ProxyOnRequestTrailers(f.contextID, int32(headerMapSize(trailers)))
		if err != nil || action != proxywasm.ActionContinue {
			log.DefaultLogger.Errorf("[proxywasm][filter] OnReceive call ProxyOnRequestTrailers err: %v", err)
			return api.StreamFilterStop
		}
	}

	return api.StreamFilterContinue
}

// Append ResponseData of filter
func (f *Filter) Append(ctx context.Context, headers api.HeaderMap, buf buffer.IoBuffer, trailers api.HeaderMap) api.StreamFilterStatus {
	f.senderFilterHandler.SetResponseData(f.responseBuffer)
	return api.StreamFilterContinue
}

// Get RootContext ID of filter's FilterConfigFactory
func (f *Filter) GetRootContextID() int32 {
	return f.factory.RootContextID
}

// Get the used WasmPlugin VmConfig of filter
func (f *Filter) GetVmConfig() common.IoBuffer {
	return f.pluginUsed.GetVmConfig()
}

// Get the used WasmPlugin config of filter
func (f *Filter) GetPluginConfig() common.IoBuffer {
	return f.pluginUsed.GetPluginConfig()
}

// Get the HttpRequest header of proxy-wasm
func (f *Filter) GetHttpRequestHeader() common.HeaderMap {
	if f.receiverFilterHandler == nil {
		return nil
	}

	return &proxywasm010.HeaderMapWrapper{HeaderMap: f.receiverFilterHandler.GetRequestHeaders()}
}

// Get the HttpRequest body of proxy-wasm
func (f *Filter) GetHttpRequestBody() common.IoBuffer {
	if f.receiverFilterHandler == nil {
		return nil
	}

	return &proxywasm010.IoBufferWrapper{IoBuffer: f.requestBuffer}
}

// Get the HttpRequest trailer of proxy-wasm
func (f *Filter) GetHttpRequestTrailer() common.HeaderMap {
	if f.receiverFilterHandler == nil {
		return nil
	}

	return &proxywasm010.HeaderMapWrapper{HeaderMap: f.receiverFilterHandler.GetRequestTrailers()}
}

// Get the HttpResponse header of proxy-wasm
func (f *Filter) GetHttpResponseHeader() common.HeaderMap {
	if f.senderFilterHandler == nil {
		return nil
	}

	return &proxywasm010.HeaderMapWrapper{HeaderMap: f.senderFilterHandler.GetResponseHeaders()}
}

// Get the HttpResponse body of proxy-wasm
func (f *Filter) GetHttpResponseBody() common.IoBuffer {
	if f.senderFilterHandler == nil {
		return nil
	}

	return &proxywasm010.IoBufferWrapper{IoBuffer: f.responseBuffer}
}

// Get the HttpResponse trailer of proxy-wasm
func (f *Filter) GetHttpResponseTrailer() common.HeaderMap {
	if f.senderFilterHandler == nil {
		return nil
	}

	return &proxywasm010.HeaderMapWrapper{HeaderMap: f.senderFilterHandler.GetResponseTrailers()}
}
