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

	"mosn.io/api"
	v2 "mosn.io/mosn/pkg/config/v2"
	"mosn.io/mosn/pkg/log"
	"mosn.io/mosn/pkg/types"
	"mosn.io/mosn/pkg/wasm"
	"mosn.io/mosn/pkg/wasm/abi"
	"mosn.io/pkg/utils"
)

const LayottoWasm = "Layotto"

func init() {
	api.RegisterStream(LayottoWasm, createProxyWasmFilterFactory)
}

// FilterConfigFactory contains multi wasm-plugin configs
// its pointer implemente api.StreamFilterChainFactory
type FilterConfigFactory struct {
	LayottoHandler

	config        []*filterConfigItem // contains multi wasm config
	RootContextID int32

	plugins []*WasmPlugin
	router  *Router
	//vmConfigBytes     buffer.IoBuffer
	//pluginConfigBytes buffer.IoBuffer
}

var (
	_ api.StreamFilterChainFactory = &FilterConfigFactory{}
	// TODO: useless for now
	_ types.WasmPluginHandler = &FilterConfigFactory{}
)

func createProxyWasmFilterFactory(confs map[string]interface{}) (api.StreamFilterChainFactory, error) {
	factory := &FilterConfigFactory{
		config:        make([]*filterConfigItem, 0, len(confs)),
		RootContextID: 1,
		plugins:       make([]*WasmPlugin, 0, len(confs)),
		router:        &Router{routes: make(map[string]Group)},
	}

	for configID, confIf := range confs {
		conf, ok := confIf.(map[string]interface{})
		if !ok {
			log.DefaultLogger.Errorf("[proxywasm][factory] createProxyWasmFilterFactory config not a map, configID: %s", configID)
			return nil, errors.New("config not a map")
		}

		config, err := parseFilterConfigItem(conf)
		if err != nil {
			log.DefaultLogger.Errorf("[proxywasm][factory] createProxyWasmFilterFactory fail to parse config, configID: %s, err: %v", configID, err)
			return nil, err
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
				log.DefaultLogger.Errorf("[proxywasm][factory] createProxyWasmFilterFactory fail to add plugin, err: %v", err)
				return nil, err
			}

			//addWatchFile(config, pluginName)
		} else {
			pluginName = config.FromWasmPlugin
		}
		config.PluginName = pluginName

		pw := wasm.GetWasmManager().GetWasmPluginWrapperByName(pluginName)
		if pw == nil {
			return nil, errors.New("plugin not found")
		}

		config.VmConfig = pw.GetConfig().VmConfig
		factory.config = append(factory.config, config)

		plugin := pw.GetPlugin()
		// an instance used for sub tasks, not for user request
		instance := plugin.GetInstance()
		defer plugin.ReleaseInstance(instance)

		// handler set instance
		factory.LayottoHandler.Instance = instance
		// plugin register handler
		pw.RegisterPluginHandler(factory)

		// get the ABI of instance
		pluginABI := abi.GetABI(instance, AbiV2)
		if pluginABI == nil {
			log.DefaultLogger.Errorf("[proxywasm][factory] createProxyWasmFilterFactory fail to get instance abi, pluginName: %s", pluginName)
			plugin.ReleaseInstance(instance)
			return nil, errors.New("abi not found in instance")
		}
		// set the imports of abi
		pluginABI.SetABIImports(factory)

		// get wasm exports
		exports := pluginABI.GetABIExports().(Exports)

		instance.Lock(pluginABI)
		defer instance.Unlock()

		id, err := exports.ProxyGetID()
		if err != nil {
			log.DefaultLogger.Errorf("[proxywasm][factory] createProxyWasmFilterFactory fail to get wasm id, PluginName: %v, err: %v",
				config.PluginName, err)
			return nil, err
		}

		wasmPlugin := &WasmPlugin{
			pluginName:    config.PluginName,
			plugin:        plugin,
			rootContextID: config.RootContextID,
		}

		factory.router.RegisterRoute(id, wasmPlugin)
		factory.plugins = append(factory.plugins, wasmPlugin)
	}

	return factory, nil
}

func (f *FilterConfigFactory) CreateFilterChain(context context.Context, callbacks api.StreamFilterChainFactoryCallbacks) {
	filter := NewFilter(context, f)
	if filter == nil {
		return
	}

	callbacks.AddStreamReceiverFilter(filter, api.BeforeRoute)
	callbacks.AddStreamSenderFilter(filter, api.BeforeSend)
}

func (f *FilterConfigFactory) GetRootContextID() int32 {
	return f.RootContextID
}

//func (f *FilterConfigFactory) GetVmConfig() common.IoBuffer {
//	if f.vmConfigBytes != nil {
//		return f.vmConfigBytes
//	}
//
//	vmConfig := *f.config.VmConfig
//	typeOf := reflect.TypeOf(vmConfig)
//	valueOf := reflect.ValueOf(&vmConfig).Elem()
//
//	if typeOf.Kind() != reflect.Struct || typeOf.NumField() == 0 {
//		return nil
//	}
//
//	m := make(map[string]string)
//	for i := 0; i < typeOf.NumField(); i++ {
//		m[typeOf.Field(i).Name] = fmt.Sprintf("%v", valueOf.Field(i).Interface())
//	}
//
//	b := proxywasm.EncodeMap(m)
//	if b == nil {
//		return nil
//	}
//
//	f.vmConfigBytes = buffer.NewIoBufferBytes(b)
//
//	return f.vmConfigBytes
//}
//
//func (f *FilterConfigFactory) GetPluginConfig() common.IoBuffer {
//	if f.pluginConfigBytes != nil {
//		return f.pluginConfigBytes
//	}
//
//	b := proxywasm.EncodeMap(f.config.UserData)
//	if b == nil {
//		return nil
//	}
//
//	f.pluginConfigBytes = buffer.NewIoBufferBytes(b)
//
//	return f.pluginConfigBytes
//}

func (f *FilterConfigFactory) OnConfigUpdate(config v2.WasmPluginConfig) {
	// TODO: update the correct wasm
	//f.config.InstanceNum = config.InstanceNum
	//f.config.VmConfig = config.VmConfig
}

func (f *FilterConfigFactory) OnPluginStart(plugin types.WasmPlugin) {
	plugin.Exec(func(instance types.WasmInstance) bool {
		a := abi.GetABI(instance, AbiV2)
		a.SetABIImports(f)
		exports := a.GetABIExports().(Exports)

		instance.Lock(a)
		defer instance.Unlock()

		err := exports.ProxyOnContextCreate(f.RootContextID, 0)
		if err != nil {
			log.DefaultLogger.Errorf("[proxywasm][factory] OnPluginStart fail to create root context id, err: %v", err)
			return true
		}

		vmConfigSize := 0
		if vmConfigBytes := f.GetVmConfig(); vmConfigBytes != nil {
			vmConfigSize = vmConfigBytes.Len()
		}

		_, err = exports.ProxyOnVmStart(f.RootContextID, int32(vmConfigSize))
		if err != nil {
			log.DefaultLogger.Errorf("[proxywasm][factory] OnPluginStart fail to create root context id, err: %v", err)
			return true
		}

		pluginConfigSize := 0
		if pluginConfigBytes := f.GetPluginConfig(); pluginConfigBytes != nil {
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

func (f *FilterConfigFactory) OnPluginDestroy(plugin types.WasmPlugin) {}
