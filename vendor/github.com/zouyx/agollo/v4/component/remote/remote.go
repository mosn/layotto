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

package remote

import (
	"github.com/zouyx/agollo/v4/env/config"
	"github.com/zouyx/agollo/v4/protocol/http"
)

// ApolloConfig apollo 配置
type ApolloConfig interface {
	// GetNotifyURLSuffix 获取异步更新路径
	GetNotifyURLSuffix(notifications string, config config.AppConfig) string
	// GetSyncURI 获取同步路径
	GetSyncURI(config config.AppConfig, namespaceName string) string
	// Sync 同步获取 Apollo 配置
	Sync(appConfigFunc func() config.AppConfig) []*config.ApolloConfig
	// CallBack 根据 namespace 获取 callback 方法
	CallBack(namespace string) http.CallBack
	// SyncWithNamespace 通过 namespace 同步 apollo 配置
	SyncWithNamespace(namespace string, appConfigFunc func() config.AppConfig) *config.ApolloConfig
}
