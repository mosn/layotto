package sequencer

import (
	"context"
	"fmt"
	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"mosn.io/layotto/components/sequencer"
	"mosn.io/layotto/components/sequencer/redis"
	"mosn.io/pkg/log"
	"sync"
	"testing"
	"time"
)

const key = "resource_xxx"
const idLimit = 6000

func TestGetNextIdFromCache(t *testing.T) {
	s, err := miniredis.Run()
	assert.NoError(t, err)
	defer s.Close()
	// construct componen
	comp := redis.NewStandaloneRedisSequencer(log.DefaultLogger)
	cfg := sequencer.Configuration{
		Properties: make(map[string]string),
	}
	cfg.Properties["redisHost"] = s.Addr()
	cfg.Properties["redisPassword"] = ""
	// init
	err = comp.Init(cfg)
	assert.NoError(t, err)

	for i := 1; i < idLimit; i++ {
		support, id, err := GetNextIdFromCache(context.Background(), comp, &sequencer.GetNextIdRequest{
			Key: key,
		})
		assert.NoError(t, err)
		assert.Equal(t, true, support)
		assert.Equal(t, id, int64(i))
	}

}

func TestConcurrentGetNextIdFromCache(t *testing.T) {
	s, err := miniredis.Run()
	assert.NoError(t, err)
	defer s.Close()
	// construct componen
	comp := redis.NewStandaloneRedisSequencer(log.DefaultLogger)
	cfg := sequencer.Configuration{
		Properties: make(map[string]string),
	}
	cfg.Properties["redisHost"] = s.Addr()
	cfg.Properties["redisPassword"] = ""
	// init
	err = comp.Init(cfg)
	assert.NoError(t, err)

	var wg sync.WaitGroup
	GRCount := 100
	wg.Add(GRCount)
	startTime := time.Now().UnixMilli()
	for g := 0; g < GRCount; g++ {
		go func() {
			for i := 0; i < idLimit; i++ {
				support, id, err := GetNextIdFromCache(context.Background(), comp, &sequencer.GetNextIdRequest{
					Key: key,
				})
				fmt.Println(id)
				assert.NoError(t, err)
				assert.Equal(t, true, support)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	support, id, err := GetNextIdFromCache(context.Background(), comp, &sequencer.GetNextIdRequest{
		Key: key,
	})
	assert.NoError(t, err)
	assert.Equal(t, true, support)
	assert.Equal(t, int64(idLimit*GRCount+1), id)
	//fmt.Printf("next id: %v \n", id)
	endTime := time.Now().UnixMilli()
	fmt.Println(endTime - startTime)
	//fmt.Printf("cost time: %v ms ", endTime-startTime)
}
