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

package persistence

import (
	"math/rand"
	"sync"
	"time"

	"mosn.io/pkg/utils"

	"mosn.io/layotto/pkg/common"
	"mosn.io/layotto/pkg/filter/network/tcpcopy/model"
)

func init() {
	dumpWorkPoolInstance = NewDefaultWorkPool(20)
}

var dumpWorkPoolInstance *DefaultWorkPool

type WorkGoroutine struct {
	tasks *sync.Map
}

func NewWorkGoroutine() *WorkGoroutine {
	worker := &WorkGoroutine{
		tasks: new(sync.Map),
	}
	return worker
}

func (g *WorkGoroutine) AddTask(key string, data *model.DumpUploadDynamicConfig) {
	g.tasks.Store(key, data)
}

func (g *WorkGoroutine) Start() {
	utils.GoWithRecover(func() {
		tick := time.NewTicker(500 * time.Millisecond)
		<-tick.C
		g.work()
	}, func(r interface{}) {
		g.Start()
	})
}

func (g *WorkGoroutine) work() {
	g.tasks.Range(func(key, value interface{}) bool {
		data := value.(*model.DumpUploadDynamicConfig)
		persistence(data)
		g.tasks.Delete(key)
		return true
	})
}

type DefaultWorkPool struct {
	size           int64
	workers        *sync.Map
	randomInstance *rand.Rand
	lock           *sync.Mutex
}

func NewDefaultWorkPool(size int64) *DefaultWorkPool {
	workPool := &DefaultWorkPool{
		size:           size,
		workers:        new(sync.Map),
		randomInstance: rand.New(rand.NewSource(time.Now().UnixNano())),
		lock:           new(sync.Mutex),
	}
	return workPool
}

func GetDumpWorkPoolInstance() *DefaultWorkPool {
	return dumpWorkPoolInstance
}

func (w *DefaultWorkPool) random() int64 {
	w.lock.Lock()
	defer w.lock.Unlock()
	return w.randomInstance.Int63n(w.size)
}

func (w *DefaultWorkPool) Schedule(data *model.DumpUploadDynamicConfig) {
	index := w.random()
	key := common.CalculateMd5(string(data.BusinessType)) + common.CalculateMd5ForBytes(data.Binary_flow_data)
	if value, ok := w.workers.Load(index); ok {
		worker := value.(*WorkGoroutine)
		worker.AddTask(key, data)
	} else {
		worker := NewWorkGoroutine()
		worker.AddTask(key, data)
		if _, ok := w.workers.LoadOrStore(index, worker); !ok {
			worker.Start()
		} else {
			worker.AddTask(key, data)
		}
	}
}
