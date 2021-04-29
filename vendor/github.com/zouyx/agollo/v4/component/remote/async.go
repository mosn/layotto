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
	"net/url"
	"path"
	"time"

	"github.com/zouyx/agollo/v4/component/log"
	"github.com/zouyx/agollo/v4/constant"
	"github.com/zouyx/agollo/v4/env"
	"github.com/zouyx/agollo/v4/env/config"
	"github.com/zouyx/agollo/v4/extension"
	"github.com/zouyx/agollo/v4/protocol/http"
	"github.com/zouyx/agollo/v4/utils"
)

const (
	//notify timeout
	notifyConnectTimeout = 10 * time.Minute //10m

	defaultContentKey = "content"
)

// CreateAsyncApolloConfig 创建异步 apollo 配置
func CreateAsyncApolloConfig() ApolloConfig {
	a := &asyncApolloConfig{}
	a.remoteApollo = a
	return a
}

type asyncApolloConfig struct {
	AbsApolloConfig
}

func (*asyncApolloConfig) GetNotifyURLSuffix(notifications string, config config.AppConfig) string {
	return fmt.Sprintf("notifications/v2?appId=%s&cluster=%s&notifications=%s",
		url.QueryEscape(config.AppID),
		url.QueryEscape(config.Cluster),
		url.QueryEscape(notifications))
}

func (*asyncApolloConfig) GetSyncURI(config config.AppConfig, namespaceName string) string {
	return fmt.Sprintf("configs/%s/%s/%s?releaseKey=%s&ip=%s",
		url.QueryEscape(config.AppID),
		url.QueryEscape(config.Cluster),
		url.QueryEscape(namespaceName),
		url.QueryEscape(config.GetCurrentApolloConfig().GetReleaseKey(namespaceName)),
		utils.GetInternal())
}

func (a *asyncApolloConfig) Sync(appConfigFunc func() config.AppConfig) []*config.ApolloConfig {
	appConfig := appConfigFunc()
	remoteConfigs, err := a.notifyRemoteConfig(appConfigFunc, utils.Empty)

	var apolloConfigs []*config.ApolloConfig
	if err != nil {
		apolloConfigs = loadBackupConfig(appConfig.NamespaceName, appConfig)
	}

	if len(remoteConfigs) == 0 || len(apolloConfigs) > 0 {
		return apolloConfigs
	}
	//只是拉去有变化的配置, 并更新拉取成功的namespace的notify ID
	for _, notifyConfig := range remoteConfigs {
		apolloConfig := a.SyncWithNamespace(notifyConfig.NamespaceName, appConfigFunc)
		if apolloConfig != nil {
			appConfig.GetNotificationsMap().UpdateNotify(notifyConfig.NamespaceName, notifyConfig.NotificationID)
			apolloConfigs = append(apolloConfigs, apolloConfig)
		}
	}
	return apolloConfigs
}

func (*asyncApolloConfig) CallBack(namespace string) http.CallBack {
	return http.CallBack{
		SuccessCallBack:   createApolloConfigWithJSON,
		NotModifyCallBack: touchApolloConfigCache,
		Namespace:         namespace,
	}
}

func (a *asyncApolloConfig) notifyRemoteConfig(appConfigFunc func() config.AppConfig, namespace string) ([]*config.Notification, error) {
	if appConfigFunc == nil {
		panic("can not find apollo config!please confirm!")
	}
	appConfig := appConfigFunc()
	notificationsMap := appConfig.GetNotificationsMap()
	urlSuffix := a.GetNotifyURLSuffix(notificationsMap.GetNotifies(namespace), appConfig)

	connectConfig := &env.ConnectConfig{
		URI:    urlSuffix,
		AppID:  appConfig.AppID,
		Secret: appConfig.Secret,
	}
	connectConfig.Timeout = notifyConnectTimeout
	notifies, err := http.RequestRecovery(appConfig, connectConfig, &http.CallBack{
		SuccessCallBack: func(responseBody []byte, callback http.CallBack) (interface{}, error) {
			return toApolloConfig(responseBody)
		},
		NotModifyCallBack: touchApolloConfigCache,
		Namespace:         namespace,
	})

	if notifies == nil {
		return nil, err
	}

	return notifies.([]*config.Notification), err
}

func touchApolloConfigCache() error {
	return nil
}

func toApolloConfig(resBody []byte) ([]*config.Notification, error) {
	remoteConfig := make([]*config.Notification, 0)

	err := json.Unmarshal(resBody, &remoteConfig)

	if err != nil {
		log.Error("Unmarshal Msg Fail,Error:", err)
		return nil, err
	}
	return remoteConfig, nil
}

func loadBackupConfig(namespace string, appConfig config.AppConfig) []*config.ApolloConfig {
	apolloConfigs := make([]*config.ApolloConfig, 0)
	config.SplitNamespaces(namespace, func(namespace string) {
		c, err := extension.GetFileHandler().LoadConfigFile(appConfig.BackupConfigPath, appConfig.AppID, namespace)
		if err != nil {
			log.Error("LoadConfigFile error, error", err)
			return
		}
		if c == nil {
			return
		}
		apolloConfigs = append(apolloConfigs, c)
	})
	return apolloConfigs
}

func createApolloConfigWithJSON(b []byte, callback http.CallBack) (o interface{}, err error) {
	apolloConfig := &config.ApolloConfig{}
	err = json.Unmarshal(b, apolloConfig)
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
	m, err := parser.Parse(apolloConfig.Configurations[defaultContentKey])
	if err != nil {
		log.Debug("GetContent fail ! error:", err)
	}

	if len(m) > 0 {
		apolloConfig.Configurations = m
	}
	return apolloConfig, nil
}
