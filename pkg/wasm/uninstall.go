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

	"mosn.io/mosn/pkg/wasm"
	"mosn.io/pkg/log"
)

func init() {
	GetDefault().AddEndpoint("uninstall", &UnInstallEndpoint{})
}

type UnInstallEndpoint struct {
}

func (e *UnInstallEndpoint) Handle(ctx context.Context, f *Filter) (map[string]interface{}, error) {
	conf := make(map[string]interface{})
	err := json.Unmarshal(f.receiverFilterHandler.GetRequestData().Bytes(), &conf)
	if err != nil {
		log.DefaultLogger.Errorf("[proxywasm][uninstall] invalid body for request /wasm/uninstall, err:%v", err)
		return nil, err
	}

	if conf["name"] == nil {
		log.DefaultLogger.Errorf("[proxywasm][uninstall] can't get name property")
		return nil, errors.New("can't get name property")
	}

	id := (conf["name"]).(string)
	wasmPlugin, _ := f.router.GetRandomPluginByID(id)
	if wasmPlugin == nil {
		log.DefaultLogger.Errorf("[proxywasm][uninstall] %v is not registered", id)
		return nil, errors.New(id + " is not registered")
	}

	if f.pluginUsed != nil && f.pluginUsed.pluginName == wasmPlugin.pluginName {
		err = f.releaseUsedInstance()
		if err != nil {
			return nil, err
		}
	}

	err = wasm.GetWasmManager().UninstallWasmPluginByName(wasmPlugin.pluginName)
	if err != nil {
		log.DefaultLogger.Errorf("[proxywasm] [uninstall] fail to uninstall plugin, err: %v", err)
		return nil, err
	}

	f.factory.config = filter(f.factory.config, func(item *filterConfigItem) bool {
		return item.PluginName != wasmPlugin.pluginName
	}).([]*filterConfigItem)
	delete(f.factory.plugins, wasmPlugin.pluginName)
	removeWatchFile(wasmPlugin.config)
	f.router.RemoveRoute(id)
	return nil, nil
}
