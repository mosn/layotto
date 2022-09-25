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
package custom

import (
	"fmt"

	"mosn.io/layotto/components/pkg/info"
)

type Registry interface {
	Register(kind string, factorys ...*ComponentFactory)
	Create(kind, compType string) (Component, error)
}

type ComponentFactory struct {
	Type          string
	FactoryMethod func() Component
}

func NewComponentFactory(compType string, f func() Component) *ComponentFactory {
	return &ComponentFactory{
		Type:          compType,
		FactoryMethod: f,
	}
}

type componentRegistry struct {
	stores map[string]map[string]func() Component
	info   *info.RuntimeInfo
}

func NewRegistry(info *info.RuntimeInfo) Registry {
	return &componentRegistry{
		stores: make(map[string]map[string]func() Component),
		info:   info,
	}
}

func (r *componentRegistry) Register(kind string, fs ...*ComponentFactory) {
	if len(fs) == 0 {
		return
	}
	r.info.AddService(kind)
	// lazy init
	if _, ok := r.stores[kind]; !ok {
		r.stores[kind] = make(map[string]func() Component)
	}
	// register FactoryMethod
	for _, f := range fs {
		r.stores[kind][f.Type] = f.FactoryMethod
		r.info.RegisterComponent(kind, f.Type)
	}
}

func (r *componentRegistry) Create(kind, compType string) (Component, error) {
	store, ok := r.stores[kind]
	if !ok {
		return nil, fmt.Errorf("custom component kind %s is not regsitered", kind)
	}
	if f, ok := store[compType]; ok {
		r.info.LoadComponent(kind, compType)
		return f(), nil
	}
	return nil, fmt.Errorf("custom component %s is not regsitered", compType)
}
