package wasm

import (
	"errors"

	"mosn.io/mosn/pkg/log"
)

type Router struct {
	routes map[string]*Group
}

type Group struct {
	weightNodes []*WeightNode
}

type WeightNode struct {
	plugin          *WasmPlugin
	weight          int
	currentWeight   int
	effectiveWeight int
}

// RegisterRoute register a group with id
// unsafe for concurrent
func (route *Router) RegisterRoute(id string, plugin *WasmPlugin) {
	node := &WeightNode{
		plugin:          plugin,
		weight:          plugin.plugin.InstanceNum(),
		effectiveWeight: plugin.plugin.InstanceNum(),
	}
	if group, found := route.routes[id]; found {
		group.weightNodes = append(group.weightNodes, node)
	} else {
		route.routes[id] = &Group{
			weightNodes: []*WeightNode{node},
		}
	}
}

func (route *Router) GetPluginByID(id string) (*WasmPlugin, error) {
	group, ok := route.routes[id]
	if !ok {
		log.DefaultLogger.Errorf("[proxywasm][dispatch] GetPluginByID id not registered, id: %s", id)
		return nil, errors.New("id is not registered")
	}

	if plugin := group.Next(); plugin != nil {
		log.DefaultLogger.Infof("[proxywasm][dispatch] GetPluginByID return plugin: %s", plugin.pluginName)
		return plugin, nil
	}
	return nil, errors.New("Next return nil")
}

func (g *Group) Next() *WasmPlugin {
	var best *WeightNode
	total := 0
	for i := 0; i < len(g.weightNodes); i++ {
		w := g.weightNodes[i]
		total += w.effectiveWeight
		w.currentWeight += w.effectiveWeight
		if w.effectiveWeight < w.weight {
			w.effectiveWeight++
		}

		if best == nil || w.currentWeight > best.currentWeight {
			best = w
		}
	}

	if best == nil {
		return nil
	}
	best.currentWeight -= total
	return best.plugin
}
