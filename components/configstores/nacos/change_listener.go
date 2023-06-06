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
