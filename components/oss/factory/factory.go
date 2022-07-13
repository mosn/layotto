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

package factory

import (
	"encoding/json"
	"sync"
)

var initFuncRegistry map[string]S3ClientInit
var mux sync.RWMutex

type S3ClientInit func(staticConf json.RawMessage, dynConf map[string]string) (map[string]interface{}, error)

func init() {
	initFuncRegistry = make(map[string]S3ClientInit)
}
func RegisterInitFunc(name string, f S3ClientInit) {
	mux.Lock()
	initFuncRegistry[name] = f
	mux.Unlock()
}

func GetInitFunc(name string) S3ClientInit {
	mux.RLock()
	defer mux.RUnlock()
	return initFuncRegistry[name]
}
