package mongo

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"mosn.io/layotto/components/lock"
	"mosn.io/pkg/log"
	"sync"
	"testing"
)

const (
	resourceId  = "resource_xxx"
	resourceId2 = "resource_xxx2"
	resourceId3 = "resource_xxx3"
	resourceId4 = "resource_xxx4"
	lockOwner   = "lockOwner1"
	lockOwner2  = "lockOwner2"
	lockOwner3  = "lockOwner3"
)

func TestMongoLock_Init(t *testing.T) {
	var err error
	var mongoUrl = "localhost:27017"
	comp := NewMongoLock(log.DefaultLogger)

	cfg := lock.Metadata{
		Properties: make(map[string]string),
	}

	err = comp.Init(cfg)
	assert.Error(t, err)

	cfg.Properties["mongoHost"] = mongoUrl
	cfg.Properties["operationTimeout"] = "a"
	err = comp.Init(cfg)
	assert.Error(t, err)

	cfg.Properties["operationTimeout"] = "2"
	err = comp.Init(cfg)
	assert.Error(t, err)
}

func TestMongoLock_TryLock(t *testing.T) {
	var err error
	var resp *lock.TryLockResponse
	var mongoUrl = "localhost:27017"

	comp := NewMongoLock(log.DefaultLogger)

	cfg := lock.Metadata{
		Properties: make(map[string]string),
	}

	cfg.Properties["mongoHost"] = mongoUrl
	err = comp.Init(cfg)
	assert.NoError(t, err)

	ownerId1 := uuid.New().String()
	resp, err = comp.TryLock(&lock.TryLockRequest{
		ResourceId: resourceId,
		LockOwner:  ownerId1,
		Expire:     10,
	})
	assert.NoError(t, err)
	assert.Equal(t, true, resp.Success)

	resp, err = comp.TryLock(&lock.TryLockRequest{
		ResourceId: resourceId,
		LockOwner:  ownerId1,
		Expire:     10,
	})
	assert.NoError(t, err)
	assert.Equal(t, false, resp.Success)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		ownerId2 := uuid.New().String()
		resp, err = comp.TryLock(&lock.TryLockRequest{
			ResourceId: resourceId,
			LockOwner:  ownerId2,
			Expire:     10,
		})
		assert.NoError(t, err)
		assert.Equal(t, false, resp.Success)
		wg.Done()
	}()

	wg.Wait()

	//another resource
	resp, err = comp.TryLock(&lock.TryLockRequest{
		ResourceId: resourceId2,
		LockOwner:  ownerId1,
		Expire:     10,
	})
	assert.NoError(t, err)
	assert.Equal(t, true, resp.Success)
}

func TestMongoLock_Unlock(t *testing.T) {
	var err error
	var resp *lock.UnlockResponse
	var lockresp *lock.TryLockResponse
	var mongoUrl = "localhost:27017"

	comp := NewMongoLock(log.DefaultLogger)

	cfg := lock.Metadata{
		Properties: make(map[string]string),
	}

	cfg.Properties["mongoHost"] = mongoUrl
	err = comp.Init(cfg)
	assert.NoError(t, err)

	ownerId1 := uuid.New().String()
	lockresp, err = comp.TryLock(&lock.TryLockRequest{
		ResourceId: resourceId3,
		LockOwner:  ownerId1,
		Expire:     10,
	})
	assert.NoError(t, err)
	assert.Equal(t, true, lockresp.Success)

	//error ownerid
	resp, err = comp.Unlock(&lock.UnlockRequest{
		ResourceId: resourceId3,
		LockOwner:  uuid.New().String(),
	})
	assert.NoError(t, err)
	assert.Equal(t, lock.LOCK_BELONG_TO_OTHERS, resp.Status)

	//error resourceid
	resp, err = comp.Unlock(&lock.UnlockRequest{
		ResourceId: resourceId4,
		LockOwner:  ownerId1,
	})
	assert.NoError(t, err)
	assert.Equal(t, lock.LOCK_UNEXIST, resp.Status)

	//success
	resp, err = comp.Unlock(&lock.UnlockRequest{
		ResourceId: resourceId3,
		LockOwner:  ownerId1,
	})
	assert.NoError(t, err)
	assert.Equal(t, lock.SUCCESS, resp.Status)
}
