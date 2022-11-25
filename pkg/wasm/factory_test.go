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
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	v2 "mosn.io/mosn/pkg/config/v2"
	"mosn.io/mosn/pkg/mock"
	"mosn.io/mosn/pkg/wasm"
)

func TestGetFactory(t *testing.T) {
	assert.Equal(t, factory, GetFactory())
}

func TestFilterConfigFactory_IsRegister(t *testing.T) {
	assert.False(t, factory.IsRegister("id_1"))
}

func TestFilterConfigFactory_Install(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	engine := mock.NewMockWasmVM(ctrl)
	wasm.RegisterWasmEngine("wasmtime", engine)

	conf := make(map[string]interface{})
	config := "{\"name\":\"id_1\",\"instance_num\":2,\"vm_config\":{\"engine\":\"wasmtime\",\"path\":\"nofile\"}}"
	err := json.Unmarshal([]byte(config), &conf)
	assert.NoError(t, err)
	err = factory.Install(conf)
	assert.NoError(t, err)
}

func TestFilterConfigFactory_Install_WithErrorConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	engine := mock.NewMockWasmVM(ctrl)
	wasm.RegisterWasmEngine("wasmtime", engine)

	conf := make(map[string]interface{})
	config := "{\"name\":\"id_1\"}"
	err := json.Unmarshal([]byte(config), &conf)
	assert.NoError(t, err)
	err = factory.Install(conf)
	assert.Equal(t, "nil vm config", err.Error())
}

func TestFilterConfigFactory_UnInstall_WithNoInstall(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	engine := mock.NewMockWasmVM(ctrl)
	wasm.RegisterWasmEngine("wasmtime", engine)

	id := "id_1"
	err := factory.UnInstall(id)
	assert.Equal(t, "id_1 is not registered", err.Error())
}

func TestFilterConfigFactory_UpdateInstanceNum_WithNoInstall(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	engine := mock.NewMockWasmVM(ctrl)
	wasm.RegisterWasmEngine("wasmtime", engine)

	id := "id_1"
	instanceNum := 1
	err := factory.UpdateInstanceNum(id, instanceNum)
	assert.Equal(t, "id_1 is not registered", err.Error())
}

func TestCreateProxyWasmFilterFactory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	engine := mock.NewMockWasmVM(ctrl)
	wasm.RegisterWasmEngine("wasmtime", engine)

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
