package actuator

import (
	"mosn.io/layotto/pkg/actuator/health"
	"sync"
)

const (
	reasonKey           = "reason"
	reasonValueStarting = "starting"
)

var runtimeReady *runtimeIndicatorImpl
var runtimeLive *runtimeIndicatorImpl
var once sync.Once

func init() {
	runtimeReady = &runtimeIndicatorImpl{
		started: false,
		health:  true,
		reason:  "",
	}
	runtimeLive = &runtimeIndicatorImpl{
		started: false,
		health:  true,
		reason:  "",
	}
}

type RuntimeIndicator interface {
	Report() (status string, details map[string]interface{})
	SetUnhealthy(reason string)
	SetHealthy(reason string)
	SetStarted()
}

func GetRuntimeReadinessIndicator() RuntimeIndicator {
	return runtimeReady
}

func GetRuntimeLivenessIndicator() RuntimeIndicator {
	return runtimeLive
}

type runtimeIndicatorImpl struct {
	mu sync.RWMutex

	started bool
	health  bool
	reason  string
}

func (idc *runtimeIndicatorImpl) SetStarted() {
	idc.mu.Lock()
	defer idc.mu.Unlock()

	idc.started = true
}

func (idc *runtimeIndicatorImpl) Report() (status string, details map[string]interface{}) {
	idc.mu.RLock()
	defer idc.mu.RUnlock()

	if !idc.health {
		h := health.NewHealth(health.DOWN)
		h.SetDetail(reasonKey, idc.reason)
		return h.Status, h.Details
	}
	if !idc.started {
		h := health.NewHealth(health.INIT)
		h.SetDetail(reasonKey, reasonValueStarting)
		return h.Status, h.Details
	}
	h := health.NewHealth(health.UP)
	h.SetDetail(reasonKey, idc.reason)
	return h.Status, h.Details
}

func (idc *runtimeIndicatorImpl) SetUnhealthy(reason string) {
	idc.mu.Lock()
	defer idc.mu.Unlock()

	idc.health = false
	idc.reason = reason
}

func (idc *runtimeIndicatorImpl) SetHealthy(reason string) {
	idc.mu.Lock()
	defer idc.mu.Unlock()

	idc.health = true
	idc.reason = reason
}
