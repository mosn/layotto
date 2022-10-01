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

package ref

import (
	"fmt"

	"github.com/dapr/components-contrib/secretstores"

	"mosn.io/layotto/components/configstores"
	"mosn.io/layotto/components/ref"
)

type DefaultInjector struct {
	Container RefContainer
}

// NewDefaultInjector return a single Inject
func NewDefaultInjector(secretStores map[string]secretstores.SecretStore, configStores map[string]configstores.Store) *DefaultInjector {
	injector := &DefaultInjector{
		Container: RefContainer{
			SecretRef: secretStores,
			ConfigRef: configStores,
		},
	}
	return injector
}

// InjectSecretRef  inject secret to metaData
// TODO: permission control
func (i *DefaultInjector) InjectSecretRef(items []*ref.SecretRefConfig, metaData map[string]string) (map[string]string, error) {
	if metaData == nil {
		metaData = make(map[string]string)
	}
	if len(items) == 0 {
		return metaData, nil
	}

	meta := make(map[string]string)
	for _, item := range items {
		store := i.Container.getSecretStore(item.StoreName)
		secret, err := store.GetSecret(secretstores.GetSecretRequest{
			Name: item.Key,
		})
		if err != nil {
			return metaData, err
		}
		for k, v := range secret.Data {
			if k != item.SubKey {
				continue
			}
			if item.InjectAs == "" {
				meta[k] = v
			} else {
				meta[item.InjectAs] = v
			}
		}
	}
	//avoid part of assign because of err
	for k, v := range meta {
		metaData[k] = v
	}
	return metaData, nil
}

func (i *DefaultInjector) GetConfigStore(cf *ref.ComponentRefConfig) (configstores.Store, error) {
	if cf == nil || cf.ConfigStore == "" {
		return nil, nil
	}
	configStore := i.Container.getConfigStore(cf.ConfigStore)
	if configStore == nil {
		return nil, fmt.Errorf("fail to get configStore:%v", cf.ConfigStore)
	}
	return configStore, nil
}

func (i *DefaultInjector) GetSecretStore(cf *ref.ComponentRefConfig) (secretstores.SecretStore, error) {
	if cf == nil || cf.SecretStore == "" {
		return nil, nil
	}
	secretStore := i.Container.getSecretStore(cf.SecretStore)
	if secretStore == nil {
		return nil, fmt.Errorf("fail to get secretStore:%v", cf.SecretStore)
	}
	return secretStore, nil
}
