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

package grpc

import (
	"github.com/dapr/components-contrib/bindings"
	"github.com/dapr/components-contrib/pubsub"
	"github.com/dapr/components-contrib/secretstores"
	"github.com/dapr/components-contrib/state"
	"mosn.io/layotto/components/pkg/common"
	"mosn.io/layotto/pkg/runtime/lifecycle"

	"mosn.io/layotto/components/configstores"
	"mosn.io/layotto/components/custom"
	"mosn.io/layotto/components/file"
	"mosn.io/layotto/components/hello"
	"mosn.io/layotto/components/lock"
	"mosn.io/layotto/components/oss"
	"mosn.io/layotto/components/rpc"
	"mosn.io/layotto/components/sequencer"
)

// ApplicationContext contains all you need to construct your GrpcAPI, such as all the components.
// For example, your `SuperState` GrpcAPI can hold the `StateStores` components and use them to implement your own `Super State API` logic.
type ApplicationContext struct {
	AppId                 string
	Hellos                map[string]hello.HelloService
	ConfigStores          map[string]configstores.Store
	Rpcs                  map[string]rpc.Invoker
	PubSubs               map[string]pubsub.PubSub
	StateStores           map[string]state.Store
	Files                 map[string]file.File
	Oss                   map[string]oss.Oss
	LockStores            map[string]lock.LockStore
	Sequencers            map[string]sequencer.Store
	SendToOutputBindingFn func(name string, req *bindings.InvokeRequest) (*bindings.InvokeResponse, error)
	SecretStores          map[string]secretstores.SecretStore
	DynamicComponents     map[lifecycle.ComponentKey]common.DynamicComponent
	CustomComponent       map[string]map[string]custom.Component
}
