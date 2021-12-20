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

package consul

import (
	"sync"
	"time"
)

type task func()

// generate a GC task which delete element in the map after specific ttl
func generateGCTask(ttl int32, m *sync.Map, key string) task {
	return func() {
		time.Sleep(time.Second * time.Duration(ttl))
		//may delete the second lock,but not affect the result
		m.Delete(key)
	}
}
