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

import "mosn.io/pkg/log"

type Actuator struct {
	endpointRegistry map[string]Endpoint
}

func New() *Actuator {
	return &Actuator{
		endpointRegistry: make(map[string]Endpoint),
	}
}

func (act *Actuator) GetEndpoint(name string) (endpoint Endpoint, ok bool) {
	e, ok := act.endpointRegistry[name]
	return e, ok
}

func (act *Actuator) AddEndpoint(name string, ep Endpoint) {
	_, ok := act.endpointRegistry[name]
	if ok {
		log.DefaultLogger.Warnf("Duplicate Endpoint name:  %v !", name)
	}
	act.endpointRegistry[name] = ep
}
