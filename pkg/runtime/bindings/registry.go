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

package bindings

import (
	"fmt"

	"github.com/dapr/components-contrib/bindings"

	"mosn.io/layotto/components/pkg/info"
)

const (
	ServiceName = "bindings"
)

type Registry interface {
	RegisterOutputBinding(fs ...*OutputBindingFactory)
	RegisterInputBinding(fs ...*InputBindingFactory)
	CreateOutputBinding(compType string) (bindings.OutputBinding, error)
	CreateInputBinding(compType string) (bindings.InputBinding, error)
}

type bindingsRegistry struct {
	outputBindingStores map[string]func() bindings.OutputBinding
	inputBindingStores  map[string]func() bindings.InputBinding
	info                *info.RuntimeInfo
}

func NewRegistry(info *info.RuntimeInfo) Registry {
	info.AddService(ServiceName)
	return &bindingsRegistry{
		outputBindingStores: make(map[string]func() bindings.OutputBinding),
		inputBindingStores:  make(map[string]func() bindings.InputBinding),
		info:                info,
	}
}

func (r *bindingsRegistry) RegisterOutputBinding(fs ...*OutputBindingFactory) {
	for _, f := range fs {
		r.outputBindingStores[f.CompType] = f.FactoryMethod
		r.info.RegisterComponent(ServiceName, f.CompType)
	}
}

func (r *bindingsRegistry) RegisterInputBinding(fs ...*InputBindingFactory) {
	for _, f := range fs {
		r.inputBindingStores[f.CompType] = f.FactoryMethod
		r.info.RegisterComponent(ServiceName, f.CompType)
	}
}

func (r *bindingsRegistry) CreateOutputBinding(compType string) (bindings.OutputBinding, error) {
	if f, ok := r.outputBindingStores[compType]; ok {
		r.info.LoadComponent(ServiceName, compType)
		return f(), nil
	}
	return nil, fmt.Errorf("service component %s is not regsitered", compType)
}

func (r *bindingsRegistry) CreateInputBinding(compType string) (bindings.InputBinding, error) {
	if f, ok := r.inputBindingStores[compType]; ok {
		r.info.LoadComponent(ServiceName, compType)
		return f(), nil
	}
	return nil, fmt.Errorf("service component %s is not regsitered", compType)
}
