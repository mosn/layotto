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
		group.count += 1
		group.plugins = append(group.plugins, plugin)
	} else {
		route.routes[id] = &Group{
			count:   1,
			plugins: []*WasmPlugin{plugin},
		}
	}
}

func (route *Router) GetRandomPluginByID(id string) (*WasmPlugin, error) {
	group, ok := route.routes[id]
	if !ok {
		log.DefaultLogger.Errorf("[proxywasm][filter] GetRandomPluginByID id not registered, id: %s", id)
		return nil, errors.New("id is not registered")
	}

	idx := rand.Intn(group.count)
	plugin := group.plugins[idx]
	log.DefaultLogger.Infof("[proxywasm][dispatch] GetRandomPluginByID return index: %d, plugin: %s", idx, plugin.pluginName)
	return plugin, nil

}
