package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
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

//A lock A unlock
func TestConsulLock_TryLock(t *testing.T) {

	comp := NewConsulLock(log.DefaultLogger)
	cfg := lock.Metadata{
		Properties: make(map[string]string),
	}
	cfg.Properties["address"] = "127.0.0.1:8500"
	err := comp.Init(cfg)

	assert.NoError(t, err)
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
func TestConsulLock_ALock_BUnlock(t *testing.T) {

	comp := NewConsulLock(log.DefaultLogger)
	cfg := lock.Metadata{
		Properties: make(map[string]string),
	}
	cfg.Properties["address"] = "127.0.0.1:8500"
	err := comp.Init(cfg)

	assert.NoError(t, err)
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
}

//func Test_RedisComplete(t *testing.T) {
//
//	//mock lock competition
//	for i := 0; i < 5; i++ {
//		go redisOption(i)
//	}
//
//	time.Sleep(time.Second * 2000)
//}
//
//func redisOption(number int) {
//	comp := NewConsulLock(log.DefaultLogger)
//	cfg := lock.Metadata{
//		Properties: make(map[string]string),
//	}
//	cfg.Properties["address"] = "127.0.0.1:8500"
//	comp.Init(cfg)
//	tryLock := &lock.TryLockResponse{}
//	lockOwner := "P_" + fmt.Sprint(number)
//	//loop to get lock
//	for !tryLock.Success {
//		tryLock, _ = comp.TryLock(&lock.TryLockRequest{
//			ResourceId: resouseId,
//			LockOwner:  lockOwner,
//			Expire:     10,
//		})
//		time.Sleep(time.Second * 2)
//	}
//
//	//in critical section----
//	fmt.Println(number, "get the lock")
//	//redis option
//	opts := &redis.Options{
//		Addr: "127.0.0.1:6379",
//	}
//	client := redis.NewClient(opts)
//	//get
//	get := client.Get(context.Background(), "r1")
//	i, _ := get.Int()
//	//sleep
//	time.Sleep(time.Second * 2)
//	//set
//	client.Set(context.Background(), "r1", i+1, 0)
//
//	//out critical section----
//	//unlock
//	unlock, _ := comp.Unlock(&lock.UnlockRequest{
//		ResourceId: resouseId,
//		LockOwner:  lockOwner,
//	})
//
//	fmt.Println(number, "release the lock", unlock)
//
//}
//
func TestConsul(t *testing.T) {
	// Get a new client
	client, err := api.NewClient(api.DefaultConfig())

	session, _, _ := client.Session().Create(&api.SessionEntry{
		TTL:       getTTL(5),
		LockDelay: 0,
		Name:      lockOwerA,
	}, nil)
	info, _, _ := client.Session().Info(session, nil)
	fmt.Println(info.Name)

	// Get a handle to the KV API
	kv := client.KV()
	// PUT a new KV pair
	p := &api.KVPair{Key: "REDIS_MAXCLIENTS", Value: []byte("1000"), Session: session}

	kv.Acquire(p, nil)

	time.Sleep(time.Second * 10)
	if err != nil {
		panic(err)
	}

	session2, _, _ := client.Session().Create(&api.SessionEntry{
		TTL:       getTTL(5),
		LockDelay: 0,
	}, nil)
	p2 := &api.KVPair{Key: "REDIS_MAXCLIENTS", Value: []byte("1000"), Session: session2}
	acquire, _, err := kv.Acquire(p2, nil)
	fmt.Println(acquire)
	// Lookup the pair
	pair, _, err := kv.Release(p, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(pair)
	fmt.Println(err)
	//fmt.Printf("KV: %v %s %s %v\n", pair.Key, pair.Value,pair.Session,pair.LockIndex)

}

//
//func TestConsulLock(t *testing.T) {
//
//	// Get a new client
//	client, err := api.NewClient(api.DefaultConfig())
//	if err != nil {
//		panic(err)
//	}
//
//	// Get a handle to the KV API
//	wg := sync.WaitGroup{}
//	wg.Add(1)
//	lock, _ := client.LockKey(resouseId)
//
//	ch, err := lock.Lock(nil)
//	fmt.Println("A 获取到锁", ch)
//	stopCh := make(chan struct{}, 1)
//	go func() {
//		client, _ := api.NewClient(api.DefaultConfig())
//		lock, _ := client.LockKey(resouseId)
//
//		fmt.Println("B 开始获取锁")
//		_, err2 := lock.Lock(stopCh)
//
//		//select {
//		//case <-ch: {
//		//   fmt.Println("B 结束了")
//		//}
//		//}
//
//		fmt.Println("1")
//		fmt.Println(err2)
//		select {}
//		wg.Done()
//	}()
//
//	go func() {
//		fmt.Println("thread B")
//		time.Sleep(time.Second * 5)
//		stopCh <- struct{}{}
//	}()
//
//	time.Sleep(time.Second * 5)
//
//	err = lock.Unlock()
//	fmt.Println("A 释放锁")
//
//	wg.Wait()
//}
//
//func TestLock(t *testing.T) {
//
//	for i := 0; i < 10; i++ {
//		go A(i)
//	}
//
//	time.Sleep(time.Second * 100)
//}
//
//func A(name int) {
//	// Get a new client
//	client, err := api.NewClient(api.DefaultConfig())
//	if err != nil {
//		panic(err)
//	}
//
//	lock, _ := client.LockKey(resouseId)
//
//	ch, err := lock.Lock(nil)
//	fmt.Println(name, "获取到锁", ch)
//
//	lock.Unlock()
//	fmt.Println(name, "释放锁", ch)
//
//}
//
//func TestSession(t *testing.T) {
//	client, err := api.NewClient(api.DefaultConfig())
//	if err != nil {
//		panic(err)
//	}
//	sessionId, _, err := client.Session().Create(&api.SessionEntry{
//		TTL: "10s",
//	}, nil)
//	if err != nil {
//		panic(err)
//	}
//	kv := client.KV()
//	p := &api.KVPair{Key: "resoure_1", Value: []byte("A"), Session: sessionId}
//
//	acquire, _, err := kv.Acquire(p, nil)
//	fmt.Println(acquire)
//	q, err := kv.Put(p, nil)
//	fmt.Println(q)
//
//	/*	acquire, _, err := kv.Release(p, nil)
//		if err != nil {
//			panic(err)
//		}
//		fmt.Println(acquire)*/
//	sessionB, _, err := client.Session().Create(&api.SessionEntry{
//		TTL: "10s",
//	}, nil)
//	pb := &api.KVPair{Key: "REDIS_MAXCLIENTS", Value: []byte("B"), Session: sessionB}
//	b, err := kv.Put(pb, nil)
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println(b)
//}
//
//func TestTxn(t *testing.T) {
//	client, err := api.NewClient(api.DefaultConfig())
//	if err != nil {
//		panic(err)
//	}
//	/*	sessionId, _, err := client.Session().Create(&api.SessionEntry{
//		TTL: "10s",
//	}, nil)*/
//	sessionB, _, err := client.Session().Create(&api.SessionEntry{
//		TTL: "10s",
//	}, nil)
//	txn, response, _, err := client.Txn().Txn(api.TxnOps{
//		/*		&api.TxnOp{
//				KV: &api.KVTxnOp{
//					Verb:    api.KVLock,
//					Key:     "REDIS_MAXCLIENTS",
//					Value:   []byte("???"),
//					Session: sessionId,
//				},
//			},*/
//		&api.TxnOp{
//			KV: &api.KVTxnOp{
//				Verb:    api.KVSet,
//				Key:     "REDIS_MAXCLIENTS",
//				Value:   []byte("BVVVVV"),
//				Session: sessionB,
//			},
//		},
//		&api.TxnOp{
//			KV: &api.KVTxnOp{
//				Verb:    api.KVLock,
//				Key:     "REDIS_MAXCLIENTS",
//				Value:   []byte("???"),
//				Session: sessionB,
//			},
//		},
//		&api.TxnOp{
//			KV: &api.KVTxnOp{
//				Verb:    api.KVGet,
//				Key:     "REDIS_MAXCLIENTS",
//				Value:   []byte("BVVVVV"),
//				Session: sessionB,
//			},
//		},
//
//		/*		&api.TxnOp{
//				KV: &api.KVTxnOp{
//					Verb:    api.KVUnlock,
//					Key:     "REDIS_MAXCLIENTS",
//					Value:   []byte("B"),
//					Session: sessionId,
//				},
//			},*/
//
//	}, nil)
//	fmt.Println(txn, response, err)
//	txn1, response1, _, err := client.Txn().Txn(api.TxnOps{
//
//		&api.TxnOp{
//			KV: &api.KVTxnOp{
//				Verb:    api.KVGet,
//				Key:     "REDIS_MAXCLIENTS",
//				Session: sessionB,
//			},
//		},
//		/*		&api.TxnOp{
//				KV: &api.KVTxnOp{
//					Verb:    api.KVDeleteCAS,
//					Key:     "REDIS_MAXCLIENTS",
//					Value:   []byte("???"),
//					Index: response.Results[0].KV.ModifyIndex,
//					Session: sessionB,
//				},
//			},*/
//		/*		&api.TxnOp{
//				KV: &api.KVTxnOp{
//					Verb:    api.KVUnlock,
//					Key:     "REDIS_MAXCLIENTS",
//					Value:   []byte("B"),
//					Session: sessionB,
//				},
//			},*/
//	}, nil)
//	fmt.Println(txn1, response1, err)
//
//}
