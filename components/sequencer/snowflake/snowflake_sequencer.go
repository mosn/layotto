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
	"errors"
	"time"

	"mosn.io/pkg/log"

	"mosn.io/layotto/components/pkg/utils"
	"mosn.io/layotto/components/sequencer"
)

type SnowFlakeSequencer struct {
	metadata   utils.SnowflakeMetadata
	ch         chan int64
	db         *sql.DB
	biggerThan map[string]int64
	logger     log.ErrorLogger
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewSnowFlakeSequencer(logger log.ErrorLogger) *SnowFlakeSequencer {
	return &SnowFlakeSequencer{
		ch:     make(chan int64, 1000),
		logger: logger,
	}
}

func (s *SnowFlakeSequencer) Init(config sequencer.Configuration) error {
	sm, err := utils.ParseSnowflakeMetadata(config.Properties)
	if err != nil {
		return err
	}
	//for unit test
	sm.MysqlMetadata.Db = s.db

	s.metadata = sm
	s.biggerThan = config.BiggerThan
	s.ctx, s.cancel = context.WithCancel(context.Background())

	var workId int64
	if workId, err = utils.NewMysqlClient(*s.metadata.MysqlMetadata); err != nil {
		return err
	}

	timestampShift := sm.SeqBits
	workidShift := sm.TimeBits + sm.SeqBits

	currentTimeStamp := time.Now().Unix() - sm.StartTime

	var sequence int64
	startId := workId<<workidShift | currentTimeStamp<<timestampShift | sequence
	maxId := (workId + 1) << workidShift

	//producer
	go func(ctx context.Context, id int64, maxId int64) {
		defer func() {
			if x := recover(); x != nil {
				log.DefaultLogger.Errorf("panic when producing id with snowflake algorithm: %v", x)
			}
		}()
		for {
			select {
			case <-ctx.Done():
				close(s.ch)
				return
			default:
				if id == maxId {
					//if id reach the max value, wait for the id to be used up
					if len(s.ch) == 0 {
						close(s.ch)
						return
					}
				} else {
					s.ch <- id
					id++
				}
			}
		}
	}(s.ctx, startId, maxId)
	return err
}

func (s *SnowFlakeSequencer) GetNextId(req *sequencer.GetNextIdRequest) (*sequencer.GetNextIdResponse, error) {
	var id int64
	var ok bool
	if id, ok = <-s.ch; !ok {
		return nil, errors.New("timeBits has been used up, please adjust the start time")
	}

	return &sequencer.GetNextIdResponse{
		NextId: id,
	}, nil
}

func (s *SnowFlakeSequencer) GetSegment(req *sequencer.GetSegmentRequest) (support bool, result *sequencer.GetSegmentResponse, err error) {
	return false, nil, nil
}

func (s *SnowFlakeSequencer) Close() error {
	s.cancel()
	return nil
}
