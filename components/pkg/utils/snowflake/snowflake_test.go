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
	"github.com/stretchr/testify/assert"
)

func TestParseSnowflakeRingBufferMetadata(t *testing.T) {
	properties := make(map[string]string)

	properties["workerBits"] = "1"
	properties["timeBits"] = "1"
	properties["seqBits"] = "1"
	_, err := ParseSnowflakeRingBufferMetadata(properties)
	assert.Error(t, err)

	properties["startTime"] = "2022.01.01"
	_, err = ParseSnowflakeRingBufferMetadata(properties)
	assert.Error(t, err)

	properties["boostPower"] = "a"
	_, err = ParseSnowflakeRingBufferMetadata(properties)
	assert.Error(t, err)

	properties["paddingFactor"] = "a"
	_, err = ParseSnowflakeRingBufferMetadata(properties)
	assert.Error(t, err)

	properties["workerBits"] = "a"
	_, err = ParseSnowflakeRingBufferMetadata(properties)
	assert.Error(t, err)

	properties["timeBits"] = "a"
	_, err = ParseSnowflakeRingBufferMetadata(properties)
	assert.Error(t, err)

	properties["seqBits"] = "a"
	_, err = ParseSnowflakeRingBufferMetadata(properties)
	assert.Error(t, err)
}

func TestParseSnowflakeMysqlMetadata(t *testing.T) {
	properties := make(map[string]string)

	_, err := ParseSnowflakeMysqlMetadata(properties)
	assert.Error(t, err)

	properties["databaseName"] = "layotto_sequence_snowflake"
	_, err = ParseSnowflakeMysqlMetadata(properties)
	assert.Error(t, err)
}

func TestGetMysqlPort(t *testing.T) {
	p := getMysqlPort()
	assert.NotEmpty(t, p)
}

func TestNewMysqlClient(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectExec("CREATE TABLE").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT HOST_NAME").WillReturnError(sql.ErrNoRows)
	mock.ExpectExec("INSERT INTO").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectQuery("SELECT ID").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	m := MysqlMetadata{}
	m.Db = db
	workId, err := NewMysqlClient(m)
	assert.NoError(t, err)
	assert.Equal(t, workId, int64(1))
}
