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
package dapr

import (
	"github.com/dapr/components-contrib/state"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetResponse2GetStateResponse(t *testing.T) {
	resp := GetResponse2GetStateResponse(&state.GetResponse{
		Data:     []byte("v"),
		ETag:     nil,
		Metadata: make(map[string]string),
	})
	assert.Equal(t, resp.Data, []byte("v"))
	assert.Equal(t, resp.Etag, "")
	assert.True(t, len(resp.Metadata) == 0)
}

func TestGetResponse2BulkStateItem(t *testing.T) {
	itm := GetResponse2BulkStateItem(&state.GetResponse{
		Data:     []byte("v"),
		ETag:     nil,
		Metadata: make(map[string]string),
	}, "key")
	assert.Equal(t, itm.Key, "key")
	assert.Equal(t, itm.Data, []byte("v"))
	assert.Equal(t, itm.Etag, "")
	assert.Equal(t, itm.Error, "")
	assert.True(t, len(itm.Metadata) == 0)
}

func TestBulkGetResponse2BulkStateItem(t *testing.T) {
	t.Run("convert nil", func(t *testing.T) {
		itm := BulkGetResponse2BulkStateItem(nil)
		assert.NotNil(t, itm)
	})
	t.Run("normal", func(t *testing.T) {
		itm := BulkGetResponse2BulkStateItem(&state.BulkGetResponse{
			Key:      "key",
			Data:     []byte("v"),
			ETag:     nil,
			Metadata: nil,
			Error:    "",
		})
		assert.Equal(t, itm.Key, "key")
		assert.Equal(t, itm.Data, []byte("v"))
		assert.Equal(t, itm.Etag, "")
		assert.Equal(t, itm.Error, "")
		assert.True(t, len(itm.Metadata) == 0)
	})
}
