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

package rpc

import (
	"fmt"

	"mosn.io/layotto/components/pkg/info"
)

const ServiceName = "rpc"

type Registry interface {
	Register(fs ...*Factory)
	Create(name string) (Invoker, error)
}

type rpcRegistry struct {
	// Key as implementing component name
	rpc  map[string]FactoryMethod
	info *info.RuntimeInfo
}

type FactoryMethod func() Invoker

func NewRegistry(info *info.RuntimeInfo) Registry {
	info.AddService(ServiceName)
	return &rpcRegistry{
		rpc:  make(map[string]FactoryMethod),
		info: info,
	}
}

func (r rpcRegistry) Register(fs ...*Factory) {
	for _, f := range fs {
		r.rpc[f.Name] = f.Fm
		r.info.RegisterComponent(ServiceName, f.Name)
	}
}

func (r rpcRegistry) Create(name string) (Invoker, error) {
	if f, ok := r.rpc[name]; ok {
		r.info.LoadComponent(ServiceName, name)
		return f(), nil
	}
	return nil, fmt.Errorf("service component %s is not registered", name)
}

type Factory struct {
	Name string
	Fm   FactoryMethod
}

func NewRpcFactory(name string, fm FactoryMethod) *Factory {
	return &Factory{
		Name: name,
		Fm:   fm,
	}
}
