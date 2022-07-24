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
package snowflake

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"mosn.io/pkg/log"

	"mosn.io/layotto/components/pkg/utils/snowflake"
	"mosn.io/layotto/components/sequencer"
)

const (
	mysqlUrl     = "localhost:3306"
	databaseName = "layotto_sequence"
	tableName    = "snowflake"
)

func TestSnowFlakeSequence_Init(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	s := NewSnowFlakeSequencer(log.DefaultLogger)
	s.db = db

	mock.ExpectQuery("SELECT TABLE_NAME").WillReturnError(sql.ErrNoRows)
	mock.ExpectExec("CREATE TABLE").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT HOST_NAME").WillReturnError(sql.ErrNoRows)
	mock.ExpectExec("INSERT INTO").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery("SELECT ID").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	cfg := sequencer.Configuration{
		Properties: make(map[string]string),
		BiggerThan: make(map[string]int64),
	}

	cfg.Properties[snowflake.DefaultMysqlUrl] = mysqlUrl
	cfg.Properties[snowflake.DefaultDatabaseName] = databaseName
	cfg.Properties[snowflake.DefaultTableName] = tableName

	err = s.Init(cfg)

	assert.NoError(t, err)
}
