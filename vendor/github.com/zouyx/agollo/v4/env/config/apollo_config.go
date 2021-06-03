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

package config

import (
	"sync"

	"github.com/zouyx/agollo/v4/utils"
)

// CurrentApolloConfig 当前 apollo 返回的配置信息
type CurrentApolloConfig struct {
	l       sync.RWMutex
	configs map[string]*ApolloConnConfig
}

// CreateCurrentApolloConfig nolint
func CreateCurrentApolloConfig() *CurrentApolloConfig {
	return &CurrentApolloConfig{
		configs: make(map[string]*ApolloConnConfig, 1),
	}
}

//SetCurrentApolloConfig 设置apollo配置
func (c *CurrentApolloConfig) Set(namespace string, connConfig *ApolloConnConfig) {
	c.l.Lock()
	defer c.l.Unlock()

	c.configs[namespace] = connConfig
}

//GetCurrentApolloConfig 获取Apollo链接配置
func (c *CurrentApolloConfig) Get() map[string]*ApolloConnConfig {
	c.l.RLock()
	defer c.l.RUnlock()

	return c.configs
}

//GetCurrentApolloConfigReleaseKey 获取release key
func (c *CurrentApolloConfig) GetReleaseKey(namespace string) string {
	c.l.RLock()
	defer c.l.RUnlock()
	config := c.configs[namespace]
	if config == nil {
		return utils.Empty
	}

	return config.ReleaseKey
}

// ApolloConnConfig apollo链接配置
type ApolloConnConfig struct {
	AppID         string `json:"appId"`
	Cluster       string `json:"cluster"`
	NamespaceName string `json:"namespaceName"`
	ReleaseKey    string `json:"releaseKey"`
	sync.RWMutex
}

// ApolloConfig apollo配置
type ApolloConfig struct {
	ApolloConnConfig
	Configurations map[string]interface{} `json:"configurations"`
}

//Init 初始化
func (a *ApolloConfig) Init(appID string, cluster string, namespace string) {
	a.AppID = appID
	a.Cluster = cluster
	a.NamespaceName = namespace
}
