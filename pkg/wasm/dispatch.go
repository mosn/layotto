// Copyright 2021 Layotto Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package wasm

import (
	"errors"
	"math/rand"

	"mosn.io/mosn/pkg/log"
)

type Group struct {
	count   int
	plugins []*WasmPlugin
}

type Router struct {
	routes map[string]*Group
}

// RegisterRoute register a group with id
// unsafe for concurrent
func (route *Router) RegisterRoute(id string, plugin *WasmPlugin) {
	if group, found := route.routes[id]; found {
		group.plugins = append(filter(group.plugins, func(item *WasmPlugin) bool {
			return item.pluginName != plugin.pluginName
		}).([]*WasmPlugin), plugin)
		group.count = len(group.plugins)
	} else {
		route.routes[id] = &Group{
			count:   1,
			plugins: []*WasmPlugin{plugin},
		}
	}
}

func (route *Router) RemoveRoute(id string) {
	delete(route.routes, id)
}

// Get random plugin with rand id
func (route *Router) GetRandomPluginByID(id string) (*WasmPlugin, error) {
	group, ok := route.routes[id]
	if !ok {
		log.DefaultLogger.Infof("[proxywasm][dispatch] GetRandomPluginByID id not registered, id: %s", id)
		return nil, errors.New("id is not registered")
	}

	idx := rand.Intn(group.count)
	plugin := group.plugins[idx]
	log.DefaultLogger.Infof("[proxywasm][dispatch] GetRandomPluginByID return index: %d, plugin: %s", idx, plugin.pluginName)
	return plugin, nil
}
