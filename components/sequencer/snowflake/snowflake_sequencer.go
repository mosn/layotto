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
	"time"

	"mosn.io/pkg/log"

	"mosn.io/layotto/components/pkg/utils/snowflake"
	"mosn.io/layotto/components/sequencer"
)

type SnowFlakeSequencer struct {
	metadata   snowflake.SnowflakeMetadata
	workId     int64
	ringBuffer *snowflake.RingBuffer
	db         *sql.DB
	biggerThan map[string]int64
	logger     log.ErrorLogger
}

func NewSnowFlakeSequencer(logger log.ErrorLogger) *SnowFlakeSequencer {
	return &SnowFlakeSequencer{
		logger: logger,
	}
}

func (s *SnowFlakeSequencer) Init(config sequencer.Configuration) error {
	mm, err := snowflake.ParseSnowflakeMysqlMetadata(config.Properties)
	if err != nil {
		return err
	}
	mm.Db = s.db

	s.metadata.MysqlMetadata = &mm
	s.biggerThan = config.BiggerThan

	var workId int64
	if workId, err = snowflake.NewMysqlClient(*s.metadata.MysqlMetadata); err != nil {
		return err
	}

	s.workId = workId

	rm, err := snowflake.ParseSnowflakeRingBufferMetadata(config.Properties)

	startTime := rm.StartTime
	seqBits := rm.SeqBits

	var maxSeq int64 = ^(-1 << seqBits) + 1
	bufferSize := maxSeq << rm.BoostPower

	s.ringBuffer = snowflake.NewRingBuffer(bufferSize)

	s.ringBuffer.WorkId = workId
	s.ringBuffer.CurrentTimeStamp = time.Now().Unix() - startTime
	s.ringBuffer.TimeBits = rm.TimeBits
	s.ringBuffer.WorkIdBits = rm.WorkIdBits
	s.ringBuffer.SeqBits = seqBits
	s.ringBuffer.PaddingFactor = rm.PaddingFactor
	s.ringBuffer.MaxSeq = ^(-1 << seqBits) + 1

	s.ringBuffer.PaddingRingBuffer()

	return err
}

func (s *SnowFlakeSequencer) GetNextId(req *sequencer.GetNextIdRequest) (*sequencer.GetNextIdResponse, error) {
	uid, err := s.ringBuffer.Take()
	if err != nil {
		return nil, err
	}

	return &sequencer.GetNextIdResponse{
		NextId: uid,
	}, nil
}

func (s *SnowFlakeSequencer) GetSegment(req *sequencer.GetSegmentRequest) (support bool, result *sequencer.GetSegmentResponse, err error) {
	return false, nil, nil
}

func (s *SnowFlakeSequencer) Close(db *sql.DB) error {
	return nil
}
