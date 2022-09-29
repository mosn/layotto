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
package redis

import (
	"fmt"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"mosn.io/pkg/log"

	"mosn.io/layotto/components/sequencer"
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

func TestStandaloneRedisSequencer_biggerThan_fail_reset(t *testing.T) {
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
	assert.NoError(t, err)

	//first request
	id, err := comp.GetNextId(&sequencer.GetNextIdRequest{
		Key: key,
	})
	assert.NoError(t, err)
	assert.Equal(t, defalutVal+5+1, id.NextId)

}

func TestStandaloneRedisSequencer_segment(t *testing.T) {
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

	// init
	err = comp.Init(cfg)
	assert.NoError(t, err)
	//first request
	_, id, err := comp.GetSegment(&sequencer.GetSegmentRequest{
		Key:  key,
		Size: 6,
	})
	assert.NoError(t, err)
	assert.Equal(t, defalutVal+1, id.From)
	assert.Equal(t, defalutVal+6, id.To)

}
