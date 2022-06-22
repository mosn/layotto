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

package apollo

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"mosn.io/mosn/pkg/log"

	"mosn.io/layotto/components/configstores"
)

const (
	prod = "prod"
)

// MockRepository implements Repository interface
type MockRepository struct {
	cfg     *repoConfig
	invoked []string
	cache   map[string]map[string]string
}

func (a *MockRepository) Connect() error {
	var err error
	return err
}

func (a *MockRepository) SetConfig(r *repoConfig) {
	a.cfg = r
}

func newMockRepository() *MockRepository {
	return &MockRepository{
		invoked: make([]string, 0, 2),
		cache:   make(map[string]map[string]string),
	}
}

func (a *MockRepository) Get(namespace string, key string) (interface{}, error) {
	if namespace == "wrong" {
		return nil, errors.New("test")
	}

	a.invoked = append(a.invoked, "Get,"+namespace+","+key)
	if _, ok := a.cache[namespace]; !ok {
		a.cache[namespace] = make(map[string]string)
	}
	v := a.cache[namespace][key]
	return v, nil
}

func (a *MockRepository) Range(namespace string, f func(key interface{}, value interface{}) bool) error {
	if namespace == "wrong" {
		return errors.New("test")
	}

	a.invoked = append(a.invoked, "Range,"+namespace)
	for k, v := range a.cache[namespace] {
		if !f(k, v) {
			return nil
		}
	}
	return nil
}

const (
	appId = "testApplication_yang"
)

func (a *MockRepository) AddChangeListener(listener *changeListener) {
	if listener == nil {
		panic("nil listener.")
	}

	a.invoked = append(a.invoked, "AddChangeListener")
}

func (a *MockRepository) Set(namespace string, key string, value string) error {
	if _, ok := a.cache[namespace]; !ok {
		a.cache[namespace] = make(map[string]string)
	}
	a.cache[namespace][key] = value
	return nil
}

func TestConfigStore_read(t *testing.T) {
	// 1. set up
	// inject the MockRepository into a ConfigStore
	store, cfg := setup(t)
	kvRepo := store.kvRepo.(*MockRepository)
	kvRepo.Set("application", "sofa@$prod", "sofa@$prod")
	kvRepo.Set("application", "apollo@$prod", "apollo@$prod")
	kvRepo.Set("dubbo", "dubbo", "dubbo")

	// 2. test the ConfigStore,which has a MockRepository in it
	// init
	err := store.Init(cfg)
	if err != nil {
		t.Error(err)
	}
	// get appid
	assert.True(t, store.GetAppId() == "testApplication_yang")
	assert.True(t, store.GetDefaultGroup() != "")
	assert.True(t, store.GetDefaultLabel() == "")

	// get storeName
	assert.True(t, store.GetStoreName() == "config_demo")

	//	get key
	var req configstores.GetRequest
	req.AppId = appId
	req.Group = defaultGroup
	req.Label = prod
	req.Keys = []string{"sofa"}
	resp, err := store.Get(context.Background(), &req)
	if err != nil || len(resp) == 0 || resp[0].Content != "sofa@$prod" {
		t.Error(err)
	}
	//	 get key under namespace
	req2 := req
	req2.Keys = nil
	resp, err = store.Get(context.Background(), &req2)
	if err != nil || len(resp) != 2 {
		t.Error(err)
	}
	if resp[0].Key != "sofa" {
		resp[0], resp[1] = resp[1], resp[0]
	}
	assert.True(t, resp[0].Content == "sofa@$prod" && resp[1].Content == "apollo@$prod")
	//	 get key under default namespace
	req0 := req
	req0.Group = ""
	req0.Keys = []string{"sofa", "apollo"}
	resp, err = store.Get(context.Background(), &req0)
	if err != nil || len(resp) != 2 {
		t.Error(err)
	}
	if resp[0].Key != "sofa" {
		resp[0], resp[1] = resp[1], resp[0]
	}
	assert.True(t, resp[0].Content == "sofa@$prod" && resp[1].Content == "apollo@$prod")
	//	getAllWithAppId
	req3 := req2
	req3.Group = ""
	resp, err = store.Get(context.Background(), &req3)
	if err != nil || len(resp) != 3 {
		t.Error(err)
	}
	//	subscribe
	var subReq configstores.SubscribeReq
	ch := make(chan *configstores.SubscribeResp)
	subReq.AppId = "testApplication_yang"
	subReq.Group = defaultGroup
	subReq.Label = prod
	subReq.Keys = []string{"sofa"}
	err = store.Subscribe(&subReq, ch)
	if err != nil {
		t.Error(err)
	}
	subReq.Group = defaultGroup
	subReq.Label = ""
	subReq.Keys = []string{}
	err = store.Subscribe(&subReq, ch)
	if err != nil {
		t.Error(err)
	}
	subReq.Group = ""
	err = store.Subscribe(&subReq, ch)
	if err != nil {
		t.Error(err)
	}
}

