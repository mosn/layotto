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

package nacos

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/stretchr/testify/assert"

	"mosn.io/layotto/components/configstores"
	"mosn.io/layotto/components/pkg/mock"
)

const (
	storeName = "test-store"
	appName   = "test-app"
	address   = "127.0.0.1:8848"
)

func getMockNacosClient(t *testing.T) *mock.MockNacosConfigClient {
	ctrl := gomock.NewController(t)
	return mock.NewMockNacosConfigClient(ctrl)
}

func setup(t *testing.T, client Client) *ConfigStore {
	t.Helper()
	store := NewStore()
	// with default namespace and timeout
	config := &configstores.StoreConfig{
		StoreName: storeName,
		Address:   []string{address},
		Metadata: map[string]string{
			appNameKey: appName,
		},
	}

	err := store.Init(config)
	if err != nil {
		t.Fatal(err)
		return nil
	}

	nacosStore := store.(*ConfigStore)
	if client != nil {
		nacosStore.client = client
	}

	return nacosStore
}

func TestNacosConfigStore_Delete(t *testing.T) {
	mockClient := getMockNacosClient(t)
	mockClient.EXPECT().DeleteConfig(gomock.Any()).Return(true, nil)

	store := setup(t, mockClient)
	params := &configstores.DeleteRequest{
		Group: "group",
		Keys:  []string{"key"},
		AppId: "test-delete-app",
	}

	t.Run("delete success", func(t *testing.T) {
		err := store.Delete(context.Background(), params)
		assert.Nil(t, err)
	})

	t.Run("delete without app_id", func(t *testing.T) {
		row := params.AppId
		params.AppId = ""
		err := store.Delete(context.Background(), params)
		assert.Error(t, err)
		params.AppId = row
	})

	t.Run("delete without group", func(t *testing.T) {
		row := params.Group
		params.Group = ""
		err := store.Delete(context.Background(), params)
		assert.Error(t, err)
		params.Group = row
	})

	t.Run("delete with empty keys", func(t *testing.T) {
		row := params.Keys
		params.Keys = []string{}
		err := store.Delete(context.Background(), params)
		assert.Error(t, err)
		params.Keys = row
	})
}

func TestNacosConfigStore_Get(t *testing.T) {
	// Only support get configs from the app_id has been set in store.
	t.Run("test get with other app id", func(t *testing.T) {
		mockClient := getMockNacosClient(t)
		params := &configstores.GetRequest{
			AppId: "test-app1", // different from app stored in the nacos instance
			Group: "test-get-group",
			Keys:  []string{"test-get-key1"},
		}

		content := "content"
		mockClient.EXPECT().GetConfig(gomock.Eq(vo.ConfigParam{
			DataId:  params.Keys[0],
			Group:   params.Group,
			AppName: appName, //  app name that stored in the store instance
		})).Return(content, nil)
		store := setup(t, mockClient)
		get, err := store.Get(context.Background(), params)
		assert.Nil(t, err)
		assert.EqualValues(t, content, get[0].Content)
	})

	t.Run("test success with key level", func(t *testing.T) {
		mockClient := getMockNacosClient(t)
		params := &configstores.GetRequest{
			AppId: appName, // different from app stored in the nacos instance
			Group: "test-get-group",
			Keys:  []string{"test-get-key1"},
		}

		content := "content"
		mockClient.EXPECT().GetConfig(gomock.Eq(vo.ConfigParam{
			DataId:  params.Keys[0],
			Group:   params.Group,
			AppName: appName, //  app name that stored in the store instance
		})).Return(content, nil)
		store := setup(t, mockClient)
		get, err := store.Get(context.Background(), params)
		assert.Nil(t, err)
		assert.EqualValues(t, content, get[0].Content)
	})

	t.Run("test success with app level", func(t *testing.T) {
		mockClient := getMockNacosClient(t)
		params := &configstores.GetRequest{
			AppId: appName, // different from app stored in the nacos instance
		}

		content := "content"
		mockClient.EXPECT().SearchConfig(gomock.Eq(vo.SearchConfigParam{
			Search:  "accurate",
			AppName: appName, //  app name that stored in the store instance
		})).Return(&model.ConfigPage{
			PageItems: []model.ConfigItem{
				{
					DataId:  "data_id",
					Group:   "group",
					Appname: appName,
					Content: content,
				},
			},
		}, nil)
		store := setup(t, mockClient)
		get, err := store.Get(context.Background(), params)
		assert.Nil(t, err)
		assert.EqualValues(t, content, get[0].Content)
		assert.EqualValues(t, "group", get[0].Group)
		assert.EqualValues(t, "data_id", get[0].Key)
	})

	t.Run("test success with group level", func(t *testing.T) {
		mockClient := getMockNacosClient(t)
		params := &configstores.GetRequest{
			AppId: appName, // different from app stored in the nacos instance
			Group: "test-get-group",
			//Keys:  []string{"test-get-key1"}, // without keys
		}

		content := "content"
		mockClient.EXPECT().SearchConfig(gomock.Eq(vo.SearchConfigParam{
			Search:  "accurate",
			AppName: appName, //  app name that stored in the store instance
			Group:   params.Group,
		})).Return(&model.ConfigPage{
			PageItems: []model.ConfigItem{
				{
					DataId:  "data_id",
					Group:   params.Group,
					Appname: appName,
					Content: content,
				},
			},
		}, nil)
		store := setup(t, mockClient)
		get, err := store.Get(context.Background(), params)
		assert.Nil(t, err)
		assert.EqualValues(t, content, get[0].Content)
		assert.EqualValues(t, params.Group, get[0].Group)
		assert.EqualValues(t, "data_id", get[0].Key)
	})

	t.Run("test success with illegal params", func(t *testing.T) {
		mockClient := getMockNacosClient(t)
		params := &configstores.GetRequest{
			AppId: appName, // different from app stored in the nacos instance
			//Group: "test-get-group", // without group
			Keys: []string{"test-get-key1"},
		}

		content := "content"
		mockClient.EXPECT().GetConfig(gomock.Eq(vo.ConfigParam{
			DataId:  params.Keys[0],
			Group:   defaultGroup, // use default group
			AppName: appName,      //  app name that stored in the store instance
		})).Return(content, nil)
		store := setup(t, mockClient)
		get, err := store.Get(context.Background(), params)
		assert.Nil(t, err)
		assert.EqualValues(t, content, get[0].Content)
	})
}

