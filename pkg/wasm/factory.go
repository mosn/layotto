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

var factory = &FilterConfigFactory{
	config:        make([]*filterConfigItem, 0),
	RootContextID: 1,
	plugins:       make(map[string]*WasmPlugin),
	router:        &Router{routes: make(map[string]*Group)},
}

var _ api.StreamFilterChainFactory = &FilterConfigFactory{}

func GetFactory() *FilterConfigFactory {
	return factory
}

// Create a proxy factory for WasmFilter
func createProxyWasmFilterFactory(confs map[string]interface{}) (api.StreamFilterChainFactory, error) {
	for configID, confIf := range confs {
		conf, ok := confIf.(map[string]interface{})
		if !ok {
			log.DefaultLogger.Errorf("[proxywasm][factory] createProxyWasmFilterFactory config not a map, configID: %s", configID)
			return nil, errors.New("config not a map")
		}
		err := factory.Install(conf)
		if err != nil {
			log.DefaultLogger.Errorf("[proxywasm][factory] createProxyWasmFilterFactory install error: %v", err)
			return nil, err
		}
	}

	return factory, nil
}

// Create the FilterChain
var filterChain *Filter

func (f *FilterConfigFactory) CreateFilterChain(context context.Context, callbacks api.StreamFilterChainFactoryCallbacks) {
	filterChain = NewFilter(context, f)
	if filterChain == nil {
		return
	}

	callbacks.AddStreamReceiverFilter(filterChain, api.BeforeRoute)
	callbacks.AddStreamSenderFilter(filterChain, api.BeforeSend)
}

func (f *FilterConfigFactory) IsRegister(id string) bool {
	plugin, err := f.router.GetRandomPluginByID(id)
	return err == nil && plugin != nil
}

func (f *FilterConfigFactory) Install(conf map[string]interface{}) error {
	config, err := parseFilterConfigItem(conf)
	if err != nil {
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

func (f *FilterConfigFactory) UpdateInstanceNum(id string, instanceNum int) error {
	wasmPlugin, _ := f.router.GetRandomPluginByID(id)
	if wasmPlugin == nil {
		log.DefaultLogger.Errorf("[proxywasm][factory] GetRandomPluginByID id not registered, id: %s", id)
		return errors.New(id + " is not registered")
	}

	var config *filterConfigItem
	for _, item := range f.config {
		if item.PluginName == wasmPlugin.pluginName {
			config = item
			break
		}
	}
	if config == nil {
		return errors.New("can't find config for " + id)
	}

	if config.InstanceNum == instanceNum {
		return nil
	}

	config.InstanceNum = instanceNum
	v2Config := v2.WasmPluginConfig{
		PluginName:  config.PluginName,
		VmConfig:    config.VmConfig,
		InstanceNum: config.InstanceNum,
	}
	err := wasm.GetWasmManager().AddOrUpdateWasm(v2Config)
	if err != nil {
		return err
	}
	pw := wasm.GetWasmManager().GetWasmPluginWrapperByName(config.PluginName)
	if pw == nil {
		return errors.New("plugin not found")
	}
	f.plugins[config.PluginName] = &WasmPlugin{
		pluginName:    config.PluginName,
		plugin:        pw.GetPlugin(),
		rootContextID: config.RootContextID,
		config:        config,
	}
	pw.RegisterPluginHandler(f)
	return nil
}

func (f *FilterConfigFactory) UnInstall(id string) error {
	wasmPlugin, _ := f.router.GetRandomPluginByID(id)
	if wasmPlugin == nil {
		log.DefaultLogger.Errorf("[proxywasm][factory] GetRandomPluginByID id not registered, id: %s", id)
		return errors.New(id + " is not registered")
	}
	err := wasm.GetWasmManager().UninstallWasmPluginByName(wasmPlugin.pluginName)
	if err != nil {
		return err
	}

	if filterChain != nil && filterChain.pluginUsed != nil && filterChain.pluginUsed.pluginName == wasmPlugin.pluginName {
		err = filterChain.releaseUsedInstance()
		if err != nil {
			return err
		}
	}

	f.config = filter(f.config, func(item *filterConfigItem) bool {
		return item.PluginName != wasmPlugin.pluginName
	}).([]*filterConfigItem)
	delete(f.plugins, wasmPlugin.pluginName)
	removeWatchFile(wasmPlugin.config)
	f.router.RemoveRoute(id)
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
