package wasm

import (
	"errors"

	"mosn.io/mosn/pkg/log"
)

type Group struct {
	plugins       []*WasmPlugin
	instanceCount int
	tPool         chan int
}

type Router struct {
	routes map[string]Group
}

// RegisterRoute register a group with id
// unsafe for concurrent
func (route *Router) RegisterRoute(id string, plugin *WasmPlugin) {
	if group, found := route.routes[id]; found {
		group.instanceCount += plugin.config.InstanceNum
		group.plugins = append(group.plugins, plugin)
		group.tPool = make(chan int, group.instanceCount)
	} else {
		route.routes[id] = Group{
			plugins:       []*WasmPlugin{plugin},
			instanceCount: plugin.config.InstanceNum,
			tPool:         make(chan int, plugin.config.InstanceNum),
		}
	}
}

func (route *Router) GetPluginByID(id string) (*WasmPlugin, error) {
	group, ok := route.routes[id]
	if !ok {
		log.DefaultLogger.Errorf("[proxywasm][filter] GetPluginByID id not registered, id: %s", id)
		return nil, errors.New("id is not registered")
	}

	if len(group.tPool) == 0 {
		for idx, plugin := range group.plugins {
			instanceNum := plugin.plugin.InstanceNum()
			for i := 0; i < instanceNum; i++ {
				group.tPool <- idx
			}
		}
	}

	return group.plugins[<-group.tPool], nil
}
