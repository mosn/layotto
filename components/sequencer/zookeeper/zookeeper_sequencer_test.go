package zookeeper

import (
	"github.com/go-zookeeper/zk"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"mosn.io/layotto/components/sequencer"
	"mosn.io/pkg/log"
	"testing"
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
	client := NewMockZKConnection(ctrl)

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
