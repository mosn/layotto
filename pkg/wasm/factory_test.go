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
	"encoding/json"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	mockwasm "mosn.io/layotto/pkg/mock/wasm"

	v2 "mosn.io/mosn/pkg/config/v2"
	"mosn.io/mosn/pkg/mock"
	"mosn.io/mosn/pkg/types"
	"mosn.io/mosn/pkg/wasm"
	"mosn.io/mosn/pkg/wasm/abi"
)

var (
	configNormal = `{
    "name": "function_1",
    "instance_num": 2,
    "vm_config": {
        "engine": "wasmtime",
        "path": "no_file"
    }
}`

	configGlobal = `{
	"from_wasm_plugin": "global_plugin",
    "instance_num": 2
}`
)

func prepare(t *testing.T) *gomock.Controller {
	ctrl := gomock.NewController(t)
	// mock wasm engine
	engine := mock.NewMockWasmVM(ctrl)
	wasm.RegisterWasmEngine("wasmtime", engine)

	return ctrl
}

func reset(ctrl *gomock.Controller) {
	ctrl.Finish()
}

func mockConfig(cfgStr string) map[string]interface{} {
	cfg := make(map[string]interface{})
	_ = json.Unmarshal([]byte(cfgStr), &cfg)
	return cfg
}

func mockWasmVmConfig(path string) v2.WasmVmConfig {
	return v2.WasmVmConfig{
		Engine: "wasmtime",
		Path:   path,
	}
}

func mockWasmConfig(pluginName string, instanceNum int) v2.WasmPluginConfig {
	vmConfig := mockWasmVmConfig("no_file")
	return v2.WasmPluginConfig{
		PluginName:  pluginName,
		VmConfig:    &vmConfig,
		InstanceNum: instanceNum,
	}
}

func mockFilterConfigItem(wasmConfig v2.WasmPluginConfig, ctxID int32, pluginName string) *filterConfigItem {
	return &filterConfigItem{
		FromWasmPlugin: "",
		VmConfig:       wasmConfig.VmConfig,
		InstanceNum:    wasmConfig.InstanceNum,
		RootContextID:  ctxID,
		PluginName:     pluginName,
	}
}

func mockLayottoWasmPlugin(pluginName string, instanceNum int, plugin *mock.MockWasmPlugin) *WasmPlugin {
	plugin.EXPECT().PluginName().Return(pluginName).AnyTimes()
	wasmConfig := mockWasmConfig(pluginName, instanceNum)
	ctxID := int32(1)

	return &WasmPlugin{
		pluginName:    pluginName,
		plugin:        plugin,
		rootContextID: ctxID,
		config:        mockFilterConfigItem(wasmConfig, ctxID, pluginName),
	}
}

func TestGetFactory(t *testing.T) {
	assert.Equal(t, factory, GetFactory())
}

func TestFilterConfigFactory_IsRegister(t *testing.T) {
	assert.False(t, factory.IsRegister("id_1"))
}

func TestFilterConfigFactory_Install(t *testing.T) {
	ctrl := prepare(t)
	defer reset(ctrl)

	conf := make(map[string]interface{})
	config := "{\"name\":\"id_1\",\"instance_num\":2,\"vm_config\":{\"engine\":\"wasmtime\",\"path\":\"nofile\"}}"
	err := json.Unmarshal([]byte(config), &conf)
	assert.NoError(t, err)
	manager := wasm.GetWasmManager()
	err = factory.Install(conf, manager)
	assert.NoError(t, err)
}

func TestFilterConfigFactory_Install_WithErrorConfig(t *testing.T) {
	ctrl := prepare(t)
	defer reset(ctrl)

	conf := make(map[string]interface{})
	config := "{\"name\":\"id_1\"}"
	err := json.Unmarshal([]byte(config), &conf)
	assert.NoError(t, err)
	manager := wasm.GetWasmManager()
	err = factory.Install(conf, manager)
	assert.Equal(t, "nil vm config", err.Error())
}

func TestFilterConfigFactory_UnInstall_WithNoInstall(t *testing.T) {
	ctrl := prepare(t)
	defer reset(ctrl)

	id := "id_1"
	manager := mock.NewMockWasmManager(ctrl)
	err := factory.UnInstall(id, manager)
	assert.Equal(t, "id_1 is not registered", err.Error())
}

