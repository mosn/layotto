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
	"go.uber.org/atomic"
	"mosn.io/layotto/components/lock"
	"sync"
	"time"
)

type InMemoryLock struct {
	features []lock.Feature
	data     *sync.Map
	wLock    sync.Mutex
}

type lockData struct {
	key        string
	owner      string
	expireTime time.Time
	lock       *atomic.Int32
}

func NewInMemoryLock() *InMemoryLock {
	return &InMemoryLock{
		features: make([]lock.Feature, 0),
		data:     &sync.Map{},
	}
}

func (s *InMemoryLock) Init(_ lock.Metadata) error {
	return nil
}

func (s *InMemoryLock) Features() []lock.Feature {
	return s.features
}

func (s *InMemoryLock) TryLock(req *lock.TryLockRequest) (*lock.TryLockResponse, error) {
	item, ok := s.data.Load(req.ResourceId)
	if !ok {
		newItem := &lockData{
			key:  req.ResourceId,
			lock: &atomic.Int32{},
		}
		s.wLock.Lock()
		item, _ = s.data.LoadOrStore(req.ResourceId, newItem)
		s.wLock.Unlock()
	}

	//0 unlock, 1 lock
	d := item.(*lockData)

	//check expire
	if d.owner != "" && time.Now().Before(d.expireTime) {
		s.wLock.Lock()
		//double check
		s.data.Delete(req.ResourceId)
		item = &lockData{
			key:  req.ResourceId,
			lock: &atomic.Int32{},
		}
		s.data.Store(req.ResourceId, item)
		s.wLock.Unlock()
	}

	if !d.lock.CAS(0, 1) {
		//lock failed
		return &lock.TryLockResponse{
			Success: false,
		}, nil
	}

	d.owner = req.LockOwner
	d.expireTime = time.Now().Add(time.Second * time.Duration(req.Expire))
	return &lock.TryLockResponse{
		Success: true,
	}, nil
}

func (s *InMemoryLock) Unlock(req *lock.UnlockRequest) (*lock.UnlockResponse, error) {
	item, ok := s.data.Load(req.ResourceId)

	if !ok {
		return &lock.UnlockResponse{
			Status: lock.LOCK_UNEXIST,
		}, nil
	}

	d := item.(*lockData)
	if d.lock.Load() != 1 {
		return &lock.UnlockResponse{
			Status: lock.SUCCESS,
		}, nil
	}

	if d.owner != req.LockOwner {
		return &lock.UnlockResponse{
			Status: lock.LOCK_BELONG_TO_OTHERS,
		}, nil
	}

	if !d.lock.CAS(1, 0) {
		return &lock.UnlockResponse{
			Status: lock.LOCK_UNEXIST,
		}, nil
	}

	d.owner = ""
	return &lock.UnlockResponse{
		Status: lock.SUCCESS,
	}, nil
}
