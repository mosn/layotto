package etcd

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"mosn.io/layotto/components/lock"
	"mosn.io/pkg/log"

	"go.etcd.io/etcd/client/v3"
)

const (
	defaultDialTimeout = 5
	defaultKeyPrefix   = "/layotto/"

	prefixKey      = "keyPrefix"
	usernameKey    = "username"
	passwordKey    = "password"
	dialTimeoutKey = "dialTimeout"
	endpointsKey   = "endpoints"
)

type EtcdLock struct {
	client   *clientv3.Client
	metadata metadata

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

func (e *EtcdLock) Init(metadata lock.Metadata) error {
	// 1. parse config
	m, err := parseEtcdMetadata(metadata)
	if err != nil {
		return err
	}
	e.metadata = m
	// 2. construct client
	if e.client, err = e.newClient(m); err != nil {
		return err
	}

	e.ctx, e.cancel = context.WithCancel(context.Background())

	return err
}

func (e *EtcdLock) Features() []lock.Feature {
	return e.features
}

func (e *EtcdLock) TryLock(req *lock.TryLockRequest) (*lock.TryLockResponse, error) {
	var leaseId clientv3.LeaseID
	//1.Create new lease
	lease := clientv3.NewLease(e.client)
	if leaseGrantResp, err := lease.Grant(e.ctx, int64(req.Expire)); err != nil {
		return &lock.TryLockResponse{}, fmt.Errorf("[etcdLock]: Create new lease returned error: %s.ResourceId: %s", err, req.ResourceId)
	} else {
		leaseId = leaseGrantResp.ID
	}

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

func (e *EtcdLock) Unlock(req *lock.UnlockRequest) (*lock.UnlockResponse, error) {
	key := e.getKey(req.ResourceId)

	kv := clientv3.NewKV(e.client)
	txn := kv.Txn(e.ctx)
	txn.If(clientv3.Compare(clientv3.Value(key), "=", req.LockOwner)).Then(
		clientv3.OpDelete(key)).Else(
		clientv3.OpGet(key))
	txnResponse, err := txn.Commit()
	if err != nil {
		return newInternalErrorUnlockResponse(), fmt.Errorf("[etcdLock]: Unlock returned error: %s.ResourceId: %s", err, req.ResourceId)
	}

	if txnResponse.Succeeded {
		return &lock.UnlockResponse{Status: lock.SUCCESS}, nil
	} else {
		resp := txnResponse.Responses[0].GetResponseRange()
		if len(resp.Kvs) == 0 {
			return &lock.UnlockResponse{Status: lock.LOCK_UNEXIST}, nil
		}

		return &lock.UnlockResponse{Status: lock.LOCK_BELONG_TO_OTHERS}, nil
	}
}

func (e *EtcdLock) Close() error {
	e.cancel()

	return e.client.Close()
}

func (e *EtcdLock) newClient(meta metadata) (*clientv3.Client, error) {

	config := clientv3.Config{
		Endpoints:   meta.endpoints,
		DialTimeout: time.Duration(meta.dialTimeout) * time.Second,
		Username:    meta.username,
		Password:    meta.password,
	}

	if client, err := clientv3.New(config); err != nil {
		return nil, err
	} else {
		return client, nil
	}
}

func (e *EtcdLock) getKey(resourceId string) string {
	return fmt.Sprintf("%s%s", e.metadata.keyPrefix, resourceId)
}

func newInternalErrorUnlockResponse() *lock.UnlockResponse {
	return &lock.UnlockResponse{
		Status: lock.INTERNAL_ERROR,
	}
}

func parseEtcdMetadata(meta lock.Metadata) (metadata, error) {
	m := metadata{}
	var err error

	if val, ok := meta.Properties[endpointsKey]; ok && val != "" {
		m.endpoints = strings.Split(val, ";")
	} else {
		return m, errors.New("etcd lock error: missing endpoints address")
	}

	if val, ok := meta.Properties[dialTimeoutKey]; ok && val != "" {
		if m.dialTimeout, err = strconv.Atoi(val); err != nil {
			return m, fmt.Errorf("etcd lock error: ncorrect dialTimeout value %s", val)
		}
	} else {
		m.dialTimeout = defaultDialTimeout
	}

	if val, ok := meta.Properties[prefixKey]; ok && val != "" {
		m.keyPrefix = val
	} else {
		m.keyPrefix = defaultKeyPrefix
	}

	if val, ok := meta.Properties[usernameKey]; ok && val != "" {
		m.username = val
	}

	if val, ok := meta.Properties[passwordKey]; ok && val != "" {
		m.password = val
	}

	return m, nil
}

type metadata struct {
	keyPrefix   string
	dialTimeout int
	endpoints   []string
	username    string
	password    string
}
