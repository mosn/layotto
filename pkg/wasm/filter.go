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
	"sync"
	"sync/atomic"

	"mosn.io/api"
	"mosn.io/mosn/pkg/log"
	"mosn.io/mosn/pkg/types"
	"mosn.io/mosn/pkg/wasm"
	"mosn.io/mosn/pkg/wasm/abi"
	"mosn.io/mosn/pkg/wasm/abi/proxywasm010"
	"mosn.io/pkg/buffer"
	"mosn.io/proxy-wasm-go-host/common"
	"mosn.io/proxy-wasm-go-host/proxywasm"
)

type Filter struct {
	LayottoHandler

	ctx context.Context

	factory *FilterConfigFactory

	router  Router
	plugins []*WasmPlugin

	receiverFilterHandler api.StreamReceiverFilterHandler
	senderFilterHandler   api.StreamSenderFilterHandler

	destroyOnce sync.Once

	buffer api.IoBuffer
}

type WasmPlugin struct {
	pluginName string             // 单个wasm文件的name
	plugin     types.WasmPlugin   // 单个wasm文件，包括多个instance
	instance   types.WasmInstance // 一个instance
	abi        types.ABI          // 单个wasm文件对应的abi
	exports    proxywasm.Exports  // 单个wasm文件导出的方法

	rootContextID int32
	contextID     int32
}

var contextIDGenerator int32

func newContextID(rootContextID int32) int32 {
	for {
		id := atomic.AddInt32(&contextIDGenerator, 1)
		if id != rootContextID {
			return id
		}
	}
}

func NewFilter(ctx context.Context, factory *FilterConfigFactory) *Filter {
	filter := &Filter{
		ctx:     ctx,
		factory: factory,
		buffer:  buffer.NewIoBuffer(100),
	}

	configs := factory.config
	plugins := make([]*WasmPlugin, 0, len(configs))
	for _, pluginConfig := range configs {
		pluginWrapper := wasm.GetWasmManager().GetWasmPluginWrapperByName(pluginConfig.PluginName)
		if pluginWrapper == nil {
			log.DefaultLogger.Errorf("[proxywasm][filter] NewFilter wasm plugin not exists, plugin name: %v", pluginConfig.PluginName)
			return nil
		}

		plugin := pluginWrapper.GetPlugin()
		instance := plugin.GetInstance()

		pluginABI := abi.GetABI(instance, AbiV2)
		if pluginABI == nil {
			log.DefaultLogger.Errorf("[proxywasm][filter] NewFilter abi not found in instance")
			plugin.ReleaseInstance(instance)
			return nil
		}

		// TODO: 确定这里做了什么事，调用顺序有影响吗
		pluginABI.SetABIImports(filter)

		exports := pluginABI.GetABIExports().(proxywasm.Exports)
		if exports == nil {
			log.DefaultLogger.Errorf("[proxywasm][filter] NewFilter fail to get exports part from abi")
			plugin.ReleaseInstance(instance)

			return nil
		}

		contextID := newContextID(pluginConfig.RootContextID)

		err := exports.ProxyOnContextCreate(contextID, pluginConfig.RootContextID)
		if err != nil {
			log.DefaultLogger.Errorf("[proxywasm][filter] NewFilter fail to create context id: %v, rootContextID: %v, err: %v",
				contextID, pluginConfig.RootContextID, err)
			return nil
		}

		instance.Lock(pluginABI)
		defer instance.Unlock()

		wasmPlugin := &WasmPlugin{
			pluginName:    pluginConfig.PluginName,
			plugin:        plugin,
			instance:      instance,
			abi:           pluginABI,
			exports:       exports,
			rootContextID: pluginConfig.RootContextID,
			contextID:     contextID,
		}
		plugins = append(plugins, wasmPlugin)

		// TODO: 获取id，注册路由
		{
			exports := pluginABI.GetABIExports().(Exports)
			if exports == nil {
				log.DefaultLogger.Errorf("[proxywasm][filter] NewFilter fail to get exports part from abi")
				plugin.ReleaseInstance(instance)

				return nil
			}

			contextID := newContextID(pluginConfig.RootContextID)

			id, err := exports.ProxyGetID()
			if err != nil {
				log.DefaultLogger.Errorf("[proxywasm][filter] NewFilter fail to create context id: %v, rootContextID: %v, err: %v",
					contextID, pluginConfig.RootContextID, err)
				return nil
			}

			RegisterRoute(id, wasmPlugin)
		}
	}
	filter.plugins = plugins

	// TODO: 确定这个的作用
	//filter.LayottoHandler.Instance = instance

	return filter
}

func (f *Filter) OnDestroy() {
	f.destroyOnce.Do(func() {
		for _, plugin := range f.plugins {
			plugin.instance.Lock(plugin.abi)

			_, err := plugin.exports.ProxyOnDone(plugin.contextID)
			if err != nil {
				log.DefaultLogger.Errorf("[proxywasm][filter] OnDestroy fail to call ProxyOnDone, err: %v", err)
			}

			err = plugin.exports.ProxyOnDelete(plugin.contextID)
			if err != nil {
				log.DefaultLogger.Errorf("[proxywasm][filter] OnDestroy fail to call ProxyOnDelete, err: %v", err)
			}

			plugin.instance.Unlock()
			plugin.plugin.ReleaseInstance(plugin.instance)
		}
	})
}

