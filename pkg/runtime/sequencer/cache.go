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
package sequencer

import (
	"context"
	"errors"
	"sync"
	"time"

	"mosn.io/pkg/log"
	"mosn.io/pkg/utils"

	"mosn.io/layotto/components/sequencer"
)

const defaultSize = 10000
const defaultLimit = 1000
const defaultRetry = 5
const waitTime = time.Second * 2

// DoubleBuffer is double segment id buffer.
// There are two buffers in DoubleBuffer: inUseBuffer is in use, BackUpBuffer is a backup buffer.
// Their default capacity is 1000. When the inUseBuffer usage exceeds 30%, the BackUpBuffer will be initialized.
// When inUseBuffer is used up, swap them.
type DoubleBuffer struct {
	Key              string
	size             int
	inUseBuffer      *Buffer
	backUpBufferChan chan *Buffer
	lock             sync.Mutex
	Store            sequencer.Store
}

type Buffer struct {
	from int64
	to   int64
}

func NewDoubleBuffer(key string, store sequencer.Store) *DoubleBuffer {

	d := &DoubleBuffer{
		Key:              key,
		size:             defaultSize,
		Store:            store,
		backUpBufferChan: make(chan *Buffer, 1),
	}

	return d
}

// init double buffer
func (d *DoubleBuffer) init() error {

	buffer, err := d.getNewBuffer()
	if err != nil {
		return err
	}

	d.inUseBuffer = buffer

	return nil
}

// getId next id
func (d *DoubleBuffer) getId() (int64, error) {

	d.lock.Lock()
	defer d.lock.Unlock()

	if d.inUseBuffer == nil {
		return 0, errors.New("[DoubleBuffer] Get error: inUseBuffer nil ")
	}
	//check swap
	if d.inUseBuffer.from > d.inUseBuffer.to {
		err := d.swap()
		if err != nil {
			return 0, err
		}
	}
	next := d.inUseBuffer.from
	d.inUseBuffer.from++

	//when inUseBuffer id more than limit used, initialize BackUpBuffer.
	//equal make sure only one thread enter
	if d.inUseBuffer.to-d.inUseBuffer.from == defaultLimit {
		utils.GoWithRecover(func() {
			//quick retry
			for i := 0; i < defaultRetry; i++ {
				buffer, err := d.getNewBuffer()
				if err != nil {
					log.DefaultLogger.Errorf("[DoubleBuffer] [getNewBuffer] error: %v", err)
					continue
				}
				d.backUpBufferChan <- buffer
				return
			}
			//slow retry
			for {
				buffer, err := d.getNewBuffer()
				if err != nil {
					log.DefaultLogger.Errorf("[DoubleBuffer] [getNewBuffer] error: %v", err)
					time.Sleep(waitTime)
					continue
				}
				d.backUpBufferChan <- buffer
				return
			}
		}, nil)
	}

	return next, nil
}

// swap inUseBuffer and BackUpBuffer, must be locked
func (d *DoubleBuffer) swap() error {

	select {
	case buffer := <-d.backUpBufferChan:
		{
			d.inUseBuffer = buffer
			return nil
		}
	//timeout, return error
	case <-time.After(waitTime):
		{
			return errors.New("[DoubleBuffer] swap error")
		}
	}
}

// getNewBuffer return a new segment
func (d *DoubleBuffer) getNewBuffer() (*Buffer, error) {
	support, result, err := d.Store.GetSegment(&sequencer.GetSegmentRequest{
		Key:  d.Key,
		Size: d.size,
	})
	if err != nil {
		return nil, err
	}
	if !support {
		return nil, errors.New("[DoubleBuffer] unSupport Segment id")
	}
	return &Buffer{
		from: result.From,
		to:   result.To,
	}, nil
}

// BufferCatch catch key and buffer
var BufferCatch = map[string]*DoubleBuffer{}

// read/write lock for BufferCatch
var rwLock sync.RWMutex

func GetNextIdFromCache(ctx context.Context, store sequencer.Store, req *sequencer.GetNextIdRequest) (bool, int64, error) {

	// 1. check support
	support, _, _ := store.GetSegment(&sequencer.GetSegmentRequest{
		Key:  req.Key,
		Size: 0,
	})

	// return if not support
	if !support {
		return false, 0, nil
	}

	// 2. find the DoubleBuffer for this store and key
	var d *DoubleBuffer
	var err error

	d = getDoubleBufferInRL(req.Key)
	if d == nil {
		d, err = getDoubleBufferInWL(req.Key, store)
	}

	if err != nil {
		return true, 0, err
	}

	// 3. get the next id.
	// The buffer should automatically load segment into cache if the cache is (nearly) empty
	id, err := d.getId()

	if err != nil {
		return true, 0, err
	}

	return true, id, nil
}

// get DoubleBuffer using write lock
func getDoubleBufferInWL(key string, store sequencer.Store) (*DoubleBuffer, error) {
	d := NewDoubleBuffer(key, store)
	rwLock.Lock()
	defer rwLock.Unlock()
	//double check
	if _, ok := BufferCatch[key]; ok {
		return BufferCatch[key], nil
	}
	err := d.init()
	if err != nil {
		return nil, err
	}
	BufferCatch[key] = d
	return d, nil
}

// get DoubleBuffer using read lock
func getDoubleBufferInRL(key string) *DoubleBuffer {
	rwLock.RLock()
	defer rwLock.RUnlock()
	if buffer, ok := BufferCatch[key]; ok {
		return buffer
	}
	return nil
}
