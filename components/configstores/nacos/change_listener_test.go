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

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_newSubscriberHolder(t *testing.T) {
	listener := newSubscriberHolder()
	assert.NotNil(t, listener)
	assert.NotNil(t, listener.keyMap)
}

func Test_subscriberHolder_AddSubscriberKey(t *testing.T) {
	listener := newSubscriberHolder()
	key := subscriberKey{key: "test-key", group: "test-group"}
	listener.AddSubscriberKey(key)
	s, ok := listener.keyMap[key]
	assert.True(t, ok)
	assert.EqualValues(t, struct{}{}, s)
}

func Test_subscriberHolder_GetSubscriberKey(t *testing.T) {
	listener := newSubscriberHolder()
	key := subscriberKey{key: "test-key", group: "test-group"}
	listener.AddSubscriberKey(key)
	getSubscriberKey := listener.GetSubscriberKey()
	assert.EqualValues(t, 1, len(getSubscriberKey))
	assert.EqualValues(t, key, getSubscriberKey[0])

	key2 := subscriberKey{key: "test-key2", group: "test-group2"}
	listener.AddSubscriberKey(key2)
	getSubscriberKey = listener.GetSubscriberKey()
	// Due to the use of go map results, the return order of map traversal
	// cannot be determined, so the testing here haven't test the key2's value.
	assert.Equal(t, 2, len(getSubscriberKey))
}

func Test_subscriberHolder_RemoveSubscriberKey(t *testing.T) {
	listener := newSubscriberHolder()
	key := subscriberKey{key: "test-key", group: "test-group"}
	listener.AddSubscriberKey(key)
	assert.EqualValues(t, 1, len(listener.keyMap))

	listener.RemoveSubscriberKey(key)
	assert.Equal(t, 0, len(listener.keyMap))
}
