package wasm

import "testing"

func TestRouter_RegisterRoute(t *testing.T) {
	// 1. Test case for appending a new plugin to an existing group when group is already present in route:
	routes := make(map[string]*Group)
	group := &Group{
		count:   1,
		plugins: []*WasmPlugin{{pluginName: "p1"}},
	}

	routes["test"] = group

	router := &Router{routes}

	router.RegisterRoute("test", &WasmPlugin{pluginName: "p2"})

	group = routes["test"]
	if group.count != 2 || len(group.plugins) != 2 {
		t.Errorf("Invalid group count or number of plugins")
	}

	// 2. Test case for creating a new group and appending plugin when group is not already present in route:
	routes = make(map[string]*Group)

	router = &Router{routes}

	router.RegisterRoute("test", &WasmPlugin{pluginName: "p1"})

	group = routes["test"]
	if group.count != 1 || len(group.plugins) != 1 {
		t.Errorf("Invalid group count or number of plugins")
	}

}

func TestRouter_RegisterRoute3(t *testing.T) {
	//3. Test case for appending the same plugin to an existing group (plugin with same pluginName already exists in group):
	routes := make(map[string]*Group)
	group := &Group{
		count:   1,
		plugins: []*WasmPlugin{{pluginName: "p1"}},
	}

	routes["test"] = group

	router := &Router{routes}

	// try adding same plugin
	router.RegisterRoute("test", &WasmPlugin{pluginName: "p1"})

	group = routes["test"]
	if group.count != 1 || len(group.plugins) != 1 {
		t.Errorf("Invalid group count or number of plugins")
	}
}
