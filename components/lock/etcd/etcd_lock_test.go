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
package etcd

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"sync"
	"testing"
	"time"

	"mosn.io/pkg/log"

	"mosn.io/layotto/components/lock"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.etcd.io/etcd/server/v3/embed"
)

const resourceId = "resource_xxx"
const resourceId2 = "resource_xxx2"

const resourceId3 = "resource_xxx3"
const resourceId4 = "resource_xxx4"

func TestEtcdLock_Init(t *testing.T) {
	var err error
	var etcdServer *embed.Etcd
	var etcdTestDir = "init.test.etcd"
	var etcdUrl = "localhost:2379"

	etcdServer, err = startEtcdServer(etcdTestDir, 2379)
	assert.NoError(t, err)
	defer func() {
		etcdServer.Server.Stop()
		os.RemoveAll(etcdTestDir)
	}()
	comp := NewEtcdLock(log.DefaultLogger)

	cfg := lock.Metadata{
		Properties: make(map[string]string),
	}
	err = comp.Init(cfg)
	assert.Error(t, err)

	cfg.Properties["endpoints"] = ""
	err = comp.Init(cfg)
	assert.Error(t, err)

	cfg.Properties["endpoints"] = etcdUrl
	cfg.Properties["dialTimeout"] = "a"
	err = comp.Init(cfg)
	assert.Error(t, err)

	cfg.Properties["dialTimeout"] = "2"
	err = comp.Init(cfg)
	assert.NoError(t, err)
	err = comp.Close()
	assert.NoError(t, err)

	//ca
	cfg.Properties["tlsCa"] = "/tmp"
	cfg.Properties["tlsCert"] = "/tmp"
	cfg.Properties["tlsCertKey"] = "/tmp"
	err = comp.Init(cfg)
	assert.Error(t, err)
}

func TestEtcdLock_CreateConnTimeout(t *testing.T) {
	var err error

	comp := NewEtcdLock(log.DefaultLogger)

	cfg := lock.Metadata{
		Properties: make(map[string]string),
	}

	cfg.Properties["endpoints"] = "localhost:18888"
	cfg.Properties["dialTimeout"] = "3"
	startTime := time.Now()
	err = comp.Init(cfg)
	endTIme := time.Now()
	d := endTIme.Sub(startTime)
	assert.Error(t, err)
	assert.Equal(t, true, d.Seconds() >= 3)
}

func TestEtcdLock_TryLock(t *testing.T) {
	var err error
	var resp *lock.TryLockResponse
	var etcdServer *embed.Etcd
	var etcdTestDir = "trylock.test.etcd"
	var etcdUrl = "localhost:23780"

	etcdServer, err = startEtcdServer(etcdTestDir, 23780)
	assert.NoError(t, err)
	defer func() {
		etcdServer.Server.Stop()
		os.RemoveAll(etcdTestDir)
	}()

	comp := NewEtcdLock(log.DefaultLogger)

	cfg := lock.Metadata{
		Properties: make(map[string]string),
	}

	cfg.Properties["endpoints"] = etcdUrl
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

	//repeat
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
		//another owner
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

func TestEtcdLock_UnLock(t *testing.T) {
	var err error
	var resp *lock.UnlockResponse
	var lockresp *lock.TryLockResponse
	var etcdServer *embed.Etcd
	var etcdTestDir = "trylock.test.etcd"
	var etcdUrl = "localhost:23790"

	etcdServer, err = startEtcdServer(etcdTestDir, 23790)
	assert.NoError(t, err)
	defer func() {
		etcdServer.Server.Stop()
		os.RemoveAll(etcdTestDir)
	}()

	comp := NewEtcdLock(log.DefaultLogger)

	cfg := lock.Metadata{
		Properties: make(map[string]string),
	}

	cfg.Properties["endpoints"] = etcdUrl
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

	// not implement LockKeepAlive
	keepAliveResp, err := comp.LockKeepAlive(context.TODO(), &lock.LockKeepAliveRequest{})
	assert.Nil(t, keepAliveResp)
	assert.Nil(t, err)
}

func startEtcdServer(dir string, port int) (*embed.Etcd, error) {
	lc, _ := url.Parse(fmt.Sprintf("http://localhost:%v", port))
	lp, _ := url.Parse(fmt.Sprintf("http://localhost:%v", port+1))

	cfg := embed.NewConfig()
	cfg.Dir = dir
	cfg.LogLevel = "error"
	cfg.LCUrls = []url.URL{*lc}
	cfg.LPUrls = []url.URL{*lp}
	e, err := embed.StartEtcd(cfg)
	if err != nil {
		return nil, err
	}
	<-e.Server.ReadyNotify()
	return e, nil
}
