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
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	nacoslog "github.com/nacos-group/nacos-sdk-go/v2/common/logger"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"mosn.io/pkg/log"

	"mosn.io/layotto/components/configstores"
)

type ConfigStore struct {
	client      config_client.IConfigClient
	storeName   string
	appId       string
	namespaceId string
	listener    sync.Map
}

func NewStore() configstores.Store {
	return &ConfigStore{}
}

// Init SetConfig the configuration store.
func (n *ConfigStore) Init(config *configstores.StoreConfig) (err error) {
	if config == nil {
		return errors.New("configuration illegal:no config data")
	}

	// store name, required
	n.storeName = config.StoreName
	if n.storeName == "" {
		return errConfigMissingField("store_mame")
	}

	n.appId = config.AppId
	if n.appId == "" {
		return errConfigMissingField("app_id")
	}

	// parse config metadata
	metadata, err := ParseNacosMetadata(config.Metadata)
	if err != nil {
		return err
	}

	// the nacos's addresses, required if not using acm mode.
	if len(config.Address) == 0 && !metadata.OpenKMS {
		return errConfigMissingField("address")
	}

	n.namespaceId = metadata.NameSpaceId
	// the timeout of connect to nacos, not required
	timeout := defaultTimeout
	if config.TimeOut != "" {
		timeout, err = strconv.Atoi(config.TimeOut)
		if err != nil {
			log.DefaultLogger.Errorf("wrong configuration for time out configuration: %+v, set default value(10s)", config.TimeOut)
			return err
		}
	}
	timeoutMs := uint64(timeout) * uint64(time.Second/time.Millisecond)

	// choose different mode to connect to the nacos server.
	var client config_client.IConfigClient
	if metadata.OpenKMS {
		client, err = n.initWithACM(timeoutMs, metadata)
	} else {
		client, err = n.init(config.Address, timeoutMs, metadata)
	}
	if err != nil {
		return err
	}

	n.client = client
	// replace nacos sdk log
	return n.setupLogger(metadata)
}

