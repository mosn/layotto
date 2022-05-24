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

package pubsub

import (
	"fmt"

	dpubsub "github.com/dapr/components-contrib/pubsub"

	"mosn.io/layotto/components/pkg/info"
)

const serviceName = "pubsub"

// Registry is the pubsub registry with pubsub name as the key
type Registry interface {
	Register(fs ...*Factory)
	Create(compType string) (dpubsub.PubSub, error)
}

type pubsubRegistry struct {
	stores map[string]func() dpubsub.PubSub
	info   *info.RuntimeInfo
}

func NewRegistry(info *info.RuntimeInfo) Registry {
	info.AddService(serviceName)
	return &pubsubRegistry{
		stores: make(map[string]func() dpubsub.PubSub),
		info:   info,
	}
}

func (r *pubsubRegistry) Register(fs ...*Factory) {
	for _, f := range fs {
		r.stores[f.CompType] = f.FactoryMethod
		r.info.RegisterComponent(serviceName, f.CompType)
	}
}

func (r *pubsubRegistry) Create(compType string) (dpubsub.PubSub, error) {
	if f, ok := r.stores[compType]; ok {
		r.info.LoadComponent(serviceName, compType)
		return f(), nil
	}
	return nil, fmt.Errorf("service component %s is not registered", compType)
}
