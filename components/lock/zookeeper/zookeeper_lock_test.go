package zookeeper

import (
	"github.com/stretchr/testify/assert"
	"mosn.io/layotto/components/lock"
	"mosn.io/pkg/log"
	"testing"
	"time"
)

const resouseId = "resoure_1"
const lockOwerA = "p1"
const lockOwerB = "p2"
const expireTime = 5

var cfg = lock.Metadata{
	Properties: make(map[string]string),
}

func TestMain(m *testing.M) {
	cfg.Properties["zookeeperHosts"] = "127.0.0.1;127.0.0.1"
	cfg.Properties["zookeeperPassword"] = ""
	m.Run()
}

// A lock ,A unlock
func TestZookeeperLock_One(t *testing.T) {
	comp := NewZookeeperLock(log.DefaultLogger)

	comp.Init(cfg)

	tryLock, err := comp.TryLock(&lock.TryLockRequest{
		ResourceId: resouseId,
		LockOwner:  lockOwerA,
		Expire:     expireTime,
	})
	assert.NoError(t, err)
	assert.Equal(t, tryLock.Success, true)
	unlock, _ := comp.Unlock(&lock.UnlockRequest{
		ResourceId: resouseId,
		LockOwner:  lockOwerA,
	})
	assert.NoError(t, err)
	assert.Equal(t, unlock.Status, lock.SUCCESS)

}

// A lock ,B unlock
func TestZookeeperLock_Two(t *testing.T) {
	comp := NewZookeeperLock(log.DefaultLogger)

	comp.Init(cfg)

	tryLock, err := comp.TryLock(&lock.TryLockRequest{
		ResourceId: resouseId,
		LockOwner:  lockOwerA,
		Expire:     expireTime,
	})
	assert.NoError(t, err)
	assert.Equal(t, tryLock.Success, true)
	unlock, err := comp.Unlock(&lock.UnlockRequest{
		ResourceId: resouseId,
		LockOwner:  lockOwerB,
	})
	assert.NoError(t, err)
	assert.Equal(t, unlock.Status, lock.LOCK_BELONG_TO_OTHERS)

}

// A lock , B lock ,A unlock ,A lock,B lock,B unlock
func TestZookeeperLock_Three(t *testing.T) {
	comp := NewZookeeperLock(log.DefaultLogger)

	comp.Init(cfg)
	//A lock
	tryLock, err := comp.TryLock(&lock.TryLockRequest{
		ResourceId: resouseId,
		LockOwner:  lockOwerA,
		Expire:     expireTime,
	})
	assert.NoError(t, err)
	assert.Equal(t, tryLock.Success, true)
	//B lock
	tryLock, err = comp.TryLock(&lock.TryLockRequest{
		ResourceId: resouseId,
		LockOwner:  lockOwerB,
		Expire:     expireTime,
	})
	assert.NoError(t, err)
	assert.Equal(t, tryLock.Success, false)
	//A unlock
	unlock, _ := comp.Unlock(&lock.UnlockRequest{
		ResourceId: resouseId,
		LockOwner:  lockOwerA,
	})
	assert.NoError(t, err)
	assert.Equal(t, unlock.Status, lock.SUCCESS)

	//B lock
	tryLock, err = comp.TryLock(&lock.TryLockRequest{
		ResourceId: resouseId,
		LockOwner:  lockOwerB,
		Expire:     expireTime,
	})
	assert.NoError(t, err)
	assert.Equal(t, tryLock.Success, true)

	//B unlock
	unlock, _ = comp.Unlock(&lock.UnlockRequest{
		ResourceId: resouseId,
		LockOwner:  lockOwerB,
	})
	assert.NoError(t, err)
	assert.Equal(t, unlock.Status, lock.SUCCESS)
}

//A lock
func TestZookeeperLock_Four(t *testing.T) {
	comp := NewZookeeperLock(log.DefaultLogger)
	comp.Init(cfg)
	go func() {
		tryLock, err := comp.TryLock(&lock.TryLockRequest{
			ResourceId: resouseId,
			LockOwner:  lockOwerA,
			Expire:     expireTime,
		})
		assert.NoError(t, err)
		assert.Equal(t, tryLock.Success, true)
	}()

	time.Sleep(time.Second * 9999)

}
