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

	clientv3 "go.etcd.io/etcd/client/v3"

	"mosn.io/layotto/components/pkg/utils"

	"mosn.io/pkg/log"

	"mosn.io/layotto/components/lock"
)

// Etcd lock store
type EtcdLock struct {
	client   *clientv3.Client
	metadata utils.EtcdMetadata

	features []lock.Feature
	logger   log.ErrorLogger

	ctx    context.Context
	cancel context.CancelFunc
}

// NewEtcdLock returns a new etcd lock
func NewEtcdLock(logger log.ErrorLogger) *EtcdLock {
	s := &EtcdLock{
		features: make([]lock.Feature, 0),
		logger:   logger,
	}

	return s
}

// Init EtcdLock
func (e *EtcdLock) Init(metadata lock.Metadata) error {
	// 1. parse config
	m, err := utils.ParseEtcdMetadata(metadata.Properties)
	if err != nil {
		return err
	}
	e.metadata = m
	// 2. construct client
	if e.client, err = utils.NewEtcdClient(m); err != nil {
		return err
	}

	e.ctx, e.cancel = context.WithCancel(context.Background())

	return err
}

// LockKeepAlive try to renewal lease
func (e *EtcdLock) LockKeepAlive(ctx context.Context, request *lock.LockKeepAliveRequest) (*lock.LockKeepAliveResponse, error) {
	//TODO: implemnt function
	return nil, nil
}

// Features is to get EtcdLock's features
func (e *EtcdLock) Features() []lock.Feature {
	return e.features
}

// Node tries to acquire a etcd lock
func (e *EtcdLock) TryLock(req *lock.TryLockRequest) (*lock.TryLockResponse, error) {
	var leaseId clientv3.LeaseID
	//1.Create new lease
	lease := clientv3.NewLease(e.client)
	leaseGrantResp, err := lease.Grant(e.ctx, int64(req.Expire))
	if err != nil {
		return &lock.TryLockResponse{}, fmt.Errorf("[etcdLock]: Create new lease returned error: %s.ResourceId: %s", err, req.ResourceId)
	}
	leaseId = leaseGrantResp.ID

	key := e.getKey(req.ResourceId)

	//2.Create new KV
	kv := clientv3.NewKV(e.client)
	//3.Create txn
	txn := kv.Txn(e.ctx)
	txn.If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0)).Then(
		clientv3.OpPut(key, req.LockOwner, clientv3.WithLease(leaseId))).Else(
		clientv3.OpGet(key))
	//4.Commit and try get lock
	txnResponse, err := txn.Commit()
	if err != nil {
		return &lock.TryLockResponse{}, fmt.Errorf("[etcdLock]: Creat lock returned error: %s.ResourceId: %s", err, req.ResourceId)
	}

	return &lock.TryLockResponse{
		Success: txnResponse.Succeeded,
	}, nil
}

// Node tries to release a etcd lock
func (e *EtcdLock) Unlock(req *lock.UnlockRequest) (*lock.UnlockResponse, error) {
	key := e.getKey(req.ResourceId)

	// 1.Create new KV
	kv := clientv3.NewKV(e.client)
	// 2.Create txn
	txn := kv.Txn(e.ctx)
	txn.If(clientv3.Compare(clientv3.Value(key), "=", req.LockOwner)).Then(
		clientv3.OpDelete(key)).Else(
		clientv3.OpGet(key))
	// 3.Commit and try release lock
	txnResponse, err := txn.Commit()
	if err != nil {
		return newInternalErrorUnlockResponse(), fmt.Errorf("[etcdLock]: Unlock returned error: %s.ResourceId: %s", err, req.ResourceId)
	}

	if txnResponse.Succeeded {
		return &lock.UnlockResponse{Status: lock.SUCCESS}, nil
	}
	resp := txnResponse.Responses[0].GetResponseRange()
	if len(resp.Kvs) == 0 {
		return &lock.UnlockResponse{Status: lock.LOCK_UNEXIST}, nil
	}

	return &lock.UnlockResponse{Status: lock.LOCK_BELONG_TO_OTHERS}, nil
}

// Close shuts down the client's etcd connections.
func (e *EtcdLock) Close() error {
	e.cancel()

	return e.client.Close()
}

// getkey is to return string of type KeyPrefix + resourceId
func (e *EtcdLock) getKey(resourceId string) string {
	return fmt.Sprintf("%s%s", e.metadata.KeyPrefix, resourceId)
}

// newInternalErrorUnlockResponse is to return lock release error
func newInternalErrorUnlockResponse() *lock.UnlockResponse {
	return &lock.UnlockResponse{
		Status: lock.INTERNAL_ERROR,
	}
}
