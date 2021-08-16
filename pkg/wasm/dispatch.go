package wasm

import (
	"errors"
	"math/rand"
)

type Group struct {
	count   int
	plugins []*WasmPlugin
}

type Router struct {
	routes map[string]Group
}

// RegisterRoute register a group with id
// unsafe for concurrent
func (route *Router) RegisterRoute(id string, plugin *WasmPlugin) {
	if group, found := route.routes[id]; found {
		group.count += 1
		group.plugins = append(group.plugins, plugin)
		route.routes[id] = group
	} else {
		route.routes[id] = Group{
			count:   1,
			plugins: []*WasmPlugin{plugin},
		}
	}
}

func (route *Router) GetRandomPluginByID(id string) (*WasmPlugin, error) {
	group, ok := route.routes[id]
	if !ok {
		return nil, errors.New("id is not registered")
	}

	idx := rand.Intn(group.count)
	return group.plugins[idx], nil
}
