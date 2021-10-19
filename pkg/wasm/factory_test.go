package wasm

import (
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	v2 "mosn.io/mosn/pkg/config/v2"
	"mosn.io/mosn/pkg/mock"
	"mosn.io/mosn/pkg/wasm"
	"testing"
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
