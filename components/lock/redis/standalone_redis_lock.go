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
	"time"

	"github.com/go-redis/redis/v8"
	"mosn.io/pkg/log"

	"mosn.io/layotto/components/lock"
	"mosn.io/layotto/components/pkg/utils"
)

// Standalone Redis lock store.Any fail-over related features are not supported,such as Sentinel and Redis Cluster.
type StandaloneRedisLock struct {
	client   *redis.Client
	metadata utils.RedisMetadata

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

// Init StandaloneRedisLock
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

// Features is to get StandaloneRedisLock's features
func (p *StandaloneRedisLock) Features() []lock.Feature {
	return p.features
}

// LockKeepAlive try to renewal lease
func (p *StandaloneRedisLock) LockKeepAlive(ctx context.Context, request *lock.LockKeepAliveRequest) (*lock.LockKeepAliveResponse, error) {
	//TODO: implemnt function
	return nil, nil
}

// Node tries to acquire a redis lock
func (p *StandaloneRedisLock) TryLock(req *lock.TryLockRequest) (*lock.TryLockResponse, error) {
	// 1.Setting redis expiration time
	nx := p.client.SetNX(p.ctx, req.ResourceId, req.LockOwner, time.Second*time.Duration(req.Expire))
	if nx == nil {
		return &lock.TryLockResponse{}, fmt.Errorf("[standaloneRedisLock]: SetNX returned nil.ResourceId: %s", req.ResourceId)
	}
	// 2. check error
	err := nx.Err()
	if err != nil {
		return &lock.TryLockResponse{}, err
	}

	return &lock.TryLockResponse{
		Success: nx.Val(),
	}, nil
}

const unlockScript = "local v = redis.call(\"get\",KEYS[1]); if v==false then return -1 end; if v~=ARGV[1] then return -2 else return redis.call(\"del\",KEYS[1]) end"

// Node tries to release a redis lock
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

// newInternalErrorUnlockResponse is to return lock release error
func newInternalErrorUnlockResponse() *lock.UnlockResponse {
	return &lock.UnlockResponse{
		Status: lock.INTERNAL_ERROR,
	}
}

// Close shuts down the client's redis connections.
func (p *StandaloneRedisLock) Close() error {
	if p.cancel != nil {
		p.cancel()
	}
	if p.client != nil {
		return p.client.Close()
	}
	return nil
}
