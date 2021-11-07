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

package client

import (
	"github.com/dapr/components-contrib/bindings"
	"github.com/dapr/components-contrib/pubsub"
	"github.com/dapr/components-contrib/state"
	"mosn.io/layotto/components/configstores"
	"mosn.io/layotto/components/file"
	"mosn.io/layotto/components/hello"
	"mosn.io/layotto/components/lock"
	"mosn.io/layotto/components/rpc"
	"mosn.io/layotto/components/sequencer"
)

var (
	appId                 = ""
	hellos                = make(map[string]hello.HelloService)
	configStores          = make(map[string]configstores.Store)
	rpcs                  = make(map[string]rpc.Invoker)
	pubSubs               = make(map[string]pubsub.PubSub)
	stateStores           = make(map[string]state.Store)
	files                 = make(map[string]file.File)
	lockStores            = make(map[string]lock.LockStore)
	sequencers            = make(map[string]sequencer.Store)
	sendToOutputBindingFn func(name string, req *bindings.InvokeRequest) (*bindings.InvokeResponse, error)
)

func RegisterAppId(myAppId string) {
	appId = myAppId
}

func RegisterHello(name string, service hello.HelloService) {
	hellos[name] = service
}

func RegisterConfigStores(name string, service configstores.Store) {
	configStores[name] = service
}

func RegisterRpcs(name string, service rpc.Invoker) {
	rpcs[name] = service
}

func RegisterPubsubs(name string, service pubsub.PubSub) {
	pubSubs[name] = service
}

func RegisterStateStores(name string, service state.Store) {
	stateStores[name] = service
}

func RegisterFiles(name string, service file.File) {
	files[name] = service
}

func RegisterLockStores(name string, service lock.LockStore) {
	lockStores[name] = service
}

func RegisterSequencers(name string, service sequencer.Store) {
	sequencers[name] = service
}

func RegisterSendToOutputBindingFn(fn func(name string, req *bindings.InvokeRequest) (*bindings.InvokeResponse, error)) {
	sendToOutputBindingFn = fn
}
