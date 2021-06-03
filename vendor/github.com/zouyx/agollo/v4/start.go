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
	"github.com/zouyx/agollo/v4/agcache"
	"github.com/zouyx/agollo/v4/cluster"
	"github.com/zouyx/agollo/v4/component/log"
	"github.com/zouyx/agollo/v4/env/file"
	"github.com/zouyx/agollo/v4/extension"
	"github.com/zouyx/agollo/v4/protocol/auth"
)

//SetSignature 设置自定义 http 授权控件
func SetSignature(auth auth.HTTPAuth) {
	if auth != nil {
		extension.SetHTTPAuth(auth)
	}
}

//SetBackupFileHandler 设置自定义备份文件处理组件
func SetBackupFileHandler(file file.FileHandler) {
	if file != nil {
		extension.SetFileHandler(file)
	}
}

//SetLoadBalance 设置自定义负载均衡组件
func SetLoadBalance(loadBalance cluster.LoadBalance) {
	if loadBalance != nil {
		extension.SetLoadBalance(loadBalance)
	}
}

//SetLogger 设置自定义logger组件
func SetLogger(loggerInterface log.LoggerInterface) {
	if loggerInterface != nil {
		log.InitLogger(loggerInterface)
	}
}

//SetCache 设置自定义cache组件
func SetCache(cacheFactory agcache.CacheFactory) {
	if cacheFactory != nil {
		extension.SetCacheFactory(cacheFactory)
	}
}