func TestFilterConfigFactory_UpdateInstanceNum_WithNoInstall(t *testing.T) {
	ctrl := prepare(t)
	defer reset(ctrl)

	id := "id_1"
	instanceNum := 1
	manager := mock.NewMockWasmManager(ctrl)
	err := factory.UpdateInstanceNum(id, instanceNum, manager)
	assert.Equal(t, "id_1 is not registered", err.Error())
}

func TestCreateProxyWasmFilterFactory(t *testing.T) {
	ctrl := prepare(t)
	defer reset(ctrl)

	config := `
				{
					  "type": "Layotto",
					  "config": {
						"function1": {
						  "name": "function1",
						  "instance_num": 1,
						  "vm_config": {
							"engine": "wasmtime",
							"path": "nofile"
						  }
						},
						"function2": {
						  "name": "function2",
						  "instance_num": 1,
						  "vm_config": {
							"engine": "wasmtime",
							"path": "nofile"
						  }
						}
					  }
            	}
			  `
	conf := &v2.Filter{}
	err := json.Unmarshal([]byte(config), conf)
	assert.NoError(t, err)
	_, err = createProxyWasmFilterFactory(conf.Config)
	assert.NoError(t, err)
}

func TestFilterConfigFactory_Install1(t *testing.T) {
	ctrl := prepare(t)
	defer reset(ctrl)

	type fields struct {
		LayottoHandler LayottoHandler
		config         []*filterConfigItem
		RootContextID  int32
		plugins        map[string]*WasmPlugin
		router         *Router
	}
	type args struct {
		conf map[string]interface{}
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantErr      assert.ErrorAssertionFunc
		mockAndCheck func(ctrl *gomock.Controller) *mock.MockWasmManager
	}{
		{
			name: "config.FromWasmPlugin is empty",
			fields: fields{
				LayottoHandler: LayottoHandler{},
				config:         []*filterConfigItem{},
				RootContextID:  1,
				plugins:        map[string]*WasmPlugin{},
				router:         &Router{routes: map[string]*Group{}},
			},
			args: args{
				conf: mockConfig(configNormal),
			},
			wantErr: assert.NoError,
			mockAndCheck: func(ctrl *gomock.Controller) *mock.MockWasmManager {
				manager := mock.NewMockWasmManager(ctrl)
				pw := mock.NewMockWasmPluginWrapper(ctrl)

				// mock WasmManager & WasmPluginWrapper
				gomock.InOrder(
					manager.EXPECT().AddOrUpdateWasm(gomock.Any()).Return(nil).Times(1),
					manager.EXPECT().GetWasmPluginWrapperByName(gomock.Any()).Return(pw).Times(1),
					pw.EXPECT().GetConfig().Return(mockWasmConfig("function_1", 2)).Times(1),
					pw.EXPECT().GetPlugin().Times(1),
					pw.EXPECT().RegisterPluginHandler(gomock.Any()).Times(1),
				)

				return manager
			},
		},
		{
			name: "config.FromWasmPlugin is not empty",
			fields: fields{
				LayottoHandler: LayottoHandler{},
				config:         []*filterConfigItem{},
				RootContextID:  1,
				plugins:        map[string]*WasmPlugin{},
				router:         &Router{routes: map[string]*Group{}},
			},
			args: args{
				conf: mockConfig(configGlobal),
			},
			wantErr: assert.NoError,
			mockAndCheck: func(ctrl *gomock.Controller) *mock.MockWasmManager {
				manager := mock.NewMockWasmManager(ctrl)
				pw := mock.NewMockWasmPluginWrapper(ctrl)

				// mock WasmManager & WasmPluginWrapper
				manager.EXPECT().AddOrUpdateWasm(gomock.Any()).Return(nil).MaxTimes(0)
				gomock.InOrder(
					manager.EXPECT().GetWasmPluginWrapperByName("global_plugin").Return(pw).Times(1),
					pw.EXPECT().GetConfig().Return(mockWasmConfig("global_plugin", 2)).Times(1),
					pw.EXPECT().GetPlugin().Times(1),
					pw.EXPECT().RegisterPluginHandler(gomock.Any()).Times(1),
				)

				return manager
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FilterConfigFactory{
				LayottoHandler: tt.fields.LayottoHandler,
				config:         tt.fields.config,
				RootContextID:  tt.fields.RootContextID,
				plugins:        tt.fields.plugins,
				router:         tt.fields.router,
			}

			manager := tt.mockAndCheck(ctrl)
			tt.wantErr(t, f.Install(tt.args.conf, manager), fmt.Sprintf("Install(%v)", tt.args.conf))
		})
	}
}

