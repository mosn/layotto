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

package configstores

import (
	"fmt"

	"mosn.io/layotto/components/pkg/info"
)

type Registry interface {
	Register(fs ...*StoreFactory)
	Create(name string) (Store, error)
}

type StoreFactory struct {
	Name          string
	FactoryMethod func() Store
}

func NewStoreFactory(name string, f func() Store) *StoreFactory {
	return &StoreFactory{
		Name:          name,
		FactoryMethod: f,
	}
}

type StoreRegistry struct {
	stores map[string]func() Store
	info   *info.RuntimeInfo
}

func NewRegistry(info *info.RuntimeInfo) Registry {
	info.AddService(ServiceName)
	return &StoreRegistry{
		stores: make(map[string]func() Store),
		info:   info,
	}
}

func (r *StoreRegistry) Register(fs ...*StoreFactory) {
	for _, f := range fs {
		r.stores[f.Name] = f.FactoryMethod
		r.info.RegisterComponent(ServiceName, f.Name)
	}
}

func (r *StoreRegistry) Create(name string) (Store, error) {
	if f, ok := r.stores[name]; ok {
		r.info.LoadComponent(ServiceName, name)
		return f(), nil
	}
	return nil, fmt.Errorf("service component %s is not regsitered", name)
}
