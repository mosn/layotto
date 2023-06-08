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
	"errors"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	nacoslog "github.com/nacos-group/nacos-sdk-go/v2/common/logger"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"mosn.io/layotto/components/configstores"
	"mosn.io/pkg/log"
	"strconv"
	"strings"
	"time"
)

type NacosConfigStore struct {
	client      config_client.IConfigClient
	storeName   string
	appName     string
	logDir      string
	cacheDir    string
	addresses   []string
	namespaceId string
	listener    *subscriberHolder
}

func NewStore() configstores.Store {
	return &NacosConfigStore{
		listener: newSubscriberHolder(),
	}
}

// Init SetConfig the configuration store.
func (n *NacosConfigStore) Init(config *configstores.StoreConfig) (err error) {
	// 1.parse the config
	if config == nil {
		return errors.New("configuration illegal:no config data")
	}

	// store name, required
	n.storeName = config.StoreName
	if n.storeName == "" {
		return errConfigMissingField("store_mame")
	}

	// the nacos's addresses, required
	if config.Address == nil || len(config.Address) == 0 {
		return errConfigMissingField("address")
	}

	// parse config metadata
	metadata, err := ParseNacosMetadata(config.Metadata)
	if err != nil {
		return err
	}
	n.appName = metadata.AppName
	n.namespaceId = metadata.NameSpaceId

	// the timeout of connect to nacos, not required
	timeout := defaultTimeout
	if config.TimeOut != "" {
		timeout, err = strconv.Atoi(config.TimeOut)
		if err != nil {
			log.DefaultLogger.Errorf("wrong configuration for time out configuration: %+v, set default value(10s)", config.TimeOut)
			timeout = defaultTimeout
		}
	}
	timeoutMs := uint64(timeout) * uint64(time.Second/time.Millisecond)

	// 2.create ServerConfigs
	serverConfigs := make([]constant.ServerConfig, 0, len(config.Address))
	for _, v := range config.Address {
		// split the addresses to ip and port
		splitAddr := strings.Split(v, ":")
		if len(splitAddr) != 2 {
			return errors.New("configuration illegal: addresses is not in the format of ip:port")
		}

		ip := splitAddr[0]
		port, err := strconv.Atoi(splitAddr[1])
		if err != nil {
			return errors.New("configuration illegal: can't convert port form string to int type")
		}
		sc := *constant.NewServerConfig(ip, uint64(port))
		serverConfigs = append(serverConfigs, sc)
		n.addresses = append(n.addresses, v)
	}

	// 3.create client config
	// TODO: support acm mode
	clientConfig := *constant.NewClientConfig(
		constant.WithTimeoutMs(timeoutMs),
		constant.WithNamespaceId(n.namespaceId),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir(defaultLogDir),
		constant.WithCacheDir(defaultCacheDir),
		constant.WithLogLevel(defaultLogLevel),
	)

	// 4.create config client
	// it only creates a client instance but not connect to nacos.
	// so if the address is wrong, the client instance will still be created successfully.
	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		return err
	}

	// 5.set default nacos go-sdk default log
	defaultLogger := NewDefaultLogger(log.DefaultLogger)
	nacoslog.SetLogger(defaultLogger)
	n.client = client

	return nil
}

// Get gets configuration from configuration store.
func (n *NacosConfigStore) Get(ctx context.Context, request *configstores.GetRequest) ([]*configstores.ConfigurationItem, error) {
	// todo: pagenation
	// use the configuration's app_name instead of the app_id in request
	// 0. check if illegal
	if request.Group == "" && len(request.Keys) > 0 {
		request.Group = defaultGroup
	}

	// 1. app level
	if request.Group == "" {
		return n.getAllWithAppId(ctx)
	}

	// 2.group level
	if len(request.Keys) == 0 {
		return n.getAllWithGroup(ctx, request.Group)
	}

	// 3.key level
	// todo: haven't support tag and label
	return n.getAllWithKeys(ctx, request.Group, request.Keys)
}

func (n *NacosConfigStore) getAllWithAppId(ctx context.Context) ([]*configstores.ConfigurationItem, error) {
	values, err := n.client.SearchConfig(vo.SearchConfigParam{
		Search:  "accurate",
		AppName: n.appName,
	})
	if err != nil {
		log.DefaultLogger.Errorf("fail get all app_id key-value,err: %+v", err)
		return nil, err
	}

	res := make([]*configstores.ConfigurationItem, 0, len(values.PageItems))
	for _, v := range values.PageItems {
		config := &configstores.ConfigurationItem{
			Content: v.Content,
			Key:     v.DataId,
			Group:   v.Group,
		}
		res = append(res, config)
	}

	return res, nil
}

func (n *NacosConfigStore) getAllWithGroup(ctx context.Context, group string) ([]*configstores.ConfigurationItem, error) {
	values, err := n.client.SearchConfig(vo.SearchConfigParam{
		Search:  "accurate",
		AppName: n.appName,
		Group:   group,
	})
	if err != nil {
		log.DefaultLogger.Errorf("fail get all group key-value,err: %+v", err)
		return nil, err
	}

	res := make([]*configstores.ConfigurationItem, 0, len(values.PageItems))
	for _, v := range values.PageItems {
		config := &configstores.ConfigurationItem{
			Content: v.Content,
			Key:     v.DataId,
			Group:   v.Group,
		}
		res = append(res, config)
	}

	return res, nil
}

