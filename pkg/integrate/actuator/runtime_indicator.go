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

package actuator

import (
	"sync"

	"mosn.io/layotto/pkg/actuator/health"
)

const (
	reasonKey           = "reason"
	reasonValueStarting = "starting"
)

var runtimeReady *runtimeIndicatorImpl
var runtimeLive *runtimeIndicatorImpl

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
	Report() (status health.Status, details map[string]interface{})
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

func (idc *runtimeIndicatorImpl) Report() (status health.Status, details map[string]interface{}) {
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
