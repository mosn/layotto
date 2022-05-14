package postgresql

/**
* @Author: azh
* @Date: 2022/5/13 22:11
* @Context:
 */

import (
	"context"
	"fmt"
	"mosn.io/layotto/components/pkg/utils"
	"mosn.io/layotto/components/sequencer/postgresql/model"
	"mosn.io/layotto/components/sequencer/postgresql/service"
	"mosn.io/pkg/log"

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
		fmt.Println("init config error")
		return err
	}
	for key, value := range p.biggerThan {
		err := p.client.InitMaxId(p.ctx, key, value, service.DEFAULT_STEP)
		if err != nil {
			fmt.Println("init max_id error")
			return err
		}
	}
	p.client = s
	p.ctx, p.cancel = context.WithCancel(context.Background())

	return nil
}

// Create 用户可以根据自定义去初始化id序列，纬度是以业务为纬度，biz_tag
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

// GetSegment 其实该runtime cache 用处不大，因为我系统已经实现双buffer模式了~
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

