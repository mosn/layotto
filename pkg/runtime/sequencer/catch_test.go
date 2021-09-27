package sequencer

import (
	"context"
	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"mosn.io/layotto/components/sequencer"
	"mosn.io/layotto/components/sequencer/redis"
	"mosn.io/pkg/log"
	"testing"
)

const key = "resource_xxx"

func TestGetNextIdFromCache(t *testing.T) {
	s, err := miniredis.Run()
	assert.NoError(t, err)
	defer s.Close()
	// construct component
	comp := redis.NewStandaloneRedisSequencer(log.DefaultLogger)
	cfg := sequencer.Configuration{
		Properties: make(map[string]string),
	}
	cfg.Properties["redisHost"] = s.Addr()
	cfg.Properties["redisPassword"] = ""
	// init
	err = comp.Init(cfg)
	assert.NoError(t, err)

	for i := 1; i < 8000; i++ {
		support, id, err := GetNextIdFromCache(context.Background(), comp, &sequencer.GetNextIdRequest{
			Key: key,
		})
		assert.NoError(t, err)
		assert.Equal(t, true, support)
		assert.Equal(t, id, int64(i))
	}

}
