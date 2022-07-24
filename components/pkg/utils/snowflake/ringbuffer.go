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
	"fmt"
	"sync"
	"sync/atomic"
)

const (
	CAN_TAKE_FLAG = 1
	CAN_PUT_FLAG  = 0
)

type PaddedInt struct {
	_     [7]int64
	value int64
}

type RingBuffer struct {
	m          sync.Mutex
	slots      []int64
	flags      []PaddedInt
	tail       PaddedInt
	cursor     PaddedInt
	bufferSize int64

	paddingThreshold int64
}

func NewRingBuffer(bufferSize int64) *RingBuffer {
	p := PaddedInt{}
	p.value = -1
	return &RingBuffer{
		slots:      make([]int64, bufferSize),
		flags:      make([]PaddedInt, bufferSize),
		tail:       p,
		cursor:     p,
		bufferSize: bufferSize,
	}
}

func (r *RingBuffer) Put(uid int64) (bool, error) {
	r.m.Lock()
	defer r.m.Unlock()
	currentTail := r.tail.value
	currentCursor := r.cursor.value
	if currentCursor == -1 {
		currentCursor = 0
	}
	distance := currentTail - currentCursor
	if distance == r.bufferSize-1 {
		r.RejectPutBuffer(uid)
		return false, fmt.Errorf("Catched!Rejected putting buffer")
	}

	nextTailIndex := (currentTail + 1) & (r.bufferSize - 1)

	if r.flags[nextTailIndex].value != CAN_PUT_FLAG {
		r.RejectPutBuffer(uid)
		return false, fmt.Errorf("Rejected putting buffer")
	}

	r.slots[nextTailIndex] = uid
	r.flags[nextTailIndex].value = CAN_TAKE_FLAG

	atomic.AddInt64(&r.tail.value, 1)

	return true, nil
}

func (r *RingBuffer) Take() (int64, error) {
	r.m.Lock()
	defer r.m.Unlock()
	currentCursor := r.cursor.value
	if r.cursor.value != r.tail.value {
		atomic.AddInt64(&r.cursor.value, 1)
	}
	nextCursor := r.cursor.value

	currentTail := r.tail.value
	if currentTail-nextCursor < r.paddingThreshold {
	}

	if currentCursor == nextCursor {
		return 0, fmt.Errorf("Buffer is empty, rejected take buffer.")
	}

	nextCursorIndex := (nextCursor) & (r.bufferSize - 1)
	if r.flags[nextCursorIndex].value != CAN_TAKE_FLAG {
		return 0, fmt.Errorf("Curosr not in can take status")
	}

	uid := r.slots[nextCursorIndex]
	r.flags[nextCursorIndex].value = CAN_PUT_FLAG
	return uid, nil
}

func (r *RingBuffer) RejectPutBuffer(uid int64) {
}

func (r *RingBuffer) RejectTakeBuffer() {

}