func TestFilterConfigFactory_OnConfigUpdate(t *testing.T) {
	vmConfig := mockWasmVmConfig("no_file")

	type fields struct {
		LayottoHandler LayottoHandler
		config         []*filterConfigItem
		RootContextID  int32
		plugins        map[string]*WasmPlugin
		router         *Router
	}
	type args struct {
		config v2.WasmPluginConfig
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		checkFunc func(t *testing.T, f *FilterConfigFactory)
	}{
		{
			name: "f.config is empty",
			fields: fields{
				LayottoHandler: LayottoHandler{},
				config:         []*filterConfigItem{},
				RootContextID:  1,
				plugins:        map[string]*WasmPlugin{},
				router:         &Router{routes: map[string]*Group{}},
			},
			args: args{
				config: mockWasmConfig("function_1", 2),
			},
			checkFunc: func(t *testing.T, f *FilterConfigFactory) {
				assert.Equal(t, 0, len(factory.config))
			},
		},
		{
			name: "hit pluginName",
			fields: fields{
				LayottoHandler: LayottoHandler{},
				config: []*filterConfigItem{
					{
						VmConfig:      &vmConfig,
						InstanceNum:   2,
						RootContextID: 1,
						PluginName:    "function_1",
					},
				},
				RootContextID: 1,
				plugins:       map[string]*WasmPlugin{},
				router:        &Router{routes: map[string]*Group{}},
			},
			args: args{
				config: mockWasmConfig("function_1", 3),
			},
			checkFunc: func(t *testing.T, f *FilterConfigFactory) {
				assert.Equal(t, 1, len(f.config))
				assert.Equal(t, 3, f.config[0].InstanceNum)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FilterConfigFactory{
				LayottoHandler: tt.fields.LayottoHandler,
				config:         tt.fields.config,
				RootContextID:  tt.fields.RootContextID,
				plugins:        tt.fields.plugins,
				router:         tt.fields.router,
			}
			f.OnConfigUpdate(tt.args.config)
			tt.checkFunc(t, f)
		})
	}
}

func TestFilterConfigFactory_OnPluginDestroy(t *testing.T) {
	f := &FilterConfigFactory{
		LayottoHandler: LayottoHandler{},
		config:         []*filterConfigItem{},
		plugins:        map[string]*WasmPlugin{},
		router:         &Router{routes: map[string]*Group{}},
	}
	// empty test case
	f.OnPluginDestroy(nil)
}

