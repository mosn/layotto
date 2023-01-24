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
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/assert"
	"mosn.io/pkg/log"

	"mosn.io/layotto/components/lock"
	"mosn.io/layotto/components/pkg/mock"
)

const resouseId = "resoure_1"
const lockOwerA = "p1"
const lockOwerB = "p2"
const expireTime = 5

// Init with wrong config
func TestConsulLock_InitWithWrongConfig(t *testing.T) {
	t.Run("when no address then error", func(t *testing.T) {
		comp := NewConsulLock(log.DefaultLogger)
		cfg := lock.Metadata{
			Properties: make(map[string]string),
		}
		//cfg.Properties["address"] = "127.0.0.1:8500"
		err := comp.Init(cfg)
		assert.Error(t, err)
	})
}

// Test features
func TestConsulLock_Features(t *testing.T) {
	comp := NewConsulLock(log.DefaultLogger)
	assert.True(t, len(comp.Features()) == 0)
}

// A lock A unlock
func TestConsulLock_TryLock(t *testing.T) {
	//mock
	ctrl := gomock.NewController(t)
	client := mock.NewMockConsulClient(ctrl)
	factory := mock.NewMockSessionFactory(ctrl)
	kv := mock.NewMockConsulKV(ctrl)

	comp := NewConsulLock(log.DefaultLogger)
	cfg := lock.Metadata{
		Properties: make(map[string]string),
	}
	cfg.Properties["address"] = "127.0.0.1:8500"
	err := comp.Init(cfg)
	assert.Nil(t, err)
	comp.client = client
	comp.sessionFactory = factory
	comp.kv = kv
	factory.EXPECT().Create(&api.SessionEntry{TTL: getTTL(expireTime), LockDelay: 0, Behavior: "delete"}, nil).
		Return("session1", nil, nil).Times(1)
	factory.EXPECT().Destroy("session1", nil).Return(nil, nil).Times(1)
	kv.EXPECT().Acquire(&api.KVPair{Key: resouseId, Value: []byte(lockOwerA), Session: "session1"}, nil).
		Return(true, nil, nil).Times(1)
	kv.EXPECT().Release(&api.KVPair{Key: resouseId, Value: []byte(lockOwerA), Session: "session1"}, nil).
		Return(true, nil, nil).Times(1)

	tryLock, err := comp.TryLock(&lock.TryLockRequest{
		ResourceId: resouseId,
		LockOwner:  lockOwerA,
		Expire:     expireTime,
	})

	assert.NoError(t, err)
	assert.Equal(t, true, tryLock.Success)

	unlock, err := comp.Unlock(&lock.UnlockRequest{
		ResourceId: resouseId,
		LockOwner:  lockOwerA,
	})

	assert.NoError(t, err)
	assert.Equal(t, lock.SUCCESS, unlock.Status)

}

// A lock B lock
func TestConsulLock_ALock_BLock(t *testing.T) {

	//mock
	ctrl := gomock.NewController(t)
	client := mock.NewMockConsulClient(ctrl)
	factory := mock.NewMockSessionFactory(ctrl)
	kv := mock.NewMockConsulKV(ctrl)

	comp := NewConsulLock(log.DefaultLogger)
	cfg := lock.Metadata{
		Properties: make(map[string]string),
	}
	cfg.Properties["address"] = "127.0.0.1:8500"
	err := comp.Init(cfg)
	comp.client = client
	comp.sessionFactory = factory
	comp.kv = kv
	factory.EXPECT().Create(&api.SessionEntry{TTL: getTTL(expireTime), LockDelay: 0, Behavior: "delete"}, nil).
		Return("session1", nil, nil).Times(1)
	factory.EXPECT().Create(&api.SessionEntry{TTL: getTTL(expireTime), LockDelay: 0, Behavior: "delete"}, nil).
		Return("session2", nil, nil).Times(1)
	kv.EXPECT().Acquire(&api.KVPair{Key: resouseId, Value: []byte(lockOwerA), Session: "session1"}, nil).
		Return(true, nil, nil).Times(1)
	kv.EXPECT().Acquire(&api.KVPair{Key: resouseId, Value: []byte(lockOwerB), Session: "session2"}, nil).
		Return(false, nil, nil).Times(1)

	tryLock, _ := comp.TryLock(&lock.TryLockRequest{
		ResourceId: resouseId,
		LockOwner:  lockOwerA,
		Expire:     expireTime,
	})

	assert.NoError(t, err)
	assert.Equal(t, true, tryLock.Success)

	bLock, _ := comp.TryLock(&lock.TryLockRequest{
		ResourceId: resouseId,
		LockOwner:  lockOwerB,
		Expire:     expireTime,
	})

	assert.NoError(t, err)
	assert.Equal(t, false, bLock.Success)

}

// A lock B unlock A unlock
func TestConsulLock_ALock_BUnlock(t *testing.T) {
	//mock
	ctrl := gomock.NewController(t)
	client := mock.NewMockConsulClient(ctrl)
	factory := mock.NewMockSessionFactory(ctrl)
	kv := mock.NewMockConsulKV(ctrl)

	comp := NewConsulLock(log.DefaultLogger)
	cfg := lock.Metadata{
		Properties: make(map[string]string),
	}
	cfg.Properties["address"] = "127.0.0.1:8500"
	err := comp.Init(cfg)
	comp.client = client
	comp.sessionFactory = factory
	comp.kv = kv
	factory.EXPECT().Create(&api.SessionEntry{TTL: getTTL(expireTime), LockDelay: 0, Behavior: "delete"}, nil).
		Return("session1", nil, nil).Times(1)
	factory.EXPECT().Destroy("session1", nil).Return(nil, nil).Times(1)
	kv.EXPECT().Acquire(&api.KVPair{Key: resouseId, Value: []byte(lockOwerA), Session: "session1"}, nil).
		Return(true, nil, nil).Times(1)
	kv.EXPECT().Release(&api.KVPair{Key: resouseId, Value: []byte(lockOwerA), Session: "session1"}, nil).
		Return(true, nil, nil).Times(1)

	tryLock, _ := comp.TryLock(&lock.TryLockRequest{
		ResourceId: resouseId,
		LockOwner:  lockOwerA,
		Expire:     expireTime,
	})

	assert.NoError(t, err)
	assert.Equal(t, true, tryLock.Success)

	unlock, _ := comp.Unlock(&lock.UnlockRequest{
		ResourceId: resouseId,
		LockOwner:  lockOwerB,
	})

	assert.NoError(t, err)
	assert.Equal(t, lock.LOCK_UNEXIST, unlock.Status)

	unlock2, err := comp.Unlock(&lock.UnlockRequest{
		ResourceId: resouseId,
		LockOwner:  lockOwerA,
	})

	assert.NoError(t, err)
	assert.Equal(t, lock.SUCCESS, unlock2.Status)

	// not implement LockKeepAlive
	keepAliveResp, err := comp.LockKeepAlive(context.TODO(), &lock.LockKeepAliveRequest{})
	assert.Nil(t, keepAliveResp)
	assert.Nil(t, err)
}
