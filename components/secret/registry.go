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
package secret

import (
	"fmt"

	"mosn.io/layotto/components/pkg/info"
)

type Registry interface {
	Register(fs ...*WrapperFactory)
	Create(compType string) (Wrapper, error)
}

type WrapperFactory struct {
	CompType      string
	FactoryMethod func() Wrapper
}

func NewWrapperFactory(compType string, f func() Wrapper) *WrapperFactory {
	return &WrapperFactory{
		CompType:      compType,
		FactoryMethod: f,
	}
}

type wrapperRegistry struct {
	stores map[string]func() Wrapper
	info   *info.RuntimeInfo
}

func NewRegistry(info *info.RuntimeInfo) Registry {
	return &wrapperRegistry{
		stores: make(map[string]func() Wrapper),
		info:   info,
	}
}

func (r *wrapperRegistry) Register(fs ...*WrapperFactory) {
	for _, f := range fs {
		r.stores[f.CompType] = f.FactoryMethod
		r.info.RegisterComponent(ServiceName, f.CompType) // 注册组件信息
	}
}

func (r *wrapperRegistry) Create(compType string) (Wrapper, error) {
	if f, ok := r.stores[compType]; ok {
		r.info.LoadComponent(ServiceName, compType) // 加载组件信息
		return f(), nil
	}
	return nil, fmt.Errorf("service component %s is not regsitered", compType)
}