func (f *Filter) SetReceiveFilterHandler(handler api.StreamReceiverFilterHandler) {
	f.receiverFilterHandler = handler
}

func (f *Filter) SetSenderFilterHandler(handler api.StreamSenderFilterHandler) {
	f.senderFilterHandler = handler
}

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

func (f *Filter) OnReceive(ctx context.Context, headers api.HeaderMap, buf buffer.IoBuffer, trailers api.HeaderMap) api.StreamFilterStatus {
	// TODO: should diapatch by header
	id, ok := headers.Get("id")
	if !ok {
		log.DefaultLogger.Errorf("[proxywasm][filter] OnReceive call ProxyOnRequestHeaders no id in headers")
		return api.StreamFilterStop
	}

	plugin, err := GetRandomPluginByID(id)
	if err != nil {
		log.DefaultLogger.Errorf("[proxywasm][filter] OnReceive call ProxyOnRequestHeaders id, err: %v", err)
		return api.StreamFilterStop
	}

	plugin.instance.Lock(plugin.abi)
	defer plugin.instance.Unlock()

	endOfStream := 1
	if (buf != nil && buf.Len() > 0) || trailers != nil {
		endOfStream = 0
	}

	action, err := plugin.exports.ProxyOnRequestHeaders(plugin.contextID, int32(headerMapSize(headers)), int32(endOfStream))
	if err != nil || action != proxywasm.ActionContinue {
		log.DefaultLogger.Errorf("[proxywasm][filter] OnReceive call ProxyOnRequestHeaders err: %v", err)
		return api.StreamFilterStop
	}

	endOfStream = 1
	if trailers != nil {
		endOfStream = 0
	}

	if buf != nil && buf.Len() > 0 {
		action, err = plugin.exports.ProxyOnRequestBody(plugin.contextID, int32(buf.Len()), int32(endOfStream))
		if err != nil || action != proxywasm.ActionContinue {
			log.DefaultLogger.Errorf("[proxywasm][filter] OnReceive call ProxyOnRequestBody err: %v", err)
			return api.StreamFilterStop
		}
	}

	if trailers != nil {
		action, err = plugin.exports.ProxyOnRequestTrailers(plugin.contextID, int32(headerMapSize(trailers)))
		if err != nil || action != proxywasm.ActionContinue {
			log.DefaultLogger.Errorf("[proxywasm][filter] OnReceive call ProxyOnRequestTrailers err: %v", err)
			return api.StreamFilterStop
		}
	}

	return api.StreamFilterContinue
}

func (f *Filter) Append(ctx context.Context, headers api.HeaderMap, buf buffer.IoBuffer, trailers api.HeaderMap) api.StreamFilterStatus {
	f.senderFilterHandler.SetResponseData(f.buffer)
	return api.StreamFilterContinue
}

// TODO: get the plugin content ID corresponding to the caller wasm plugin
func (f *Filter) GetRootContextID() int32 {
	return f.factory.RootContextID
}

// TODO: get the plugin vm config corresponding to the caller wasm plugin
func (f *Filter) GetVmConfig() common.IoBuffer {
	return f.factory.GetVmConfig()
}

// TODO: get the plugin config corresponding to the caller wasm plugin
func (f *Filter) GetPluginConfig() common.IoBuffer {
	return f.factory.GetPluginConfig()
}

func (f *Filter) GetHttpRequestHeader() common.HeaderMap {
	if f.receiverFilterHandler == nil {
		return nil
	}

	return &proxywasm010.HeaderMapWrapper{HeaderMap: f.receiverFilterHandler.GetRequestHeaders()}
}

func (f *Filter) GetHttpRequestBody() common.IoBuffer {
	if f.receiverFilterHandler == nil {
		return nil
	}

	return &proxywasm010.IoBufferWrapper{IoBuffer: f.receiverFilterHandler.GetRequestData()}
}

func (f *Filter) GetHttpRequestTrailer() common.HeaderMap {
	if f.receiverFilterHandler == nil {
		return nil
	}

	return &proxywasm010.HeaderMapWrapper{HeaderMap: f.receiverFilterHandler.GetRequestTrailers()}
}

func (f *Filter) GetHttpResponseHeader() common.HeaderMap {
	if f.senderFilterHandler == nil {
		return nil
	}

	return &proxywasm010.HeaderMapWrapper{HeaderMap: f.senderFilterHandler.GetResponseHeaders()}
}

func (f *Filter) GetHttpResponseBody() common.IoBuffer {
	if f.senderFilterHandler == nil {
		return nil
	}

	return &proxywasm010.IoBufferWrapper{IoBuffer: f.buffer}
}

func (f *Filter) GetHttpResponseTrailer() common.HeaderMap {
	if f.senderFilterHandler == nil {
		return nil
	}

	return &proxywasm010.HeaderMapWrapper{HeaderMap: f.senderFilterHandler.GetResponseTrailers()}
}
