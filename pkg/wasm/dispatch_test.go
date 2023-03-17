package wasm

import "testing"

// Generated using Cursor
// TestRouter tests the Router struct methods
func TestRouter(t *testing.T) {
	router := &Router{
		routes: make(map[string]*Group),
	}

	plugin1 := &WasmPlugin{pluginName: "plugin1"}
	plugin2 := &WasmPlugin{pluginName: "plugin2"}

	// Test RegisterRoute
	router.RegisterRoute("route1", plugin1)
	router.RegisterRoute("route1", plugin2)

	if len(router.routes["route1"].plugins) != 2 {
		t.Errorf("Expected 2 plugins in route1, got %d", len(router.routes["route1"].plugins))
	}

	// Test RemoveRoute
	router.RemoveRoute("route1")

	if _, found := router.routes["route1"]; found {
		t.Error("Expected route1 to be removed")
	}

	// Test GetRandomPluginByID
	router.RegisterRoute("route2", plugin1)
	router.RegisterRoute("route2", plugin2)

	_, err := router.GetRandomPluginByID("route2")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	_, err = router.GetRandomPluginByID("nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent route")
	}
}

// Generated using chatGPT
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

// TestRegisterRoute adds a new route and checks if it exists in the map
func TestRegisterRoute(t *testing.T) {
	route := &Router{routes: make(map[string]*Group)}
	plugin := &WasmPlugin{pluginName: "test"}
	route.RegisterRoute("test", plugin)
	_, ok := route.routes["test"]
	if !ok {
		t.Errorf("Expected route with id: test to exist")
	}
}

// TestRemoveRoute adds a new route, removes it and checks if it does not exist in the map
func TestRemoveRoute(t *testing.T) {
	route := &Router{routes: make(map[string]*Group)}
	plugin := &WasmPlugin{pluginName: "test"}
	route.RegisterRoute("test", plugin)
	route.RemoveRoute("test")
	_, ok := route.routes["test"]
	if ok {
		t.Errorf("Expected route with id: test to be removed")
	}
}

// TestGetRandomPluginByID adds a new route with a plugin, retrieves a random plugin and checks if it is in the group
func TestGetRandomPluginByID(t *testing.T) {
	route := &Router{routes: make(map[string]*Group)}
	plugin := &WasmPlugin{pluginName: "test"}
	route.RegisterRoute("test", plugin)
	randPlugin, err := route.GetRandomPluginByID("test")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	var found bool
	for _, p := range route.routes["test"].plugins {
		if p == randPlugin {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected random plugin to be in the group")
	}
}

// TestGetRandomPluginByIDInvalidID checks if an error is returned for an invalid id
func TestGetRandomPluginByIDInvalidID(t *testing.T) {
	route := &Router{routes: make(map[string]*Group)}
	_, err := route.GetRandomPluginByID("test")
	if err == nil {
		t.Errorf("Expected an error for an invalid id")
	}
}
