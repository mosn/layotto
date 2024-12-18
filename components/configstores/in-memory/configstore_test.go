// Copyright 2021 Layotto Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package in_memory

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"mosn.io/layotto/components/configstores"

	"github.com/stretchr/testify/assert"
)

func TestImMemoryConfigStore_Set(t *testing.T) {
	t.Run("set success", func(t *testing.T) {
		params := &configstores.SetRequest{
			AppId: "test-set-app",
			Items: []*configstores.ConfigurationItem{
				{
					Group:   "test-set-group",
					Content: "content",
					Key:     "test-set-key",
				},
			},
		}

		store := NewStore()
		err := store.Set(context.Background(), params)
		assert.Nil(t, err)
	})

	t.Run("set with empty items", func(t *testing.T) {
		store := NewStore()
		params := &configstores.SetRequest{
			AppId: "test-set-app",
			Items: []*configstores.ConfigurationItem{},
		}
		err := store.Set(context.Background(), params)
		assert.EqualError(t, fmt.Errorf("params illegal:item is empty"), err.Error())
	})

	t.Run("test notify listener success", func(t *testing.T) {
		store := &InMemoryConfigStore{
			data:     &sync.Map{},
			listener: &sync.Map{},
		}
		ch := make(chan *configstores.SubscribeResp, 2)

		config := &configstores.StoreConfig{
			AppId:     "test-app",
			StoreName: "test-store",
		}

		err := store.Init(config)
		assert.Nil(t, err)

		subReqParams := &configstores.SubscribeReq{
			AppId: "test-app",
			Group: "group",
			Keys:  []string{"data_id"},
		}

		setReqParams := &configstores.SetRequest{
			AppId:     "test-app",
			StoreName: "test-store",
			Items: []*configstores.ConfigurationItem{
				{
					Group:   "group",
					Content: "content",
					Key:     "data_id",
				},
			},
		}

		err = store.Subscribe(subReqParams, ch)
		assert.Nil(t, err)

		expected := &configstores.SubscribeResp{
			StoreName: store.storeName,
			AppId:     store.appId,
			Items: []*configstores.ConfigurationItem{
				{
					Key:     "data_id",
					Content: "content",
					Group:   "group",
				},
			},
		}

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			i := 0
			for v := range ch {
				i++
				assert.EqualValues(t, expected, v)
			}
			assert.EqualValues(t, i, 3)
		}()
		err = store.Set(context.Background(), setReqParams)
		assert.Nil(t, err)

		err = store.Set(context.Background(), setReqParams)
		assert.Nil(t, err)

		err = store.Set(context.Background(), setReqParams)
		assert.Nil(t, err)
		close(ch)
		wg.Wait()
	})

}

func TestInMemoryConfigStore_Get(t *testing.T) {

	t.Run("get success", func(t *testing.T) {
		store := NewStore()

		params := &configstores.SetRequest{
			AppId: "test-set-app",
			Items: []*configstores.ConfigurationItem{
				{
					Group:   "test-group",
					Content: "content",
					Key:     "test-key",
				},
			},
		}
		err := store.Set(context.Background(), params)
		assert.Nil(t, err)

		getRequestParams := &configstores.GetRequest{
			Group: "test-group",
			Keys:  []string{"test-key"},
		}

		v, err := store.Get(context.Background(), getRequestParams)
		assert.Nil(t, err)
		expect := []*configstores.ConfigurationItem{
			{
				Key:     getRequestParams.Keys[0],
				Group:   getRequestParams.Group,
				Content: "content",
			},
		}
		assert.EqualValues(t, expect, v)
	})

}

func TestInMemoryConfigStore_Subscribe(t *testing.T) {
	t.Run("test subscribe success", func(t *testing.T) {
		store := &InMemoryConfigStore{
			data:     &sync.Map{},
			listener: &sync.Map{},
		}
		ch := make(chan *configstores.SubscribeResp, 2)

		config := &configstores.StoreConfig{
			AppId:     "test-app",
			StoreName: "test-store",
		}

		err := store.Init(config)
		assert.Nil(t, err)

		subReqParams := &configstores.SubscribeReq{
			AppId: "test-app",
			Group: "group",
			Keys:  []string{"data_id"},
		}

		err = store.Subscribe(subReqParams, ch)
		assert.Nil(t, err)
		f, ok := store.listener.Load("data_id")
		assert.True(t, ok)
		assert.NotNil(t, f)

		expected := &configstores.SubscribeResp{
			StoreName: store.storeName,
			AppId:     store.appId,
			Items: []*configstores.ConfigurationItem{
				{
					Key:     "data_id",
					Content: "content",
					Group:   "group",
				},
			},
		}

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			i := 0
			for v := range ch {
				i++
				assert.EqualValues(t, expected, v)
			}
			assert.EqualValues(t, i, 3)
		}()
		f.(OnChangeFunc)("group", "data_id", "content")
		f.(OnChangeFunc)("group", "data_id", "content")
		f.(OnChangeFunc)("group", "data_id", "content")
		close(ch)
		wg.Wait()
	})
}

func TestInMemoryConfigStore_StopSubscribe(t *testing.T) {
	t.Run("test stop subscribe success", func(t *testing.T) {
		store := &InMemoryConfigStore{
			data:     &sync.Map{},
			listener: &sync.Map{},
		}
		ch := make(chan *configstores.SubscribeResp, 2)

		config := &configstores.StoreConfig{
			AppId:     "test-app",
			StoreName: "test-store",
		}

		err := store.Init(config)
		assert.Nil(t, err)

		subReqParams := &configstores.SubscribeReq{
			AppId: "test-app",
			Group: "group",
			Keys:  []string{"data_id"},
		}

		err = store.Subscribe(subReqParams, ch)
		assert.Nil(t, err)

		store.StopSubscribe()
		f, ok := store.listener.Load("data_id")
		assert.False(t, ok)
		assert.Nil(t, f)

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			i := 0
			for range ch {
				i++
			}
			assert.EqualValues(t, 0, i)
		}()

		setReqParams := &configstores.SetRequest{
			AppId:     "test-app",
			StoreName: "test-store",
			Items: []*configstores.ConfigurationItem{
				{
					Group:   "group",
					Content: "content",
					Key:     "data_id",
				},
			},
		}
		err = store.Set(context.Background(), setReqParams)
		assert.Nil(t, err)
		close(ch)
		wg.Wait()
	})
}