func TestNacosConfigStore_GetDefaultGroup(t *testing.T) {
	store := setup(t, nil)
	group := store.GetDefaultGroup()
	assert.EqualValues(t, defaultGroup, group)
}

func TestNacosConfigStore_GetDefaultLabel(t *testing.T) {
	store := setup(t, nil)
	label := store.GetDefaultLabel()
	assert.EqualValues(t, defaultLabel, label)
}

func TestNacosConfigStore_Init(t *testing.T) {
	const (
		namespace = "test-namespace"
		storeName = "test-store"
		address   = "127.0.0.1:8848"
		appName   = "test-app"
		timeout   = "10" // seconds
	)

	store := NewStore()
	// test success without all params
	config := &configstores.StoreConfig{}
	config.Metadata = map[string]string{
		namespaceIdKey: namespace,
		appNameKey:     appName,
	}
	config.StoreName = storeName
	config.Address = []string{address}
	config.TimeOut = timeout

	t.Run("test success", func(t *testing.T) {
		err := store.Init(config)
		assert.Nil(t, err)
		// check config params
		nacosStore := store.(*ConfigStore)
		assert.EqualValues(t, config.Metadata[namespaceIdKey], nacosStore.namespaceId)
		assert.EqualValues(t, config.Metadata[appNameKey], nacosStore.appName)
		assert.EqualValues(t, config.StoreName, nacosStore.storeName)
	})

	t.Run("test without config", func(t *testing.T) {
		store := NewStore()
		err := store.Init(nil)
		assert.Error(t, err)
	})

	t.Run("test without store name", func(t *testing.T) {
		store := NewStore()
		row := config.StoreName
		config.StoreName = ""
		err := store.Init(config)
		assert.Error(t, err)
		config.StoreName = row
	})

	t.Run("test address", func(t *testing.T) {
		store := NewStore()
		row := config.Address

		// with nil address
		config.Address = nil
		err := store.Init(config)
		assert.Error(t, err)

		// with empty address
		config.Address = []string{}
		err = store.Init(config)
		assert.Error(t, err)

		// with wrong address (address haven't port)
		config.Address = []string{"127.0.0.1"}
		err = store.Init(config)
		assert.Error(t, err)

		config.Address = row
	})

	t.Run("test metadata", func(t *testing.T) {
		store := NewStore()
		// test without metadata (the app name required has be tested in config_test.go)
		row := config.Metadata
		config.Metadata = nil
		err := store.Init(config)
		assert.Error(t, err)
		config.Metadata = row
	})
}

