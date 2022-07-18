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

package wasm

import (
	"mosn.io/layotto/pkg/filter/stream/common/http"

	"mosn.io/pkg/log"
)

type Wasm struct {
	endpointRegistry map[string]http.Endpoint
}

// New init a Wasm.
func New() *Wasm {
	return &Wasm{
		endpointRegistry: make(map[string]http.Endpoint),
	}
}

// GetEndpoint get an Endpoint from Wasm with name.
func (wasm *Wasm) GetEndpoint(name string) (endpoint http.Endpoint, ok bool) {
	e, ok := wasm.endpointRegistry[name]
	return e, ok
}

// AddEndpoint add an Endpoint to Wasm。
func (wasm *Wasm) AddEndpoint(name string, ep http.Endpoint) {
	_, ok := wasm.endpointRegistry[name]
	if ok {
		log.DefaultLogger.Warnf("Duplicate Endpoint name:  %v !", name)
	}
	wasm.endpointRegistry[name] = ep
}

var singleton = New()

func GetDefault() *Wasm {
	return singleton
}
