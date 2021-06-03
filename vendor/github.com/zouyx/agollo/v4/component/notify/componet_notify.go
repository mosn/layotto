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

package notify

import (
	"github.com/zouyx/agollo/v4/component/remote"
	"github.com/zouyx/agollo/v4/storage"
	"time"

	"github.com/zouyx/agollo/v4/env/config"
)

const (
	longPollInterval = 2 * time.Second //2s
)

//ConfigComponent 配置组件
type ConfigComponent struct {
	appConfigFunc func() config.AppConfig
	cache         *storage.Cache
}

// SetAppConfig nolint
func (c *ConfigComponent) SetAppConfig(appConfigFunc func() config.AppConfig) {
	c.appConfigFunc = appConfigFunc
}

// SetCache nolint
func (c *ConfigComponent) SetCache(cache *storage.Cache) {
	c.cache = cache
}

//Start 启动配置组件定时器
func (c *ConfigComponent) Start() {
	t2 := time.NewTimer(longPollInterval)
	instance := remote.CreateAsyncApolloConfig()
	//long poll for sync
	for {
		select {
		case <-t2.C:
			configs := instance.Sync(c.appConfigFunc)
			for _, apolloConfig := range configs {
				c.cache.UpdateApolloConfig(apolloConfig, c.appConfigFunc)
			}
			t2.Reset(longPollInterval)
		}
	}
}
