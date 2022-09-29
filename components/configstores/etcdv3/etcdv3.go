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
	"strconv"
	"strings"
	"sync"
	"time"

	"mosn.io/pkg/utils"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mosn.io/pkg/log"

	"mosn.io/layotto/components/configstores"
	"mosn.io/layotto/components/trace"
)

const (
	defaultGroup = "default"
	defaultLabel = "default"
)

type EtcdV3ConfigStore struct {
	client *clientv3.Client
	sync.RWMutex
	subscribeKey map[string]string
	appIdKey     string
	storeName    string
	// cancel is the func, call cancel will stop watching on the appIdKey
	cancel       context.CancelFunc
	watchStarted bool
	watchRespCh  chan *configstores.SubscribeResp
}

func (c *EtcdV3ConfigStore) GetDefaultGroup() string {
	return defaultGroup
}

func (c *EtcdV3ConfigStore) GetDefaultLabel() string {
	return defaultLabel
}

func NewStore() configstores.Store {
	return &EtcdV3ConfigStore{subscribeKey: make(map[string]string), watchRespCh: make(chan *configstores.SubscribeResp)}
}

// Init init the configuration store.
func (c *EtcdV3ConfigStore) Init(config *configstores.StoreConfig) error {
	t, err := strconv.Atoi(config.TimeOut)
	if err != nil {
		log.DefaultLogger.Errorf("wrong configuration for time out configuration: %+v, set default value(10s)", config.TimeOut)
		t = 10
	}
	c.client, err = clientv3.New(clientv3.Config{
		Endpoints:   config.Address,
		DialTimeout: time.Duration(t) * time.Second,
	})
	c.storeName = config.StoreName
	return err
}

func (c *EtcdV3ConfigStore) GetPrimaryKeyWithoutTag(s string) string {
	//key no tag
	if strings.Count(s, "/") == configstores.Tag {
		return s
	}
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '/' {
			return s[:i]
		}
	}
	return s
}

func (c *EtcdV3ConfigStore) GetItemsFromAllKeys(kvs []*mvccpb.KeyValue, targetString []string) []*configstores.ConfigurationItem {
	res := make([]*configstores.ConfigurationItem, 0, 10)
	resMap := make(map[string]int)

	for _, kv := range kvs {
		isTarget := true
		k := strings.Split(string(kv.Key), "/")[1:]
		if len(k) < configstores.Tag {
			continue
		}
		prmKey := c.GetPrimaryKeyWithoutTag(string(kv.Key))
		for i := 0; i < configstores.Tag; i++ {
			if targetString[i] != k[i] && targetString[i] != configstores.All {
				isTarget = false
				break
			}
		}
		if isTarget {
			if index, ok := resMap[prmKey]; !ok {
				item := &configstores.ConfigurationItem{}
				item.Group = k[configstores.Group]
				item.Label = k[configstores.Label]
				item.Key = k[configstores.Key]
				item.Content = string(kv.Value)
				item.Tags = make(map[string]string)
				resMap[prmKey] = len(res)
				res = append(res, item)
			} else {
				res[index].Tags[k[configstores.Tag]] = string(kv.Value)
			}
		}
	}
	return res
}

// Get gets configuration from configuration store.
func (c *EtcdV3ConfigStore) Get(ctx context.Context, req *configstores.GetRequest) ([]*configstores.ConfigurationItem, error) {
	targetString := []string{req.AppId, "*", "*", "*"}
	//TODO: the imp read all keys under app, then do match operation, should change later.
	keyValues, err := c.client.Get(ctx, "/"+req.AppId, clientv3.WithPrefix())
	res := make([]*configstores.ConfigurationItem, 0)
	if err != nil {
		log.DefaultLogger.Errorf("fail get all group key-value,err: %+v", err)
		return nil, err
	}
	targetString[configstores.Group] = req.Group
	targetString[configstores.Label] = req.Label
	if len(req.Keys) == 0 {
		return c.GetItemsFromAllKeys(keyValues.Kvs, targetString), nil
	}
	for _, key := range req.Keys {
		targetString[configstores.Key] = key
		res = append(res, c.GetItemsFromAllKeys(keyValues.Kvs, targetString)...)
	}
	trace.SetExtraComponentInfo(ctx, fmt.Sprintf("method: %+v, store: %+v", "Get", "etcd"))
	return res, nil
}

