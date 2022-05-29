package postgresql

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

/**
* @Author: azh
* @Date: 2022/5/13 22:11
* @Context:
 */

import (
	"context"
	"mosn.io/pkg/log"

	"mosn.io/layotto/components/pkg/utils"
	"mosn.io/layotto/components/sequencer/postgresql/model"
	"mosn.io/layotto/components/sequencer/postgresql/service"

	"mosn.io/layotto/components/sequencer"
)

type PostgresqlSequencer struct {
	biggerThan map[string]int64
	client     *service.PostgresqlService

	logger log.ErrorLogger

	ctx    context.Context
	cancel context.CancelFunc
}

var PostgresqlConfigFilePath = "D:\\goTest\\test\\layotto\\components\\sequencer\\postgresql\\conf\\postgresql.yaml"

// NewPostgresqlSequencer returns a new postgresql sequencer
func NewPostgresqlSequencer(logger log.ErrorLogger) *PostgresqlSequencer {
	s := &PostgresqlSequencer{
		logger: logger,
	}
	return s
}

func (p *PostgresqlSequencer) Init(config sequencer.Configuration) error {
	s, err := utils.InitPostgresql(config.Properties)
	if err != nil {
		p.logger.Infof("init config error")
		return err
	}
	for key, value := range p.biggerThan {
		err := p.client.InitMaxId(p.ctx, key, value, service.DEFAULT_STEP)
		if err != nil {
			p.logger.Infof("init max_id error")
			return err
		}
	}
	p.client = s
	p.ctx, p.cancel = context.WithCancel(context.Background())

	return nil
}

// Create The user can initialize the ID sequence according to the customization.
//The dimension takes the business as the dimension biz_tag
func (p *PostgresqlSequencer) Create(model *model.PostgresqlModel) error {
	err := p.client.Create(p.ctx, model)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresqlSequencer) GetNextId(req *sequencer.GetNextIdRequest) (*sequencer.GetNextIdResponse, error) {
	id, err := p.client.GetId(p.ctx, req.Key)
	if err != nil {
		return nil, err
	}

	return &sequencer.GetNextIdResponse{
		NextId: int64(id),
	}, nil
}

// GetSegment In fact, the runtime cache is not very useful, because my system has implemented the dual buffer mode~
func (p *PostgresqlSequencer) GetSegment(req *sequencer.GetSegmentRequest) (bool, *sequencer.GetSegmentResponse, error) {

	if req.Size == 0 {
		return true, nil, nil
	}

	by, err := p.client.GetId(p.ctx, req.Key)
	if err != nil {
		return true, nil, err
	}
	by = by + uint64(req.Size)

	return true, &sequencer.GetSegmentResponse{
		From: int64(by) - int64(req.Size) + 1,
		To:   int64(by),
	}, nil
}

func (p *PostgresqlSequencer) Close() error {
	p.cancel()
	return nil
}
