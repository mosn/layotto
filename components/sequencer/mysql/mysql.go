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
	"database/sql"
	"fmt"

	"mosn.io/pkg/log"

	"mosn.io/layotto/components/pkg/utils"
	"mosn.io/layotto/components/sequencer"
)

type MySQLSequencer struct {
	metadata   utils.MySQLMetadata
	biggerThan map[string]int64
	logger     log.ErrorLogger
	db         *sql.DB
}

func NewMySQLSequencer(logger log.ErrorLogger) *MySQLSequencer {
	s := &MySQLSequencer{
		logger: logger,
	}

	return s
}

func (e *MySQLSequencer) Init(config sequencer.Configuration) error {

	m, err := utils.ParseMySQLMetadata(config.Properties)

	if err != nil {
		return err
	}
	e.metadata = m
	e.metadata.Db = e.db
	e.biggerThan = config.BiggerThan

	if err = utils.NewMySQLClient(e.metadata); err != nil {
		return err
	}

	if len(e.biggerThan) > 0 {
		var Key string
		var Value int64
		for k, bt := range e.biggerThan {
			if bt > 0 {
				m.Db.QueryRow(fmt.Sprintf("SELECT sequencer_key,sequencer_value FROM %s WHERE sequencer_key = ?", m.TableName), k).Scan(&Key, &Value)
				if Value < bt {
					return fmt.Errorf("mysql sequencer error: can not satisfy biggerThan guarantee.key: %s,key in MySQL: %s", k, Key)
				}
			}
		}
	}
	return nil
}

func (e *MySQLSequencer) GetNextId(req *sequencer.GetNextIdRequest) (*sequencer.GetNextIdResponse, error) {

	metadata, err := utils.ParseMySQLMetadata(req.Metadata)
	metadata.Db = e.db
	var Key string
	var Value int64
	var Version int64
	var oldVersion int64
	if err != nil {
		return nil, err
	}
	begin, err := metadata.Db.Begin()
	if err != nil {
		return nil, err
	}
	err = begin.QueryRow("SELECT sequencer_key, sequencer_value, version FROM ? WHERE sequencer_key = ?", metadata.TableName, req.Key).Scan(&Key, &Value, &oldVersion)

	if err == sql.ErrNoRows {
		Value = 1
		Version = 1
		_, err := begin.Exec("INSERT INTO ?(sequencer_key, sequencer_value, version) VALUES(?,?,?)", metadata.TableName, req.Key, Value, Version)
		if err != nil {
			begin.Rollback()
			return nil, err
		}

	} else {
		_, err := begin.Exec("UPDATE ? SET sequencer_value +=1, version += 1 WHERE sequencer_key = ? and version = ?", metadata.TableName, req.Key, oldVersion)
		if err != nil {
			begin.Rollback()
			return nil, err
		}
	}
	defer e.Close(metadata.Db)

	return &sequencer.GetNextIdResponse{
		NextId: Value,
	}, nil
}

func (e *MySQLSequencer) GetSegment(req *sequencer.GetSegmentRequest) (support bool, result *sequencer.GetSegmentResponse, err error) {

	if req.Size == 0 {
		return true, nil, nil
	}

	metadata, err := utils.ParseMySQLMetadata(req.Metadata)

	var Key string
	var Value int64
	var Version int64
	var oldVersion int64
	if err != nil {
		return false, nil, err
	}

	metadata.Db = e.db

	begin, err := metadata.Db.Begin()
	if err != nil {
		return false, nil, err
	}
	err = begin.QueryRow("SELECT sequencer_key, sequencer_value, version FROM ? WHERE sequencer_key == ?", metadata.TableName, req.Key).Scan(&Key, &Value, &oldVersion)
	if err == sql.ErrNoRows {
		Value = int64(req.Size)
		Version = 1
		_, err := begin.Exec("INSERT INTO ?(sequencer_key, sequencer_value, version) VALUES(?,?,?)", metadata.TableName, req.Key, Value, Version)
		if err != nil {
			begin.Rollback()
			return false, nil, err
		}

	} else {
		Value += int64(req.Size)
		_, err1 := begin.Exec("UPDATE ? SET sequencer_value = ?, version += 1 WHERE sequencer_key = ? AND version = ?", metadata.TableName, Value, req.Key, oldVersion)
		if err1 != nil {
			begin.Rollback()
			return false, nil, err1
		}
	}
	defer e.Close(metadata.Db)

	return false, &sequencer.GetSegmentResponse{
		From: Value - int64(req.Size) + 1,
		To:   Value,
	}, nil
}

func (e *MySQLSequencer) Close(db *sql.DB) error {

	return db.Close()
}
