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

	injector := DefaultInjector{Container: *container}
	meta := make(map[string]string)
	var items []*ref.Item
	items = append(items, &ref.Item{
		ComponentType: "fake_secret_store",
		Key:           "good-key",
	})
	injector.InjectSecretRef(items, meta)
	assert.Equal(t, meta["good-key"], "life is good")

}
