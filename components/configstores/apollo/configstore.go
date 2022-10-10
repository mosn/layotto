/*
 * Copyright 2021 Layotto Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package apollo

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"mosn.io/layotto/components/pkg/actuators"

	"mosn.io/pkg/log"

	"mosn.io/layotto/components/configstores"
)

var (
	openAPIClientSingleton = &http.Client{}
	once                   sync.Once
	readinessIndicator     *actuators.HealthIndicator
	livenessIndicator      *actuators.HealthIndicator
)

const (
	defaultGroup  = "application"
	componentName = "apollo"
)

func init() {
	readinessIndicator = actuators.NewHealthIndicator()
	livenessIndicator = actuators.NewHealthIndicator()
}

type ConfigStore struct {
	tagsNamespace  string
	delimiter      string
	openAPIToken   string
	openAPIAddress string
	openAPIUser    string
	env            string
	listener       *changeListener
	kvRepo         Repository
	tagsRepo       Repository
	kvConfig       *repoConfig
	tagsConfig     *repoConfig
	openAPIClient  httpClient
}
type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type httpClientImpl struct {
	client *http.Client
}

func (c *httpClientImpl) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}

func (c *ConfigStore) GetDefaultGroup() string {
	return defaultGroup
}

func (c *ConfigStore) GetDefaultLabel() string {
	return ""
}

func NewStore() configstores.Store {
	registerActuator()
	return &ConfigStore{
		tagsNamespace: defaultTagsNamespace,
		delimiter:     defaultDelimiter,
		env:           defaultEnv,
		kvRepo:        newAgolloRepository(),
		tagsRepo:      newAgolloRepository(),
		openAPIClient: newHttpClient(),
	}
}

func registerActuator() {
	once.Do(func() {
		indicators := &actuators.ComponentsIndicator{ReadinessIndicator: readinessIndicator, LivenessIndicator: livenessIndicator}
		actuators.SetComponentsIndicator(componentName, indicators)
	})
}

func newHttpClient() httpClient {
	return &httpClientImpl{
		client: openAPIClientSingleton,
	}
}

// Init SetConfig the configuration store.
func (c *ConfigStore) Init(config *configstores.StoreConfig) error {
	err := c.doInit(config)
	if err != nil {
		readinessIndicator.ReportError(err.Error())
		livenessIndicator.ReportError(err.Error())
	}
	readinessIndicator.SetStarted()
	livenessIndicator.SetStarted()
	return err
}

func (c *ConfigStore) doInit(config *configstores.StoreConfig) error {
	// 1. validate and parse config
	if config == nil {
		return ErrNoConfig
	}
	// Metadata,required
	metadata := config.Metadata
	if len(metadata) == 0 {
		return errConfigMissingField("metadata")
	}
	// Address,required
	if len(config.Address) == 0 || config.Address[0] == "" {
		return errConfigMissingField("address")
	}
	addr := config.Address[0]
	// is_backup_config,not required
	// whether backup config after fetch config from apollo
	s, ok := metadata["is_backup_config"]
	var isBackupConfig = defaultIsBackupConfig
	var err error
	if ok && s != "" {
		isBackupConfig, err = strconv.ParseBool(s)
		if err != nil {
			return err
		}
	}
	// app_id,required
	appId, ok := metadata[configKeyAppId]
	if !ok || appId == "" {
		return errConfigMissingField(configKeyAppId)
	}
	// open_api_token,required
	c.openAPIToken = metadata["open_api_token"]
	if c.openAPIToken == "" {
		return errConfigMissingField("open_api_token")
	}
	// open_api_address,not required
	c.openAPIAddress = metadata["open_api_address"]
	if c.openAPIAddress == "" {
		return errConfigMissingField("open_api_address")
	}
	// open_api_user,required
	c.openAPIUser = metadata["open_api_user"]
	if c.openAPIUser == "" {
		return errConfigMissingField("open_api_user")
	}
	// TODO make 'env' configurable
	// 2. SetConfig client
	kvRepoConfig := &repoConfig{
		addr:           addr,
		appId:          appId,
		storeName:      config.StoreName,
		env:            c.env,
		cluster:        metadata["cluster"],
		namespaceName:  metadata["namespace_name"],
		isBackupConfig: isBackupConfig,
		// secret,not required
		secret: metadata["secret"],
	}
	c.kvConfig = kvRepoConfig
	c.kvRepo.SetConfig(kvRepoConfig)
	err = c.kvRepo.Connect()
	if err != nil {
		return err
	}

	// 3. SetConfig client for tags query
	tagsRepoConfig := *kvRepoConfig
	tagsRepoConfig.namespaceName = c.tagsNamespace
	c.tagsConfig = &tagsRepoConfig
	err = c.initTagsClient(&tagsRepoConfig)
	if err != nil {
		return err
	}
	// 4. SetConfig listener
	listener := newChangeListener(c)
	c.listener = listener
	c.kvRepo.AddChangeListener(listener)
	return nil
}

func (c *ConfigStore) GetAppId() string {
	if c.kvConfig == nil {
		return ""
	}
	return c.kvConfig.appId
}

func (c *ConfigStore) GetStoreName() string {
	if c.kvConfig == nil {
		return ""
	}
	return c.kvConfig.storeName
}

// Get gets configuration from configuration store.
func (c *ConfigStore) Get(ctx context.Context, req *configstores.GetRequest) ([]*configstores.ConfigurationItem, error) {
	// TODO forced pagination
	// 0. check if illegal
	if len(req.Keys) > 0 && req.Group == "" {
		req.Group = defaultNamespace
	}
	// 1. app level
	if req.Group == "" {
		return c.getAllWithAppId()
	}
	// 2. group level
	if len(req.Keys) == 0 {
		return c.getAllWithNamespace(req.Group)
	}
	// 3. group+key+label level
	return c.getKeys(req.Group, req.Keys, req.Label)
}

func (c *ConfigStore) getAllTags(group string, keyWithLabel string) (tags map[string]string, err error) {
	res := make(map[string]string)
	//	1. concatenate group+key+label
	k := c.concatenateKeyForTag(group, keyWithLabel)
	//	2. query
	value, err := c.tagsRepo.Get(c.tagsNamespace, k)
	if err != nil || value == nil || value == "" {
		//	it means no tag
		return res, nil
	}
	//	3. convert
	err = json.Unmarshal([]byte(fmt.Sprintf("%v", value)), &res)
	return res, err
}

// Set saves configuration into configuration store.
func (c *ConfigStore) Set(ctx context.Context, req *configstores.SetRequest) error {
	// 1. check params
	if req.AppId == "" {
		return errParamsMissingField("AppId")
	}
	if len(req.Items) == 0 {
		return errParamsMissingField("Items")
	}
	// 2. loop set
	groupMap := make(map[string]struct{})
	for _, itm := range req.Items {
		// 2.1. set kv
		err := c.setItem(req.AppId, itm)
		if err != nil {
			return err
		}
		groupMap[itm.Group] = struct{}{}
		// 2.2. set tags
		if len(itm.Tags) == 0 {
			continue
		}
		tagItm := &configstores.ConfigurationItem{}
		tagItm.Key = c.concatenateKeyForTag(itm.Group, itm.Key)
		tagItm.Group = c.tagsNamespace
		tagItm.Label = itm.Label
		// set tags value
		data, err := json.Marshal(itm.Tags)
		if err != nil {
			return nil
		}
		tagItm.Content = string(data)
		c.setItem(req.AppId, tagItm)
	}
	// 3. commit tagsNamespace
	err := c.commit(c.env, req.AppId, c.tagsConfig.cluster, c.tagsNamespace)
	if err != nil {
		return err
	}
	// 4. commit kv namespace
	for g := range groupMap {
		err := c.commit(c.env, req.AppId, c.kvConfig.cluster, g)
		if err != nil {
			return err
		}
	}
	// TODO 5. write cache
	return nil
}

// Delete deletes configuration from configuration store.
func (c *ConfigStore) Delete(ctx context.Context, req *configstores.DeleteRequest) error {
	// 1. check params
	if req.AppId == "" {
		return errParamsMissingField("AppId")
	}
	if len(req.Keys) == 0 {
		return errParamsMissingField("Keys")
	}
	if req.Group == "" {
		req.Group = defaultNamespace
	}
	// 2. loop delete
	for _, k := range req.Keys {
		// 2.1. delete item
		err := c.deleteItem(c.env, req.AppId, c.kvConfig.cluster, req.Group, k, req.Label)
		if err != nil {
			return err
		}
		//	2.2. delete tags
		groupAndKey := c.concatenateKeyForTag(req.Group, k)
		err = c.deleteItem(c.env, req.AppId, c.kvConfig.cluster, c.tagsNamespace, groupAndKey, req.Label)
		if err != nil {
			return err
		}
	}
	// 3. commit tagsNamespace
	err := c.commit(c.env, req.AppId, c.tagsConfig.cluster, c.tagsNamespace)
	if err != nil {
		return err
	}
	// 4. commit kv namespace
	err = c.commit(c.env, req.AppId, c.kvConfig.cluster, req.Group)
	// TODO 5. write cache
	return err
}

// Subscribe gets configuration from configuration store and subscribe the updates.
func (c *ConfigStore) Subscribe(req *configstores.SubscribeReq, ch chan *configstores.SubscribeResp) error {
	// 0. check if illegal
	if len(req.Keys) > 0 && req.Group == "" {
		req.Group = defaultNamespace
	}
	// 1. app level
	if len(req.Keys) == 0 && req.Group == "" {
		split := strings.Split(c.kvConfig.namespaceName, ",")
		// loop every namespace in config
		for _, ns := range split {
			if ns == "" {
				continue
			}
			err := c.listener.addByTopic(ns, "", ch)
			if err != nil {
				return err
			}
		}
		return nil
	}
	// 2. group level
	if len(req.Keys) == 0 {
		err := c.listener.addByTopic(req.Group, "", ch)
		return err
	}
	// 3. key level
	for _, k := range req.Keys {
		err := c.listener.addByTopic(req.Group, c.concatenateKey(k, req.Label), ch)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *ConfigStore) StopSubscribe() {
	//	TODO  Now the api layer only supports single connection and does not support multi-connection.
	//	 If it supports multiple connections in the future, we can use a context to cancel specific connections
	c.listener.reset()
}

func (c *ConfigStore) getKeys(group string, keys []string, label string) ([]*configstores.ConfigurationItem, error) {
	log.DefaultLogger.Debugf("getKeys start.namespace : %v, keys : %v, label : %v", group, keys, label)
	// 1. prepare suffix
	suffix := ""
	if label != "" {
		suffix = c.delimiter + label
	}
	res := make([]*configstores.ConfigurationItem, 0, 10)
	// 2. loop query
	for _, k := range keys {
		keyWithLabel := k + suffix
		//query value
		value, err := c.kvRepo.Get(group, keyWithLabel)
		if err != nil {
			//log error and ignore this key
			log.DefaultLogger.Errorf("error when querying configuration :%v", err)
			continue
		}
		item := &configstores.ConfigurationItem{}
		item.Group = group
		item.Label = label
		item.Key = k
		item.Content = fmt.Sprintf("%v", value)
		// query tags
		item.Tags, err = c.getAllTags(group, keyWithLabel)
		if err != nil {
			log.DefaultLogger.Errorf("error when querying tags :%v", err)
		}
		res = append(res, item)
	}

	return res, nil
}

func (c *ConfigStore) getAllWithAppId() ([]*configstores.ConfigurationItem, error) {
	log.DefaultLogger.Debugf("getAllWithAppId start.namespace:%v", c.kvConfig.namespaceName)
	split := strings.Split(c.kvConfig.namespaceName, ",")
	res := make([]*configstores.ConfigurationItem, 0, 10)
	// loop every namespace in config
	for _, ns := range split {
		items, err := c.getAllWithNamespace(ns)
		res = append(res, items...)
		if err != nil {
			return res, err
		}
	}
	return res, nil
}

func (c *ConfigStore) getAllWithNamespace(group string) ([]*configstores.ConfigurationItem, error) {
	log.DefaultLogger.Debugf("getAllWithNamespace start.namespace:%v", group)
	res := make([]*configstores.ConfigurationItem, 0, 10)
	// 1. loop query
	err := c.kvRepo.Range(group, func(key, value interface{}) bool {
		// 1.1. convert
		item := &configstores.ConfigurationItem{}
		item.Group = group
		k := key.(string)
		if k == "" {
			//	never happen
			log.DefaultLogger.Errorf("find configuration item with blank key under namespace:%v", group)
		} else {
			split := strings.Split(k, defaultDelimiter)
			item.Key = split[0]
			if len(split) > 1 {
				item.Label = split[1]
			}
			// 1.2. query tags
			tags, err := c.getAllTags(group, k)
			if err != nil {
				log.DefaultLogger.Errorf("error when querying tags :%v", err)
			} else {
				item.Tags = tags
			}
		}
		item.Content = fmt.Sprintf("%v", value)
		// 1.3. append result.
		res = append(res, item)
		//continue
		return true
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *ConfigStore) setItem(appId string, item *configstores.ConfigurationItem) error {
	// 1. put request
	keyWithLabel := c.concatenateKey(item.Key, item.Label)
	setUrl := fmt.Sprintf(setUrlTpl, c.openAPIAddress, c.env, appId, c.kvConfig.cluster, item.Group, keyWithLabel)
	// add body
	reqBody := map[string]string{
		"key":                      keyWithLabel,
		"value":                    item.Content,
		"dataChangeCreatedBy":      c.openAPIUser,
		"dataChangeLastModifiedBy": c.openAPIUser,
	}
	reqBodyJson, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", setUrl, strings.NewReader(string(reqBodyJson)))
	if err != nil {
		return err
	}
	// add params
	q := req.URL.Query()
	q.Add("createIfNotExists", "true")
	req.URL.RawQuery = q.Encode()
	// add headers
	c.addHeaderForOpenAPI(req)
	// do put request
	resp, err := c.openAPIClient.Do(req)
	// 2. parse
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	return err
}

func (c *ConfigStore) addHeaderForOpenAPI(req *http.Request) {
	//https://www.apolloconfig.com/#/zh/usage/apollo-open-api-platform?id=_3211-%e4%bf%ae%e6%94%b9%e9%85%8d%e7%bd%ae%e6%8e%a5%e5%8f%a3
	//Http Header中增加一个Authorization字段，字段值为申请的token
	//Http Header的Content-Type字段需要设置成application/json;charset=UTF-8
	req.Header.Add("Authorization", c.openAPIToken)
	req.Header.Add("Content-Type", `application/json;charset=UTF-8`)
}

func (c *ConfigStore) concatenateKey(key, label string) string {
	if label == "" {
		return key
	}
	return key + c.delimiter + label
}

func (c *ConfigStore) concatenateKeyForTag(group, keyWithLabel string) string {
	if keyWithLabel == "" {
		return ""
	}
	return group + c.delimiter + keyWithLabel
}

func (c *ConfigStore) splitKey(keyWithLabel string) (key, label string) {
	if keyWithLabel == "" {
		return "", ""
	}
	res := strings.Split(keyWithLabel, c.delimiter)
	if len(res) < 2 {
		return res[0], ""
	}
	return res[0], res[1]
}

func (c *ConfigStore) commit(env string, appId string, cluster string, namespace string) error {
	// 1. post request
	commitUrl := fmt.Sprintf(commitUrlTpl, c.openAPIAddress, env, appId, cluster, namespace)
	// add body
	reqBody := map[string]string{
		"releaseTitle": time.Now().Format("01-02-2006"),
		"releasedBy":   c.openAPIUser,
	}
	reqBodyJson, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", commitUrl, strings.NewReader(string(reqBodyJson)))
	if err != nil {
		return err
	}
	// add headers
	c.addHeaderForOpenAPI(req)
	// do request
	resp, err := c.openAPIClient.Do(req)
	// 2. parse
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	return err
}

func (c *ConfigStore) deleteItem(env string, appId string, cluster string, group string, key string, label string) error {
	// 1. delete request
	keyWithLabel := c.concatenateKey(key, label)
	deleteUrl := fmt.Sprintf(deleteUrlTpl, c.openAPIAddress, env, appId, cluster, group, keyWithLabel)
	req, err := http.NewRequest("DELETE", deleteUrl, nil)
	if err != nil {
		return err
	}
	// add params
	q := req.URL.Query()
	q.Add("key", keyWithLabel)
	q.Add("operator", c.openAPIUser)
	req.URL.RawQuery = q.Encode()
	// add headers
	c.addHeaderForOpenAPI(req)
	// do request
	resp, err := c.openAPIClient.Do(req)
	// 2. parse
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	return err
}

func (c *ConfigStore) initTagsClient(tagCfg *repoConfig) error {
	// 1. create if not exist
	err := c.createNamespace(c.env, tagCfg.appId, tagCfg.cluster, c.tagsNamespace)
	if err != nil {
		return err
	}
	// 2. Connect
	c.tagsRepo.SetConfig(tagCfg)
	return c.tagsRepo.Connect()
}

// refer to https://www.apolloconfig.com/#/zh/usage/apollo-open-api-platform?id=_327-%e5%88%9b%e5%bb%banamespace
func (c *ConfigStore) createNamespace(env string, appId string, cluster string, namespace string) error {
	// 1. request
	url := fmt.Sprintf(createNamespaceUrlTpl, c.openAPIAddress, appId)
	// add body
	reqBody := map[string]string{
		"name":                namespace,
		"appId":               appId,
		"format":              "properties",
		"isPublic":            "false",
		"dataChangeCreatedBy": c.openAPIUser,
	}
	reqBodyJson, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", url, strings.NewReader(string(reqBodyJson)))
	if err != nil {
		return err
	}
	// add headers
	c.addHeaderForOpenAPI(req)
	log.DefaultLogger.Debugf("createNamespace url: %v, request body: %s, request: %+v", url, reqBodyJson, req)
	// do request
	resp, err := c.openAPIClient.Do(req)
	// 2. parse
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		// 4. commit
		return c.commit(env, appId, cluster, namespace)
	}
	// if the namespace already exists, the status code will be 400
	if resp.StatusCode == http.StatusBadRequest {
		// log debug information
		if log.DefaultLogger.GetLogLevel() >= log.DEBUG {
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.DefaultLogger.Errorf("An error occurred when parsing createNamespace response. statusCode: %v ,error: %v", resp.StatusCode, err)
				return err
			}
			log.DefaultLogger.Debugf("createNamespace not ok. StatusCode: %v, response body: %s", resp.StatusCode, b)
		}
		return nil
	}
	// Fail fast and take it as an startup error if the status code is neither 200 nor 400
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.DefaultLogger.Errorf("An error occurred when parsing createNamespace response. statusCode: %v ,error: %v", resp.StatusCode, err)
		return err
	}
	return fmt.Errorf("createNamespace error. StatusCode: %v, response body: %s", resp.StatusCode, b)
}
