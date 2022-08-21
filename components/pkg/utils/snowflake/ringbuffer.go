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
package snowflake

import (
	"errors"
	"runtime"
	"sync"
	"sync/atomic"

	"mosn.io/pkg/log"
)

const (
	CAN_TAKE_FLAG = 1
	CAN_PUT_FLAG  = 0
)

//avoid false sharing
type PaddedInt struct {
	_     [7]int64
	Value int64
}

type RingBuffer struct {
	m      sync.Mutex
	slots  []int64
	flags  []PaddedInt
	tail   PaddedInt
	cursor PaddedInt

	MaxSeq     int64
	GoNum      int32
	TimeBits   int64
	WorkIdBits int64
	SeqBits    int64
	bufferSize int64
	WorkId     int64

	PaddingFactor    int64
	CurrentTimeStamp int64
}

var cores int = runtime.NumCPU()

func NewRingBuffer(bufferSize int64) *RingBuffer {
	p := PaddedInt{}
	p.Value = -1
	return &RingBuffer{
		slots:      make([]int64, bufferSize),
		flags:      make([]PaddedInt, bufferSize),
		tail:       p,
		cursor:     p,
		bufferSize: bufferSize,
		GoNum:      1,
	}
}

func (r *RingBuffer) Put(uid int64) (bool, error) {
	r.m.Lock()
	defer r.m.Unlock()
	currentTail := r.tail.Value
	currentCursor := r.cursor.Value
	if currentCursor == -1 {
		currentCursor = 0
	}
	distance := currentTail - currentCursor
	if distance == r.bufferSize-1 {
		return false, errors.New("catched!Rejected putting buffer")
	}

	//(currentTail + 1) mod r.bufferSize
	nextTailIndex := (currentTail + 1) & (r.bufferSize - 1)

	if r.flags[nextTailIndex].Value != CAN_PUT_FLAG {
		return false, errors.New("tail not in can put status")
	}

	r.slots[nextTailIndex] = uid
	r.flags[nextTailIndex].Value = CAN_TAKE_FLAG
	r.tail.Value++
	return true, nil
}

func (r *RingBuffer) Take() (int64, error) {
	r.m.Lock()
	defer r.m.Unlock()
	var uid int64
	currentCursor := r.cursor.Value
	if r.cursor.Value != r.tail.Value {
		r.cursor.Value++
	}
	nextCursor := r.cursor.Value
	currentTail := r.tail.Value

	if currentTail-nextCursor < r.PaddingFactor*r.bufferSize/100 {
		//limit the numbers of goroutine
		if int(r.GoNum) <= 2*cores {
			r.GoNum++
			go r.PaddingRingBuffer()
		}
	}

	if currentCursor == nextCursor {
		return uid, errors.New("buffer is empty, rejected take buffer")
	}

	//check next slot flag is CAN_TAKE_FLAG
	nextCursorIndex := (nextCursor) & (r.bufferSize - 1)
	if r.flags[nextCursorIndex].Value != CAN_TAKE_FLAG {
		return uid, errors.New("curosr not in can take status")
	}

	uid = r.slots[nextCursorIndex]
	r.flags[nextCursorIndex].Value = CAN_PUT_FLAG
	return uid, nil
}

func (r *RingBuffer) Allocator() int64 {
	var sequence int64
	timestampShift := r.SeqBits
	workidShift := r.TimeBits + r.SeqBits
	workid := r.WorkId
	r.m.Lock()
	timestamp := r.CurrentTimeStamp
	r.CurrentTimeStamp++
	r.m.Unlock()
	return timestamp<<timestampShift | (workid << workidShift) | sequence
}

func (r *RingBuffer) GenerateUid(cur int64) (bool, error) {
	maxSeq := r.MaxSeq
	var offset int64
	for offset = 0; offset < maxSeq; offset++ {
		if ok, err := r.Put(cur + offset); !ok {
			return false, err
		}
	}
	return true, nil

}

func (r *RingBuffer) PaddingRingBuffer() {
	defer func() {
		if x := recover(); x != nil {
			log.DefaultLogger.Errorf("panic when generatoring uid with snowflake algorithm and padding ringbuffer: %v", x)
		}
	}()
	for {
		u := r.Allocator()
		if ok, err := r.GenerateUid(u); !ok || err != nil {
			log.DefaultLogger.Warnf("%v", err)
			break
		}
	}
	atomic.AddInt32(&r.GoNum, -1)
}
