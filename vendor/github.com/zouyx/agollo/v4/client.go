/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package agollo

import (
	"container/list"
	"github.com/zouyx/agollo/v4/agcache"
	"github.com/zouyx/agollo/v4/agcache/memory"
	"github.com/zouyx/agollo/v4/cluster/roundrobin"
	"github.com/zouyx/agollo/v4/component"
	"github.com/zouyx/agollo/v4/component/log"
	"github.com/zouyx/agollo/v4/component/notify"
	"github.com/zouyx/agollo/v4/component/remote"
	"github.com/zouyx/agollo/v4/component/serverlist"
	"github.com/zouyx/agollo/v4/constant"
	"github.com/zouyx/agollo/v4/env"
	"github.com/zouyx/agollo/v4/env/config"
	jsonFile "github.com/zouyx/agollo/v4/env/file/json"
	"github.com/zouyx/agollo/v4/extension"
	"github.com/zouyx/agollo/v4/protocol/auth/sign"
	"github.com/zouyx/agollo/v4/storage"
	"github.com/zouyx/agollo/v4/utils"
	"github.com/zouyx/agollo/v4/utils/parse/normal"
	"github.com/zouyx/agollo/v4/utils/parse/properties"
	"github.com/zouyx/agollo/v4/utils/parse/yaml"
	"github.com/zouyx/agollo/v4/utils/parse/yml"
	"strconv"
)

var (
	//next try connect period - 60 second
	nextTryConnectPeriod int64 = 60
)

func init() {
	extension.SetCacheFactory(&memory.DefaultCacheFactory{})
	extension.SetLoadBalance(&roundrobin.RoundRobin{})
	extension.SetFileHandler(&jsonFile.FileHandler{})
	extension.SetHTTPAuth(&sign.AuthSignature{})

	// file parser
	extension.AddFormatParser(constant.DEFAULT, &normal.Parser{})
	extension.AddFormatParser(constant.Properties, &properties.Parser{})
	extension.AddFormatParser(constant.YML, &yml.Parser{})
	extension.AddFormatParser(constant.YAML, &yaml.Parser{})
}

var syncApolloConfig = remote.CreateSyncApolloConfig()

// Client apollo 客户端实例
type Client struct {
	initAppConfigFunc func() (*config.AppConfig, error)
	appConfig         *config.AppConfig
	cache             *storage.Cache
}

func (c *Client) getAppConfig() config.AppConfig {
	return *c.appConfig
}

func create() *Client {

	appConfig := env.InitFileConfig()
	return &Client{
		appConfig: appConfig,
	}
}

// Start 根据默认文件启动
func Start() (*Client, error) {
	return StartWithConfig(nil)
}

// StartWithConfig 根据配置启动
func StartWithConfig(loadAppConfig func() (*config.AppConfig, error)) (*Client, error) {
	// 有了配置之后才能进行初始化
	appConfig, err := env.InitConfig(loadAppConfig)
	if err != nil {
		return nil, err
	}

	c := create()
	if appConfig != nil {
		c.appConfig = appConfig
	}

	c.cache = storage.CreateNamespaceConfig(appConfig.NamespaceName)
	appConfig.Init()

	serverlist.InitSyncServerIPList(c.getAppConfig)

	//first sync
	configs := syncApolloConfig.Sync(c.getAppConfig)
	if len(configs) > 0 {
		for _, apolloConfig := range configs {
			c.cache.UpdateApolloConfig(apolloConfig, c.getAppConfig)
		}
	}

	log.Debug("init notifySyncConfigServices finished")

	//start long poll sync config
	configComponent := &notify.ConfigComponent{}
	configComponent.SetAppConfig(c.getAppConfig)
	configComponent.SetCache(c.cache)
	go component.StartRefreshConfig(configComponent)

	log.Info("agollo start finished ! ")

	return c, nil
}

//GetConfig 根据namespace获取apollo配置
func (c *Client) GetConfig(namespace string) *storage.Config {
	return c.GetConfigAndInit(namespace)
}

