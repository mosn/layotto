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
	"encoding/json"
	"fmt"
	"github.com/zouyx/agollo/v4/component/log"
	"github.com/zouyx/agollo/v4/constant"
	"github.com/zouyx/agollo/v4/env/config"
	"github.com/zouyx/agollo/v4/extension"
	"github.com/zouyx/agollo/v4/protocol/http"
	"github.com/zouyx/agollo/v4/utils"
	"net/url"
	"path"
)

// CreateSyncApolloConfig 创建同步获取 Apollo 配置
func CreateSyncApolloConfig() ApolloConfig {
	a := &syncApolloConfig{}
	a.remoteApollo = a
	return a
}

type syncApolloConfig struct {
	AbsApolloConfig
}

func (*syncApolloConfig) GetNotifyURLSuffix(notifications string, config config.AppConfig) string {
	return ""
}

func (*syncApolloConfig) GetSyncURI(config config.AppConfig, namespaceName string) string {
	return fmt.Sprintf("configfiles/json/%s/%s/%s?&ip=%s",
		url.QueryEscape(config.AppID),
		url.QueryEscape(config.Cluster),
		url.QueryEscape(namespaceName),
		utils.GetInternal())
}

func (*syncApolloConfig) CallBack(namespace string) http.CallBack {
	return http.CallBack{
		SuccessCallBack:   processJSONFiles,
		NotModifyCallBack: touchApolloConfigCache,
		Namespace:         namespace,
	}
}

func processJSONFiles(b []byte, callback http.CallBack) (o interface{}, err error) {
	apolloConfig := &config.ApolloConfig{}
	apolloConfig.NamespaceName = callback.Namespace

	configurations := make(map[string]interface{}, 0)
	apolloConfig.Configurations = configurations
	err = json.Unmarshal(b, &apolloConfig.Configurations)

	if utils.IsNotNil(err) {
		return nil, err
	}

	parser := extension.GetFormatParser(constant.ConfigFileFormat(path.Ext(apolloConfig.NamespaceName)))
	if parser == nil {
		parser = extension.GetFormatParser(constant.DEFAULT)
	}

	if parser == nil {
		return apolloConfig, nil
	}
	m, err := parser.Parse(configurations[defaultContentKey])
	if err != nil {
		log.Debug("GetContent fail ! error:", err)
	}

	if len(m) > 0 {
		apolloConfig.Configurations = m
	}
	return apolloConfig, nil
}

func (a *syncApolloConfig) Sync(appConfigFunc func() config.AppConfig) []*config.ApolloConfig {
	appConfig := appConfigFunc()
	configs := make([]*config.ApolloConfig, 0, 8)
	config.SplitNamespaces(appConfig.NamespaceName, func(namespace string) {
		apolloConfig := a.SyncWithNamespace(namespace, appConfigFunc)
		if apolloConfig != nil {
			configs = append(configs, apolloConfig)
			return
		}
		configs = append(configs, loadBackupConfig(appConfig.NamespaceName, appConfig)...)
	})
	return configs
}
