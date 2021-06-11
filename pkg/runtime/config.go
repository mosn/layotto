package runtime

import (
	"encoding/json"
	"github.com/layotto/layotto/components/configstores"
	"github.com/layotto/layotto/components/hello"
	"github.com/layotto/layotto/components/rpc"
	"github.com/layotto/layotto/pkg/runtime/pubsub"
	"github.com/layotto/layotto/pkg/runtime/state"
)

type AppConfig struct {
	AppId            string `json:"app_id"`
	GrpcCallbackPort int    `json:"grpc_callback_port"`
}

type MosnRuntimeConfig struct {
	AppManagement          AppConfig                           `json:"app"`
	HelloServiceManagement map[string]hello.HelloConfig        `json:"hellos"`
	ConfigStoreManagement  map[string]configstores.StoreConfig `json:"config_stores"`
	RpcManagement          map[string]rpc.RpcConfig            `json:"rpcs"`
	PubSubManagement       map[string]pubsub.Config            `json:"pub_subs"`
	StateManagement        map[string]state.Config             `json:"state"`
}

func ParseRuntimeConfig(data json.RawMessage) (*MosnRuntimeConfig, error) {
	cfg := &MosnRuntimeConfig{}
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
