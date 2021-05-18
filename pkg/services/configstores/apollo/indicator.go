package apollo

import (
	"github.com/layotto/layotto/pkg/filter/stream/actuator/health"
	"sync"
)

const (
	reasonKey = "reason"
)

var (
	readinessIndicator *healthIndicator
)

func init() {
	readinessIndicator = newHealthIndicator()
	health.AddReadinessIndicator("apollo", readinessIndicator)
	health.AddLivenessIndicator("apollo", readinessIndicator)
}

func newHealthIndicator() *healthIndicator {
	return &healthIndicator{}
}

func getReadinessIndicator() *healthIndicator {
	return readinessIndicator
}

type healthIndicator struct {
	mu sync.Mutex

	isErr     bool
	errReason string
}

func (idc *healthIndicator) Report() health.Health {
	idc.mu.Lock()
	defer idc.mu.Unlock()

	var h health.Health

	if idc.isErr {
		h = health.NewHealth(health.DOWN)
		h.Details[reasonKey] = idc.errReason
		return h
	}

	h = health.NewHealth(health.UP)
	return h
}

func (idc *healthIndicator) ReportError(reason string) {
	idc.mu.Lock()
	defer idc.mu.Unlock()

	if idc.isErr {
		return
	}

	idc.isErr = true
	idc.errReason = reason
}
