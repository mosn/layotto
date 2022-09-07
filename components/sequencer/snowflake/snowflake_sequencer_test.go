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
	"runtime"
	"sync"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	mapset "github.com/deckarep/golang-set"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"mosn.io/pkg/log"

	"mosn.io/layotto/components/sequencer"
)

const (
	mysqlUrl     = "localhost:3306"
	databaseName = "layotto_sequence"
	tableName    = "layotto_sequence_snowflake"
	userName     = "root"
	password     = "123456"

	boostPower       = "3"
	paddingFactor    = "50"
	scheduleInterval = "5"
	timeBits         = "28"
	workerBits       = "22"
	seqBits          = "13"
	startTime        = "2022-01-01"

	key = "resource_xxx"
)

func TestSnowFlakeSequence_GetNextId(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

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

	cfg.Properties["mysql"] = mysqlUrl
	cfg.Properties["databaseName"] = databaseName
	cfg.Properties["tableName"] = tableName
	cfg.Properties["userName"] = userName
	cfg.Properties["password"] = password

	cfg.Properties["boostPower"] = boostPower
	cfg.Properties["paddingFactor"] = paddingFactor
	cfg.Properties["scheduleInterval"] = scheduleInterval
	cfg.Properties["timeBits"] = timeBits
	cfg.Properties["workerBits"] = workerBits
	cfg.Properties["seqBits"] = seqBits
	cfg.Properties["startTime"] = startTime

	err = s.Init(cfg)

	assert.NoError(t, err)

	mset := mapset.NewSet()

	var size int = 10000
	for i := 0; i < size; i++ {
		resp, err := s.GetNextId(&sequencer.GetNextIdRequest{
			Key: key,
		})
		assert.NoError(t, err)
		mset.Add(resp.NextId)
	}
	assert.Equal(t, size, mset.Cardinality())
}

func TestSnowFlakeSequence_ParallelGetNextId(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

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

	cfg.Properties["mysql"] = mysqlUrl
	cfg.Properties["databaseName"] = databaseName
	cfg.Properties["tableName"] = tableName
	cfg.Properties["userName"] = userName
	cfg.Properties["password"] = password

	err = s.Init(cfg)

	assert.NoError(t, err)

	mset := mapset.NewSet()

	var size int = 1000
	var cores int = runtime.NumCPU()

	var wg sync.WaitGroup
	wg.Add(cores)

	for i := 0; i < cores; i++ {
		go func() {
			defer func() {
				if x := recover(); x != nil {
					log.DefaultLogger.Errorf("panic when testing parallel generatoring uid with snowflake algorithm: %v", x)
				}
			}()
			for j := 0; j < size; j++ {
				resp, err := s.GetNextId(&sequencer.GetNextIdRequest{
					Key: key,
				})
				assert.NoError(t, err)
				mset.Add(resp.NextId)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	assert.Equal(t, size*cores, mset.Cardinality())
}
