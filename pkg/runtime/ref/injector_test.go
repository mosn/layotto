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
	"testing"

	"github.com/stretchr/testify/assert"

	"mosn.io/layotto/components/ref"
	"mosn.io/layotto/pkg/mock"
	"mosn.io/layotto/pkg/mock/components/secret"
)

func TestInject(t *testing.T) {

	container := NewRefContainer()

	ss := &secret.FakeSecretStore{}
	container.SecretRef["fake_secret_store"] = ss
	cf := &mock.MockStore{}
	container.ConfigRef["mock_config_store"] = cf

	injector := NewDefaultInjector(container.SecretRef, container.ConfigRef)
	meta := make(map[string]string)

	var items []*ref.SecretRefConfig
	secretRef, err := injector.InjectSecretRef(nil, meta)
	assert.Nil(t, err)
	assert.Equal(t, len(secretRef), 0)
	secretRef, err = injector.InjectSecretRef(items, nil)
	assert.Nil(t, err)
	assert.Equal(t, len(secretRef), 0)
	items = append(items, &ref.SecretRefConfig{
		StoreName: "fake_secret_store",
		Key:       "good-key",
		SubKey:    "good-key",
		InjectAs:  "ref-key",
	})
	items = append(items, &ref.SecretRefConfig{
		StoreName: "fake_secret_store",
		Key:       "good-key",
		SubKey:    "good-key",
	})
	injector.InjectSecretRef(items, meta)
	assert.Equal(t, meta["ref-key"], "life is good")
	assert.Equal(t, meta["good-key"], "life is good")
	secretStoreRef, err := injector.GetSecretStore(&ref.ComponentRefConfig{
		SecretStore: "fake_secret_store",
	})
	assert.Nil(t, err)
	assert.Equal(t, secretStoreRef, ss)

	_, err = injector.GetSecretStore(&ref.ComponentRefConfig{
		SecretStore: "null",
	})
	assert.NotNil(t, err)

	configStoreRef, err := injector.GetConfigStore(&ref.ComponentRefConfig{
		ConfigStore: "mock_config_store",
	})
	assert.Nil(t, err)
	assert.Equal(t, configStoreRef, cf)

	_, err = injector.GetConfigStore(&ref.ComponentRefConfig{
		ConfigStore: "null",
	})
	assert.NotNil(t, err)
}
