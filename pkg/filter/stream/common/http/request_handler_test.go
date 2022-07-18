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

package http

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetRequestData_isError(t *testing.T) {
	ctx := context.Background()
	data, err := GetRequestData(ctx)
	assert.Nil(t, data)
	assert.Equal(t, "invalid request body", err.Error())

	ctx = context.WithValue(ctx, ContextKeyRequestData{}, []byte("invalid json"))
	data, err = GetRequestData(ctx)
	assert.Nil(t, data)
	assert.Equal(t, "invalid character 'i' looking for beginning of value", err.Error())
}

func Test_GetRequestData_isOk(t *testing.T) {
	expectedByte := []byte("{\"name\":\"id_1\",\"instance_num\":2}")
	ctx := context.Background()
	ctx = context.WithValue(ctx, ContextKeyRequestData{}, expectedByte)
	expected := make(map[string]interface{})
	_ = json.Unmarshal(expectedByte, &expected)
	data, err := GetRequestData(ctx)
	assert.Nil(t, err)
	assert.Equal(t, expected, data)
}
