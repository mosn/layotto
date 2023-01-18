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
	"context"
	"strings"
	"sync"
	"testing"

	miniredis "github.com/alicebob/miniredis/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"mosn.io/pkg/log"

	"mosn.io/layotto/components/lock"
)

const (
	redisHosts  = "127.0.0.1"
	cResourceId = "resource_red_lock"
)

func TestClusterRedisLock_InitError(t *testing.T) {
	t.Run("error when connection fail", func(t *testing.T) {
		// construct component
		comp := NewClusterRedisLock(log.DefaultLogger)

		cfg := lock.Metadata{
			Properties: make(map[string]string),
		}
		cfg.Properties["redisHosts"] = redisHosts
		cfg.Properties["redisPassword"] = ""

		// init
		err := comp.Init(cfg)
		assert.Error(t, err)
	})

	t.Run("error when no host", func(t *testing.T) {
		// construct component
		comp := NewClusterRedisLock(log.DefaultLogger)

		cfg := lock.Metadata{
			Properties: make(map[string]string),
		}
		cfg.Properties["redisHosts"] = ""
		cfg.Properties["redisPassword"] = ""

		// init
		err := comp.Init(cfg)
		assert.Error(t, err)
	})

	t.Run("error when wrong MaxRetries", func(t *testing.T) {
		// construct component
		comp := NewClusterRedisLock(log.DefaultLogger)

		cfg := lock.Metadata{
			Properties: make(map[string]string),
		}
		cfg.Properties["redisHosts"] = redisHosts
		cfg.Properties["redisPassword"] = ""
		cfg.Properties["maxRetries"] = "1 "

		// init
		err := comp.Init(cfg)
		assert.Error(t, err)
	})

}

func TestClusterRedisLock_TryLock(t *testing.T) {
	// start 5 miniredis instances
	redisAddrs := make([]string, 0, 5)
	var err error
	for i := 0; i < 5; i++ {
		redis, err := miniredis.Run()
		assert.NoError(t, err)
		redisAddrs = append(redisAddrs, redis.Addr())
	}
	// construct component
	comp := NewClusterRedisLock(log.DefaultLogger)
	cfg := lock.Metadata{
		Properties: make(map[string]string),
	}
	cfg.Properties["redisHosts"] = strings.Join(redisAddrs, ",")
	cfg.Properties["redisPassword"] = ""
	// init
	err = comp.Init(cfg)
	assert.NoError(t, err)
	// 1. client1 trylock
	ownerId1 := uuid.New().String()
	resp, err := comp.TryLock(&lock.TryLockRequest{
		ResourceId: cResourceId,
		LockOwner:  ownerId1,
		Expire:     10,
	})
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	var wg sync.WaitGroup
	wg.Add(1)
	//	2. Client2 tryLock fail
	go func() {
		owner2 := uuid.New().String()
		resp2, err2 := comp.TryLock(&lock.TryLockRequest{
			ResourceId: cResourceId,
			LockOwner:  owner2,
			Expire:     10,
		})
		assert.NoError(t, err2)
		assert.False(t, resp2.Success)
		wg.Done()
	}()
	wg.Wait()
	// 3. client 1 unlock
	unlockResp, err := comp.Unlock(&lock.UnlockRequest{
		ResourceId: cResourceId,
		LockOwner:  ownerId1,
	})
	assert.NoError(t, err)
	assert.True(t, unlockResp.Status == 0, "client1 failed to unlock!")
	// 4. client 2 get lock
	wg.Add(1)
	go func() {
		owner2 := uuid.New().String()
		resp2, err2 := comp.TryLock(&lock.TryLockRequest{
			ResourceId: cResourceId,
			LockOwner:  owner2,
			Expire:     10,
		})
		assert.NoError(t, err2)
		assert.True(t, resp2.Success, "client2 failed to get lock?!")
		// 5. client2 unlock
		unlockResp, err := comp.Unlock(&lock.UnlockRequest{
			ResourceId: cResourceId,
			LockOwner:  owner2,
		})
		assert.NoError(t, err)
		assert.True(t, unlockResp.Status == 0, "client2 failed to unlock!")
		wg.Done()
	}()
	wg.Wait()
	// not implement LockKeepAlive
	keepAliveResp, err := comp.LockKeepAlive(context.TODO(), &lock.LockKeepAliveRequest{})
	assert.Nil(t, keepAliveResp)
	assert.Nil(t, err)
}
