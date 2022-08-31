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
	"google.golang.org/grpc"
	"mosn.io/pkg/log"

	"mosn.io/layotto/components/oss"

	"mosn.io/layotto/components/configstores"
	"mosn.io/layotto/components/custom"
	"mosn.io/layotto/components/file"
	"mosn.io/layotto/components/hello"
	"mosn.io/layotto/components/rpc"
	rgrpc "mosn.io/layotto/pkg/grpc"
	mbindings "mosn.io/layotto/pkg/runtime/bindings"
	runtime_lock "mosn.io/layotto/pkg/runtime/lock"
	"mosn.io/layotto/pkg/runtime/pubsub"
	msecretstores "mosn.io/layotto/pkg/runtime/secretstores"
	runtime_sequencer "mosn.io/layotto/pkg/runtime/sequencer"
	"mosn.io/layotto/pkg/runtime/state"
)

// services encapsulates the service to include in the runtime
type services struct {
	hellos        []*hello.HelloFactory
	configStores  []*configstores.StoreFactory
	rpcs          []*rpc.Factory
	files         []*file.FileFactory
	oss           []*oss.Factory
	pubSubs       []*pubsub.Factory
	states        []*state.Factory
	locks         []*runtime_lock.Factory
	sequencers    []*runtime_sequencer.Factory
	outputBinding []*mbindings.OutputBindingFactory
	inputBinding  []*mbindings.InputBindingFactory
	secretStores  []*msecretstores.SecretStoresFactory
	// Custom components.
	// The key is component kind
	custom map[string][]*custom.ComponentFactory
	extensionComponentFactorys
}

type runtimeOptions struct {
	// services
	services services
	// other config options
	srvMaker rgrpc.NewServer
	errInt   ErrInterceptor
	options  []grpc.ServerOption
	// new grpc api
	apiFactorys []rgrpc.NewGrpcAPI
}

func newRuntimeOptions() *runtimeOptions {
	return &runtimeOptions{
		services: services{
			custom: make(map[string][]*custom.ComponentFactory),
		},
	}
}

type Option func(o *runtimeOptions)

func WithNewServer(f rgrpc.NewServer) Option {
	return func(o *runtimeOptions) {
		o.srvMaker = f
	}
}

func WithGrpcOptions(options ...grpc.ServerOption) Option {
	return func(o *runtimeOptions) {
		o.options = append(o.options, options...)
	}
}

func WithGrpcAPI(apiFuncs ...rgrpc.NewGrpcAPI) Option {
	return func(o *runtimeOptions) {
		o.apiFactorys = append(o.apiFactorys, apiFuncs...)
	}
}

type ErrInterceptor func(err error, format string, args ...interface{})

func WithErrInterceptor(i ErrInterceptor) Option {
	return func(o *runtimeOptions) {
		if o.errInt != nil {
			log.DefaultLogger.Fatalf("the error interceptor was already setted")
		}
		o.errInt = i
	}
}

func WithCustomComponentFactory(kind string, factorys ...*custom.ComponentFactory) Option {
	return func(o *runtimeOptions) {
		if len(factorys) == 0 {
			return
		}
		o.services.custom[kind] = append(o.services.custom[kind], factorys...)
	}
}

func WithHelloFactory(hellos ...*hello.HelloFactory) Option {
	return func(o *runtimeOptions) {
		o.services.hellos = append(o.services.hellos, hellos...)
	}
}

func WithConfigStoresFactory(configStores ...*configstores.StoreFactory) Option {
	return func(o *runtimeOptions) {
		o.services.configStores = append(o.services.configStores, configStores...)
	}
}

func WithRpcFactory(rpcs ...*rpc.Factory) Option {
	return func(o *runtimeOptions) {
		o.services.rpcs = append(o.services.rpcs, rpcs...)
	}
}

func WithOssFactory(oss ...*oss.Factory) Option {
	return func(o *runtimeOptions) {
		o.services.oss = append(o.services.oss, oss...)
	}
}

func WithFileFactory(files ...*file.FileFactory) Option {
	return func(o *runtimeOptions) {
		o.services.files = append(o.services.files, files...)
	}
}

func WithPubSubFactory(factorys ...*pubsub.Factory) Option {
	return func(o *runtimeOptions) {
		o.services.pubSubs = append(o.services.pubSubs, factorys...)
	}
}

func WithLockFactory(factorys ...*runtime_lock.Factory) Option {
	return func(o *runtimeOptions) {
		o.services.locks = append(o.services.locks, factorys...)
	}
}

func WithStateFactory(factorys ...*state.Factory) Option {
	return func(o *runtimeOptions) {
		o.services.states = append(o.services.states, factorys...)
	}
}

// WithInputBindings adds input binding components to the runtime.
func WithInputBindings(factorys ...*mbindings.InputBindingFactory) Option {
	return func(o *runtimeOptions) {
		o.services.inputBinding = append(o.services.inputBinding, factorys...)
	}
}

// WithOutputBindings adds output binding components to the runtime.
func WithOutputBindings(factorys ...*mbindings.OutputBindingFactory) Option {
	return func(o *runtimeOptions) {
		o.services.outputBinding = append(o.services.outputBinding, factorys...)
	}
}

func WithSequencerFactory(factorys ...*runtime_sequencer.Factory) Option {
	return func(o *runtimeOptions) {
		o.services.sequencers = append(o.services.sequencers, factorys...)
	}
}

func WithSecretStoresFactory(factorys ...*msecretstores.SecretStoresFactory) Option {
	return func(o *runtimeOptions) {
		o.services.secretStores = append(o.services.secretStores, factorys...)
	}
}
