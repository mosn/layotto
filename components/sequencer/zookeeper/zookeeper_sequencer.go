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
package zookeeper

import (
	"context"
	"fmt"

	"github.com/go-zookeeper/zk"
	"mosn.io/pkg/log"

	"mosn.io/layotto/components/pkg/utils"
	"mosn.io/layotto/components/sequencer"
)

const maxInt32 = 2147483647

type ZookeeperSequencer struct {
	client     utils.ZKConnection
	metadata   utils.ZookeeperMetadata
	BiggerThan map[string]int64
	logger     log.ErrorLogger
	factory    utils.ConnectionFactory
	ctx        context.Context
	cancel     context.CancelFunc
}

// NewZookeeperSequencer returns a new zookeeper sequencer
func NewZookeeperSequencer(logger log.ErrorLogger) *ZookeeperSequencer {
	s := &ZookeeperSequencer{
		logger: logger,
	}

	return s
}

func (s *ZookeeperSequencer) Init(config sequencer.Configuration) error {
	m, err := utils.ParseZookeeperMetadata(config.Properties)
	if err != nil {
		return err
	}
	//init
	s.metadata = m
	s.BiggerThan = config.BiggerThan
	s.factory = &utils.ConnectionFactoryImpl{}
	connection, err := s.factory.NewConnection(0, s.metadata)
	if err != nil {
		return err
	}
	s.client = connection
	s.ctx, s.cancel = context.WithCancel(context.Background())

	//check biggerThan
	for k, needV := range s.BiggerThan {
		if needV <= 0 {
			continue
		}

		if needV >= maxInt32 {
			return fmt.Errorf("the maximum value of zookeeper version cannot exceed int32")
		}
		_, stat, err := s.client.Get("/" + k)
		if err != nil {
			//key not exist
			if err == zk.ErrNoNode {
				return fmt.Errorf("zookeeper sequencer error: can not satisfy biggerThan guarantee.key: %s, current key does not exist", k)
			}
			//other error
			return err
		}
		realV := int64(stat.Version)

		if realV < needV {
			return fmt.Errorf("zookeeper sequencer error: can not satisfy biggerThan guarantee.key: %s,current id:%v", k, realV)
		}

	}
	return err

}

func (s *ZookeeperSequencer) GetNextId(req *sequencer.GetNextIdRequest) (*sequencer.GetNextIdResponse, error) {

	stat, err := s.client.Set("/"+req.Key, []byte(""), -1)

	if err != nil {
		if err == zk.ErrNoNode {
			_, errCreate := s.client.Create("/"+req.Key, []byte(""), zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
			if errCreate != nil {
				return nil, errCreate
			}
			return s.GetNextId(req)
		}
		return nil, err
	}
	// create node version=0, every time we set node  will result in version+1
	// so if version=0, an overflow int32 has occurred
	//but this time return error ,what to do next time ？
	if stat.Version <= 0 {
		s.logger.Errorf("an overflow int32 has occurred in zookeeper , the key is %s", req.Key)
		return nil, fmt.Errorf("an overflow int32 has occurred in zookeeper, the key is %s", req.Key)
	}

	return &sequencer.GetNextIdResponse{
		NextId: int64(stat.Version),
	}, nil
}

func (s *ZookeeperSequencer) GetSegment(req *sequencer.GetSegmentRequest) (support bool, result *sequencer.GetSegmentResponse, err error) {
	return false, nil, nil
}
func (s *ZookeeperSequencer) Close() error {
	s.cancel()
	s.client.Close()
	return nil
}
