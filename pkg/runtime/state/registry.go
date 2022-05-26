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

package state

import (
	"fmt"

	"github.com/dapr/components-contrib/state"

	"mosn.io/layotto/components/pkg/info"
)

const (
	ServiceName = "state"
)

type Registry interface {
	Register(fs ...*Factory)
	Create(compType string) (state.Store, error)
}

type stateRegistry struct {
	stores map[string]func() state.Store
	info   *info.RuntimeInfo
}

// Create a new Registry type variable
func NewRegistry(info *info.RuntimeInfo) Registry {
	info.AddService(ServiceName)
	return &stateRegistry{
		stores: make(map[string]func() state.Store),
		info:   info,
	}
}

// Registration for multiple Factories
func (r *stateRegistry) Register(fs ...*Factory) {
	for _, f := range fs {
		r.stores[f.CompType] = f.FactoryMethod
		r.info.RegisterComponent(ServiceName, f.CompType)
	}
}

// Loading components for a registered Factory
func (r *stateRegistry) Create(compType string) (state.Store, error) {
	if f, ok := r.stores[compType]; ok {
		r.info.LoadComponent(ServiceName, compType)
		return f(), nil
	}
	return nil, fmt.Errorf("service component %s is not regsitered", compType)
}
