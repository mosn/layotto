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

package apollo

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"mosn.io/layotto/components/configstores"
)

// Test CRUD operations of subscriberHolder
func Test_subscriberHolder_crud(t *testing.T) {
	h := newSubscriberHolder()
	// add
	ch := make(chan *configstores.SubscribeResp)
	err := h.addByTopic("application", "key1", ch)
	if err != nil {
		t.Errorf("addByTopic() error = %v", err)
	}
	for k, v := range h.chanMap {
		assert.Equal(t, k.group, "application")
		assert.Equal(t, k.keyWithLabel, "key1")
		assert.True(t, v[0].respChan == ch)
	}
	//	find
	topic := h.findByTopic("application", "key1")
	assert.True(t, len(topic) == 1)
	assert.True(t, topic[0].respChan == ch)

	//	 add another item
	ch2 := make(chan *configstores.SubscribeResp)
	err = h.addByTopic("application", "key2", ch2)
	if err != nil {
		t.Errorf("addByTopic() error = %v", err)
	}
	topic = h.findByTopic("application", "key2")
	assert.True(t, len(topic) == 1)
	assert.True(t, topic[0].respChan == ch2)
	s := topic[0]
	// when remove nil then ok
	h.remove(nil)
	topic = h.findByTopic("application", "key2")
	assert.True(t, len(topic) == 1)
	// remove
	h.remove(s)
	topic = h.findByTopic("application", "key2")
	assert.True(t, len(topic) == 0)
	// when removed key not exist then ok
	s.subscriberKey.group = "asdasddasda"
	h.remove(s)
	// reset
	topic = h.findByTopic("application", "key1")
	assert.True(t, len(topic) == 1)
	assert.True(t, topic[0].respChan == ch)
	h.reset()
	topic = h.findByTopic("application", "key1")
	assert.True(t, len(topic) == 0)
}

func Test_addByTopic_whenKeyNotExist_thenReturnEmptySlice(t *testing.T) {
	h := newSubscriberHolder()
	ch := make(chan *configstores.SubscribeResp)
	err := h.addByTopic("application", "key1", ch)
	if err != nil {
		t.Errorf("addByTopic() error = %v", err)
	}
	topic := h.findByTopic("application", "key2")
	assert.True(t, len(topic) == 0)
}

func Test_addByTopic_whenChanNil_thenError(t *testing.T) {
	h := newSubscriberHolder()
	err := h.addByTopic("application", "key1", nil)
	if notNil := assert.NotNil(t, err); notNil {
		assert.True(t, err.Error() != "")
	}
	topic := h.findByTopic("application", "key1")
	assert.True(t, len(topic) == 0)
}
