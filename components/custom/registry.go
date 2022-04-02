//
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
	Register(componentType string, factorys ...*ComponentFactory)
	Create(componentType, name string) (Component, error)
}

type ComponentFactory struct {
	Name          string
	FactoryMethod func() Component
}

func NewComponentFactory(name string, f func() Component) *ComponentFactory {
	return &ComponentFactory{
		Name:          name,
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

func (r *componentRegistry) Register(componentType string, fs ...*ComponentFactory) {
	if len(fs) == 0 {
		return
	}
	r.info.AddService(componentType)
	// lazy init
	if _, ok := r.stores[componentType]; !ok {
		r.stores[componentType] = make(map[string]func() Component)
	}
	// register FactoryMethod
	for _, f := range fs {
		r.stores[componentType][f.Name] = f.FactoryMethod
		r.info.RegisterComponent(componentType, f.Name)
	}
}

func (r *componentRegistry) Create(componentType, name string) (Component, error) {
	store, ok := r.stores[componentType]
	if !ok {
		return nil, fmt.Errorf("custom component type %s is not regsitered", componentType)
	}
	if f, ok := store[name]; ok {
		r.info.LoadComponent(componentType, name)
		return f(), nil
	}
	return nil, fmt.Errorf("custom component %s is not regsitered", name)
}
