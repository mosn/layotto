package zookeeper

import (
	"errors"
	"github.com/go-zookeeper/zk"
	"github.com/golang/mock/gomock"
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
func TestZookeeperLock_ALock_AUnlock(t *testing.T) {

	comp := NewZookeeperLock(log.DefaultLogger)
	comp.Init(cfg)

	//mock
	ctrl := gomock.NewController(t)
	connection := NewMockZKConnection(ctrl)
	path := "/" + resouseId
	connection.EXPECT().NewConnection(time.Duration(expireTime), comp.metadata).Return(connection, nil).Times(2)
	connection.EXPECT().Create(path, []byte(lockOwerA), int32(zk.FlagEphemeral), zk.WorldACL(zk.PermAll)).Return("", nil).Times(1)
	connection.EXPECT().Close().Return().Times(1)
	connection.EXPECT().Get(path).Return([]byte(lockOwerA), &zk.Stat{Version: 123}, nil).Times(1)
	connection.EXPECT().Delete(path, int32(123)).Return(nil).Times(1)
	comp.conn = connection

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
func TestZookeeperLock_ALock_BUnlock(t *testing.T) {

	comp := NewZookeeperLock(log.DefaultLogger)
	comp.Init(cfg)

	//mock
	ctrl := gomock.NewController(t)
	connection := NewMockZKConnection(ctrl)
	path := "/" + resouseId
	connection.EXPECT().NewConnection(time.Duration(expireTime), comp.metadata).Return(connection, nil).Times(2)
	connection.EXPECT().Create(path, []byte(lockOwerA), int32(zk.FlagEphemeral), zk.WorldACL(zk.PermAll)).Return("", nil).Times(1)
	connection.EXPECT().Close().Return().Times(1)
	connection.EXPECT().Get(path).Return([]byte(lockOwerA), &zk.Stat{Version: 123}, nil).Times(1)
	connection.EXPECT().Delete(path, int32(123)).Return(nil).Times(1)
	comp.conn = connection

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

// A lock , B lock ,A unlock ,B lock,B unlock
func TestZookeeperLock_ALock_BLock_AUnlock_BLock_BUnlock(t *testing.T) {

	comp := NewZookeeperLock(log.DefaultLogger)
	comp.Init(cfg)

	//mock
	ctrl := gomock.NewController(t)
	connection := NewMockZKConnection(ctrl)
	path := "/" + resouseId
	connection.EXPECT().NewConnection(time.Duration(expireTime), comp.metadata).Return(connection, nil).Times(5)
	connection.EXPECT().Create(path, []byte(lockOwerA), int32(zk.FlagEphemeral), zk.WorldACL(zk.PermAll)).Return("", nil).Times(1)
	connection.EXPECT().Create(path, []byte(lockOwerB), int32(zk.FlagEphemeral), zk.WorldACL(zk.PermAll)).Return("", errors.New("")).Times(1)
	connection.EXPECT().Create(path, []byte(lockOwerB), int32(zk.FlagEphemeral), zk.WorldACL(zk.PermAll)).Return("", nil).Times(1)

	connection.EXPECT().Close().Return().Times(5)
	connection.EXPECT().Get(path).Return([]byte(lockOwerA), &zk.Stat{Version: 123}, nil).Times(1)
	connection.EXPECT().Get(path).Return([]byte(lockOwerB), &zk.Stat{Version: 124}, nil).Times(1)
	connection.EXPECT().Delete(path, int32(123)).Return(nil).Times(2)
	connection.EXPECT().Delete(path, int32(124)).Return(nil).Times(2)

	comp.conn = connection

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

//test redis compter
/*func Test_RedisComplete(t *testing.T) {

	//mock lock competition
	for i := 0; i < 30; i++ {
		go redisOperate(i)
	}

	time.Sleep(time.Second * 2000)
}

func redisOperate(number int) {
	comp := NewZookeeperLock(log.DefaultLogger)
	comp.Init(cfg)
	tryLock := &lock.TryLockResponse{}

	lockOwner := "P_" + fmt.Sprint(number)
	//loop to get lock
	for !tryLock.Success {
		time.Sleep(time.Second * 2)
		tryLock, _ = comp.TryLock(&lock.TryLockRequest{
			ResourceId: resouseId,
			LockOwner:  lockOwner,
			Expire:     10,
		})
	}

	//in critical section----
	fmt.Println(number, "get the lock")
	//redis option
	opts := &redis.Options{
		Addr: "127.0.0.1:6379",
	}
	client := redis.NewClient(opts)
	//get
	get := client.Get(context.Background(), "r1")
	i, _ := get.Int()
	//sleep
	time.Sleep(time.Second * 2)
	//set
	client.Set(context.Background(), "r1", i+1, 0)

	//out critical section----
	//unlock
	unlock, _ := comp.Unlock(&lock.UnlockRequest{
		ResourceId: resouseId,
		LockOwner:  lockOwner,
	})

	fmt.Println(number, "release the lock", unlock)

}
*/