func TestFilterConfigFactory_OnPluginStart(t *testing.T) {
	ctrl := prepare(t)
	defer reset(ctrl)

	commonMockFunc := func(ctrl *gomock.Controller, instance *mock.MockWasmInstance) *mock.MockWasmPlugin {
		plugin := mock.NewMockWasmPlugin(ctrl)
		plugin.EXPECT().Exec(gomock.Any()).Times(1).
			Do(func(lambda func(instance types.WasmInstance) bool) {
				lambda(instance)
			})
		return plugin
	}

	mockAbiFunc := func(ctrl *gomock.Controller, instance *mock.MockWasmInstance) *mock.MockABI {
		a := mock.NewMockABI(ctrl)
		abiFactory := func(instance types.WasmInstance) types.ABI {
			return a
		}
		abi.RegisterABI(AbiV2, abiFactory)
		return a
	}

	type fields struct {
		LayottoHandler LayottoHandler
		config         []*filterConfigItem
		RootContextID  int32
		plugins        map[string]*WasmPlugin
		router         *Router
	}
	tests := []struct {
		name         string
		fields       fields
		mockAndCheck func(ctrl *gomock.Controller, plugin *mock.MockWasmPlugin, instance *mock.MockWasmInstance, f *FilterConfigFactory)
	}{
		{
			name: "f.plugins empty",
			fields: fields{
				LayottoHandler: LayottoHandler{},
				config:         []*filterConfigItem{},
				RootContextID:  1,
				plugins:        map[string]*WasmPlugin{},
				router:         &Router{routes: map[string]*Group{}},
			},
			mockAndCheck: func(ctrl *gomock.Controller, plugin *mock.MockWasmPlugin, instance *mock.MockWasmInstance, f *FilterConfigFactory) {
				plugin.EXPECT().PluginName().Return("not_exist").Times(2)
			},
		},
		{
			name: "exports.ProxyGetID error",
			fields: fields{
				LayottoHandler: LayottoHandler{},
				config:         []*filterConfigItem{},
				RootContextID:  1,
				plugins: map[string]*WasmPlugin{
					"plugin_1": mockLayottoWasmPlugin("plugin_1", 2, mock.NewMockWasmPlugin(ctrl)),
				},
				router: &Router{routes: map[string]*Group{}},
			},
			mockAndCheck: func(ctrl *gomock.Controller, plugin *mock.MockWasmPlugin, instance *mock.MockWasmInstance, f *FilterConfigFactory) {
				a := mockAbiFunc(ctrl, instance)
				exports := mockwasm.NewMockExports(ctrl)
				module := mock.NewMockWasmModule(ctrl)
				gomock.InOrder(
					plugin.EXPECT().PluginName().Return("plugin_1").Times(1),
					instance.EXPECT().GetModule().Return(module).Times(1),
					module.EXPECT().GetABINameList().Return([]string{AbiV2}).Times(1),
					a.EXPECT().SetABIImports(gomock.Any()).Times(1),
					a.EXPECT().GetABIExports().Return(exports).Times(1),
					instance.EXPECT().Lock(gomock.Any()).Times(1),
					exports.EXPECT().ProxyGetID().Return("", errors.New("exports.ProxyGetID error")).Times(1),
					plugin.EXPECT().PluginName().Return("plugin_1").Times(1),
					instance.EXPECT().Unlock().Times(1),
				)
			},
		},
		{
			name: "exports.ProxyOnContextCreate error",
			fields: fields{
				LayottoHandler: LayottoHandler{},
				config:         []*filterConfigItem{},
				RootContextID:  1,
				plugins: map[string]*WasmPlugin{
					"plugin_1": mockLayottoWasmPlugin("plugin_1", 2, mock.NewMockWasmPlugin(ctrl)),
				},
				router: &Router{routes: map[string]*Group{}},
			},
			mockAndCheck: func(ctrl *gomock.Controller, plugin *mock.MockWasmPlugin, instance *mock.MockWasmInstance, f *FilterConfigFactory) {
				a := mockAbiFunc(ctrl, instance)
				exports := mockwasm.NewMockExports(ctrl)
				module := mock.NewMockWasmModule(ctrl)
				gomock.InOrder(
					plugin.EXPECT().PluginName().Return("plugin_1").Times(1),
					instance.EXPECT().GetModule().Return(module).Times(1),
					module.EXPECT().GetABINameList().Return([]string{AbiV2}).Times(1),
					a.EXPECT().SetABIImports(gomock.Any()).Times(1),
					a.EXPECT().GetABIExports().Return(exports).Times(1),
					instance.EXPECT().Lock(gomock.Any()).Times(1),
					exports.EXPECT().ProxyGetID().Return("id_1", nil).Times(1),
					exports.EXPECT().ProxyOnContextCreate(f.RootContextID, int32(0)).
						Return(errors.New("exports.ProxyOnContextCreate error")).Times(1),
					instance.EXPECT().Unlock().Times(1),
				)
			},
		},
		// TODO: Identify the reason why plugin.GetVmConfig() cannot be mocked
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FilterConfigFactory{
				LayottoHandler: tt.fields.LayottoHandler,
				config:         tt.fields.config,
				RootContextID:  tt.fields.RootContextID,
				plugins:        tt.fields.plugins,
				router:         tt.fields.router,
			}

			instance := mock.NewMockWasmInstance(ctrl)
			plugin := commonMockFunc(ctrl, instance)
			tt.mockAndCheck(ctrl, plugin, instance, f)
			f.OnPluginStart(plugin)
		})
	}
}