// Connect to self built nacos services
func (n *ConfigStore) init(address []string, timeoutMs uint64, metadata *Metadata) (config_client.IConfigClient, error) {
	// 1.create ServerConfigs
	serverConfigs := make([]constant.ServerConfig, 0, len(address))
	for _, v := range address {
		// split the addresses to ip and port
		splitAddr := strings.Split(v, ":")
		if len(splitAddr) != 2 {
			return nil, errors.New("configuration illegal: addresses is not in the format of ip:port")
		}

		ip := splitAddr[0]
		port, err := strconv.Atoi(splitAddr[1])
		if err != nil {
			return nil, errors.New("configuration illegal: can't convert port form string to int type")
		}
		// default use http schema and use nacos as the context
		sc := *constant.NewServerConfig(ip, uint64(port))
		serverConfigs = append(serverConfigs, sc)
	}

	// 2.create client config
	clientConfig := *constant.NewClientConfig(
		constant.WithTimeoutMs(timeoutMs),
		constant.WithNamespaceId(metadata.NameSpaceId),
		constant.WithUsername(metadata.Username),
		constant.WithPassword(metadata.Password),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithCacheDir(metadata.CacheDir),
	)

	// 3.create config client
	// it only creates a client instance but not connect to nacos.
	// so if the address is wrong, the client instance will still be created successfully.
	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// Connect to the nacos service provided by Alibaba Cloud
func (n *ConfigStore) initWithACM(timeoutMs uint64, metadata *Metadata) (config_client.IConfigClient, error) {
	cc := constant.ClientConfig{
		Endpoint:            metadata.Endpoint,
		NamespaceId:         metadata.NameSpaceId,
		RegionId:            metadata.RegionId,
		AccessKey:           metadata.AccessKey,
		SecretKey:           metadata.SecretKey,
		OpenKMS:             true,
		TimeoutMs:           timeoutMs,
		NotLoadCacheAtStart: true,
		CacheDir:            metadata.CacheDir,
	}

	// a more graceful way to create config client
	client, err := clients.CreateConfigClient(map[string]interface{}{
		constant.KEY_CLIENT_CONFIG: cc,
	})

	if err != nil {
		return nil, err
	}

	return client, nil
}

func (n *ConfigStore) setupLogger(metadata *Metadata) error {
	roller := log.DefaultRoller()
	logFilePath := filepath.Join(metadata.LogDir, defaultLogFileName)
	logger, err := log.GetOrCreateLogger(logFilePath, roller)
	if err != nil {
		return err
	}

	errLogger := &log.SimpleErrorLog{
		Logger: logger,
	}

	switch metadata.LogLevel {
	case DEBUG:
		errLogger.Level = log.DEBUG
	case INFO:
		errLogger.Level = log.INFO
	case WARN:
		errLogger.Level = log.WARN
	case ERROR:
		errLogger.Level = log.ERROR
	default:
		return errors.New("unknown log level")
	}

	nacoslog.SetLogger(NewDefaultLogger(errLogger))
	return nil
}

// Get gets configuration from configuration store.
func (n *ConfigStore) Get(ctx context.Context, request *configstores.GetRequest) ([]*configstores.ConfigurationItem, error) {
	// use the configuration's app_name instead of the app_id in request
	// 0. check if illegal
	if request.Group == "" && len(request.Keys) > 0 {
		request.Group = defaultGroup
	}

	// 1. get pagination information
	pagination := n.getPagination(request.Metadata)

	// 2. app level
	if request.Group == "" {
		return n.getAllWithAppId(ctx, pagination)
	}

	// 3.group level
	if len(request.Keys) == 0 {
		return n.getAllWithGroup(ctx, request.Group, pagination)
	}

	// 4.key level
	return n.getAllWithKeys(ctx, request.Group, request.Keys)
}

const (
	PageNo   = "PageNo"
	PageSize = "PageSize"
)

type Pagination struct {
	PageNo   int
	PageSize int
}

func (n *ConfigStore) getPagination(metadata map[string]string) *Pagination {
	res := &Pagination{}
	if v, ok := metadata[PageNo]; ok {
		pageNo, err := strconv.Atoi(v)
		if err != nil {
			return &Pagination{0, 0}
		}
		res.PageNo = pageNo
	}
	if v, ok := metadata[PageSize]; ok {
		pageSize, err := strconv.Atoi(v)
		if err != nil {
			return &Pagination{0, 0}
		}
		res.PageSize = pageSize
	}

	return res
}

func (n *ConfigStore) getAllWithAppId(ctx context.Context, pagination *Pagination) ([]*configstores.ConfigurationItem, error) {
	values, err := n.client.SearchConfig(vo.SearchConfigParam{
		Search:   "accurate",
		AppName:  n.appId,
		PageNo:   pagination.PageNo,
		PageSize: pagination.PageSize,
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

func (n *ConfigStore) getAllWithGroup(ctx context.Context, group string, pagination *Pagination) ([]*configstores.ConfigurationItem, error) {
	values, err := n.client.SearchConfig(vo.SearchConfigParam{
		Search:   "accurate",
		AppName:  n.appId,
		Group:    group,
		PageNo:   pagination.PageNo,
		PageSize: pagination.PageSize,
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

func (n *ConfigStore) getAllWithKeys(ctx context.Context, group string, keys []string) ([]*configstores.ConfigurationItem, error) {
	res := make([]*configstores.ConfigurationItem, 0, len(keys))
	// todo: make more goroutine to search the configurations.
	for _, key := range keys {
		value, err := n.client.GetConfig(vo.ConfigParam{
			DataId:  key,
			Group:   group,
			AppName: n.appId,
		})
		if err != nil {
			log.DefaultLogger.Errorf("fail get key-value,err: %+v", err)
			return nil, err
		}

		// config is not exist
		// nacos does not support an empty content.
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

func (n *ConfigStore) Set(ctx context.Context, request *configstores.SetRequest) error {
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
		if err != nil {
			log.DefaultLogger.Errorf("set key[%+v] failed with error: %+v", configItem.Key, err)
			return err
		}
		if !ok {
			return IllegalParam
		}
	}

	return nil
}

func (n *ConfigStore) Delete(ctx context.Context, request *configstores.DeleteRequest) error {
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
		if err != nil {
			log.DefaultLogger.Errorf("delete key[%+v] failed with error: %+v", key, err)
			return err
		}
		if !ok {
			return IllegalParam
		}

		// remove the config change listening
		n.listener.Delete(subscriberKey{
			group: request.Group,
			key:   key,
		})
	}

	return nil
}

func (n *ConfigStore) Subscribe(request *configstores.SubscribeReq, ch chan *configstores.SubscribeResp) error {
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

	items, err := n.Get(ctx, req)
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

type subscriberKey struct {
	group string
	key   string
}

func (n *ConfigStore) subscribeKey(item *configstores.ConfigurationItem, ch chan *configstores.SubscribeResp) error {
	err := n.client.ListenConfig(vo.ConfigParam{
		DataId:   item.Key,
		Group:    item.Group,
		AppName:  n.appId,
		OnChange: n.subscribeOnChange(ch),
	})

	if err != nil {
		return err
	}

	n.listener.Store(subscriberKey{key: item.Key, group: item.Group}, struct{}{})
	return nil
}

type OnChangeFunc func(namespace, group, dataId, data string)

func (n *ConfigStore) subscribeOnChange(ch chan *configstores.SubscribeResp) OnChangeFunc {
	return func(namespace, group, dataId, data string) {
		// package the listening data.
		resp := &configstores.SubscribeResp{
			StoreName: n.storeName,
			AppId:     n.appId,
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

func (n *ConfigStore) StopSubscribe() {
	// stop listening all subscribed configs
	n.listener.Range(func(key, value any) bool {
		subscribe := key.(subscriberKey)
		if err := n.client.CancelListenConfig(vo.ConfigParam{
			DataId:  subscribe.key,
			Group:   subscribe.group,
			AppName: n.appId,
		}); err != nil {
			log.DefaultLogger.Errorf("nacos StopSubscribe key %s-%s-%s failed", n.appId, subscribe.group, subscribe.key)
			return false
		}

		n.listener.Delete(subscribe)
		return true
	})
}

func (n *ConfigStore) GetDefaultGroup() string {
	return defaultGroup
}

func (n *ConfigStore) GetDefaultLabel() string {
	return defaultLabel
}
