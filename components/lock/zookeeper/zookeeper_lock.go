//
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

package zookeeper

import (
	"context"
	"time"

	"github.com/go-zookeeper/zk"
	"mosn.io/pkg/log"
	util "mosn.io/pkg/utils"

	"mosn.io/layotto/components/lock"
	"mosn.io/layotto/components/pkg/utils"
)

var closeConn = func(conn utils.ZKConnection, expireInSecond int32) {
	//can also
	//time.Sleep(time.Second * time.Duration(expireInSecond))
	<-time.After(time.Second * time.Duration(expireInSecond))
	// make sure close connecion
	conn.Close()
}

// ZookeeperLock lock store
type ZookeeperLock struct {
	//trylock reestablish connection  every time
	factory utils.ConnectionFactory
	//unlock reuse this conneciton
	unlockConn utils.ZKConnection
	metadata   utils.ZookeeperMetadata
	logger     log.ErrorLogger
}

// NewZookeeperLock Create ZookeeperLock
func NewZookeeperLock(logger log.ErrorLogger) *ZookeeperLock {
	lock := &ZookeeperLock{
		logger: logger,
	}
	return lock
}

// Init ZookeeperLock
func (p *ZookeeperLock) Init(metadata lock.Metadata) error {

	m, err := utils.ParseZookeeperMetadata(metadata.Properties)
	if err != nil {
		return err
	}

	p.metadata = m
	p.factory = &utils.ConnectionFactoryImpl{}

	//init unlock connection
	zkConn, err := p.factory.NewConnection(p.metadata.SessionTimeout, p.metadata)
	if err != nil {
		return err
	}
	p.unlockConn = zkConn
	return nil
}

// Features is to get ZookeeperLock's features
func (p *ZookeeperLock) Features() []lock.Feature {
	return nil
}

// LockKeepAlive try to renewal lease
func (p *ZookeeperLock) LockKeepAlive(ctx context.Context, request *lock.LockKeepAliveRequest) (*lock.LockKeepAliveResponse, error) {
	//TODO: implemnt function
	return nil, nil
}

// TryLock Node tries to acquire a zookeeper lock
func (p *ZookeeperLock) TryLock(req *lock.TryLockRequest) (*lock.TryLockResponse, error) {

	conn, err := p.factory.NewConnection(time.Duration(req.Expire)*time.Second, p.metadata)
	if err != nil {
		return &lock.TryLockResponse{}, err
	}
	//1.create zk ephemeral node
	_, err = conn.Create("/"+req.ResourceId, []byte(req.LockOwner), zk.FlagEphemeral, zk.WorldACL(zk.PermAll))

	//2.1 create node fail ,indicates lock fail
	if err != nil {
		defer conn.Close()
		//the node exists,lock fail
		if err == zk.ErrNodeExists {
			return &lock.TryLockResponse{
				Success: false,
			}, nil
		}
		//other err
		return nil, err
	}

	//2.2 create node success, asyn  to make sure zkclient alive for need time
	util.GoWithRecover(func() {
		closeConn(conn, req.Expire)
	}, nil)

	return &lock.TryLockResponse{
		Success: true,
	}, nil

}

// Unlock Node tries to release a zookeeper lock
func (p *ZookeeperLock) Unlock(req *lock.UnlockRequest) (*lock.UnlockResponse, error) {

	conn := p.unlockConn

	path := "/" + req.ResourceId
	owner, state, err := conn.Get(path)

	if err != nil {
		//node does not exist, indicates this lock has expired
		if err == zk.ErrNoNode {
			return &lock.UnlockResponse{Status: lock.LOCK_UNEXIST}, nil
		}
		//other err
		return nil, err
	}
	//node exist ,but owner not this, indicates this lock has occupied or wrong unlock
	if string(owner) != req.LockOwner {
		return &lock.UnlockResponse{Status: lock.LOCK_BELONG_TO_OTHERS}, nil
	}
	err = conn.Delete(path, state.Version)
	//owner is this, but delete fail
	if err != nil {
		// delete no node , indicates this lock has expired
		if err == zk.ErrNoNode {
			return &lock.UnlockResponse{Status: lock.LOCK_UNEXIST}, nil
			// delete version error , indicates this lock has occupied by others
		} else if err == zk.ErrBadVersion {
			return &lock.UnlockResponse{Status: lock.LOCK_BELONG_TO_OTHERS}, nil
			//other error
		} else {
			return nil, err
		}
	}
	//delete success, unlock success
	return &lock.UnlockResponse{Status: lock.SUCCESS}, nil
}
