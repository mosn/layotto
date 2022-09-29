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
	"testing"

	"github.com/go-zookeeper/zk"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"mosn.io/pkg/log"

	"mosn.io/layotto/components/pkg/mock"
	"mosn.io/layotto/components/sequencer"
)

const key = "resoure_1"

func TestZookeeperSequencer_GetNextId(t *testing.T) {
	cfg := sequencer.Configuration{
		BiggerThan: nil,
		Properties: map[string]string{
			"zookeeperHosts": "127.0.0.1",
		},
	}

	comp := NewZookeeperSequencer(log.DefaultLogger)
	comp.Init(cfg)

	//mock
	ctrl := gomock.NewController(t)
	client := mock.NewMockZKConnection(ctrl)

	path := "/" + key
	var val int32 = 1
	client.EXPECT().Set(path, []byte(""), int32(-1)).Return(&zk.Stat{Version: val}, nil).Times(1)
	val++
	client.EXPECT().Set(path, []byte(""), int32(-1)).Return(&zk.Stat{Version: val}, nil).Times(1)
	comp.client = client
	//first
	resp, err := comp.GetNextId(&sequencer.GetNextIdRequest{
		Key: key,
	})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), resp.NextId)

	//repeat
	resp, err = comp.GetNextId(&sequencer.GetNextIdRequest{
		Key: key,
	})
	assert.NoError(t, err)
	assert.Equal(t, int64(2), resp.NextId)

}
