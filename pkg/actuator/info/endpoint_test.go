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

package info

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockContributor struct {
}

func (m MockContributor) GetInfo() (info interface{}, err error) {
	return map[string]string{
		"k":  "v",
		"k1": "v1",
	}, nil
}

func TestEndpoint_Handle(t *testing.T) {
	ep := NewEndpoint()
	handle, err := ep.Handle(context.Background(), nil)
	assert.True(t, err == nil)
	assert.True(t, len(handle) == 0)

	AddInfoContributorFunc("test", nil)
	handle, err = ep.Handle(context.Background(), nil)
	assert.True(t, err == nil)
	assert.True(t, len(handle) == 0)

	AddInfoContributor("test", nil)
	handle, err = ep.Handle(context.Background(), nil)
	assert.True(t, err == nil)
	assert.True(t, len(handle) == 0)

	AddInfoContributor("test", MockContributor{})
	handle, err = ep.Handle(context.Background(), nil)
	assert.True(t, err == nil)
	assert.True(t, len(handle) == 1)
	assert.True(t, handle["test"].(map[string]string)["k"] == "v")
	assert.True(t, handle["test"].(map[string]string)["k1"] == "v1")
}
