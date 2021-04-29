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
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/zouyx/agollo/v4/component/log"
	"github.com/zouyx/agollo/v4/utils"
)

//ConfigFile json文件读写
type ConfigFile struct {
}

//Load json文件读
func (t *ConfigFile) Load(fileName string, unmarshal func([]byte) (interface{}, error)) (interface{}, error) {
	fs, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, errors.New("Fail to read config file:" + err.Error())
	}

	config, loadErr := unmarshal(fs)

	if utils.IsNotNil(loadErr) {
		return nil, errors.New("Load Json Config fail:" + loadErr.Error())
	}

	return config, nil
}

//Write json文件写
func (t *ConfigFile) Write(content interface{}, configPath string) error {
	if content == nil {
		log.Error("content is null can not write backup file")
		return errors.New("content is null can not write backup file")
	}
	file, e := os.Create(configPath)
	if e != nil {
		log.Errorf("writeConfigFile fail,error:", e)
		return e
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(content)
}
