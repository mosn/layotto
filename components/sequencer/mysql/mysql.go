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
	"context"
	"database/sql"
	"fmt"
	"mosn.io/layotto/components/pkg/utils"
	"mosn.io/layotto/components/sequencer"
	"mosn.io/pkg/log"
)

type MySQLSequencer struct {
	metadata   utils.MySQLMetadata
	biggerThan map[string]int64
	logger     log.ErrorLogger
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewMySQLSequencer(logger log.ErrorLogger) *MySQLSequencer {
	s := &MySQLSequencer{
		logger: logger,
	}

	return s
}

func (e *MySQLSequencer) Init(config sequencer.Configuration, db *sql.DB) error {

	m, err := utils.ParseMySQLMetadata(config.Properties, db)

	if err != nil {
		return err
	}
	e.metadata = m
	e.biggerThan = config.BiggerThan

	if err = utils.NewMySQLClient(e.metadata); err != nil {
		return err
	}
	e.ctx, e.cancel = context.WithCancel(context.Background())

	if len(e.biggerThan) > 0 {
		var Key string
		var Value int64
		for k, bt := range e.biggerThan {
			if bt <= 0 {
				continue
			} else {
				m.Db.QueryRow(fmt.Sprintf("SELECT * FROM %s WHERE sequencer_key = ?", m.TableName), k).Scan(&Key, &Value)
				if Value < bt {
					return fmt.Errorf("MySQL sequencer error: can not satisfy biggerThan guarantee.key: %s,key in MySQL: %s", k, Key)
				}
			}
		}
	}
	return nil
}

func (e *MySQLSequencer) GetNextId(req *sequencer.GetNextIdRequest, db *sql.DB) (*sequencer.GetNextIdResponse, error) {

	metadata, err := utils.ParseMySQLMetadata(req.Metadata, db)

	var Key string
	var Value int64
	if err != nil {
		return nil, err
	}

	begin, err := metadata.Db.Begin()
	if err != nil {
		return nil, err
	}

	err = begin.QueryRow("SELECT sequencer_key,sequencer_value FROM ? WHERE sequencer_key = ?", metadata.TableName, req.Key).Scan(&Key, &Value)
	if err != nil {
		return nil, err
	}

	Value += 1
	_, err1 := begin.Exec("UPDATE ? SET sequencer_value = ? WHERE sequencer_key = ?", metadata.TableName, Value, req.Key)
	if err1 != nil {
		return nil, err1
	}

	Value += 1
	if _, err2 := begin.Exec("INSERT INTO ?(sequencer_key, sequencer_value) VALUES(?,?)", metadata.TableName, req.Key, Value); err2 != nil {
		return nil, err2
	}

	e.Close(metadata.Db)

	return &sequencer.GetNextIdResponse{
		NextId: Value,
	}, nil
}

func (e *MySQLSequencer) GetSegment(req *sequencer.GetSegmentRequest, db *sql.DB) (support bool, result *sequencer.GetSegmentResponse, err error) {

	if req.Size == 0 {
		return true, nil, nil
	}

	metadata, err := utils.ParseMySQLMetadata(req.Metadata, db)

	var Key string
	var Value int64
	if err != nil {
		return false, nil, err
	}

	begin, err := metadata.Db.Begin()
	if err != nil {
		return false, nil, err
	}

	err = begin.QueryRow("SELECT sequencer_key,sequencer_value FROM ? WHERE sequencer_key = ?", metadata.TableName, req.Key).Scan(&Key, &Value)
	if err != nil {
		return false, nil, err
	}
	Value += int64(req.Size)

	_, err1 := begin.Exec("UPDATE ? SET sequencer_value = ? WHERE sequencer_key = ?", metadata.TableName, Value, req.Key)
	if err1 != nil {
		return false, nil, err1
	}

	if _, err2 := begin.Exec("INSERT INTO ?(sequencer_key, sequencer_value) VALUES(?,?)", metadata.TableName, req.Key, Value+1); err2 != nil {
		return false, nil, err2
	}

	e.Close(metadata.Db)

	return false, &sequencer.GetSegmentResponse{
		From: Value - int64(req.Size) + 1,
		To:   Value,
	}, nil
}

func (e *MySQLSequencer) Close(db *sql.DB) error {

	return db.Close()
}
