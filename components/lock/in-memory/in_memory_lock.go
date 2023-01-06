/*
 * Copyright 2021 Layotto Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package in_memory

import (
	"context"
	"sync"
	"time"

	"mosn.io/layotto/components/lock"
)

type InMemoryLock struct {
	features []lock.Feature
	data     *lockMap
}

// memoryLock is a lock holder
type memoryLock struct {
	key        string
	owner      string
	expireTime time.Time
	lock       int
}

type lockMap struct {
	sync.Mutex
	locks map[string]*memoryLock
}

func NewInMemoryLock() *InMemoryLock {
	return &InMemoryLock{
		features: make([]lock.Feature, 0),
		data: &lockMap{
			locks: make(map[string]*memoryLock),
		},
	}
}

func (s *InMemoryLock) Init(_ lock.Metadata) error {
	return nil
}

// LockKeepAlive try to renewal lease
func (s *InMemoryLock) LockKeepAlive(ctx context.Context, request *lock.LockKeepAliveRequest) (*lock.LockKeepAliveResponse, error) {
	//TODO: implemnt function
	return nil, nil
}

func (s *InMemoryLock) Features() []lock.Feature {
	return s.features
}

// Try to add a lock. Currently this is a non-reentrant lock
func (s *InMemoryLock) TryLock(req *lock.TryLockRequest) (*lock.TryLockResponse, error) {
	s.data.Lock()
	defer s.data.Unlock()
	// 1. Find the memoryLock for this resourceId
	item, ok := s.data.locks[req.ResourceId]
	if !ok {
		item = &memoryLock{
			key: req.ResourceId,
			//0 unlock, 1 lock
			lock: 0,
		}
		s.data.locks[req.ResourceId] = item
	}

	// 2. Construct a new one if the lockData has expired
	//check expire
	if item.owner != "" && time.Now().After(item.expireTime) {
		item = &memoryLock{
			key:  req.ResourceId,
			lock: 0,
		}
		s.data.locks[req.ResourceId] = item
	}

	// 3. Check if it has been locked by others.
	// Currently this is a non-reentrant lock
	if item.lock == 1 {
		//lock failed
		return &lock.TryLockResponse{
			Success: false,
		}, nil
	}

	// 4. Update owner information
	item.lock = 1
	item.owner = req.LockOwner
	item.expireTime = time.Now().Add(time.Second * time.Duration(req.Expire))

	return &lock.TryLockResponse{
		Success: true,
	}, nil
}

func (s *InMemoryLock) Unlock(req *lock.UnlockRequest) (*lock.UnlockResponse, error) {
	s.data.Lock()
	defer s.data.Unlock()
	// 1. Find the memoryLock for this resourceId
	item, ok := s.data.locks[req.ResourceId]

	if !ok {
		return &lock.UnlockResponse{
			Status: lock.LOCK_UNEXIST,
		}, nil
	}
	// 2. check the owner information
	if item.lock != 1 {
		return &lock.UnlockResponse{
			Status: lock.LOCK_UNEXIST,
		}, nil
	}
	if item.owner != req.LockOwner {
		return &lock.UnlockResponse{
			Status: lock.LOCK_BELONG_TO_OTHERS,
		}, nil
	}
	// 3. unlock and reset the owner information
	item.owner = ""
	item.lock = 0
	return &lock.UnlockResponse{
		Status: lock.SUCCESS,
	}, nil
}
