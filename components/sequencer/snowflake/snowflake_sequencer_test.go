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
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"mosn.io/pkg/log"

	"mosn.io/layotto/components/sequencer"
)

const (
	mysqlHostUrl = "localhost:3306"
	databaseName = "layotto_sequencer"
	userName     = "root"
	password     = "123456"
	key          = "resource_xxx"
)

func TestSnowflakeSequence_GetNextId(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	s := NewSnowFlakeSequencer(log.DefaultLogger)
	s.db = db

	mock.ExpectExec("CREATE TABLE").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("CREATE TABLE").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT HOST_NAME").WillReturnError(sql.ErrNoRows)
	mock.ExpectExec("INSERT INTO").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery("SELECT ID").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()
	mock.ExpectQuery("SELECT WORKER_ID").WillReturnError(sql.ErrNoRows)

	cfg := sequencer.Configuration{
		Properties: make(map[string]string),
		BiggerThan: make(map[string]int64),
	}

	cfg.Properties["mysqlHost"] = mysqlHostUrl
	cfg.Properties["databaseName"] = databaseName
	cfg.Properties["tableName"] = defaultMysqlTableName
	cfg.Properties["userName"] = userName
	cfg.Properties["password"] = password

	cfg.Properties["timeBits"] = "28"
	cfg.Properties["workerBits"] = "22"
	cfg.Properties["seqBits"] = "13"
	cfg.Properties["startTime"] = "2022-01-01"

	err = s.Init(cfg)

	assert.NoError(t, err)

	var falseUidNum int
	var preUid int64
	var size int = 10000
	for i := 0; i < size; i++ {
		resp, err := s.GetNextId(&sequencer.GetNextIdRequest{
			Key: key,
		})
		//assert next uid is bigger than previous
		if preUid != 0 && resp.NextId <= preUid {
			falseUidNum++
		}
		preUid = resp.NextId
		assert.NoError(t, err)
	}
	assert.Equal(t, falseUidNum, 0)
}

func TestSnowflakeSequence_ParallelGetNextId(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	s := NewSnowFlakeSequencer(log.DefaultLogger)
	s.db = db

	mock.ExpectExec("CREATE TABLE").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("CREATE TABLE").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT HOST_NAME").WillReturnError(sql.ErrNoRows)
	mock.ExpectExec("INSERT INTO").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery("SELECT ID").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()
	mock.ExpectQuery("SELECT WORKER_ID").WillReturnError(sql.ErrNoRows)
	mock.ExpectQuery("SELECT WORKER_ID").WillReturnError(sql.ErrNoRows)

	cfg := sequencer.Configuration{
		Properties: make(map[string]string),
		BiggerThan: make(map[string]int64),
	}

	cfg.Properties["mysqlHost"] = mysqlHostUrl
	cfg.Properties["databaseName"] = databaseName
	cfg.Properties["tableName"] = defaultMysqlTableName
	cfg.Properties["userName"] = userName
	cfg.Properties["password"] = password

	err = s.Init(cfg)

	assert.NoError(t, err)

	var size int = 20000
	var falseUidNum int

	var wg sync.WaitGroup
	wg.Add(2)
	var key string
	for i := 0; i < 2; i++ {
		key = strconv.Itoa(i)
		go func(key string) {
			defer func() {
				if x := recover(); x != nil {
					log.DefaultLogger.Errorf("panic when testing parallel generatoring uid with snowflake algorithm: %v", x)
				}
			}()
			var preUid int64
			for j := 0; j < size; j++ {
				//assert next uid is bigger than previous
				resp, err := s.GetNextId(&sequencer.GetNextIdRequest{
					Key: key,
				})
				if preUid != 0 && resp.NextId <= preUid {
					falseUidNum++
				}
				preUid = resp.NextId
				assert.NoError(t, err)
			}
			wg.Done()
		}(key)
	}
	wg.Wait()
}

func TestKeyTimeout(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	s := NewSnowFlakeSequencer(log.DefaultLogger)
	s.db = db

	mock.ExpectExec("CREATE TABLE").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("CREATE TABLE").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT HOST_NAME").WillReturnError(sql.ErrNoRows)
	mock.ExpectExec("INSERT INTO").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery("SELECT ID").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	mock.ExpectQuery("SELECT WORKER_ID").WillReturnError(sql.ErrNoRows)
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT WORKER_ID").WillReturnError(sql.ErrNoRows)
	mock.ExpectExec("INSERT INTO").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	cfg := sequencer.Configuration{
		Properties: make(map[string]string),
		BiggerThan: make(map[string]int64),
	}

	cfg.Properties["mysqlHost"] = mysqlHostUrl
	cfg.Properties["databaseName"] = databaseName
	cfg.Properties["tableName"] = defaultMysqlTableName
	cfg.Properties["userName"] = userName
	cfg.Properties["password"] = password

	cfg.Properties["keyTimeout"] = "0"

	err = s.Init(cfg)

	assert.NoError(t, err)

	_, err = s.GetNextId(&sequencer.GetNextIdRequest{
		Key: key,
	})
	time.Sleep(time.Second)
	assert.NoError(t, err)
}
