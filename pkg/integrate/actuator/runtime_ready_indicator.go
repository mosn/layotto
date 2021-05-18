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
	started bool
	health  bool
	reason  string
}

func (s *runtimeReadyIndicatorImpl) SetStarted() {
	s.started = true
}

func (s *runtimeReadyIndicatorImpl) Report() health.Health {
	if !s.health {
		h := health.NewHealth(health.DOWN)
		h.SetDetail(reasonKey, s.reason)
		return h
	}
	if !s.started {
		h := health.NewHealth(health.DOWN)
		h.SetDetail(reasonKey, reasonValueStarting)
		return h
	}
	return health.NewHealth(health.UP)
}

func (c *runtimeReadyIndicatorImpl) SetUnhealth(reason string) {
	c.health = false
	c.reason = reason
}

func (c *runtimeReadyIndicatorImpl) SetHealth(reason string) {
	c.health = true
	c.reason = reason
}