//GetConfigAndInit 根据namespace获取apollo配置
func (c *Client) GetConfigAndInit(namespace string) *storage.Config {
	if namespace == "" {
		return nil
	}

	config := c.cache.GetConfig(namespace)

	if config == nil {
		//init cache
		storage.CreateNamespaceConfig(namespace)

		//sync config
		syncApolloConfig.SyncWithNamespace(namespace, c.getAppConfig)
	}

	config = c.cache.GetConfig(namespace)

	return config
}

//GetConfigCache 根据namespace获取apollo配置的缓存
func (c *Client) GetConfigCache(namespace string) agcache.CacheInterface {
	config := c.GetConfigAndInit(namespace)
	if config == nil {
		return nil
	}

	return config.GetCache()
}

//GetDefaultConfigCache 获取默认缓存
func (c *Client) GetDefaultConfigCache() agcache.CacheInterface {
	config := c.GetConfigAndInit(storage.GetDefaultNamespace())
	if config != nil {
		return config.GetCache()
	}
	return nil
}

//GetApolloConfigCache 获取默认namespace的apollo配置
func (c *Client) GetApolloConfigCache() agcache.CacheInterface {
	return c.GetDefaultConfigCache()
}

//GetValue 获取配置
func (c *Client) GetValue(key string) string {
	value := c.getConfigValue(key)
	if value == nil {
		return utils.Empty
	}

	return value.(string)
}

//GetStringValue 获取string配置值
func (c *Client) GetStringValue(key string, defaultValue string) string {
	value := c.GetValue(key)
	if value == utils.Empty {
		return defaultValue
	}

	return value
}

//GetIntValue 获取int配置值
func (c *Client) GetIntValue(key string, defaultValue int) int {
	value := c.GetValue(key)

	i, err := strconv.Atoi(value)
	if err != nil {
		log.Debug("convert to int fail!error:", err)
		return defaultValue
	}

	return i
}

//GetFloatValue 获取float配置值
func (c *Client) GetFloatValue(key string, defaultValue float64) float64 {
	value := c.GetValue(key)

	i, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Debug("convert to float fail!error:", err)
		return defaultValue
	}

	return i
}

//GetBoolValue 获取bool 配置值
func (c *Client) GetBoolValue(key string, defaultValue bool) bool {
	value := c.GetValue(key)

	b, err := strconv.ParseBool(value)
	if err != nil {
		log.Debug("convert to bool fail!error:", err)
		return defaultValue
	}

	return b
}

//GetStringSliceValue 获取[]string 配置值
func (c *Client) GetStringSliceValue(key string, defaultValue []string) []string {
	value := c.getConfigValue(key)

	if value == nil {
		return defaultValue
	}
	s, ok := value.([]string)
	if !ok {
		return defaultValue
	}
	return s
}

//GetIntSliceValue 获取[]int 配置值
func (c *Client) GetIntSliceValue(key string, defaultValue []int) []int {
	value := c.getConfigValue(key)

	if value == nil {
		return defaultValue
	}
	s, ok := value.([]int)
	if !ok {
		return defaultValue
	}
	return s
}

func (c *Client) getConfigValue(key string) interface{} {
	cache := c.GetDefaultConfigCache()
	if cache == nil {
		return utils.Empty
	}

	value, err := cache.Get(key)
	if err != nil {
		log.Errorf("get config value fail!key:%s,err:%s", key, err)
		return utils.Empty
	}

	return value
}

// AddChangeListener 增加变更监控
func (c *Client) AddChangeListener(listener storage.ChangeListener) {
	c.cache.AddChangeListener(listener)
}

// RemoveChangeListener 增加变更监控
func (c *Client) RemoveChangeListener(listener storage.ChangeListener) {
	c.cache.RemoveChangeListener(listener)
}

// GetChangeListeners 获取配置修改监听器列表
func (c *Client) GetChangeListeners() *list.List {
	return c.cache.GetChangeListeners()
}

// UseEventDispatch  添加为某些key分发event功能
func (c *Client) UseEventDispatch() {
	c.AddChangeListener(storage.UseEventDispatch())
}