func TestFilterConfigFactory_UnInstall(t *testing.T) {
	ctrl := prepare(t)
	defer reset(ctrl)

	type fields struct {
		LayottoHandler LayottoHandler
		config         []*filterConfigItem
		RootContextID  int32
		plugins        map[string]*WasmPlugin
		router         *Router
	}
	type args struct {
		id string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantErr   assert.ErrorAssertionFunc
		mockFunc  func(ctrl *gomock.Controller) *mock.MockWasmManager
		checkFunc func(t *testing.T, f *FilterConfigFactory)
	}{
		{
			name: "manager.UnInstall error",
			fields: fields{
				LayottoHandler: LayottoHandler{},
				config: []*filterConfigItem{
					mockFilterConfigItem(mockWasmConfig("function_1", 1), 1, "function_1"),
				},
				RootContextID: 1,
				plugins: map[string]*WasmPlugin{
					"function_1": mockLayottoWasmPlugin("function_1", 1, mock.NewMockWasmPlugin(ctrl)),
				},
				router: &Router{
					routes: map[string]*Group{
						"id_1": {
							count: 1,
							plugins: []*WasmPlugin{
								mockLayottoWasmPlugin("function_1", 1, mock.NewMockWasmPlugin(ctrl)),
							},
						},
					},
				},
			},
			args: args{
				id: "id_1",
			},
			wantErr: assert.Error,
			mockFunc: func(ctrl *gomock.Controller) *mock.MockWasmManager {
				manager := mock.NewMockWasmManager(ctrl)
				manager.EXPECT().UninstallWasmPluginByName(gomock.Any()).Return(fmt.Errorf("UnInstall fail")).Times(1)
				return manager
			},
			checkFunc: func(t *testing.T, f *FilterConfigFactory) {
				assert.Equal(t, 1, len(f.config))
				assert.Equal(t, 1, len(f.plugins))
				assert.Equal(t, 1, len(f.router.routes))
			},
		},
		{
			name: "UnInstall success",
			fields: fields{
				LayottoHandler: LayottoHandler{},
				config: []*filterConfigItem{
					mockFilterConfigItem(mockWasmConfig("function_1", 1), 1, "function_1"),
				},
				RootContextID: 1,
				plugins: map[string]*WasmPlugin{
					"function_1": mockLayottoWasmPlugin("function_1", 1, mock.NewMockWasmPlugin(ctrl)),
				},
				router: &Router{
					routes: map[string]*Group{
						"id_1": {
							count: 1,
							plugins: []*WasmPlugin{
								mockLayottoWasmPlugin("function_1", 1, mock.NewMockWasmPlugin(ctrl)),
							},
						},
					},
				},
			},
			args: args{
				id: "id_1",
			},
			wantErr: assert.NoError,
			mockFunc: func(ctrl *gomock.Controller) *mock.MockWasmManager {
				manager := mock.NewMockWasmManager(ctrl)
				manager.EXPECT().UninstallWasmPluginByName(gomock.Any()).Return(nil).Times(1)
				return manager
			},
			checkFunc: func(t *testing.T, f *FilterConfigFactory) {
				assert.Equal(t, 0, len(f.config))
				assert.Equal(t, 0, len(f.plugins))
				assert.Equal(t, 0, len(f.router.routes))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FilterConfigFactory{
				LayottoHandler: tt.fields.LayottoHandler,
				config:         tt.fields.config,
				RootContextID:  tt.fields.RootContextID,
				plugins:        tt.fields.plugins,
				router:         tt.fields.router,
			}

			manager := tt.mockFunc(ctrl)
			tt.wantErr(t, f.UnInstall(tt.args.id, manager), fmt.Sprintf("UnInstall(%v)", tt.args.id))
			tt.checkFunc(t, f)
		})
	}
}

