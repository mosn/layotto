package config

//
// Copyright 2021 Layotto Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/**
* @Author: azh
* @Date: 2022/5/13 22:12
* @Context:
 */

import (
	"io/ioutil"

	_ "github.com/lib/pq"
	"gopkg.in/yaml.v2"
)

type Server struct {
	Postgresql Postgresql `json:"postgresql" yaml:"postgresql"`
}

type Postgresql struct {
	Host     string `json:"host" yaml:"host"`
	Port     string `json:"port" yaml:"port"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	Db       string `json:"db" yaml:"db"`
}

// The file is not loaded using this function now
func (m *Server) Load(filePath string) error {
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlFile, m)
	if err != nil {
		return err
	}
	return nil
}
