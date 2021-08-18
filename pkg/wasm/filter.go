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
	pluginName string
	plugin     types.WasmPlugin
	instance   types.WasmInstance
	abi        types.ABI
	exports    proxywasm.Exports

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
	configs := factory.config

	filter := &Filter{
		ctx:     ctx,
		factory: factory,
		router: Router{
			routes: make(map[string]Group),
		},
		buffer: buffer.NewIoBuffer(100),
	}

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
		pluginABI.SetABIImports(filter)

		exports := pluginABI.GetABIExports().(Exports)
		if exports == nil {
			log.DefaultLogger.Errorf("[proxywasm][filter] NewFilter fail to get exports part from abi")
			plugin.ReleaseInstance(instance)
			return nil
		}

		contextID := newContextID(pluginConfig.RootContextID)
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

		instance.Lock(pluginABI)
		defer instance.Unlock()

		err := exports.ProxyOnContextCreate(contextID, pluginConfig.RootContextID)
		if err != nil {
			log.DefaultLogger.Errorf("[proxywasm][filter] NewFilter fail to create context id: %v, rootContextID: %v, err: %v",
				contextID, pluginConfig.RootContextID, err)
			return nil
		}

		// TODO: 获取id，注册路由
		id, err := exports.ProxyGetID()
		if err != nil {
			log.DefaultLogger.Errorf("[proxywasm][filter] NewFilter fail to get id context id: %v, rootContextID: %v, err: %v",
				contextID, pluginConfig.RootContextID, err)
			return nil
		}
		filter.router.RegisterRoute(id, wasmPlugin)
	}
	filter.plugins = plugins

	// TODO: 确定这个的作用
	if len(plugins) > 0 {
		filter.LayottoHandler.Instance = plugins[0].instance
	}

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

	plugin, err := f.router.GetRandomPluginByID(id)
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
