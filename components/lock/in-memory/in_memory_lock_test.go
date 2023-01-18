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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"mosn.io/layotto/components/lock"
)

func TestNew(t *testing.T) {
	s := NewInMemoryLock()
	assert.NotNil(t, s)
}

func TestInit(t *testing.T) {
	s := NewInMemoryLock()
	assert.NotNil(t, s)

	err := s.Init(lock.Metadata{})
	assert.NoError(t, err)
}

func TestFeatures(t *testing.T) {
	s := NewInMemoryLock()
	assert.NotNil(t, s)

	f := s.Features()
	assert.NotNil(t, f)
	assert.Equal(t, 0, len(f))
}

func TestTryLock(t *testing.T) {
	s := NewInMemoryLock()
	assert.NotNil(t, s)

	req := &lock.TryLockRequest{
		ResourceId: "key111",
		LockOwner:  "own",
		Expire:     3,
	}

	var err error
	var resp *lock.TryLockResponse
	resp, err = s.TryLock(req)
	assert.NoError(t, err)
	assert.NotNil(t, req)
	assert.True(t, resp.Success)

	resp, err = s.TryLock(req)
	assert.NoError(t, err)
	assert.NotNil(t, req)
	assert.False(t, resp.Success)

	req = &lock.TryLockRequest{
		ResourceId: "key112",
		LockOwner:  "own",
		Expire:     1,
	}

	resp, err = s.TryLock(req)
	assert.NoError(t, err)
	assert.NotNil(t, req)
	assert.True(t, resp.Success)

	req = &lock.TryLockRequest{
		ResourceId: "key112",
		LockOwner:  "own",
		Expire:     1,
	}

	resp, err = s.TryLock(req)
	assert.NoError(t, err)
	assert.NotNil(t, req)
	assert.False(t, resp.Success)

	s.data.locks["key112"].expireTime = time.Now().Add(-2 * time.Second)

	resp, err = s.TryLock(req)
	assert.NoError(t, err)
	assert.NotNil(t, req)
	assert.True(t, resp.Success)

}

func TestUnLock(t *testing.T) {
	s := NewInMemoryLock()
	assert.NotNil(t, s)

	req := &lock.UnlockRequest{
		ResourceId: "key111",
		LockOwner:  "own",
	}

	var err error
	var resp *lock.UnlockResponse
	resp, err = s.Unlock(req)
	assert.NoError(t, err)
	assert.NotNil(t, req)
	assert.Equal(t, lock.LOCK_UNEXIST, resp.Status)

	lockReq := &lock.TryLockRequest{
		ResourceId: "key111",
		LockOwner:  "own",
		Expire:     10,
	}

	var lockResp *lock.TryLockResponse
	lockResp, err = s.TryLock(lockReq)
	assert.NoError(t, err)
	assert.NotNil(t, req)
	assert.True(t, lockResp.Success)

	resp, err = s.Unlock(req)
	assert.NoError(t, err)
	assert.NotNil(t, req)
	assert.Equal(t, lock.SUCCESS, resp.Status)

	lockResp, err = s.TryLock(lockReq)
	assert.NoError(t, err)
	assert.NotNil(t, req)
	assert.True(t, lockResp.Success)

	req.LockOwner = "1"

	resp, err = s.Unlock(req)
	assert.NoError(t, err)
	assert.NotNil(t, req)
	assert.Equal(t, lock.LOCK_BELONG_TO_OTHERS, resp.Status)

	req.ResourceId = "11"
	lockReq.ResourceId = "11"
	req.LockOwner = "own1"
	lockReq.LockOwner = "own1"
	lockResp, err = s.TryLock(lockReq)
	assert.NoError(t, err)
	assert.NotNil(t, req)
	assert.True(t, lockResp.Success)

	resp, err = s.Unlock(req)
	assert.NoError(t, err)
	assert.NotNil(t, req)
	assert.Equal(t, lock.SUCCESS, resp.Status)

	resp, err = s.Unlock(req)
	assert.NoError(t, err)
	assert.NotNil(t, req)
	assert.Equal(t, lock.LOCK_UNEXIST, resp.Status)

	// not implement LockKeepAlive
	keepAliveResp, err := s.LockKeepAlive(context.TODO(), &lock.LockKeepAliveRequest{})
	assert.Nil(t, keepAliveResp)
	assert.Nil(t, err)
}