func TestConfigStore_Init(t *testing.T) {
	t.Run("when token invalid then error", func(t *testing.T) {
		// 1. set up
		// inject the MockRepository into a ConfigStore
		store, cfg := setup(t)
		store.openAPIClient = newMockHttpClient(http.StatusUnauthorized)
		kvRepo := store.kvRepo.(*MockRepository)
		kvRepo.Set("application", "sofa@$prod", "sofa@$prod")
		kvRepo.Set("application", "apollo@$prod", "apollo@$prod")
		kvRepo.Set("dubbo", "dubbo", "dubbo")

		// 2. test the ConfigStore,which has a MockRepository in it
		// init
		log.DefaultLogger.SetLogLevel(log.DEBUG)
		err := store.Init(cfg)
		assert.NotNil(t, err)
	})
	t.Run("when open_api_token blank then error", func(t *testing.T) {
		// 1. set up
		// inject the MockRepository into a ConfigStore
		store, cfg := setup(t)
		cfg.Metadata["open_api_token"] = ""
		store.openAPIClient = newMockHttpClient(http.StatusBadRequest)
		kvRepo := store.kvRepo.(*MockRepository)
		kvRepo.Set("application", "sofa@$prod", "sofa@$prod")
		kvRepo.Set("application", "apollo@$prod", "apollo@$prod")
		kvRepo.Set("dubbo", "dubbo", "dubbo")

		// 2. test the ConfigStore,which has a MockRepository in it
		// init
		log.DefaultLogger.SetLogLevel(log.DEBUG)
		err := store.Init(cfg)
		assert.Error(t, err)
	})

	t.Run("when namespace exist then succeed with debug information", func(t *testing.T) {
		// 1. set up
		// inject the MockRepository into a ConfigStore
		store, cfg := setup(t)
		store.openAPIClient = newMockHttpClient(http.StatusBadRequest)
		kvRepo := store.kvRepo.(*MockRepository)
		kvRepo.Set("application", "sofa@$prod", "sofa@$prod")
		kvRepo.Set("application", "apollo@$prod", "apollo@$prod")
		kvRepo.Set("dubbo", "dubbo", "dubbo")

		// 2. test the ConfigStore,which has a MockRepository in it
		// init
		log.DefaultLogger.SetLogLevel(log.DEBUG)
		err := store.Init(cfg)
		assert.Nil(t, err)
	})
}
func setup(t *testing.T) (*ConfigStore, *configstores.StoreConfig) {
	store := NewStore().(*ConfigStore)
	//mock read client
	kvRepo := newMockRepository()
	store.kvRepo = kvRepo
	store.tagsRepo = newMockRepository()
	// mock write client
	store.openAPIClient = newMockHttpClient(http.StatusOK)
	// prepare config
	cfgJson := `{
    "address": [
        "http://106.54.227.205:8080"
    ],
	"store_name": "config_demo",
    "metadata": {
        "app_id": "testApplication_yang",
        "cluster": "default",
        "namespace_name": "dubbo,product.joe,application",
        "is_backup_config": "true",
        "secret": "6ce3ff7e96a24335a9634fe9abca6d51",
        "open_api_token": "947b0db097d2931ba5bf503f1e33c10394f90d11",
        "open_api_address": "http://106.54.227.205",
        "open_api_user": "apollo"
    }
	}`
	cfg := &configstores.StoreConfig{}
	if err := json.Unmarshal([]byte(cfgJson), &cfg); err != nil {
		t.Error(err)
	}
	return store, cfg
}

func newMockHttpClient(statusCode int) *MockHttpClient {
	return &MockHttpClient{
		statusCode: statusCode,
	}
}

type MockHttpClient struct {
	count      int
	invokedUrl []string
	statusCode int
}

func (m *MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	m.count++
	m.invokedUrl = append(m.invokedUrl, req.URL.String())
	return &http.Response{
		StatusCode: m.statusCode,
		Body:       ioutil.NopCloser(bytes.NewReader(nil)),
	}, nil
}

func (m *MockHttpClient) reset() {
	m.count = 0
	m.invokedUrl = make([]string, 0)
}

