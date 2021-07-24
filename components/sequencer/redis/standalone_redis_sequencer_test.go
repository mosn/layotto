package redis

import (
	"fmt"
	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"mosn.io/layotto/components/sequencer"
	"mosn.io/pkg/log"
	"testing"
)

const key = "resource_xxx"

func TestStandaloneRedisSequencer(t *testing.T) {
	s, err := miniredis.Run()
	assert.NoError(t, err)
	defer s.Close()
	// construct component
	comp := NewStandaloneRedisSequencer(log.DefaultLogger)
	cfg := sequencer.Configuration{
		Properties: make(map[string]string),
	}
	cfg.Properties["redisHost"] = s.Addr()
	cfg.Properties["redisPassword"] = ""
	// init
	err = comp.Init(cfg)
	assert.NoError(t, err)
	//first request
	id, err := comp.GetNextId(&sequencer.GetNextIdRequest{
		Key: key,
	})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), id.NextId)

	//again
	id, err = comp.GetNextId(&sequencer.GetNextIdRequest{
		Key: key,
	})
	assert.NoError(t, err)
	assert.Equal(t, int64(2), id.NextId)
}

func TestStandaloneRedisSequencer_biggerThan_success(t *testing.T) {
	s, err := miniredis.Run()
	assert.NoError(t, err)
	defer s.Close()
	// construct component
	comp := NewStandaloneRedisSequencer(log.DefaultLogger)
	cfg := sequencer.Configuration{
		Properties: make(map[string]string),
	}
	defalutVal := int64(20)
	//init kv
	s.Set(key, fmt.Sprint(defalutVal))
	cfg.Properties["redisHost"] = s.Addr()
	cfg.Properties["redisPassword"] = ""
	cfg.BiggerThan = map[string]int64{
		key: defalutVal,
	}

	// init
	err = comp.Init(cfg)
	assert.NoError(t, err)
	//first request
	id, err := comp.GetNextId(&sequencer.GetNextIdRequest{
		Key: key,
	})
	assert.NoError(t, err)
	assert.Equal(t, defalutVal+1, id.NextId)

	//again
	id, err = comp.GetNextId(&sequencer.GetNextIdRequest{
		Key: key,
	})
	assert.NoError(t, err)
	assert.Equal(t, defalutVal+2, id.NextId)
}

func TestStandaloneRedisSequencer_biggerThan_fail(t *testing.T) {
	s, err := miniredis.Run()
	assert.NoError(t, err)
	defer s.Close()
	// construct component
	comp := NewStandaloneRedisSequencer(log.DefaultLogger)
	cfg := sequencer.Configuration{
		Properties: make(map[string]string),
	}
	defalutVal := int64(20)
	//init kv
	s.Set(key, fmt.Sprint(defalutVal))
	// init config
	cfg.Properties["redisHost"] = s.Addr()
	cfg.Properties["redisPassword"] = ""
	cfg.BiggerThan = map[string]int64{
		key: defalutVal + 5,
	}
	err = comp.Init(cfg)
	assert.Error(t, err)

}