func TestFilterConfigFactory_UpdateInstanceNum(t *testing.T) {
	ctrl := prepare(t)
	defer reset(ctrl)

	type fields struct {
		LayottoHandler LayottoHandler
		config         []*filterConfigItem
		RootContextID  int32
		plugins        map[string]*WasmPlugin
		router         *Router
	}
	type args struct {
		id          string
		instanceNum int
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantErr      assert.ErrorAssertionFunc
		mockAndCheck func(ctrl *gomock.Controller) *mock.MockWasmManager
	}{
		{
			name: "f.config is empty",
			fields: fields{
				LayottoHandler: LayottoHandler{},
				config:         []*filterConfigItem{},
				RootContextID:  1,
				plugins: map[string]*WasmPlugin{
					"function_1": mockLayottoWasmPlugin("function_1", 1, mock.NewMockWasmPlugin(ctrl)),
				},
				router: &Router{
					routes: map[string]*Group{
						"id_1": {
							count: 1,
							plugins: []*WasmPlugin{
								mockLayottoWasmPlugin("function_1", 1, mock.NewMockWasmPlugin(ctrl)),
							},
						},
					},
				},
			},
			args: args{
				id:          "id_1",
				instanceNum: 2,
			},
			wantErr: assert.Error,
			mockAndCheck: func(ctrl *gomock.Controller) *mock.MockWasmManager {
				return mock.NewMockWasmManager(ctrl)
			},
		},
		{
			name: "not update",
			fields: fields{
				LayottoHandler: LayottoHandler{},
				config: []*filterConfigItem{
					mockFilterConfigItem(mockWasmConfig("function_1", 1), 1, "function_1"),
				},
				RootContextID: 1,
				plugins: map[string]*WasmPlugin{
					"function_1": mockLayottoWasmPlugin("function_1", 1, mock.NewMockWasmPlugin(ctrl)),
				},
				router: &Router{
					routes: map[string]*Group{
						"id_1": {
							count: 1,
							plugins: []*WasmPlugin{
								mockLayottoWasmPlugin("function_1", 1, mock.NewMockWasmPlugin(ctrl)),
							},
						},
					},
				},
			},
			args: args{
				id:          "id_1",
				instanceNum: 1,
			},
			wantErr: assert.NoError,
			mockAndCheck: func(ctrl *gomock.Controller) *mock.MockWasmManager {
				return mock.NewMockWasmManager(ctrl)
			},
		},
		{
			name: "update",
			fields: fields{
				LayottoHandler: LayottoHandler{},
				config: []*filterConfigItem{
					mockFilterConfigItem(mockWasmConfig("function_1", 1), 1, "function_1"),
				},
				RootContextID: 1,
				plugins: map[string]*WasmPlugin{
					"function_1": mockLayottoWasmPlugin("function_1", 1, mock.NewMockWasmPlugin(ctrl)),
				},
				router: &Router{
					routes: map[string]*Group{
						"id_1": {
							count: 1,
							plugins: []*WasmPlugin{
								mockLayottoWasmPlugin("function_1", 1, mock.NewMockWasmPlugin(ctrl)),
							},
						},
					},
				},
			},
			args: args{
				id:          "id_1",
				instanceNum: 2,
			},
			wantErr: assert.NoError,
			mockAndCheck: func(ctrl *gomock.Controller) *mock.MockWasmManager {
				manager := mock.NewMockWasmManager(ctrl)
				pw := mock.NewMockWasmPluginWrapper(ctrl)

				gomock.InOrder(
					manager.EXPECT().AddOrUpdateWasm(gomock.Any()).Return(nil).Times(1),
					manager.EXPECT().GetWasmPluginWrapperByName(gomock.Any()).Return(pw).Times(1),
					pw.EXPECT().GetPlugin().Times(1),
					pw.EXPECT().RegisterPluginHandler(gomock.Any()).Times(1),
				)

				return manager
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FilterConfigFactory{
				LayottoHandler: tt.fields.LayottoHandler,
				config:         tt.fields.config,
				RootContextID:  tt.fields.RootContextID,
				plugins:        tt.fields.plugins,
				router:         tt.fields.router,
			}

			manager := tt.mockAndCheck(ctrl)
			tt.wantErr(t, f.UpdateInstanceNum(tt.args.id, tt.args.instanceNum, manager), fmt.Sprintf("UpdateInstanceNum(%v, %v)", tt.args.id, tt.args.instanceNum))
		})
	}
}