func (n *NacosConfigStore) getAllWithKeys(ctx context.Context, group string, keys []string) ([]*configstores.ConfigurationItem, error) {
	res := make([]*configstores.ConfigurationItem, 0, len(keys))
	// todo: make more goroutine to search the configurations.
	for _, key := range keys {
		value, err := n.client.GetConfig(vo.ConfigParam{
			DataId:  key,
			Group:   group,
			AppName: n.appName,
		})
		if err != nil {
			log.DefaultLogger.Errorf("fail get key-value,err: %+v", err)
			return nil, err
		}

		// config is not exist
		// nacos dose not support an empty content.
		if value == "" {
			continue
		}

		config := &configstores.ConfigurationItem{
			Content: value,
			Key:     key,
			Group:   group,
		}

		res = append(res, config)
	}

	return res, nil
}

func (n *NacosConfigStore) Set(ctx context.Context, request *configstores.SetRequest) error {
	if request.AppId == "" {
		return errParamsMissingField("AppId")
	}

	if len(request.Items) == 0 {
		return errParamsMissingField("Items")
	}

	for _, configItem := range request.Items {
		if configItem.Group == "" {
			return errParamsMissingField("Group")
		}
		ok, err := n.client.PublishConfig(vo.ConfigParam{
			DataId:  configItem.Key,
			Group:   configItem.Group,
			AppName: request.AppId,
			Content: configItem.Content,
		})

		// If the config does not exist, deleting the config will not result in an error.
		if err != nil || !ok {
			log.DefaultLogger.Errorf("set key[%+v] failed with error: %+v", configItem.Key, err)
			return err
		}
	}

	return nil
}

func (n *NacosConfigStore) Delete(ctx context.Context, request *configstores.DeleteRequest) error {
	if request.AppId == "" {
		return errParamsMissingField("AppId")
	}

	if request.Group == "" {
		return errParamsMissingField("Group")
	}

	if len(request.Keys) == 0 {
		return errParamsMissingField("Keys")
	}

	for _, key := range request.Keys {
		ok, err := n.client.DeleteConfig(vo.ConfigParam{
			DataId:  key,
			Group:   request.Group,
			AppName: request.AppId,
		})
		if err != nil || !ok {
			log.DefaultLogger.Errorf("delete key[%+v] failed with error: %+v", key, err)
			return err
		}

		// remove the config change listening
		n.listener.RemoveSubscriberKey(subscriberKey{
			group: request.Group,
			key:   key,
		})
	}

	return nil
}

func (n *NacosConfigStore) Subscribe(request *configstores.SubscribeReq, ch chan *configstores.SubscribeResp) error {
	if request.Group == "" && len(request.Keys) > 0 {
		request.Group = defaultGroup
	}

	ctx := context.Background()
	var err error
	var items []*configstores.ConfigurationItem

	if request.Group == "" { // 1.app level
		items, err = n.getAllWithAppId(ctx)
	} else if len(request.Keys) == 0 { // 2.group level
		items, err = n.getAllWithGroup(ctx, request.Group)
	} else { // 3.key level
		items, err = n.getAllWithKeys(ctx, request.Group, request.Keys)
	}

	if err != nil {
		return err
	}

	for _, item := range items {
		// todo: use errgroup to deal with it concurrently.
		if err := n.subscribeKey(item, ch); err != nil {
			return err
		}
	}

	return nil
}

func (n *NacosConfigStore) subscribeKey(item *configstores.ConfigurationItem, ch chan *configstores.SubscribeResp) error {
	err := n.client.ListenConfig(vo.ConfigParam{
		DataId:  item.Key,
		Group:   item.Group,
		AppName: n.appName,
		OnChange: func(namespace, group, dataId, data string) {
			// package the listening data.
			resp := &configstores.SubscribeResp{
				StoreName: n.storeName,
				AppId:     n.appName,
				Items: []*configstores.ConfigurationItem{
					{
						Key:     dataId,
						Content: data,
						Group:   group,
					},
				},
			}

			ch <- resp
		},
	})

	if err != nil {
		return err
	}

	n.listener.AddSubscriberKey(subscriberKey{key: item.Key, group: item.Group})
	return nil
}

func (n *NacosConfigStore) StopSubscribe() {
	// stop listening all subscribed configs
	keys := n.listener.GetSubscriberKey()
	for _, key := range keys {
		err := n.client.CancelListenConfig(vo.ConfigParam{
			DataId:  key.key,
			Group:   key.group,
			AppName: n.appName,
		})

		if err != nil {
			log.DefaultLogger.Errorf("nacos StopSubscribe key %s-%s-%s failed", n.appName, key.group, key.key)
			return
		}

		n.listener.RemoveSubscriberKey(key)
	}
}

func (n NacosConfigStore) GetDefaultGroup() string {
	return defaultGroup
}

func (n NacosConfigStore) GetDefaultLabel() string {
	return defaultLabel
}
