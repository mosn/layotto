/*
 * Copyright 2021 Layotto Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package mock_state

import (
	"encoding/json"
	"fmt"
	"github.com/dapr/components-contrib/state"
	"sync"
)

type inMemStateStoreItem struct {
	data []byte
	etag *string
}

type inMemoryStateStore struct {
	items map[string]*inMemStateStoreItem
	lock  *sync.RWMutex
}

func NewInMemoryStateStore() state.Store {
	return &inMemoryStateStore{
		items: map[string]*inMemStateStoreItem{},
		lock:  &sync.RWMutex{},
	}
}

func (store *inMemoryStateStore) newItem(data []byte, etagString *string) *inMemStateStoreItem {
	return &inMemStateStoreItem{
		data: data,
		etag: etagString,
	}
}

func (store *inMemoryStateStore) Init(metadata state.Metadata) error {
	return nil
}

func (store *inMemoryStateStore) Ping() error {
	return nil
}

func (store *inMemoryStateStore) Features() []state.Feature {
	return []state.Feature{state.FeatureETag, state.FeatureTransactional}
}

func (store *inMemoryStateStore) Delete(req *state.DeleteRequest) error {
	store.lock.Lock()
	defer store.lock.Unlock()
	delete(store.items, req.Key)

	return nil
}

func (store *inMemoryStateStore) BulkDelete(req []state.DeleteRequest) error {
	if req == nil || len(req) == 0 {
		return nil
	}
	for _, dr := range req {
		err := store.Delete(&dr)
		if err != nil {
			return err
		}
	}
	return nil
}

func (store *inMemoryStateStore) Get(req *state.GetRequest) (*state.GetResponse, error) {
	store.lock.RLock()
	defer store.lock.RUnlock()
	item := store.items[req.Key]

	if item == nil {
		return &state.GetResponse{Data: nil, ETag: nil}, nil
	}

	return &state.GetResponse{Data: item.data, ETag: item.etag}, nil
}

func (store *inMemoryStateStore) BulkGet(req []state.GetRequest) (bool, []state.BulkGetResponse, error) {
	res := []state.BulkGetResponse{}
	for _, oneRequest := range req {
		oneResponse, err := store.Get(&state.GetRequest{
			Key:      oneRequest.Key,
			Metadata: oneRequest.Metadata,
			Options:  oneRequest.Options,
		})
		if err != nil {
			return false, nil, err
		}

		res = append(res, state.BulkGetResponse{
			Key:  oneRequest.Key,
			Data: oneResponse.Data,
			ETag: oneResponse.ETag,
		})
	}

	return true, res, nil
}

func (store *inMemoryStateStore) Set(req *state.SetRequest) error {
	b, _ := Marshal(req.Value, json.Marshal)
	store.lock.Lock()
	defer store.lock.Unlock()
	store.items[req.Key] = store.newItem(b, req.ETag)

	return nil
}

func (store *inMemoryStateStore) BulkSet(req []state.SetRequest) error {
	for _, r := range req {
		err := store.Set(&r)
		if err != nil {
			return err
		}
	}
	return nil
}

func (store *inMemoryStateStore) Multi(request *state.TransactionalStateRequest) error {
	store.lock.Lock()
	defer store.lock.Unlock()
	// First we check all eTags
	for _, o := range request.Operations {
		var eTag *string
		key := ""
		if o.Operation == state.Upsert {
			key = o.Request.(state.SetRequest).Key
			eTag = o.Request.(state.SetRequest).ETag
		} else if o.Operation == state.Delete {
			key = o.Request.(state.DeleteRequest).Key
			eTag = o.Request.(state.DeleteRequest).ETag
		}
		item := store.items[key]
		if eTag != nil && item != nil {
			if *eTag != *item.etag {
				return fmt.Errorf("etag does not match for key %v", key)
			}
		}
		if eTag != nil && item == nil {
			return fmt.Errorf("etag does not match for key not found %v", key)
		}
	}

	// Now we can perform the operation.
	for _, o := range request.Operations {
		if o.Operation == state.Upsert {
			req := o.Request.(state.SetRequest)
			b, _ := json.Marshal(req.Value)
			store.items[req.Key] = store.newItem(b, req.ETag)
		} else if o.Operation == state.Delete {
			req := o.Request.(state.DeleteRequest)
			delete(store.items, req.Key)
		}
	}

	return nil
}

func Marshal(val interface{}, marshaler func(interface{}) ([]byte, error)) ([]byte, error) {
	var err error = nil
	bt, ok := val.([]byte)
	if !ok {
		bt, err = marshaler(val)
	}

	return bt, err
}
