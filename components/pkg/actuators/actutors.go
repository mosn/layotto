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

package actuators

import (
	"sync"
)

// Status is the enumeration value of component health status.
type Status = string

var (
	// INIT means it is starting
	INIT = Status("INIT")
	// UP means it is healthy
	UP = Status("UP")
	// DOWN means it is unhealthy
	DOWN = Status("DOWN")
)

var allActiveComponents sync.Map
var componentsActutors sync.Map

type Indicator interface {
	Report() (status Status, details map[string]interface{})
}

type ComponentsIndicator struct {
	ReadinessIndicator Indicator
	LivenessIndicator  Indicator
}

func GetIndicatorWithName(name string) *ComponentsIndicator {
	if v, ok := componentsActutors.Load(name); ok {
		return v.(*ComponentsIndicator)
	}
	return nil
}

func SetComponentsActuators(name string, indicator *ComponentsIndicator) {
	componentsActutors.Store(name, indicator)
}

//here return map for json marshal directly
func GetComponentsActiveStatus() (interface{}, error) {
	res := make(map[string]bool)
	allActiveComponents.Range(func(key, value interface{}) bool {
		res[key.(string)] = value.(bool)
		return true
	})
	return res, nil
}

func SetComponentActive(component string) {
	allActiveComponents.LoadOrStore(component, true)
}
