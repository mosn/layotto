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

package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"mosn.io/pkg/log"

	_ "github.com/lib/pq"

	"mosn.io/layotto/components/pkg/utils"
	"mosn.io/layotto/components/sequencer"
)

// PostgresqlModel postgresql model
type PostgresqlModel struct {
	ID         int64  `json:"id" form:"id"` // Primary key ID
	ValueId    int64  `json:"value_id" from:"value_id"`
	BizTag     string `json:"biz_tag" form:"biz_tag"`         // Differentiated business
	CreateTime int64  `json:"create_time" from:"create_time"` // create time
	UpdateTime int64  `json:"update_time" form:"update_time"` // update time
}

type PostgresqlSequencer struct {
	biggerThan map[string]int64

	metadata utils.PostgresqlMetaData
	db       *sql.DB
	logger   log.ErrorLogger

	ctx    context.Context
	cancel context.CancelFunc
}

func NewPostgresqlSequencer(logger log.ErrorLogger) *PostgresqlSequencer {
	p := &PostgresqlSequencer{
		logger: logger,
	}
	return p
}

func (p *PostgresqlSequencer) Init(config sequencer.Configuration) error {
	meta, err := utils.ParsePostgresqlMetaData(config.Properties)
	if err != nil {
		p.logger.Errorf("init properties error: %v", err)
		return err
	}
	p.metadata = meta
	p.db = utils.NewPostgresqlCli(meta)
	p.biggerThan = config.BiggerThan

	p.ctx, p.cancel = context.WithCancel(context.Background())

	for key, value := range p.biggerThan {
		if value <= 0 {
			continue
		}
		var model PostgresqlModel
		queryParams := fmt.Sprintf(`select id, value_id, biz_tag, create_time, update_time from %s where biz_tag = $1`, p.metadata.TableName)
		err := p.db.QueryRow(queryParams, key).Scan(&model.ID, &model.ValueId, &model.BizTag, &model.CreateTime, &model.UpdateTime)
		if err != nil {
			p.logger.Errorf("get nextId error, biz_tag: %s, err: %v", key, err)
			return err
		}
		if model.ID < value {
			return fmt.Errorf("postgresql sequenccer error: can not satisfy biggerThan guarantee.key: %s,key in postgres: %s", key, p.metadata.TableName)
		}
	}
	return nil
}

func (p *PostgresqlSequencer) GetNextId(req *sequencer.GetNextIdRequest) (*sequencer.GetNextIdResponse, error) {

	var model PostgresqlModel
	queryParams := fmt.Sprintf(`select id, value_id, biz_tag, create_time, update_time from %s where biz_tag = $1`, p.metadata.TableName)
	err := p.db.QueryRow(queryParams, req.Key).Scan(&model.ID, &model.ValueId, &model.BizTag, &model.CreateTime, &model.UpdateTime)
	if err != nil {
		return nil, err
	}
	var res sql.Result
	updateParams := fmt.Sprintf(`update %s set value_id = value_id + 1, update_time = $1 where biz_tag = $2`, p.metadata.TableName)
	res, err = p.db.Exec(updateParams, time.Now().Unix(), req.Key)
	if err != nil {
		return nil, err
	}

	rowsId, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsId == 0 {
		return nil, errors.New("no update")
	}

	if err != nil {
		p.logger.Errorf("get nextId error, biz_tag: %s, err: %v", req.Key, err)
		return nil, err
	}

	return &sequencer.GetNextIdResponse{
		NextId: model.ValueId,
	}, nil
}

func (p *PostgresqlSequencer) GetSegment(req *sequencer.GetSegmentRequest) (bool, *sequencer.GetSegmentResponse, error) {
	if req.Size == 0 {
		return false, nil, nil
	}
	var model PostgresqlModel
	queryParams := fmt.Sprintf(`select id, value_id, biz_tag, create_time, update_time from %s where biz_tag = $1`, p.metadata.TableName)
	err := p.db.QueryRow(queryParams, req.Key).Scan(&model.ID, &model.ValueId, &model.BizTag, &model.CreateTime, &model.UpdateTime)
	if err != nil {
		return false, nil, err
	}
	var res sql.Result
	model.ValueId += int64(req.Size)
	updateParams := fmt.Sprintf(`update %v set value_id = $1, update_time = $2 where biz_tag = $3`, p.metadata.TableName)
	res, err = p.db.Exec(updateParams, model.ValueId, time.Now().Unix(), req.Key)
	if err != nil {
		return false, nil, err
	}

	rowsId, err := res.RowsAffected()
	if err != nil {
		return false, nil, err
	}
	if rowsId == 0 {
		return false, nil, errors.New("no update")
	}

	if err != nil {
		return false, nil, err
	}

	return false, &sequencer.GetSegmentResponse{
		From: model.ValueId - int64(req.Size) + 1,
		To:   model.ValueId,
	}, nil
}

func (p *PostgresqlSequencer) Close() error {
	p.cancel()
	return p.db.Close()
}
