package nacos

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
