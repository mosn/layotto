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

package secretstores

import (
	"github.com/dapr/components-contrib/secretstores"
	"github.com/pkg/errors"

	"mosn.io/layotto/components/pkg/info"
)

const ServiceName = "secretStore"

type (
	// Registry is used to get registered secret store implementations.
	Registry interface {
		Register(ss ...*SecretStoresFactory)
		Create(compType string) (secretstores.SecretStore, error)
	}

	secretStoreRegistry struct {
		secretStores map[string]func() secretstores.SecretStore
		info         *info.RuntimeInfo
	}
)

// NewRegistry returns a new secret store registry.
func NewRegistry(info *info.RuntimeInfo) Registry {
	info.AddService(ServiceName)
	return &secretStoreRegistry{
		secretStores: map[string]func() secretstores.SecretStore{},
		info:         info,
	}
}

// Register adds one or many new secret stores to the registry.
func (s *secretStoreRegistry) Register(ss ...*SecretStoresFactory) {
	for _, component := range ss {
		s.secretStores[component.CompType] = component.FactoryMethod
		s.info.RegisterComponent(ServiceName, component.CompType)
	}
}

// Create instantiates a secret store based on `name`.
func (s *secretStoreRegistry) Create(compType string) (secretstores.SecretStore, error) {
	if method, ok := s.secretStores[compType]; ok {
		s.info.LoadComponent(ServiceName, compType)
		return method(), nil
	}

	return nil, errors.Errorf("couldn't find secret store %s", compType)
}
