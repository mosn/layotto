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
	"sync"
	"time"

	"mosn.io/pkg/log"

	"mosn.io/layotto/components/sequencer"
)

type SnowFlakeSequencer struct {
	metadata   SnowflakeMetadata
	workerId   int64
	db         *sql.DB
	mu         sync.Mutex
	smap       map[string]chan int64
	biggerThan map[string]int64
	logger     log.ErrorLogger
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewSnowFlakeSequencer(logger log.ErrorLogger) *SnowFlakeSequencer {
	return &SnowFlakeSequencer{
		logger: logger,
		smap:   make(map[string]chan int64),
	}
}

func (s *SnowFlakeSequencer) Init(config sequencer.Configuration) error {
	var err error
	s.metadata, err = ParseSnowflakeMetadata(config.Properties)
	if err != nil {
		return err
	}
	//for unit test
	s.metadata.MysqlMetadata.Db = s.db

	s.biggerThan = config.BiggerThan
	s.ctx, s.cancel = context.WithCancel(context.Background())

	if s.workerId, err = NewMysqlClient(&s.metadata.MysqlMetadata); err != nil {
		return err
	}
	return err
}

func (s *SnowFlakeSequencer) GetNextId(req *sequencer.GetNextIdRequest) (*sequencer.GetNextIdResponse, error) {
	s.mu.Lock()
	ch, ok := s.smap[req.Key]
	//If the key appears for the first time, start a new goroutine for it. If the key doesn't appear for a long time, close the goroutine
	if !ok {
		ch = make(chan int64, 1000)
		s.smap[req.Key] = ch

		var oldWorkerId int64
		var oldTimeStamp int64

		timestamp := time.Now().Unix() - s.metadata.StartTime

		err := s.metadata.MysqlMetadata.Db.QueryRow("SELECT WORKER_ID, TIMESTAMP FROM "+s.metadata.MysqlMetadata.KeyTableName+" WHERE SEQUENCER_KEY = ?", req.Key).Scan(&oldWorkerId, &oldTimeStamp)
		if err == nil {
			if oldWorkerId == s.workerId {
				timestamp = oldTimeStamp + 1
			}
		} else if err != sql.ErrNoRows {
			return nil, err
		}
		startId := timestamp<<s.metadata.TimestampShift | s.workerId<<s.metadata.WorkidShift

		go s.producer(startId, timestamp, ch, req.Key)
	}
	s.mu.Unlock()

	timeout := time.NewTimer(s.metadata.ReqTimeout)
	defer timeout.Stop()

	var id int64
	select {
	case id, ok = <-ch:
		if !ok {
			return nil, errors.New("please try again or adjust the start time")
		}
		return &sequencer.GetNextIdResponse{
			NextId: id,
		}, nil
	case <-timeout.C:
		return nil, errors.New("request id time out")
	}
}

func (s *SnowFlakeSequencer) GetSegment(req *sequencer.GetSegmentRequest) (support bool, result *sequencer.GetSegmentResponse, err error) {
	return false, nil, nil
}

func (s *SnowFlakeSequencer) producer(id, currentTimeStamp int64, ch chan int64, key string) {
	defer func() {
		if x := recover(); x != nil {
			log.DefaultLogger.Errorf("panic when producing id with snowflake algorithm: %v", x)
		}
	}()

	timeout := time.NewTimer(s.metadata.KeyTimeout)
	defer timeout.Stop()

	var maxTimeStamp int64
	var maxSeqId int64

	maxTimeStamp = 1 << s.metadata.TimeBits
	maxSeqId = 1<<s.metadata.SeqBits - 1
	for {
		timeout.Reset(s.metadata.KeyTimeout)
		select {
		case <-s.ctx.Done():
			close(ch)
			return
		//if timeout, remove key from map and record key, workerId, timestamp to mysql
		case <-timeout.C:
			s.mu.Lock()
			delete(s.smap, key)
			close(ch)

			err := MysqlRecord(s.metadata.MysqlMetadata.Db, s.metadata.MysqlMetadata.KeyTableName, key, s.workerId, currentTimeStamp)
			if err != nil {
				s.logger.Errorf("%v", err)
			}
			s.mu.Unlock()
			return
		case ch <- id:
			if currentTimeStamp == maxTimeStamp {
				close(ch)
				return
			}
			if id&maxSeqId != maxSeqId {
				id++
			} else {
				currentTimeStamp++
				id = currentTimeStamp<<s.metadata.TimestampShift | s.workerId<<s.metadata.WorkidShift
			}
		}
	}
}

func (s *SnowFlakeSequencer) Close() error {
	s.cancel()
	s.metadata.MysqlMetadata.Db.Close()
	return nil
}
