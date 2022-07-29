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
	"context"
	"database/sql"

	"mosn.io/pkg/log"

	"mosn.io/layotto/components/pkg/utils/snowflake"
	"mosn.io/layotto/components/sequencer"
)

type SnowFlakeSequencer struct {
	metadata   snowflake.SnowflakeMetadata
	workId     int64
	db         *sql.DB
	biggerThan map[string]int64
	logger     log.ErrorLogger
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewSnowFlakeSequencer(logger log.ErrorLogger) *SnowFlakeSequencer {
	return &SnowFlakeSequencer{
		logger: logger,
	}
}

func (s *SnowFlakeSequencer) Init(config sequencer.Configuration) error {

	m, err := snowflake.ParseSnowflakeMysqlMetadata(config.Properties)
	m.Db = s.db

	if err != nil {
		return err
	}

	s.metadata.MysqlMetadata = &m
	s.biggerThan = config.BiggerThan

	var workId int64
	if workId, err = snowflake.NewMysqlClient(*s.metadata.MysqlMetadata); err != nil {
		return err
	}

	s.workId = workId
	return err
}

func (s *SnowFlakeSequencer) GetNextId(req *sequencer.GetNextIdRequest) (*sequencer.GetNextIdResponse, error) {

	return &sequencer.GetNextIdResponse{}, nil
}

func (s *SnowFlakeSequencer) GetSegment(req *sequencer.GetSegmentRequest) (support bool, result *sequencer.GetSegmentResponse, err error) {
	return false, nil, nil
}

func (s *SnowFlakeSequencer) Close() error {
	return nil
}
