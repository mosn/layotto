package apollo

import (
	"sync"

	"mosn.io/layotto/components/pkg/actuators"

	"mosn.io/layotto/components/pkg/common"
)

const (
	reasonKey = "reason"
)

var (
	readinessIndicator *healthIndicator
	livenessIndicator  *healthIndicator
)

func init() {
	readinessIndicator = newHealthIndicator()
	livenessIndicator = newHealthIndicator()
	indicators := &actuators.ComponentsIndicator{ReadinessIndicator: readinessIndicator, LivenessIndicator: livenessIndicator}
	actuators.SetComponentsActuators("apollo", indicators)
}

func newHealthIndicator() *healthIndicator {
	return &healthIndicator{
		started: false,
		isErr:   false,
	}
}

func GetReadinessIndicator() *healthIndicator {
	return readinessIndicator
}

func GetLivenessIndicator() *healthIndicator {
	return livenessIndicator
}

type healthIndicator struct {
	mu sync.Mutex

	started   bool
	isErr     bool
	errReason string
}

func (idc *healthIndicator) Report() (status string, details map[string]interface{}) {
	idc.mu.Lock()
	defer idc.mu.Unlock()
	statusDetail := make(map[string]interface{})
	status = common.INIT
	if idc.isErr {
		status = common.DOWN
		statusDetail[reasonKey] = idc.errReason
	}
	if idc.started {
		status = common.UP
	}

	return status, statusDetail
}

func (idc *healthIndicator) reportError(reason string) {
	idc.mu.Lock()
	defer idc.mu.Unlock()

	if idc.isErr {
		return
	}
	idc.isErr = true
	idc.errReason = reason
}

func (idc *healthIndicator) setStarted() {
	idc.mu.Lock()
	defer idc.mu.Unlock()

	idc.started = true
}
