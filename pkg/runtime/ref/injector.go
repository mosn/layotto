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
	"context"
	"github.com/dapr/components-contrib/secretstores"
	"mosn.io/layotto/components/configstores"
	"mosn.io/layotto/components/ref"
)

const (
	TypeConfig = "config_store"
	TypeSecret = "secret_store"
)

type DefaultInjector struct {
	Container RefContainer
}

//InjectSecretCmpRef inject secret Component
func (i *DefaultInjector) InjectSecretCmpRef(componentRefs []*ref.ComponentRef) map[string]secretstores.SecretStore {
	cmpRefs := make(map[string]secretstores.SecretStore)
	for _, ref := range componentRefs {
		if ref.Type != TypeSecret {
			continue
		}
		store := i.Container.GetSecretStore(ref.Name)
		if store == nil {
			continue
		}
		cmpRefs[ref.Name] = store
	}
	return cmpRefs
}

func (i *DefaultInjector) InjectConfigCmpRef(componentRefs []*ref.ComponentRef) map[string]configstores.Store {
	cmpRefs := make(map[string]configstores.Store)
	for _, ref := range componentRefs {
		if ref.Type != TypeConfig {
			continue
		}
		store := i.Container.GetConfigStore(ref.Name)
		if store == nil {
			continue
		}
		cmpRefs[ref.Name] = store
	}
	return cmpRefs
}

//InjectSecretRef  inject secret to metaData
func (i *DefaultInjector) InjectSecretRef(items []*ref.RefItem, metaData map[string]string) error {

	meta := make(map[string]string)
	for _, item := range items {
		store := i.Container.GetSecretStore(item.Name)
		secret, err := store.GetSecret(secretstores.GetSecretRequest{
			Name: item.Key,
		})
		if err != nil {
			return err
		}
		for k, v := range secret.Data {
			meta[k] = v
		}
	}
	//avoid part of assign because of err
	for k, v := range meta {
		metaData[k] = v
	}
	return nil
}

//InjectConfigRef inject config
func (i *DefaultInjector) InjectConfigRef(items []*ref.RefItem, metaData map[string]string) error {
	//TODO: how to inject config  subscribe or once ?  It can be configured by the user
	meta := make(map[string]string)
	for _, item := range items {
		store := i.Container.GetConfigStore(item.Name)
		cfgS, err := store.Get(context.Background(), &configstores.GetRequest{
			Keys: []string{item.Key},
		})
		if err != nil {
			return err
		}

		//TODO： we can put the code of subscribe here

		for _, cfg := range cfgS {
			meta[cfg.Key] = cfg.Content
		}
	}
	for k, v := range meta {
		metaData[k] = v
	}
	return nil
}
