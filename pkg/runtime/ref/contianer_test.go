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

	"mosn.io/layotto/pkg/mock"
	"mosn.io/layotto/pkg/mock/components/secret"
)

func TestRefContainer(t *testing.T) {

	container := NewRefContainer()

	ss := &secret.FakeSecretStore{}
	container.SecretRef["fake"] = ss
	cf := &mock.MockStore{}
	container.ConfigRef["mock"] = cf
	assert.Equal(t, ss, container.getSecretStore("fake"))
	assert.Equal(t, cf, container.getConfigStore("mock"))

}
