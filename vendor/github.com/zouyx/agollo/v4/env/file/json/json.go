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

package json

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/zouyx/agollo/v4/env/config"
	"sync"

	"github.com/zouyx/agollo/v4/component/log"
	jsonConfig "github.com/zouyx/agollo/v4/env/config/json"
)

//Suffix 默认文件保存类型
const Suffix = ".json"

var (
	//jsonFileConfig 处理文件的json格式存取
	jsonFileConfig = &jsonConfig.ConfigFile{}
	//configFileMap 存取namespace文件地址
	configFileMap = make(map[string]string, 1)
	//configFileMapLock configFileMap的读取锁
	configFileMapLock sync.Mutex
)

// FileHandler 默认备份文件读写
type FileHandler struct {
}

// WriteConfigFile write config to file
func (fileHandler *FileHandler) WriteConfigFile(config *config.ApolloConfig, configPath string) error {
	return jsonFileConfig.Write(config, fileHandler.GetConfigFile(configPath, config.AppID, config.NamespaceName))
}

// GetConfigFile get real config file
func (fileHandler *FileHandler) GetConfigFile(configDir string, appID string, namespace string) string {
	key := fmt.Sprintf("%s-%s", appID, namespace)
	configFileMapLock.Lock()
	defer configFileMapLock.Unlock()
	fullPath := configFileMap[key]
	if fullPath == "" {
		filePath := fmt.Sprintf("%s%s", key, Suffix)
		if configDir != "" {
			configFileMap[namespace] = fmt.Sprintf("%s/%s", configDir, filePath)
		} else {
			configFileMap[namespace] = filePath
		}
	}
	return configFileMap[namespace]
}

//LoadConfigFile load config from file
func (fileHandler *FileHandler) LoadConfigFile(configDir string, appID string, namespace string) (*config.ApolloConfig, error) {
	configFilePath := fileHandler.GetConfigFile(configDir, appID, namespace)
	log.Info("load config file from :", configFilePath)
	c, e := jsonFileConfig.Load(configFilePath, func(b []byte) (interface{}, error) {
		config := &config.ApolloConfig{}
		e := json.NewDecoder(bytes.NewBuffer(b)).Decode(config)
		return config, e
	})

	if c == nil || e != nil {
		log.Errorf("loadConfigFile fail,error:", e)
		return nil, e
	}

	return c.(*config.ApolloConfig), e
}
