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
	"mosn.io/pkg/log"

	"mosn.io/layotto/components/pkg/utils"
	"mosn.io/layotto/components/sequencer"
)

type EtcdSequencer struct {
	client     *clientv3.Client
	metadata   utils.EtcdMetadata
	biggerThan map[string]int64

	logger log.ErrorLogger

	ctx    context.Context
	cancel context.CancelFunc
}

// EtcdSequencer returns a new etcd sequencer
func NewEtcdSequencer(logger log.ErrorLogger) *EtcdSequencer {
	s := &EtcdSequencer{
		logger: logger,
	}

	return s
}

func (e *EtcdSequencer) Init(config sequencer.Configuration) error {
	// 1. parse config
	m, err := utils.ParseEtcdMetadata(config.Properties)
	if err != nil {
		return err
	}
	e.metadata = m
	e.biggerThan = config.BiggerThan

	// 2. construct client
	if e.client, err = utils.NewEtcdClient(m); err != nil {
		return err
	}
	e.ctx, e.cancel = context.WithCancel(context.Background())

	// 3. check biggerThan
	if len(e.biggerThan) > 0 {
		kv := clientv3.NewKV(e.client)
		for k, bt := range e.biggerThan {
			if bt <= 0 {
				continue
			}
			actualKey := e.getKeyInEtcd(k)
			get, err := kv.Get(e.ctx, actualKey)
			if err != nil {
				return err
			}
			var cur int64 = 0
			if get.Count > 0 && len(get.Kvs) > 0 {
				cur = get.Kvs[0].Version
			}
			if cur < bt {
				return fmt.Errorf("etcd sequencer error: can not satisfy biggerThan guarantee.key: %s,key in etcd: %s,current id:%v", k, actualKey, cur)
			}
		}
	}
	// TODO close component?
	return nil
}

func (e *EtcdSequencer) GetNextId(req *sequencer.GetNextIdRequest) (*sequencer.GetNextIdResponse, error) {
	key := e.getKeyInEtcd(req.Key)
	// Create new KV
	kv := clientv3.NewKV(e.client)
	// Create txn
	txn := kv.Txn(e.ctx)
	txn.If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0)).Then(
		clientv3.OpPut(key, ""),
		clientv3.OpGet(key),
	).Else(
		clientv3.OpPut(key, ""),
		clientv3.OpGet(key),
	)
	// Commit
	txnResp, err := txn.Commit()
	if err != nil {
		return nil, err
	}
	return &sequencer.GetNextIdResponse{
		NextId: txnResp.Responses[1].GetResponseRange().Kvs[0].Version,
	}, nil
}

func (e *EtcdSequencer) GetSegment(req *sequencer.GetSegmentRequest) (support bool, result *sequencer.GetSegmentResponse, err error) {
	return false, nil, nil
}

func (e *EtcdSequencer) Close() error {
	e.cancel()

	return e.client.Close()
}

func (e *EtcdSequencer) getKeyInEtcd(key string) string {
	return fmt.Sprintf("%s%s", e.metadata.KeyPrefix, key)
}

func addPathSeparator(p string) string {
	if p == "" {
		return "/"
	}
	if p[0] != '/' {
		p = "/" + p
	}
	if p[len(p)-1] != '/' {
		p = p + "/"
	}
	return p
}