func TestNacosConfigStore_Set(t *testing.T) {
	mockClient := getMockNacosClient(t)
	mockClient.EXPECT().PublishConfig(gomock.Any()).Return(true, nil)

	store := setup(t, mockClient)
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

	t.Run("set success", func(t *testing.T) {
		err := store.Set(context.Background(), params)
		assert.Nil(t, err)
	})

	t.Run("set without app_id", func(t *testing.T) {
		row := params.AppId
		params.AppId = ""
		err := store.Set(context.Background(), params)
		assert.Error(t, err)
		params.AppId = row
	})

	t.Run("set with empty items", func(t *testing.T) {
		row := params.Items
		params.Items = nil
		err := store.Set(context.Background(), params)
		assert.Error(t, err)
		params.Items = row
	})

	t.Run("set without group", func(t *testing.T) {
		row := params.Items[0].Group
		params.Items[0].Group = ""
		err := store.Set(context.Background(), params)
		assert.Error(t, err)
		params.Items[0].Group = row
	})
}

func TestNacosConfigStore_StopSubscribe(t *testing.T) {
	// todo test channel closed.

	req := &configstores.SubscribeReq{
		AppId: appName,
		Group: "test-stop-subscribe-group",
		Keys:  []string{"1", "2", "3"},
	}

	client := getMockNacosClient(t)
	client.EXPECT().CancelListenConfig(gomock.Any()).Return(nil).MaxTimes(len(req.Keys))
	client.EXPECT().ListenConfig(gomock.Any()).Return(nil).MaxTimes(len(req.Keys))
	client.EXPECT().GetConfig(gomock.Any()).Return("content", nil).MaxTimes(len(req.Keys))
	store := setup(t, client)
	// listening for some configs
	ch := make(chan *configstores.SubscribeResp, 10)

	err := store.Subscribe(req, ch)
	assert.Nil(t, err)
	assert.EqualValues(t, len(req.Keys), len(store.listener.keyMap))

	// stop all listening
	store.StopSubscribe()
	assert.EqualValues(t, 0, len(store.listener.keyMap))
}

// 由于vo.ConfigParam中存在函数指针的影响，所以自定义方法去进行 matcher 比较
func EqConfigParam(param vo.ConfigParam) gomock.Matcher {
	return &configParamMatcher{expected: param}
}

type configParamMatcher struct {
	expected vo.ConfigParam
}

func (c *configParamMatcher) Matches(x interface{}) bool {
	v, ok := x.(vo.ConfigParam)
	if !ok {
		return false
	}

	if v.DataId != c.expected.DataId || v.Group != c.expected.Group || v.AppName != c.expected.AppName {
		return false
	}

	return true
}

func (c *configParamMatcher) String() string {
	return fmt.Sprintf("is equal to %v", c.expected)
}

