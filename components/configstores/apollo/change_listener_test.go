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
	"sync"
	"testing"
	"time"

	"github.com/apolloconfig/agollo/v4/storage"
	"github.com/stretchr/testify/assert"

	"mosn.io/layotto/components/configstores"
)

const (
	testAppId     = "test_app"
	testStoreName = "test_storename"
)

type MockRepo struct {
	c *ConfigStore
}

func (m *MockRepo) splitKey(keyWithLabel string) (key string, label string) {
	return m.c.splitKey(keyWithLabel)
}

func (m *MockRepo) getAllTags(group string, keyWithLabel string) (tags map[string]string, err error) {
	return nil, nil
}

func (m *MockRepo) GetAppId() string {
	return testAppId
}

func (m *MockRepo) GetStoreName() string {
	return testStoreName
}

const ns = "application"

func setupChangeListener() *changeListener {
	mockRepo := &MockRepo{
		c: NewStore().(*ConfigStore),
	}
	return newChangeListener(mockRepo)
}

// Test modified
func Test_changeListener_OnChange(t *testing.T) {
	lis := setupChangeListener()
	ch := make(chan *configstores.SubscribeResp)
	// add subscriber
	err := lis.addByTopic(ns, "key1", ch)
	if err != nil {
		t.Error(err)
	}
	// assert consume
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		select {
		case c2 := <-ch:
			assert.Equal(t, c2.StoreName, testStoreName)
			assert.Equal(t, c2.AppId, testAppId)
			assert.True(t, len(c2.Items) == 1)
			assert.True(t, c2.Items[0].Key == "key1")
			assert.True(t, c2.Items[0].Content == "v2")
			wg.Done()
		case <-time.After(time.Second * 2):
			t.Error("consume timeout")
			close(ch)
			wg.Done()
		}
	}()
	// change
	changes := make(map[string]*storage.ConfigChange)
	c := &storage.ConfigChange{
		OldValue:   "v1",
		NewValue:   "v2",
		ChangeType: storage.MODIFIED,
	}
	changes["key1"] = c
	event := &storage.ChangeEvent{
		Changes: changes,
	}
	event.Namespace = ns
	lis.OnChange(event)
	lis.reset()
	assert.True(t, len(lis.subscribers.chanMap) == 0)

	wg.Wait()
}

func Test_changeListener_timeout(t *testing.T) {
	// 1. setup
	lis := setupChangeListener()
	ch := make(chan *configstores.SubscribeResp)
	// add subscriber
	err := lis.addByTopic(ns, "key1", ch)
	if err != nil {
		t.Error(err)
	}
	// 2 assert consuming a closed chan
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		time.Sleep(2 * time.Second)
		select {
		case c2, ok := <-ch:
			assert.Nil(t, c2)
			assert.False(t, ok)
			wg.Done()
		case <-time.After(time.Second * 2):
			t.Error("consume timeout")
			close(ch)
			wg.Done()
		}
	}()
	// 3 change
	changes := make(map[string]*storage.ConfigChange)
	c := &storage.ConfigChange{
		OldValue:   "v1",
		NewValue:   "v2",
		ChangeType: storage.MODIFIED,
	}
	changes["key1"] = c
	event := &storage.ChangeEvent{
		Changes: changes,
	}
	event.Namespace = ns
	lis.OnChange(event)

	wg.Wait()
}

func Test_changeListener_writeToClosedChan(t *testing.T) {
	// 1. setup
	lis := setupChangeListener()
	ch := make(chan *configstores.SubscribeResp)
	// add subscriber
	err := lis.addByTopic(ns, "key1", ch)
	if err != nil {
		t.Error(err)
	}
	// close before write
	close(ch)
	// 3 change
	changes := make(map[string]*storage.ConfigChange)
	c := &storage.ConfigChange{
		OldValue:   "v1",
		NewValue:   "v2",
		ChangeType: storage.MODIFIED,
	}
	changes["key1"] = c
	event := &storage.ChangeEvent{
		Changes: changes,
	}
	event.Namespace = ns
	// execute
	lis.OnChange(event)
	//	 assert no panic
}
