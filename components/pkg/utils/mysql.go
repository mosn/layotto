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
package utils

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var (
	Db *sql.DB
)

const (
	defaultTableName    = "tableName"
	tableNameKey        = "tableName"
	connectionStringKey = "connectionString"
	dataBaseName        = "dataBaseName"
)

type MySQLMetadata struct {
	TableName        string
	ConnectionString string
	DataBaseName     string
	Db               *sql.DB
}

func ParseMySQLMetadata(properties map[string]string, db *sql.DB) (MySQLMetadata, error) {
	m := MySQLMetadata{}

	if val, ok := properties[tableNameKey]; ok && val != "" {
		m.TableName = val
	}

	if val, ok := properties[connectionStringKey]; ok && val != "" {
		m.ConnectionString = val
	}

	if val, ok := properties[dataBaseName]; ok && val != "" {
		m.DataBaseName = val
	}
	m.Db = db
	return m, nil
}

func NewMySQLClient(meta MySQLMetadata) error {

	val := meta

	if val.TableName == "" {
		val.TableName = defaultTableName
	}

	return tableExists(meta)
}

func tableExists(meta MySQLMetadata) error {
	exists := ""

	query := `SELECT EXISTS (
		SELECT * FROM ? WHERE TABLE_NAME = ?
		) AS 'exists'`
	err := meta.Db.QueryRow(query, meta.DataBaseName, meta.TableName).Scan(&exists)

	if exists == "1" {
		return nil
	} else {
		return err
	}
}
