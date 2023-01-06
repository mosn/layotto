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
package consul

import (
	"context"
	"runtime"
	"strconv"
	"sync"

	"github.com/hashicorp/consul/api"
	msync "mosn.io/mosn/pkg/sync"
	"mosn.io/pkg/log"

	"mosn.io/layotto/components/lock"
	"mosn.io/layotto/components/pkg/utils"
)

type ConsulLock struct {
	metadata       utils.ConsulMetadata
	logger         log.ErrorLogger
	client         utils.ConsulClient
	sessionFactory utils.SessionFactory
	kv             utils.ConsulKV
	sMap           sync.Map
	workPool       msync.WorkerPool
}

func NewConsulLock(logger log.ErrorLogger) *ConsulLock {
	consulLock := &ConsulLock{logger: logger}
	return consulLock
}

func (c *ConsulLock) Init(metadata lock.Metadata) error {
	consulMetadata, err := utils.ParseConsulMetadata(metadata)
	if err != nil {
		return err
	}
	c.metadata = consulMetadata
	client, err := api.NewClient(&api.Config{
		Address: consulMetadata.Address,
		Scheme:  consulMetadata.Scheme,
	})
	if err != nil {
		return err
	}
	c.client = client
	c.sessionFactory = client.Session()
	c.kv = client.KV()
	c.workPool = msync.NewWorkerPool(runtime.NumCPU())
	return nil
}
func (c *ConsulLock) Features() []lock.Feature {
	return nil
}

// LockKeepAlive try to renewal lease
func (c *ConsulLock) LockKeepAlive(ctx context.Context, request *lock.LockKeepAliveRequest) (*lock.LockKeepAliveResponse, error) {
	//TODO: implemnt function
	return nil, nil
}

func getTTL(expire int32) string {
	//session TTL must be between [10s=24h0m0s]
	if expire < 10 {
		expire = 10
	}
	return strconv.Itoa(int(expire)) + "s"
}

func (c *ConsulLock) TryLock(req *lock.TryLockRequest) (*lock.TryLockResponse, error) {

	// create a session TTL
	session, _, err := c.sessionFactory.Create(&api.SessionEntry{
		TTL:       getTTL(req.Expire),
		LockDelay: 0,
		Behavior:  "delete", //Controls the behavior to delete when a session is invalidated.
	}, nil)

	if err != nil {
		return nil, err
	}

	// put a new KV pair with ttl session
	p := &api.KVPair{Key: req.ResourceId, Value: []byte(req.LockOwner), Session: session}
	//acquire lock
	acquire, _, err := c.kv.Acquire(p, nil)

	if err != nil {
		return nil, err
	}

	if acquire {
		//bind lockOwner+resourceId and session
		c.sMap.Store(req.LockOwner+"-"+req.ResourceId, session)
		c.workPool.Schedule(generateGCTask(req.Expire, &c.sMap, req.LockOwner+"-"+req.ResourceId))
		return &lock.TryLockResponse{
			Success: true,
		}, nil
	}
	return &lock.TryLockResponse{
		Success: false,
	}, nil
}
func (c *ConsulLock) Unlock(req *lock.UnlockRequest) (*lock.UnlockResponse, error) {

	session, ok := c.sMap.Load(req.LockOwner + "-" + req.ResourceId)

	if !ok {
		return &lock.UnlockResponse{Status: lock.LOCK_UNEXIST}, nil
	}
	// put a new KV pair with ttl session
	p := &api.KVPair{Key: req.ResourceId, Value: []byte(req.LockOwner), Session: session.(string)}
	//release lock
	release, _, err := c.kv.Release(p, nil)

	if err != nil {
		return &lock.UnlockResponse{Status: lock.INTERNAL_ERROR}, nil
	}

	if release {
		c.sMap.Delete(req.LockOwner + "-" + req.ResourceId)
		_, err = c.sessionFactory.Destroy(session.(string), nil)
		if err != nil {
			c.logger.Errorf("consul lock session destroy error: %v", err)
		}
		return &lock.UnlockResponse{Status: lock.
			SUCCESS}, nil
	}
	return &lock.UnlockResponse{Status: lock.LOCK_BELONG_TO_OTHERS}, nil
}
