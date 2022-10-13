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

	"mosn.io/layotto/components/pkg/utils"
	"mosn.io/layotto/components/sequencer"
)

type SnowFlakeSequencer struct {
	metadata   utils.SnowflakeMetadata
	workerId   int64
	mu         sync.Mutex
	db         *sql.DB
	smap       map[string]chan int64
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
	var err error
	s.metadata, err = utils.ParseSnowflakeMetadata(config.Properties)
	if err != nil {
		return err
	}
	//for unit test
	s.metadata.MysqlMetadata.Db = s.db

	s.smap = make(map[string]chan int64)
	s.biggerThan = config.BiggerThan
	s.ctx, s.cancel = context.WithCancel(context.Background())

	if s.workerId, err = utils.NewMysqlClient(&s.metadata.MysqlMetadata); err != nil {
		return err
	}
	return err
}

func (s *SnowFlakeSequencer) GetNextId(req *sequencer.GetNextIdRequest) (*sequencer.GetNextIdResponse, error) {
	s.mu.Lock()
	var ok bool
	seed, ok := s.smap[req.Key]
	//If the key appears for the first time, start a new goroutine for it. If the key doesn't appear for a long time, close the goroutine
	if !ok {
		seed = make(chan int64, 1000)
		s.smap[req.Key] = seed

		var oldWorkerId int64
		var oldTimeStamp int64
		currentTimeStamp := time.Now().Unix() - s.metadata.StartTime
		cts := currentTimeStamp

		err := s.metadata.MysqlMetadata.Db.QueryRow("SELECT WORKER_ID, TIMESTAMP FROM "+s.metadata.MysqlMetadata.KeyTableName+" WHERE SEQUENCER_KEY = ?", req.Key).Scan(&oldWorkerId, &oldTimeStamp)
		if err == nil {
			if oldWorkerId == s.workerId {
				cts = oldTimeStamp + 1
			}
		}

		var sequence int64
		var workerId = s.workerId
		startId := cts<<s.metadata.TimestampShift | workerId<<s.metadata.WorkidShift | sequence

		//producer
		go s.producer(startId, currentTimeStamp, seed, req.Key)
	}
	s.mu.Unlock()

	var id int64
	for {
		select {
		case id, ok = <-seed:
			if !ok {
				return nil, errors.New("please try again or adjust the start time")
			}
			return &sequencer.GetNextIdResponse{
				NextId: id,
			}, nil
		case <-time.After(s.metadata.ReqTimeout):
			return nil, errors.New("request id time out")
		}
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

	var maxTimeStamp int64
	var maxSeqId int64
	var seq int64

	maxTimeStamp = 2 << s.metadata.TimeBits
	maxSeqId = 2 << s.metadata.SeqBits
	for {
		select {
		case <-s.ctx.Done():
			close(ch)
			return
		//if timeout, remove key from map and record key, workerId, timestamp to mysql
		case <-time.After(s.metadata.KeyTimeout):
			delete(s.smap, key)
			close(ch)

			err := mysqlRecord(s.metadata.MysqlMetadata.Db, s.metadata.MysqlMetadata.KeyTableName, key, s.workerId, currentTimeStamp)
			if err != nil {
				s.logger.Errorf("%v", err)
			}
			return
		case ch <- id:
			if currentTimeStamp == maxTimeStamp {
				//if id reach the max value, wait for the id to be used up
				if len(ch) == 0 {
					close(ch)
					return
				}
			} else {
				if seq < maxSeqId {
					seq++
					id++
				} else {
					var workerId int64 = s.workerId
					var sequence int64

					currentTimeStamp++
					cts := currentTimeStamp
					id = cts<<s.metadata.TimestampShift | workerId<<s.metadata.WorkidShift | sequence
					seq = 0
				}
			}
		}
	}
}

func mysqlRecord(db *sql.DB, keyTableName, key string, workerId, timestamp int64) error {
	var mysqlWorkerId int64
	var mysqlTimestamp int64

	begin, err := db.Begin()
	if err != nil {
		return err
	}
	err = begin.QueryRow("SELECT WORKER_ID, TIMESTAMP FROM "+keyTableName+" WHERE SEQUENCER_KEY = ?", key).Scan(&mysqlWorkerId, &mysqlTimestamp)
	if err == sql.ErrNoRows {
		_, err = begin.Exec("INSERT INTO "+keyTableName+"(SEQUENCER_KEY, WORKER_ID, TIMESTAMP) VALUES(?,?,?)", key, workerId, timestamp)
		if err != nil {
			return err
		}
	} else {
		_, err = begin.Exec("UPDATE INTO "+keyTableName+"(SEQUENCER_KEY, WORKER_ID, TIMESTAMP) VALUES(?,?,?)", key, workerId, timestamp)
		if err != nil {
			return err
		}
	}
	if err = begin.Commit(); err != nil {
		begin.Rollback()
		return err
	}
	return err
}

func (s *SnowFlakeSequencer) Close() error {
	s.cancel()
	s.metadata.MysqlMetadata.Db.Close()
	return nil
}
