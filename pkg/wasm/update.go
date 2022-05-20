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
	"errors"

	v2 "mosn.io/mosn/pkg/config/v2"
	"mosn.io/mosn/pkg/wasm"
	"mosn.io/pkg/log"
)

func init() {
	GetDefault().AddEndpoint("update", &UpdateEndpoint{})
}

type UpdateEndpoint struct {
}

func (e *UpdateEndpoint) Handle(ctx context.Context, f *Filter) (map[string]interface{}, error) {
	conf := make(map[string]interface{})
	err := json.Unmarshal(f.receiverFilterHandler.GetRequestData().Bytes(), &conf)
	if err != nil {
		log.DefaultLogger.Errorf("[proxywasm][update] invalid body for request /wasm/update, err:%v", err)
		return nil, err
	}

	if conf["name"] == nil {
		log.DefaultLogger.Errorf("[proxywasm][update] can't get name property")
		return nil, errors.New("can't get name property")
	}

	if conf["instance_num"] == nil {
		log.DefaultLogger.Errorf("[proxywasm][update] can't get instance_num property")
		return nil, errors.New("can't get instance_num property")
	}

	instanceNum := int(conf["instance_num"].(float64))
	if instanceNum <= 0 {
		log.DefaultLogger.Errorf("[proxywasm][update] instance_num should be greater than 0")
		return nil, errors.New("instance_num should be greater than 0")
	}

	id := (conf["name"]).(string)
	wasmPlugin, _ := f.router.GetRandomPluginByID(id)
	if wasmPlugin == nil {
		log.DefaultLogger.Errorf("[proxywasm][update] %v is not registered", id)
		return nil, errors.New(id + " is not registered")
	}
	var config *filterConfigItem
	for _, item := range f.factory.config {
		if item.PluginName == wasmPlugin.pluginName {
			config = item
			break
		}
	}
	if config == nil {
		log.DefaultLogger.Errorf("[proxywasm][update] can't find config for %v", id)
		return nil, errors.New("can't find config for " + id)
	}

	if config.InstanceNum == instanceNum {
		return nil, nil
	}

	config.InstanceNum = instanceNum
	v2Config := v2.WasmPluginConfig{
		PluginName:  config.PluginName,
		VmConfig:    config.VmConfig,
		InstanceNum: config.InstanceNum,
	}
	err = wasm.GetWasmManager().AddOrUpdateWasm(v2Config)
	if err != nil {
		log.DefaultLogger.Errorf("[proxywasm] [update] fail to update plugin, err: %v", err)
		return nil, err
	}
	pw := wasm.GetWasmManager().GetWasmPluginWrapperByName(config.PluginName)
	if pw == nil {
		log.DefaultLogger.Errorf("[proxywasm] [update] plugin not found")
		return nil, errors.New("plugin not found")
	}
	f.factory.plugins[config.PluginName] = &WasmPlugin{
		pluginName:    config.PluginName,
		plugin:        pw.GetPlugin(),
		rootContextID: config.RootContextID,
		config:        config,
	}
	pw.RegisterPluginHandler(f.factory)
	log.DefaultLogger.Infof("[proxywasm] [update] wasm instance number updated success, id: %v, num: %v", id, instanceNum)
	return nil, nil
}