func TestNacosConfigStore_Subscribe(t *testing.T) {
	// Only support get configs from the app_id has been set in store.
	t.Run("test subscribe with other app id", func(t *testing.T) {
		mockClient := getMockNacosClient(t)
		params := &configstores.SubscribeReq{
			AppId: "test-app1", // different from app stored in the nacos instance
			Group: "test-get-group",
			Keys:  []string{"test-get-key1"},
		}

		ch := make(chan *configstores.SubscribeResp)
		content := "content"
		mockClient.EXPECT().GetConfig(gomock.Eq(vo.ConfigParam{
			DataId:  params.Keys[0],
			Group:   params.Group,
			AppName: appName, //  app name that stored in the store instance
		})).Return(content, nil)

		mockClient.EXPECT().ListenConfig(EqConfigParam(vo.ConfigParam{
			DataId:   params.Keys[0],
			Group:    params.Group,
			AppName:  appName, //  app name that stored in the store instance
			OnChange: nil,     // Ignore the impact of the OnChange function
		})).Return(nil)
		store := setup(t, mockClient)
		err := store.Subscribe(params, ch)
		assert.Nil(t, err)
	})

	t.Run("test success with key level", func(t *testing.T) {
		mockClient := getMockNacosClient(t)
		params := &configstores.SubscribeReq{
			AppId: "test-app1", // different from app stored in the nacos instance
			Group: "test-get-group",
			Keys:  []string{"test-get-key1"},
		}

		ch := make(chan *configstores.SubscribeResp)
		content := "content"
		mockClient.EXPECT().GetConfig(gomock.Eq(vo.ConfigParam{
			DataId:  params.Keys[0],
			Group:   params.Group,
			AppName: appName, //  app name that stored in the store instance
		})).Return(content, nil)

		mockClient.EXPECT().ListenConfig(EqConfigParam(vo.ConfigParam{
			DataId:   params.Keys[0],
			Group:    params.Group,
			AppName:  appName, //  app name that stored in the store instance
			OnChange: nil,     // Ignore the impact of the OnChange function
		})).Return(nil)
		store := setup(t, mockClient)
		err := store.Subscribe(params, ch)
		assert.Nil(t, err)
	})

	t.Run("test success with app level", func(t *testing.T) {
		mockClient := getMockNacosClient(t)
		params := &configstores.SubscribeReq{
			AppId: appName, // different from app stored in the nacos instance
			//Group: "test-get-group", // without group
			//Keys:  []string{"test-get-key1"}, // without keys
		}

		mockClient.EXPECT().SearchConfig(gomock.Eq(vo.SearchConfigParam{
			Search:  "accurate",
			AppName: appName, //  app name that stored in the store instance
		})).Return(&model.ConfigPage{
			PageItems: []model.ConfigItem{
				{
					DataId:  "data_id",
					Group:   "group",
					Appname: appName,
					Content: "content",
				},
			},
		}, nil)

		ch := make(chan *configstores.SubscribeResp)
		mockClient.EXPECT().ListenConfig(EqConfigParam(vo.ConfigParam{
			DataId:   "data_id",
			Group:    "group",
			AppName:  appName, //  app name that stored in the store instance
			OnChange: nil,     // Ignore the impact of the OnChange function
		})).Return(nil)
		store := setup(t, mockClient)
		err := store.Subscribe(params, ch)
		assert.Nil(t, err)
	})

	t.Run("test success with group level", func(t *testing.T) {
		mockClient := getMockNacosClient(t)
		params := &configstores.SubscribeReq{
			AppId: appName, // different from app stored in the nacos instance
			Group: "test-get-group",
			//Keys:  []string{"test-get-key1"}, // without keys
		}

		mockClient.EXPECT().SearchConfig(gomock.Eq(vo.SearchConfigParam{
			Search:  "accurate",
			AppName: appName, //  app name that stored in the store instance
			Group:   params.Group,
		})).Return(&model.ConfigPage{
			PageItems: []model.ConfigItem{
				{
					DataId:  "data_id",
					Group:   params.Group,
					Appname: appName,
					Content: "content",
				},
			},
		}, nil)

		ch := make(chan *configstores.SubscribeResp)
		mockClient.EXPECT().ListenConfig(EqConfigParam(vo.ConfigParam{
			DataId:   "data_id",
			Group:    params.Group,
			AppName:  appName, //  app name that stored in the store instance
			OnChange: nil,     // Ignore the impact of the OnChange function
		})).Return(nil)
		store := setup(t, mockClient)
		err := store.Subscribe(params, ch)
		assert.Nil(t, err)
	})

	t.Run("test success with illegal params", func(t *testing.T) {
		mockClient := getMockNacosClient(t)
		params := &configstores.SubscribeReq{
			AppId: appName, // different from app stored in the nacos instance
			//Group: "test-get-group", // without group
			Keys: []string{"test-get- key1"},
		}

		ch := make(chan *configstores.SubscribeResp)
		content := "content"
		mockClient.EXPECT().GetConfig(gomock.Eq(vo.ConfigParam{
			DataId:  params.Keys[0],
			Group:   defaultGroup,
			AppName: appName, //  app name that stored in the store instance
		})).Return(content, nil)

		mockClient.EXPECT().ListenConfig(EqConfigParam(vo.ConfigParam{
			DataId:   params.Keys[0],
			Group:    defaultGroup,
			AppName:  appName, //  app name that stored in the store instance
			OnChange: nil,     // Ignore the impact of the OnChange function
		})).Return(nil)
		store := setup(t, mockClient)
		err := store.Subscribe(params, ch)
		assert.Nil(t, err)
	})

	t.Run("test default subscribe", func(t *testing.T) {
		store := setup(t, nil)
		ch := make(chan *configstores.SubscribeResp, 2)
		fn := store.subscribeOnChange(ch)
		expected := &configstores.SubscribeResp{
			StoreName: store.storeName,
			AppId:     store.appName,
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
		fn(store.namespaceId, "group", "data_id", "content")
		fn(store.namespaceId, "group", "data_id", "content")
		fn(store.namespaceId, "group", "data_id", "content")
		close(ch)
		wg.Wait()
	})
}

func TestNewStore(t *testing.T) {
	store := NewStore()
	assert.NotNil(t, store)
}
