//
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
	"fmt"
	"github.com/go-redis/redis/v8"
	"mosn.io/layotto/components/lock"
	"mosn.io/layotto/components/pkg/utils"
	"mosn.io/pkg/log"
	"time"
)

// Standalone Redis lock store.Any fail-over related features are not supported,such as Sentinel and Redis Cluster.
type StandaloneRedisLock struct {
	client   *redis.Client
	metadata utils.RedisMetadata
	replicas int

	features []lock.Feature
	logger   log.ErrorLogger

	ctx    context.Context
	cancel context.CancelFunc
}

// NewStandaloneRedisLock returns a new redis lock store
func NewStandaloneRedisLock(logger log.ErrorLogger) *StandaloneRedisLock {
	s := &StandaloneRedisLock{
		features: make([]lock.Feature, 0),
		logger:   logger,
	}

	return s
}

func (p *StandaloneRedisLock) Init(metadata lock.Metadata) error {
	// 1. parse config
	m, err := utils.ParseRedisMetadata(metadata.Properties)
	if err != nil {
		return err
	}
	p.metadata = m
	// 2. construct client
	p.client = utils.NewRedisClient(m)
	p.ctx, p.cancel = context.WithCancel(context.Background())
	// 3. connect to redis
	if _, err = p.client.Ping(p.ctx).Result(); err != nil {
		return fmt.Errorf("[standaloneRedisLock]: error connecting to redis at %s: %s", m.Host, err)
	}
	return err
}

func (p *StandaloneRedisLock) Features() []lock.Feature {
	return p.features
}

func (p *StandaloneRedisLock) TryLock(req *lock.TryLockRequest) (*lock.TryLockResponse, error) {
	nx := p.client.SetNX(p.ctx, req.ResourceId, req.LockOwner, time.Second*time.Duration(req.Expire))
	if nx == nil {
		return &lock.TryLockResponse{}, fmt.Errorf("[standaloneRedisLock]: SetNX returned nil.ResourceId: %s", req.ResourceId)
	}
	err := nx.Err()
	if err != nil {
		return &lock.TryLockResponse{}, err
	}

	return &lock.TryLockResponse{
		Success: nx.Val(),
	}, nil
}

const unlockScript = "local v = redis.call(\"get\",KEYS[1]); if v==false then return -1 end; if v~=ARGV[1] then return -2 else return redis.call(\"del\",KEYS[1]) end"

func (p *StandaloneRedisLock) Unlock(req *lock.UnlockRequest) (*lock.UnlockResponse, error) {
	// 1. delegate to client.eval lua script
	eval := p.client.Eval(p.ctx, unlockScript, []string{req.ResourceId}, req.LockOwner)
	// 2. check error
	if eval == nil {
		return newInternalErrorUnlockResponse(), fmt.Errorf("[standaloneRedisLock]: Eval unlock script returned nil.ResourceId: %s", req.ResourceId)
	}
	err := eval.Err()
	if err != nil {
		return newInternalErrorUnlockResponse(), err
	}
	// 3. parse result
	i, err := eval.Int()
	status := lock.INTERNAL_ERROR
	if err != nil {
		return &lock.UnlockResponse{
			Status: status,
		}, err
	}
	if i >= 0 {
		status = lock.SUCCESS
	} else if i == -1 {
		status = lock.LOCK_UNEXIST
	} else if i == -2 {
		status = lock.LOCK_BELONG_TO_OTHERS
	}
	return &lock.UnlockResponse{
		Status: status,
	}, nil
}

func newInternalErrorUnlockResponse() *lock.UnlockResponse {
	return &lock.UnlockResponse{
		Status: lock.INTERNAL_ERROR,
	}
}

func (p *StandaloneRedisLock) Close() error {
	p.cancel()

	return p.client.Close()
}
