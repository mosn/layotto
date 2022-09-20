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
	"sync"
	"sync/atomic"
	"time"

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

//there are two ringbuffers, one is to store uid, another is to store flag, flag represents the slot is readable or writable
type RingBuffer struct {
	m sync.Mutex
	//store uid
	slots []int64
	//store flag
	flags []PaddedInt
	//write pointer, wp is readable, next is writable
	wp PaddedInt
	//read pointer, rp is writable, next is readable
	rp PaddedInt

	//1 + WorkIdBits + TimeBits + SeqBits = 64
	TimeBits   int64
	WorkIdBits int64
	SeqBits    int64

	//ringbuffer's size
	bufferSize int64
	//when readable slots nums <= bufferSize * PaddingFactor / 100, start padding goroutine
	PaddingFactor int64

	//get id from Mysql at startup
	WorkId int64
	//get current timestamp at startup
	CurrentTimeStamp int64

	//whether there is a running padding goroutine
	running int32
}

func NewRingBuffer(bufferSize int64) *RingBuffer {
	p := PaddedInt{}
	p.Value = -1
	return &RingBuffer{
		slots:      make([]int64, bufferSize),
		flags:      make([]PaddedInt, bufferSize),
		wp:         p,
		rp:         p,
		bufferSize: bufferSize,
		running:    0,
	}
}

func (r *RingBuffer) Put(uid int64) (bool, error) {
	r.m.Lock()
	defer r.m.Unlock()
	currentWritePointer := r.wp.Value
	currentReadPointer := r.rp.Value
	if currentReadPointer == -1 {
		currentReadPointer = 0
	}

	distance := currentWritePointer - currentReadPointer
	//write pointer catches read pointer, ringbuffer is full
	if distance == r.bufferSize-1 {
		return false, errors.New("ringbuffer is full! Rejected putting buffer")
	}

	nextWriteIndex := (currentWritePointer + 1) & (r.bufferSize - 1)
	if r.flags[nextWriteIndex].Value != CAN_PUT_FLAG {
		//The flag at this location may not have been updated
		time.Sleep(time.Microsecond)
		if r.flags[nextWriteIndex].Value != CAN_PUT_FLAG {
			return false, errors.New("slot is not in writable status")
		}
	}

	r.slots[nextWriteIndex] = uid
	r.flags[nextWriteIndex].Value = CAN_TAKE_FLAG
	r.wp.Value++
	return true, nil
}

//taking uid is a lock free operation to ensure the only padding goroutine can get enough time slices
func (r *RingBuffer) Take() (int64, error) {
	var uid int64
	var currentReadPointer int64
	var nextReadPointer int64

	//atomically update r.rp.value if r.rp.Value != r.wp.Value
	for {
		currentReadPointer = r.rp.Value
		nextReadPointer = currentReadPointer
		if currentReadPointer < r.wp.Value {
			nextReadPointer++
		}
		if atomic.CompareAndSwapInt64(&r.rp.Value, currentReadPointer, nextReadPointer) {
			break
		}
	}

	if r.wp.Value-currentReadPointer < r.bufferSize*r.PaddingFactor/100 {
		go r.PaddingRingBuffer()
	}

	//if buffer is empty, wait a moment
	if currentReadPointer == nextReadPointer {
		running := false
		for !running {
			running = r.PaddingRingBuffer()
		}
		if currentReadPointer == nextReadPointer {
			log.DefaultLogger.Warnf("buffer is empty")
		}
	}

	//check next slot flag is CAN_TAKE_FLAG
	nextReadIndex := nextReadPointer & (r.bufferSize - 1)
	if atomic.LoadInt64(&r.flags[nextReadIndex].Value) != CAN_TAKE_FLAG {
		return uid, errors.New("buffer is empty! Please request again")
	}

	uid = r.slots[nextReadIndex]
	r.flags[nextReadIndex].Value = CAN_PUT_FLAG
	return uid, nil
}

//allocate workid and timestamp for uid
func (r *RingBuffer) Allocator() int64 {
	timestampShift := r.SeqBits
	workidShift := r.TimeBits + r.SeqBits

	workid := r.WorkId
	timestamp := r.CurrentTimeStamp
	r.CurrentTimeStamp++

	var sequence int64
	return timestamp<<timestampShift | (workid << workidShift) | sequence
}

//put uid into ringbuffer
//uid： workid + timestamp + (0 ~ maxSeq)
func (r *RingBuffer) GenerateUid(cur int64) (bool, error) {
	var maxSeq int64 = 1 << r.SeqBits
	var offset int64
	for offset = 0; offset < maxSeq; offset++ {
		if ok, err := r.Put(cur + offset); !ok {
			return false, err
		}
	}
	return true, nil
}

//async padding ringbuffer
func (r *RingBuffer) PaddingRingBuffer() bool {
	defer func() {
		if x := recover(); x != nil {
			log.DefaultLogger.Errorf("panic when generatoring uid with snowflake algorithm and padding ringbuffer: %v", x)
		}
	}()

	//there can only be one padding goroutine at a time, to ensure uids are strictly increasing
	if !atomic.CompareAndSwapInt32(&r.running, 0, 1) {
		return false
	}

	for {
		u := r.Allocator()
		if ok, err := r.GenerateUid(u); !ok || err != nil {
			log.DefaultLogger.Warnf("%v", err)
			break
		}
	}

	if !atomic.CompareAndSwapInt32(&r.running, 1, 0) {
		panic("running num is false")
	}
	return true
}
