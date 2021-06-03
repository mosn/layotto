package apollo

import (
	"sync"
	"testing"
	"time"

	"github.com/layotto/layotto/components/configstores"
	"github.com/stretchr/testify/assert"
	"github.com/zouyx/agollo/v4/storage"
)

const testAppId = "test_app"

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

const ns = "application"

func setupChangeListener() *changeListener {
	mockRepo := &MockRepo{
		c: NewStore().(*ConfigStore),
	}
	return newChangeListener(mockRepo)
}

//Test modified
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
			assert.Equal(t, c2.StoreName, "apollo")
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
	lis.OnChange(event)
	//	 assert no panic
}
