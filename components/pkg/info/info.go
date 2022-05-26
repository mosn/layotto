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

package info

// Runtimeinfo
type RuntimeInfo struct {
	Services ServiceInfo `json:"services"`
}

// ServiceInfo
type ServiceInfo map[string]*ComponentInfo

// ComponentInfo
type ComponentInfo struct {
	// Registered Component
	Registered []string `json:"registered"`
	// Loaded Component
	Loaded []string `json:"loaded"`
}

func NewRuntimeInfo() *RuntimeInfo {
	return &RuntimeInfo{
		Services: ServiceInfo{},
	}
}

func (info *RuntimeInfo) AddService(service string) {
	info.Services[service] = &ComponentInfo{}
}

func (info *RuntimeInfo) RegisterComponent(service string, compType string) {
	if c, ok := info.Services[service]; ok {
		c.Registered = append(c.Registered, compType)
	}
}

func (info *RuntimeInfo) LoadComponent(service string, compType string) {
	if c, ok := info.Services[service]; ok {
		c.Loaded = append(c.Loaded, compType)
	}
}
