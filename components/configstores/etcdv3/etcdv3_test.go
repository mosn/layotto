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
package etcdv3

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.etcd.io/etcd/server/v3/embed"

	"mosn.io/layotto/components/configstores"
)

const (
	appId                = "mosn"
	defaultEtcdV3WorkDir = "/tmp/default-dubbo-go-remote.etcd"
)

var etcd EtcdV3ConfigStore

func TestGetPrimaryKeyWithoutTag(t *testing.T) {
	s1 := "/mosn/group1/label1/sofa"
	assert.Equal(t, etcd.GetPrimaryKeyWithoutTag(s1), s1)
	s2 := "/mosn/group1/label1/sofa"
	assert.Equal(t, etcd.GetPrimaryKeyWithoutTag(s2), s1)
}

func TestGetDetailInfoFromResult(t *testing.T) {
	kvs := make([]*mvccpb.KeyValue, 0, 10)
	kv1 := &mvccpb.KeyValue{Key: []byte("/mosn/group1/label1/sofa"), Value: []byte("value1")}
	kv2 := &mvccpb.KeyValue{Key: []byte("/mosn/group1/label1/sofa/tag1"), Value: []byte("tag1")}
	kv3 := &mvccpb.KeyValue{Key: []byte("/mosn/group1/label1/sofa/tag2"), Value: []byte("tag2")}
	kv4 := &mvccpb.KeyValue{Key: []byte("/mosn/group1/label2/sofa"), Value: []byte("value1")}
	kv5 := &mvccpb.KeyValue{Key: []byte("/mosn/group1/label2/sofa/tag2"), Value: []byte("tag2")}
	kvs = append(kvs, kv1)
	kvs = append(kvs, kv2)
	kvs = append(kvs, kv3)
	kvs = append(kvs, kv4)
	kvs = append(kvs, kv5)
	targetStr := []string{appId, "group1", "label1", "sofa"}
	res := etcd.GetItemsFromAllKeys(kvs, targetStr)
	for _, value := range res {
		assert.Equal(t, value.Group, "group1")
		assert.Equal(t, value.Label, "label1")
		assert.Equal(t, value.Key, "sofa")
		assert.Equal(t, value.Content, "value1")
		assert.Equal(t, value.Tags, map[string]string{"tag1": "tag1", "tag2": "tag2"})
	}
	targetStr2 := []string{appId, "*", "label1", "sofa"}
	res = etcd.GetItemsFromAllKeys(kvs, targetStr2)
	for _, value := range res {
		assert.Equal(t, value.Group, "group1")
		assert.Equal(t, value.Label, "label1")
		assert.Equal(t, value.Key, "sofa")
		assert.Equal(t, value.Content, "value1")
		assert.Equal(t, value.Tags, map[string]string{"tag1": "tag1", "tag2": "tag2"})
	}

	targetStr3 := []string{appId, "*", "*", "sofa"}
	res = etcd.GetItemsFromAllKeys(kvs, targetStr3)
	for _, value := range res {
		assert.Equal(t, value.Group, "group1")
		if value.Label == "label1" {
			assert.Equal(t, value.Tags, map[string]string{"tag1": "tag1", "tag2": "tag2"})
		}
		if value.Label == "label2" {
			assert.Equal(t, value.Tags, map[string]string{"tag2": "tag2"})
		}
		assert.Equal(t, value.Key, "sofa")
		assert.Equal(t, value.Content, "value1")

	}
}

var wg sync.WaitGroup

type ClientTestSuite struct {
	suite.Suite

	etcdConfig struct {
		name      string
		endpoints []string
		timeout   time.Duration
		heartbeat int
	}

	etcd  *embed.Etcd
	store configstores.Store
}

// start etcd server
func (suite *ClientTestSuite) SetupSuite() {

	t := suite.T()
	DefaultListenPeerURLs := "http://localhost:2382"
	DefaultListenClientURLs := "http://localhost:2381"
	lpurl, _ := url.Parse(DefaultListenPeerURLs)
	lcurl, _ := url.Parse(DefaultListenClientURLs)
	cfg := embed.NewConfig()
	cfg.LPUrls = []url.URL{*lpurl}
	cfg.LCUrls = []url.URL{*lcurl}
	cfg.Dir = defaultEtcdV3WorkDir
	e, err := embed.StartEtcd(cfg)
	if err != nil {
		t.Fatal(err)
	}
	select {
	case <-e.Server.ReadyNotify():
		t.Log("Server is ready!")
	case <-time.After(60 * time.Second):
		e.Server.Stop() // trigger a shutdown
		t.Logf("Server took too long to start!")
	}

	suite.etcd = e
}

