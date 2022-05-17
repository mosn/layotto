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

package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaveConfiguration(t *testing.T) {
	ctx := context.Background()
	item1 := &ConfigurationItem{Key: "hello1", Content: "world1"}
	item2 := &ConfigurationItem{Key: "hello2", Content: "world2"}
	saveRequest := &SaveConfigurationRequest{StoreName: "etcd", AppId: "sofa"}
	saveRequest.Items = append(saveRequest.Items, item1)
	saveRequest.Items = append(saveRequest.Items, item2)
	t.Run("save configuration", func(t *testing.T) {
		err := testClient.SaveConfiguration(ctx, saveRequest)
		assert.Nil(t, err)
	})
}
func TestGetConfiguration(t *testing.T) {
	getRequest := &ConfigurationRequestItem{StoreName: "etcd", AppId: "sofa", Keys: []string{"hello1", "hello2"}}
	t.Run("get configuration", func(t *testing.T) {
		resp, err := testClient.GetConfiguration(context.Background(), getRequest)
		assert.Nil(t, err)
		assert.Equal(t, resp[0].Key, "hello1")
		assert.Equal(t, resp[0].Content, "world1")
		assert.Equal(t, resp[1].Key, "hello2")
		assert.Equal(t, resp[1].Content, "world2")
	})
}

func TestDeleteConfiguration(t *testing.T) {
	ctx := context.Background()
	deleteRequest := &ConfigurationRequestItem{StoreName: "etcd", AppId: "sofa", Keys: []string{"hello1", "hello2"}}
	t.Run("delete configuration", func(t *testing.T) {
		err := testClient.DeleteConfiguration(ctx, deleteRequest)
		assert.Nil(t, err)
	})

	getRequest := &ConfigurationRequestItem{StoreName: "etcd", AppId: "sofa", Keys: []string{"hello1", "hello2"}}
	t.Run("get configuration", func(t *testing.T) {
		resp, err := testClient.GetConfiguration(context.Background(), getRequest)
		assert.Nil(t, err)
		assert.Equal(t, 0, len(resp))
	})
}

func TestSubscribeConfiguration(t *testing.T) {
	item := &ConfigurationRequestItem{StoreName: "etcd", AppId: "sofa", Keys: []string{"hello1"}}
	ch := testClient.SubscribeConfiguration(context.Background(), item)
	for wc := range ch {
		assert.Equal(t, wc.Item.Items[0].Key, "hello1")
		assert.Equal(t, wc.Item.Items[0].Content, "Test")
	}
}
