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
