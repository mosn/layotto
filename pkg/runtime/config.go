/*
 * Copyright 2021 Layotto Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package runtime

import (
	"encoding/json"

	"mosn.io/layotto/components/lock"

	"mosn.io/layotto/components/oss"

	"mosn.io/layotto/pkg/runtime/secretstores"

	"mosn.io/layotto/components/custom"

	"mosn.io/layotto/pkg/runtime/bindings"

	"mosn.io/layotto/components/file"

	"mosn.io/layotto/components/configstores"
	"mosn.io/layotto/components/hello"
	"mosn.io/layotto/components/rpc"
	"mosn.io/layotto/components/sequencer"
	"mosn.io/layotto/pkg/runtime/pubsub"
	"mosn.io/layotto/pkg/runtime/state"
)

type AppConfig struct {
	AppId            string `json:"app_id"`
	GrpcCallbackPort int    `json:"grpc_callback_port"`
	GrpcCallbackHost string `json:"grpc_callback_host"`
}

type MosnRuntimeConfig struct {
	AppManagement          AppConfig                           `json:"app"`
	HelloServiceManagement map[string]hello.HelloConfig        `json:"hellos"`
	ConfigStoreManagement  map[string]configstores.StoreConfig `json:"config_store"`
	RpcManagement          map[string]rpc.RpcConfig            `json:"rpcs"`
	PubSubManagement       map[string]pubsub.Config            `json:"pub_subs"`
	StateManagement        map[string]state.Config             `json:"state"`
	Files                  map[string]file.FileConfig          `json:"file"`
	Oss                    map[string]oss.Config               `json:"oss"`
	LockManagement         map[string]lock.Config              `json:"lock"`
	SequencerManagement    map[string]sequencer.Config         `json:"sequencer"`
	Bindings               map[string]bindings.Metadata        `json:"bindings"`
	SecretStoresManagement map[string]secretstores.Metadata    `json:"secret_store"`
	// <component type,component name,config>
	// e.g. <"super_pubsub","etcd",config>
	CustomComponent map[string]map[string]custom.Config `json:"custom_component,omitempty"`
	Extends         map[string]json.RawMessage          `json:"extends,omitempty"` // extend config
	ExtensionComponentConfig
}

func ParseRuntimeConfig(data json.RawMessage) (*MosnRuntimeConfig, error) {
	cfg := &MosnRuntimeConfig{}
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
