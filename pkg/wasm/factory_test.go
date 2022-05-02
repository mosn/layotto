//
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

func TestCreateProxyWasmFilterFactory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	engine := mock.NewMockWasmVM(ctrl)
	wasm.RegisterWasmEngine("wasmer", engine)

	config := `
				{
					  "type": "Layotto",
					  "config": {
						"function1": {
						  "name": "function1",
						  "instance_num": 1,
						  "vm_config": {
							"engine": "wasmer",
							"path": "nofile"
						  }
						},
						"function2": {
						  "name": "function2",
						  "instance_num": 1,
						  "vm_config": {
							"engine": "wasmer",
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
