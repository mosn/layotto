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
	"fmt"
	"github.com/zouyx/agollo/v4/env/config"
	"os"
	"sync"

	"github.com/zouyx/agollo/v4/env/file"
)

var (
	raw     file.FileHandler
	rawOnce sync.Once
)

//rawFileHandler 写入备份文件时，同时写入原始内容和namespace类型
type rawFileHandler struct {
	*FileHandler
}

func writeWithRaw(config *config.ApolloConfig, configDir string) error {
	filePath := ""
	if configDir != "" {
		filePath = fmt.Sprintf("%s/%s", configDir, config.NamespaceName)
	} else {
		filePath = config.NamespaceName
	}

	file, e := os.Create(filePath)
	if e != nil {
		return e
	}
	defer file.Close()
	if config.Configurations["content"] != nil {
		_, e = file.WriteString(config.Configurations["content"].(string))
		if e != nil {
			return e
		}
	}
	return nil
}

//WriteConfigFile write config to file
func (fileHandler *rawFileHandler) WriteConfigFile(config *config.ApolloConfig, configPath string) error {
	writeWithRaw(config, configPath)
	return jsonFileConfig.Write(config, fileHandler.GetConfigFile(configPath, config.AppID, config.NamespaceName))
}

// GetRawFileHandler 获取 rawFileHandler 实例
func GetRawFileHandler() file.FileHandler {
	rawOnce.Do(func() {
		raw = &rawFileHandler{}
	})
	return raw
}
