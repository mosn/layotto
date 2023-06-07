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

package nacos

import "sync"

// store the group and keys that client subscribe
type subscriberHolder struct {
	rw     sync.RWMutex
	keyMap map[subscriberKey]struct{}
}

type subscriberKey struct {
	group string
	key   string
}

func newSubscriberHolder() *subscriberHolder {
	return &subscriberHolder{
		keyMap: make(map[subscriberKey]struct{}),
	}
}

func (s *subscriberHolder) AddSubscriberKey(key subscriberKey) {
	s.rw.Lock()
	defer s.rw.Unlock()
	s.keyMap[key] = struct{}{}
}

func (s *subscriberHolder) RemoveSubscriberKey(key subscriberKey) {
	s.rw.Lock()
	defer s.rw.Unlock()
	delete(s.keyMap, key)
}

func (s *subscriberHolder) GetSubscriberKey() []subscriberKey {
	s.rw.RLock()
	defer s.rw.RUnlock()

	res := make([]subscriberKey, 0, len(s.keyMap))
	for k := range s.keyMap {
		res = append(res, k)
	}

	return res
}
