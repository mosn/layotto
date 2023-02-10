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
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	msync "mosn.io/mosn/pkg/sync"
	"mosn.io/pkg/log"

	"mosn.io/layotto/components/lock"
	"mosn.io/layotto/components/pkg/utils"
)

// RedLock
// it will be best to use at least 5 hosts
type ClusterRedisLock struct {
	clients  []*redis.Client
	metadata utils.RedisClusterMetadata
	workpool msync.WorkerPool

	features []lock.Feature
	logger   log.ErrorLogger

	ctx    context.Context
	cancel context.CancelFunc
}

// NewClusterRedisLock returns a new redis lock store
func NewClusterRedisLock(logger log.ErrorLogger) *ClusterRedisLock {
	s := &ClusterRedisLock{
		features: make([]lock.Feature, 0),
		logger:   logger,
	}

	return s
}

type resultMsg struct {
	error        error
	host         string
	lockStatus   bool
	unlockStatus lock.LockStatus
}

func (c *ClusterRedisLock) Init(metadata lock.Metadata) error {

	m, err := utils.ParseRedisClusterMetadata(metadata.Properties)
	if err != nil {
		return err
	}
	c.metadata = m
	c.clients = utils.NewClusterRedisClient(m)
	c.ctx, c.cancel = context.WithCancel(context.Background())
	c.workpool = msync.NewWorkerPool(m.Concurrency)
	for i, client := range c.clients {
		if _, err = client.Ping(c.ctx).Result(); err != nil {
			return fmt.Errorf("[ClusterRedisLock]: error connecting to redis at %s: %s", c.metadata.Hosts[i], err)
		}
	}
	return err
}

func (c *ClusterRedisLock) Features() []lock.Feature {
	return c.features
}

// LockKeepAlive try to renewal lease
func (c *ClusterRedisLock) LockKeepAlive(ctx context.Context, request *lock.LockKeepAliveRequest) (*lock.LockKeepAliveResponse, error) {
	//TODO: implemnt function
	return nil, nil
}

func (c *ClusterRedisLock) TryLock(req *lock.TryLockRequest) (*lock.TryLockResponse, error) {
	//try to get lock on all redis nodes
	intervalStart := utils.GetMiliTimestamp(time.Now().UnixNano())
	//intervalLimit must be 1/10 of expire time to make sure time of lock far less than expire time
	intervalLimit := int64(req.Expire) * 1000 / 10
	wg := sync.WaitGroup{}
	wg.Add(len(c.clients))

	//resultChan will be used to collect results of getting lock
	resultChan := make(chan resultMsg, len(c.clients))

	//getting lock concurrently
	for i := range c.clients {
		clientIndex := i
		c.workpool.Schedule(func() {
			c.LockSingleRedis(clientIndex, req, &wg, resultChan)
		})
	}
	wg.Wait()
	intervalEnd := utils.GetMiliTimestamp(time.Now().UnixNano())

	//make sure time interval of locking far less than expire time
	if intervalLimit < intervalEnd-intervalStart {
		_, _ = c.UnlockAllRedis(&lock.UnlockRequest{
			ResourceId: req.ResourceId,
			LockOwner:  req.LockOwner,
		}, &wg)
		return &lock.TryLockResponse{
			Success: false,
		}, fmt.Errorf("[ClusterRedisLock]: lock timeout. ResourceId: %s", req.ResourceId)
	}
	close(resultChan)

	successCount := 0
	errorStrs := make([]string, 0, len(c.clients))
	for msg := range resultChan {
		if msg.error != nil {
			errorStrs = append(errorStrs, msg.error.Error())
			continue
		}
		if msg.lockStatus {
			successCount++
		}
	}
	var err error
	if len(errorStrs) > 0 {
		err = fmt.Errorf(strings.Join(errorStrs, "\n"))
	}
	//getting lock on majority of redis cluster will be regarded as locking success
	if successCount*2 > len(c.clients) {
		return &lock.TryLockResponse{
			Success: true,
		}, err
	}

	_, unlockErr := c.UnlockAllRedis(&lock.UnlockRequest{
		ResourceId: req.ResourceId,
		LockOwner:  req.LockOwner,
	}, &wg)
	if unlockErr != nil {
		errorStrs = append(errorStrs, unlockErr.Error())
		err = fmt.Errorf(strings.Join(errorStrs, "\n"))
	}
	return &lock.TryLockResponse{
		Success: false,
	}, err
}

