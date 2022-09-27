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
package utils

import (
	"database/sql"
	"fmt"
)

var (
	Db *sql.DB
)

const (
	defaultTableName    = "layotto_sequencer"
	defaultTableNameKey = "tableName"
	dataBaseName        = "dataBaseName"
	userName            = "userName"
	defaultPassword     = "password"
	mysqlUrl            = "mysqlUrl"
)

type MySQLMetadata struct {
	TableName    string
	DataBaseName string
	UserName     string
	Password     string
	MysqlUrl     string
	Db           *sql.DB
}

func ParseMySQLMetadata(properties map[string]string) (MySQLMetadata, error) {
	m := MySQLMetadata{}

	if val, ok := properties[defaultTableNameKey]; ok && val != "" {
		m.TableName = val
	}
	if val, ok := properties[dataBaseName]; ok && val != "" {
		m.DataBaseName = val
	}
	if val, ok := properties[userName]; ok && val != "" {
		m.UserName = val
	}
	if val, ok := properties[defaultPassword]; ok && val != "" {
		m.Password = val
	}
	if val, ok := properties[mysqlUrl]; ok && val != "" {
		m.MysqlUrl = val
	}
	return m, nil
}

func NewMySQLClient(meta MySQLMetadata) error {

	val := meta
	if val.TableName == "" {
		val.TableName = defaultTableName
	}
	meta.Db.Begin()
	createTable := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
			sequencer_key VARCHAR(255),
			sequencer_value INT,
			UNIQUE INDEX (sequencer_key));`, val.TableName)
	_, err := meta.Db.Exec(createTable)
	defer meta.Db.Close()
	if err != nil {
		return err
	}
	return nil
}
