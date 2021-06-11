package persistence

import (
	"math/rand"
	"mosn.io/layotto/pkg/common"
	"mosn.io/layotto/pkg/filter/network/tcpcopy/model"
	"mosn.io/pkg/utils"
	"sync"
	"time"
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
		for {
			select {
			case <-tick.C:
				g.work()
			}
		}
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
