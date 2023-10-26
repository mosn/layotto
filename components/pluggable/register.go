// Copyright 2021 Layotto Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pluggable

func init() {
	onServiceDiscovered = make(map[string]CallbackFunc)
}

var (
	// register the callback function of pluggable component, using grpc connection creates factory and registers into MosnRuntime
	onServiceDiscovered map[string]CallbackFunc
)

type Component interface{}

type CallbackFunc func(compType string, dialer GRPCConnectionDialer) Component

// AddServiceDiscoveryCallback register callback function, not concurrent secure
func AddServiceDiscoveryCallback(serviceName string, callback CallbackFunc) {
	onServiceDiscovered[serviceName] = callback
}

func GetServiceDiscoveryMapper() map[string]CallbackFunc {
	return onServiceDiscovered
}
