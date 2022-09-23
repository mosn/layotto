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
package sequencer

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"mosn.io/pkg/log"

	"mosn.io/layotto/components/sequencer"
	"mosn.io/layotto/components/sequencer/redis"
)

const keyXx = "resource_xxx"
const idLimit = 12000

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
			Key: keyXx,
		})
		assert.NoError(t, err)
		assert.Equal(t, true, support)
		assert.Equal(t, id, int64(i))
	}
}
