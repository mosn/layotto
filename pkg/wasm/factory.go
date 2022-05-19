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
	"errors"
	"mosn.io/mosn/pkg/wasm"
	"mosn.io/pkg/utils"

	"mosn.io/api"
	v2 "mosn.io/mosn/pkg/config/v2"
	"mosn.io/mosn/pkg/log"
	"mosn.io/mosn/pkg/types"
	"mosn.io/mosn/pkg/wasm/abi"
)

const LayottoWasm = "Layotto"

func init() {
	api.RegisterStream(LayottoWasm, createProxyWasmFilterFactory)
}

// FilterConfigFactory contains multi wasm-plugin configs
// its pointer implement api.StreamFilterChainFactory
type FilterConfigFactory struct {
	LayottoHandler

	config        []*filterConfigItem // contains multi wasm config
	RootContextID int32

	// map[pluginName]*WasmPlugin
	plugins map[string]*WasmPlugin
	router  *Router
}

var _ api.StreamFilterChainFactory = &FilterConfigFactory{}

// Create a proxy factory for WasmFilter
func createProxyWasmFilterFactory(confs map[string]interface{}) (api.StreamFilterChainFactory, error) {
	factory := &FilterConfigFactory{
		config:        make([]*filterConfigItem, 0, len(confs)),
		RootContextID: 1,
		plugins:       make(map[string]*WasmPlugin),
		router:        &Router{routes: make(map[string]*Group)},
	}

	for configID, confIf := range confs {
		conf, ok := confIf.(map[string]interface{})
		if !ok {
			log.DefaultLogger.Errorf("[proxywasm][factory] createProxyWasmFilterFactory config not a map, configID: %s", configID)
			return nil, errors.New("config not a map")
		}
		err := factory.register(conf)
		if err != nil {
			return nil, err
		}
	}

	return factory, nil
}

// Create the FilterChain
func (f *FilterConfigFactory) CreateFilterChain(context context.Context, callbacks api.StreamFilterChainFactoryCallbacks) {
	filter := NewFilter(context, f)
	if filter == nil {
		return
	}

	callbacks.AddStreamReceiverFilter(filter, api.BeforeRoute)
	callbacks.AddStreamSenderFilter(filter, api.BeforeSend)
}

func (f *FilterConfigFactory) register(conf map[string]interface{}) error {
	config, err := parseFilterConfigItem(conf)
	if err != nil {
		log.DefaultLogger.Errorf("[proxywasm][factory] register fail to parse config, err: %v", err)
		return err
	}
	var pluginName string
	if config.FromWasmPlugin == "" {
		pluginName = utils.GenerateUUID()
		v2Config := v2.WasmPluginConfig{
			PluginName:  pluginName,
			VmConfig:    config.VmConfig,
			InstanceNum: config.InstanceNum,
		}
		err = wasm.GetWasmManager().AddOrUpdateWasm(v2Config)
		if err != nil {
			config.PluginName = pluginName
			addWatchFile(config, f)
			return nil
		}
		addWatchFile(config, f)
	} else {
		pluginName = config.FromWasmPlugin
	}
	config.PluginName = pluginName
	pw := wasm.GetWasmManager().GetWasmPluginWrapperByName(pluginName)
	if pw == nil {
		log.DefaultLogger.Errorf("[proxywasm][factory] register plugin not found")
		return errors.New("plugin not found")
	}
	config.VmConfig = pw.GetConfig().VmConfig
	f.config = append(filter(f.config, func(item *filterConfigItem) bool {
		return item.PluginName != config.PluginName
	}).([]*filterConfigItem), config)
	wasmPlugin := &WasmPlugin{
		pluginName:    config.PluginName,
		plugin:        pw.GetPlugin(),
		rootContextID: config.RootContextID,
		config:        config,
	}
	f.plugins[config.PluginName] = wasmPlugin
	pw.RegisterPluginHandler(f)
	return nil
}

// Get RootContext's ID
func (f *FilterConfigFactory) GetRootContextID() int32 {
	return f.RootContextID
}

// FilterConfigFactory implement types.WasmPluginHandler
// for `pw.RegisterPluginHandler(factory)`
var _ types.WasmPluginHandler = &FilterConfigFactory{}

// update config of FilterConfigFactory
func (f *FilterConfigFactory) OnConfigUpdate(config v2.WasmPluginConfig) {
	for _, plugin := range f.config {
		if plugin.PluginName == config.PluginName {
			plugin.InstanceNum = config.InstanceNum
			plugin.VmConfig = config.VmConfig
		}
	}
}

// Execute the plugin of FilterConfigFactory
func (f *FilterConfigFactory) OnPluginStart(plugin types.WasmPlugin) {
	plugin.Exec(func(instance types.WasmInstance) bool {
		wasmPlugin, ok := f.plugins[plugin.PluginName()]
		if !ok {
			log.DefaultLogger.Errorf("[proxywasm][factory] createProxyWasmFilterFactory fail to get wasm plugin, PluginName: %s",
				plugin.PluginName())
			return true
		}

		a := abi.GetABI(instance, AbiV2)
		a.SetABIImports(f)
		exports := a.GetABIExports().(Exports)
		f.LayottoHandler.Instance = instance

		instance.Lock(a)
		defer instance.Unlock()

		// get the ID of wasm, register route
		id, err := exports.ProxyGetID()
		if err != nil {
			log.DefaultLogger.Errorf("[proxywasm][factory] createProxyWasmFilterFactory fail to get wasm id, PluginName: %s, err: %v",
				plugin.PluginName(), err)
			return true
		}
		f.router.RegisterRoute(id, wasmPlugin)

		err = exports.ProxyOnContextCreate(f.RootContextID, 0)
		if err != nil {
			log.DefaultLogger.Errorf("[proxywasm][factory] OnPluginStart fail to create root context id, err: %v", err)
			return true
		}

		vmConfigSize := 0
		if vmConfigBytes := wasmPlugin.GetVmConfig(); vmConfigBytes != nil {
			vmConfigSize = vmConfigBytes.Len()
		}

		_, err = exports.ProxyOnVmStart(f.RootContextID, int32(vmConfigSize))
		if err != nil {
			log.DefaultLogger.Errorf("[proxywasm][factory] OnPluginStart fail to create root context id, err: %v", err)
			return true
		}

		pluginConfigSize := 0
		if pluginConfigBytes := wasmPlugin.GetPluginConfig(); pluginConfigBytes != nil {
			pluginConfigSize = pluginConfigBytes.Len()
		}

		_, err = exports.ProxyOnConfigure(f.RootContextID, int32(pluginConfigSize))
		if err != nil {
			log.DefaultLogger.Errorf("[proxywasm][factory] OnPluginStart fail to create root context id, err: %v", err)
			return true
		}

		return true
	})
}

// Destroy the plugin of FilterConfigFactory
func (f *FilterConfigFactory) OnPluginDestroy(types.WasmPlugin) {}
