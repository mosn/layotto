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

package runtime

import "mosn.io/layotto/components/lock"

type extensionComponents struct {
	locks map[string]lock.LockStore
}

func newExtensionComponents() *extensionComponents {
	return &extensionComponents{
		locks: make(map[string]lock.LockStore),
	}
}

func (m *MosnRuntime) initExtensionComponent(s services) error {
	if err := m.initLocks(s.locks...); err != nil {
		return err
	}
	return nil
}
