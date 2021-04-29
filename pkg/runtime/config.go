package runtime

import (
	"encoding/json"

	"github.com/layotto/layotto/pkg/services/configstores"
	"github.com/layotto/layotto/pkg/services/hello"
)

type MosnRuntimeConfig struct {
	HelloServiceManagement map[string]hello.HelloConfig        `json:"hellos"`
	ConfigStoreManagement  map[string]configstores.StoreConfig `json:"config_stores"`
}

func ParseRuntimeConfig(data json.RawMessage) (*MosnRuntimeConfig, error) {
	cfg := &MosnRuntimeConfig{}
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