func (c *ClusterRedisLock) Unlock(req *lock.UnlockRequest) (*lock.UnlockResponse, error) {
	wg := sync.WaitGroup{}
	//err means there were some internal errors,then the status must be INTERNAL_ERROR
	//the LOCK_UNEXIST and LOCK_BELONG_TO_OTHERS status codes can be ignore
	//becauce they means the lock of the current redis
	//returned the status code don't need to be unlocked by current invoking
	_, err := c.UnlockAllRedis(req, &wg)
	if err != nil {
		return newInternalErrorUnlockResponse(), err
	}
	return &lock.UnlockResponse{
		Status: lock.SUCCESS,
	}, nil
}

func (c *ClusterRedisLock) UnlockAllRedis(req *lock.UnlockRequest, wg *sync.WaitGroup) (lock.LockStatus, error) {
	wg.Add(len(c.clients))
	ch := make(chan resultMsg, len(c.clients))

	//unlock concurrently
	for i := range c.clients {
		clientIndex := i
		c.workpool.Schedule(func() {
			c.UnlockSingleRedis(clientIndex, req, wg, ch)
		})
	}
	wg.Wait()
	close(ch)
	errorStrs := make([]string, 0, len(c.clients))
	status := lock.SUCCESS

	//collect result of unlocking
	for msg := range ch {
		if msg.unlockStatus == lock.INTERNAL_ERROR {
			status = msg.unlockStatus
			errorStrs = append(errorStrs, msg.error.Error())
		}
	}
	if len(errorStrs) > 0 {
		return status, fmt.Errorf(strings.Join(errorStrs, "\n"))
	}
	return status, nil
}

func (c *ClusterRedisLock) LockSingleRedis(clientIndex int, req *lock.TryLockRequest, wg *sync.WaitGroup, ch chan resultMsg) {
	defer wg.Done()
	msg := resultMsg{
		host: c.metadata.Hosts[clientIndex],
	}
	nx := c.clients[clientIndex].SetNX(c.ctx, req.ResourceId, req.LockOwner, time.Second*time.Duration(req.Expire))
	if nx == nil {
		msg.error = fmt.Errorf("[ClusterRedisLock]: SetNX returned nil. host: %s \n ResourceId: %s", c.clients[clientIndex], req.ResourceId)
		ch <- msg
		return
	}
	if nx.Err() != nil {
		msg.error = fmt.Errorf("[ClusterRedisLock]: %s host: %s \n ResourceId: %s", nx.Err().Error(), c.clients[clientIndex], req.ResourceId)
	}
	msg.lockStatus = nx.Val()
	ch <- msg
}

func (c *ClusterRedisLock) UnlockSingleRedis(clientIndex int, req *lock.UnlockRequest, wg *sync.WaitGroup, ch chan resultMsg) {
	defer wg.Done()
	eval := c.clients[clientIndex].Eval(c.ctx, unlockScript, []string{req.ResourceId}, req.LockOwner)
	msg := resultMsg{}
	msg.unlockStatus = lock.INTERNAL_ERROR
	if eval == nil {
		msg.error = fmt.Errorf("[ClusterRedisLock]: Eval unlock script returned nil. host: %s \n ResourceId: %s", c.clients[clientIndex], req.ResourceId)
		ch <- msg
		return
	}
	if eval.Err() != nil {
		msg.error = fmt.Errorf("[ClusterRedisLock]: %s host: %s \n ResourceId: %s", eval.Err().Error(), c.clients[clientIndex], req.ResourceId)
		ch <- msg
		return
	}
	i, err := eval.Int()
	if err != nil {
		msg.error = err
		ch <- msg
		return
	}
	if i >= 0 {
		msg.unlockStatus = lock.SUCCESS
	} else if i == -1 {
		msg.unlockStatus = lock.LOCK_UNEXIST
	} else if i == -2 {
		msg.unlockStatus = lock.LOCK_BELONG_TO_OTHERS
	}
	ch <- msg
}
