package actuator

import (
	"github.com/layotto/layotto/pkg/filter/stream/actuator/health"
	"sync"
)

const (
	reasonKey           = "reason"
	reasonValueStarting = "starting"
)

var singleton *runtimeReadyIndicatorImpl
var once sync.Once

func init() {
	singleton = &runtimeReadyIndicatorImpl{
		started: false,
		health:  true,
		reason:  "",
	}
}

type RuntimeReadyIndicator interface {
	Report() health.Health
	SetUnhealth(reason string)
	SetHealth(reason string)
	SetStarted()
}

func GetRuntimeReadyIndicator() RuntimeReadyIndicator {
	return singleton
}

type runtimeReadyIndicatorImpl struct {
	mu sync.RWMutex

	started bool
	health  bool
	reason  string
}

func (idc *runtimeReadyIndicatorImpl) SetStarted() {
	idc.mu.Lock()
	defer idc.mu.Unlock()

	idc.started = true
}

func (idc *runtimeReadyIndicatorImpl) Report() health.Health {
	idc.mu.RLock()
	defer idc.mu.RUnlock()

	if !idc.health {
		h := health.NewHealth(health.DOWN)
		h.SetDetail(reasonKey, idc.reason)
		return h
	}
	if !idc.started {
		h := health.NewHealth(health.INIT)
		h.SetDetail(reasonKey, reasonValueStarting)
		return h
	}
	h := health.NewHealth(health.UP)
	h.SetDetail(reasonKey, idc.reason)
	return h
}

func (idc *runtimeReadyIndicatorImpl) SetUnhealth(reason string) {
	idc.mu.Lock()
	defer idc.mu.Unlock()

	idc.health = false
	idc.reason = reason
}

func (idc *runtimeReadyIndicatorImpl) SetHealth(reason string) {
	idc.mu.Lock()
	defer idc.mu.Unlock()

	idc.health = true
	idc.reason = reason
}
