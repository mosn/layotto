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

package actuator

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockEndpoint struct {
}

func (m *MockEndpoint) Handle(ctx context.Context, params ParamsScanner) (map[string]interface{}, error) {
	return nil, nil
}

// TestActuator test AddEndpoint and GetEndpoint.
func TestActuator(t *testing.T) {
	// get singleton Actuator
	act := GetDefault()
	// reset before test
	act.AddEndpoint("health", nil)

	endpoint, ok := act.GetEndpoint("health")
	assert.True(t, ok)
	assert.Nil(t, endpoint)

	act.AddEndpoint("", nil)
	endpoint, ok = act.GetEndpoint("")
	assert.True(t, ok)
	assert.Nil(t, endpoint)

	ep := &MockEndpoint{}
	act.AddEndpoint("health", ep)
	// reset
	defer func() {
		act.AddEndpoint("health", nil)
	}()
	endpoint, ok = act.GetEndpoint("health")
	assert.True(t, ok)
	assert.Equal(t, endpoint, ep)

}