func TestConfigStore_write(t *testing.T) {
	// 1. set up
	store, cfg := setup(t)
	// init
	err := store.Init(cfg)
	if err != nil {
		t.Error(err)
	}
	client := store.openAPIClient.(*MockHttpClient)
	client.reset()
	// 2. test set
	var req configstores.SetRequest
	var item configstores.ConfigurationItem
	req.AppId = appId
	item.Key = "sofa"
	item.Content = "v1"
	item.Group = defaultGroup
	item.Label = prod
	req.StoreName = "apollo"
	req.Items = append(req.Items, &item)

	err = store.Set(context.Background(), &req)
	if err != nil {
		t.Error(err)
	}
	// 3. validate set
	assert.True(t, client.count == 3, "client.count is %v", client.count)
	assertUrl := []string{
		"http://106.54.227.205/openapi/v1/envs/DEV/apps/testApplication_yang/clusters/default/namespaces/application/items/sofa@$prod?createIfNotExists=true",
		"http://106.54.227.205/openapi/v1/envs/DEV/apps/testApplication_yang/clusters/default/namespaces/sidecar_config_tags/releases",
		"http://106.54.227.205/openapi/v1/envs/DEV/apps/testApplication_yang/clusters/default/namespaces/application/releases",
	}
	assert.True(t, reflect.DeepEqual(client.invokedUrl, assertUrl))
	// set with tags
	client.reset()
	item.Tags = map[string]string{
		"version": "1.0.0",
		"feature": "nothing",
	}
	err = store.Set(context.Background(), &req)
	if err != nil {
		t.Error(err)
	}
	// validate set with tags
	assert.True(t, client.count == 4, "client.count is %v", client.count)
	assertUrl = []string{
		"http://106.54.227.205/openapi/v1/envs/DEV/apps/testApplication_yang/clusters/default/namespaces/application/items/sofa@$prod?createIfNotExists=true",
		"http://106.54.227.205/openapi/v1/envs/DEV/apps/testApplication_yang/clusters/default/namespaces/sidecar_config_tags/items/application@$sofa@$prod?createIfNotExists=true",
		"http://106.54.227.205/openapi/v1/envs/DEV/apps/testApplication_yang/clusters/default/namespaces/sidecar_config_tags/releases",
		"http://106.54.227.205/openapi/v1/envs/DEV/apps/testApplication_yang/clusters/default/namespaces/application/releases",
	}
	//fmt.Printf("%v", client.invokedUrl)
	assert.True(t, reflect.DeepEqual(client.invokedUrl, assertUrl))

	// set with wrong args
	req2 := req
	req2.AppId = ""
	err = store.Set(context.Background(), &req2)
	assert.NotNil(t, err)
	req3 := req
	req3.Items = nil
	err = store.Set(context.Background(), &req3)
	assert.NotNil(t, err)

	//	 test delete
	client.reset()
	var delReq configstores.DeleteRequest
	delReq.AppId = appId
	delReq.Keys = []string{"sofa"}
	delReq.Group = defaultGroup
	delReq.Label = prod
	err = store.Delete(context.Background(), &delReq)
	if err != nil {
		t.Error(err)
	}
	assert.True(t, client.count == 4, "client.count is %v", client.count)
	//fmt.Printf("%v", client.invokedUrl)
	assertUrl = []string{
		"http://106.54.227.205/openapi/v1/envs/DEV/apps/testApplication_yang/clusters/default/namespaces/application/items/sofa@$prod?key=sofa%40%24prod&operator=apollo",
		"http://106.54.227.205/openapi/v1/envs/DEV/apps/testApplication_yang/clusters/default/namespaces/sidecar_config_tags/items/application@$sofa@$prod?key=application%40%24sofa%40%24prod&operator=apollo",
		"http://106.54.227.205/openapi/v1/envs/DEV/apps/testApplication_yang/clusters/default/namespaces/sidecar_config_tags/releases",
		"http://106.54.227.205/openapi/v1/envs/DEV/apps/testApplication_yang/clusters/default/namespaces/application/releases",
	}
	assert.True(t, reflect.DeepEqual(client.invokedUrl, assertUrl))
}

func TestConfigStore_Init_fail(t *testing.T) {
	// ok
	store, cfg := setup(t)
	err := store.Init(cfg)
	assert.Nil(t, err)
	// no config
	store, _ = setup(t)
	err = store.Init(nil)
	if notNil := assert.NotNil(t, err); notNil {
		assert.True(t, err.Error() != "")
	}
	// no Address
	store, cfg = setup(t)
	cfg.Address = []string{}
	err = store.Init(cfg)
	if notNil := assert.NotNil(t, err); notNil {
		assert.True(t, err.Error() != "")
	}
	// no Metadata
	store, cfg = setup(t)
	cfg.Metadata = nil
	err = store.Init(cfg)
	if notNil := assert.NotNil(t, err); notNil {
		assert.True(t, err.Error() != "")
	}
	// no openAPIAddress
	store, cfg = setup(t)
	cfg.Metadata["open_api_address"] = ""
	err = store.Init(cfg)
	if notNil := assert.NotNil(t, err); notNil {
		assert.True(t, err.Error() != "")
	}
	// no openAPIUser
	store, cfg = setup(t)
	cfg.Metadata["open_api_user"] = ""
	err = store.Init(cfg)
	if notNil := assert.NotNil(t, err); notNil {
		assert.True(t, err.Error() != "")
	}
	// no app_id
	store, cfg = setup(t)
	cfg.Metadata["app_id"] = ""
	err = store.Init(cfg)
	if notNil := assert.NotNil(t, err); notNil {
		assert.True(t, err.Error() != "")
	}
	// is_backup_config wrong
	store, cfg = setup(t)
	cfg.Metadata["is_backup_config"] = "typo"
	err = store.Init(cfg)
	if notNil := assert.NotNil(t, err); notNil {
		assert.True(t, err.Error() != "")
	}
}
