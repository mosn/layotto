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
	"github.com/dapr/components-contrib/secretstores"

	"mosn.io/layotto/components/ref"
)

type DefaultInjector struct {
	Container RefContainer
}

//InjectSecretRef  inject secret to metaData
// TODO: permission control
func (i *DefaultInjector) InjectSecretRef(items []*ref.Item, metaData map[string]string) (map[string]string, error) {
	if metaData == nil {
		metaData = make(map[string]string)
	}
	if len(items) == 0 {
		return metaData, nil
	}

	meta := make(map[string]string)
	for _, item := range items {
		store := i.Container.GetSecretStore(item.ComponentType)
		secret, err := store.GetSecret(secretstores.GetSecretRequest{
			Name: item.Key,
		})
		if err != nil {
			return metaData, err
		}
		for k, v := range secret.Data {
			if item.RefKey == "" {
				meta[k] = v
			} else {
				meta[item.RefKey] = v
			}
		}
	}
	//avoid part of assign because of err
	for k, v := range meta {
		metaData[k] = v
	}
	return metaData, nil
}
