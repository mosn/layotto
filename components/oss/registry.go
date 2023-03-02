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

package oss

import (
	"fmt"

	"mosn.io/layotto/components/pkg/info"
)

type Registry interface {
	Register(fs ...*Factory)
	Create(compType string) (Oss, error)
}

type Factory struct {
	CompType      string
	FactoryMethod func() Oss
}

func NewFactory(compType string, f func() Oss) *Factory {
	return &Factory{
		CompType:      compType,
		FactoryMethod: f,
	}
}

type registry struct {
	oss  map[string]func() Oss
	info *info.RuntimeInfo
}

func NewRegistry(info *info.RuntimeInfo) Registry {
	info.AddService(ServiceName)
	return &registry{
		oss:  make(map[string]func() Oss),
		info: info,
	}
}

func (r *registry) Register(fs ...*Factory) {
	for _, f := range fs {
		r.oss[f.CompType] = f.FactoryMethod
		r.info.RegisterComponent(ServiceName, f.CompType)
	}
}

func (r *registry) Create(compType string) (Oss, error) {
	if f, ok := r.oss[compType]; ok {
		r.info.LoadComponent(ServiceName, compType)
		return f(), nil
	}
	return nil, fmt.Errorf("service component %s is not registered", compType)
}
