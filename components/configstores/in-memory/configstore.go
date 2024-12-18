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

	"mosn.io/layotto/components/configstores"
	"mosn.io/layotto/components/pkg/actuators"
	"mosn.io/layotto/components/trace"
)

var (
	once               sync.Once
	readinessIndicator *actuators.HealthIndicator
	livenessIndicator  *actuators.HealthIndicator
)

const (
	componentName = "configstore-memory"
	defaultGroup  = "default"
	defaultLabel  = "default"
)

func init() {
	readinessIndicator = actuators.NewHealthIndicator()
	livenessIndicator = actuators.NewHealthIndicator()
}

type InMemoryConfigStore struct {
	data      *sync.Map
	listener  *sync.Map
	storeName string
	appId     string
}

func NewStore() configstores.Store {
	once.Do(func() {
		indicators := &actuators.ComponentsIndicator{ReadinessIndicator: readinessIndicator, LivenessIndicator: livenessIndicator}
		actuators.SetComponentsIndicator(componentName, indicators)
	})
	return &InMemoryConfigStore{
		data:     &sync.Map{},
		listener: &sync.Map{},
	}
}

func (m *InMemoryConfigStore) Init(config *configstores.StoreConfig) error {
	m.appId = config.AppId
	m.storeName = config.StoreName
	readinessIndicator.SetStarted()
	livenessIndicator.SetStarted()
	return nil
}

// Get gets configuration from configuration store.
func (m *InMemoryConfigStore) Get(ctx context.Context, req *configstores.GetRequest) ([]*configstores.ConfigurationItem, error) {

	res := make([]*configstores.ConfigurationItem, 0, len(req.Keys))

	for _, key := range req.Keys {
		value, ok := m.data.Load(key)
		if ok {
			config := &configstores.ConfigurationItem{
				Content: value.(string),
				Key:     key,
				Group:   req.Group,
			}
			res = append(res, config)
		}
	}
	trace.SetExtraComponentInfo(ctx, fmt.Sprintf("method: %+v, store: %+v", "Get", "memory"))
	return res, nil
}

// Set saves configuration into configuration store.
func (m *InMemoryConfigStore) Set(ctx context.Context, req *configstores.SetRequest) error {
	if len(req.Items) == 0 {
		return fmt.Errorf("params illegal:item is empty")
	}
	for _, item := range req.Items {
		m.data.Store(item.Key, item.Content)
		m.notifyChanged(item)
	}
	return nil
}

// Delete deletes configuration from configuration store.
func (m *InMemoryConfigStore) Delete(ctx context.Context, req *configstores.DeleteRequest) error {
	for _, key := range req.Keys {
		m.data.Delete(key)
	}
	return nil
}

// Subscribe gets configuration from configuration store and subscribe the updates.
func (m *InMemoryConfigStore) Subscribe(request *configstores.SubscribeReq, ch chan *configstores.SubscribeResp) error {
	if request.Group == "" && len(request.Keys) > 0 {
		request.Group = defaultGroup
	}

	ctx := context.Background()
	req := &configstores.GetRequest{
		AppId:    request.AppId,
		Group:    request.Group,
		Label:    request.Label,
		Keys:     request.Keys,
		Metadata: request.Metadata,
	}

	for _, key := range req.Keys {
		m.listener.Store(key, m.subscribeOnChange(ch))
	}

	items, err := m.Get(ctx, req)
	if err != nil {
		return err
	}

	for _, item := range items {
		m.notifyChanged(item)
	}

	return nil
}

func (m *InMemoryConfigStore) notifyChanged(item *configstores.ConfigurationItem) {
	f, ok := m.listener.Load(item.Key)
	if ok {
		f.(OnChangeFunc)(item.Group, item.Key, item.Content)
	}
}

type OnChangeFunc func(group, dataId, data string)

func (m *InMemoryConfigStore) subscribeOnChange(ch chan *configstores.SubscribeResp) OnChangeFunc {
	return func(group, dataId, data string) {
		resp := &configstores.SubscribeResp{
			StoreName: m.storeName,
			AppId:     m.appId,
			Items: []*configstores.ConfigurationItem{
				{
					Key:     dataId,
					Content: data,
					Group:   group,
				},
			},
		}

		ch <- resp
	}
}

func (m *InMemoryConfigStore) StopSubscribe() {
	// stop listening all subscribed configs
	m.listener.Range(func(key, value any) bool {
		m.listener.Delete(key)
		return true
	})
}

func (m *InMemoryConfigStore) GetDefaultGroup() string {
	return defaultGroup
}

func (m *InMemoryConfigStore) GetDefaultLabel() string {
	return defaultLabel
}
