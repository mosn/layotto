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
package mysql

import (
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"mosn.io/layotto/components/sequencer"
	"mosn.io/pkg/log"
	"testing"
)

const (
	MySQLUrl = "localhost:xxxxx"
	Value    = 1
	Key      = "sequenceKey"
	Size     = 50
)

func TestMySQLSequencer_Init(t *testing.T) {

	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	var MySQLUrl = MySQLUrl

	comp := NewMySQLSequencer(log.DefaultLogger)

	cfg := sequencer.Configuration{
		Properties: make(map[string]string),
		BiggerThan: make(map[string]int64),
	}

	cfg.Properties["MySQLHost"] = MySQLUrl
	cfg.BiggerThan["test"] = 1000

	err = comp.Init(cfg, db)

	assert.Error(t, err)
}

func TestMySQLSequencer_GetNextId(t *testing.T) {

	comp := NewMySQLSequencer(log.DefaultLogger)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"sequencer_key", "sequencer_value"}).AddRow([]driver.Value{Key, Value}...)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO").WithArgs("", "sequenceKey", 3).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	properties := make(map[string]string)

	req := &sequencer.GetNextIdRequest{Key: Key, Options: sequencer.SequencerOptions{AutoIncrement: sequencer.STRONG}, Metadata: properties}

	_, err = comp.GetNextId(req, db)
	if err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	assert.NoError(t, err)
}

func TestMySQLSequencer_GetSegment(t *testing.T) {

	comp := NewMySQLSequencer(log.DefaultLogger)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"sequencer_key", "sequencer_value"}).AddRow([]driver.Value{Key, Value}...)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	properties := make(map[string]string)

	req := &sequencer.GetSegmentRequest{Size: Size, Key: Key, Options: sequencer.SequencerOptions{AutoIncrement: sequencer.STRONG}, Metadata: properties}

	_, _, err = comp.GetSegment(req, db)
	if err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}

	assert.NoError(t, err)
}

func TestMySQLSequencer_Close(t *testing.T) {

	var MySQLUrl = MySQLUrl

	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	comp := NewMySQLSequencer(log.DefaultLogger)

	cfg := sequencer.Configuration{
		BiggerThan: nil,
		Properties: make(map[string]string),
	}

	cfg.Properties["MySQLHost"] = MySQLUrl

	_ = comp.Init(cfg, db)

	comp.Close(db)
}