// stop etcd server
func (suite *ClientTestSuite) TearDownSuite() {
	suite.etcd.Close()
	if err := os.RemoveAll(defaultEtcdV3WorkDir); err != nil {
		suite.FailNow(err.Error())
	}
}

func (suite *ClientTestSuite) setUpClient() configstores.Store {
	store := NewStore()
	err := store.Init(&configstores.StoreConfig{StoreName: "etcd", Address: suite.etcdConfig.endpoints, TimeOut: fmt.Sprintf("%d", suite.etcdConfig.timeout)})
	if err != nil {
		suite.T().Fatal(err)
	}
	return store
}

// set up a client for suite
func (suite *ClientTestSuite) SetupTest() {
	c := suite.setUpClient()
	suite.store = c
}

func (suite *ClientTestSuite) Delete() {
	var delReq configstores.DeleteRequest
	delReq.AppId = appId
	delReq.Keys = []string{"sofa"}
	delReq.Group = defaultGroup
	delReq.Label = defaultLabel
	err := suite.store.Delete(context.Background(), &delReq)
	if err != nil {
		suite.T().Fatal(err)
	}

	var req configstores.GetRequest
	req.AppId = appId
	req.Keys = []string{"sofa"}

	resp, err := suite.store.Get(context.Background(), &req)
	assert.Equal(suite.T(), len(resp), 0)
	assert.Nil(suite.T(), err)
}

func (suite *ClientTestSuite) Subscribe() {
	var subReq configstores.SubscribeReq
	var i int
	ch := make(chan *configstores.SubscribeResp)
	subReq.AppId = appId
	subReq.Group = defaultGroup
	subReq.Label = defaultLabel
	subReq.Keys = []string{"sofa"}
	wg.Add(1)
	suite.store.Subscribe(&subReq, ch)
	for event := range ch {
		if i == 0 {
			assert.Equal(suite.T(), event.Items[0].Key, "sofa")
			assert.Equal(suite.T(), event.Items[0].Content, "v1")
			assert.Equal(suite.T(), event.StoreName, "etcd")
			i++
			continue
		}
		assert.Equal(suite.T(), event.Items[0].Key, "sofa")
		assert.Equal(suite.T(), event.Items[0].Content, "")
		assert.Equal(suite.T(), event.StoreName, "etcd")
	}
	wg.Done()
}

func (suite *ClientTestSuite) Get() {
	var req configstores.GetRequest
	req.AppId = appId
	req.Group = defaultGroup
	req.Label = defaultLabel
	req.Keys = []string{"sofa"}
	resp, err := suite.store.Get(context.Background(), &req)
	if err != nil || len(resp) == 0 {
		suite.T().Fatal(err)
	}
	for _, value := range resp {
		assert.Equal(suite.T(), value.Key, "sofa")
		assert.Equal(suite.T(), value.Content, "v1")
		assert.Equal(suite.T(), value.Group, defaultGroup)
		assert.Equal(suite.T(), value.Label, defaultLabel)
	}

	req.Keys = []string{}
	resp, err = suite.store.Get(context.Background(), &req)
	if err != nil || len(resp) == 0 {
		suite.T().Fatal(err)
	}
	for _, value := range resp {
		assert.Equal(suite.T(), value.Key, "sofa")
		assert.Equal(suite.T(), value.Content, "v1")
		assert.Equal(suite.T(), value.Group, defaultGroup)
		assert.Equal(suite.T(), value.Label, defaultLabel)
	}
}

func (suite *ClientTestSuite) Set() {
	var req configstores.SetRequest
	var item configstores.ConfigurationItem
	item.Key = "sofa"
	item.Content = "v1"
	item.Group = defaultGroup
	item.Label = defaultLabel
	req.StoreName = "etcd"
	req.AppId = appId
	req.Items = append(req.Items, &item)
	err := suite.store.Set(context.Background(), &req)
	if err != nil {
		suite.T().Fatal(err)
	}
}
func (suite *ClientTestSuite) TestEtcd() {
	suite.Set()
	suite.Get()
	go suite.Subscribe()
	time.Sleep(1 * time.Second)
	suite.Set()
	suite.Delete()
	suite.store.StopSubscribe()
}

func TestClientSuite(t *testing.T) {
	suite.Run(t, &ClientTestSuite{
		etcdConfig: struct {
			name      string
			endpoints []string
			timeout   time.Duration
			heartbeat int
		}{
			name:      "test",
			endpoints: []string{"localhost:2381"},
			timeout:   time.Second,
			heartbeat: 1,
		},
	})
	wg.Wait()
}