// Set saves configuration into configuration store.
func (c *EtcdV3ConfigStore) Set(ctx context.Context, req *configstores.SetRequest) error {
	for _, item := range req.Items {
		for _, key := range c.ParseKey(req.AppId, item) {
			_, err := c.client.Put(ctx, key, item.Content)
			if err != nil {
				log.DefaultLogger.Errorf("set key[%+v] failed with error: %+v", key, err)
				return err
			}
		}
	}
	return nil
}

// Delete deletes configuration from configuration store.
func (c *EtcdV3ConfigStore) Delete(ctx context.Context, req *configstores.DeleteRequest) error {
	for _, key := range req.Keys {
		res := "/" + req.AppId + "/" + req.Group + "/" + req.Label + "/" + key
		_, err := c.client.Delete(ctx, res, clientv3.WithPrefix())
		if err != nil {
			log.DefaultLogger.Errorf("delete key[%+v] failed with error: %+v", key, err)
			return err
		}
	}
	return nil
}

func (c *EtcdV3ConfigStore) processWatchResponse(resp *clientv3.WatchResponse) {
	res := &configstores.SubscribeResp{StoreName: c.storeName, AppId: c.appIdKey}
	item := &configstores.ConfigurationItem{}
	if len(resp.Events) == 0 {
		return
	}
	c.RLock()
	defer c.RUnlock()
	if !c.watchStarted {
		return
	}
	for _, events := range resp.Events {
		s := strings.Split(string(events.Kv.Key), "/")[1:]
		if key, ok := c.subscribeKey[string(events.Kv.Key)]; ok {
			item.Group = s[configstores.Group]
			item.Label = s[configstores.Label]
			item.Key = key
			item.Content = string(events.Kv.Value)
			res.Items = append(res.Items, item)
		}
	}
	c.watchRespCh <- res
}
func (c *EtcdV3ConfigStore) watch() {
	// Add watch for propertyKey from lastUpdatedRevision updated after Initializing
	ctx, cancel := context.WithCancel(context.Background())
	c.cancel = cancel
	ch := c.client.Watch(ctx, "/"+c.appIdKey, clientv3.WithCreatedNotify(), clientv3.WithPrefix())
	for wc := range ch {
		c.processWatchResponse(&wc)
	}
}

// Subscribe gets configuration from configuration store and subscribe the updates.
func (c *EtcdV3ConfigStore) Subscribe(req *configstores.SubscribeReq, ch chan *configstores.SubscribeResp) error {
	c.appIdKey = req.AppId
	c.watchRespCh = ch
	c.Lock()
	defer c.Unlock()
	for _, key := range req.Keys {
		s := "/" + req.AppId + "/" + req.Group + "/" + req.Label + "/" + key
		c.subscribeKey[s] = key
	}
	if !c.watchStarted {
		utils.GoWithRecover(func() {
			c.watch()
		}, nil)
		c.watchStarted = true
	}
	return nil
}

func (c *EtcdV3ConfigStore) StopSubscribe() {
	c.Lock()
	defer c.Unlock()
	if !c.watchStarted {
		return
	}
	c.watchStarted = false
	c.cancel()
	close(c.watchRespCh)
}

func (c *EtcdV3ConfigStore) ParseKey(appId string, req *configstores.ConfigurationItem) []string {
	res := make([]string, 0, len(req.Tags))
	res = append(res, "/"+appId+"/"+req.Group+"/"+req.Label+"/"+req.Key)
	for _, tag := range req.Tags {
		res = append(res, "/"+appId+"/"+req.Group+"/"+req.Label+"/"+req.Key+"/"+tag)
	}
	return res
}
