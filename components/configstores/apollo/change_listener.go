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
	"time"

	"github.com/apolloconfig/agollo/v4/storage"
	"mosn.io/pkg/log"

	"mosn.io/layotto/components/configstores"
)

type changeListener struct {
	subscribers *subscriberHolder
	timeout     time.Duration
	store       RepoForListener
}

type RepoForListener interface {
	splitKey(keyWithLabel string) (key string, label string)
	getAllTags(group string, keyWithLabel string) (tags map[string]string, err error)
	GetAppId() string
	GetStoreName() string
}

func newChangeListener(c RepoForListener) *changeListener {
	return &changeListener{
		subscribers: newSubscriberHolder(),
		timeout:     time.Duration(defaultTimeoutWhenResponse) * time.Millisecond,
		store:       c,
	}
}

func (lis *changeListener) OnChange(changeEvent *storage.ChangeEvent) {
	// 1. find related subscribers
	ns := changeEvent.Namespace
	groupLevel := lis.subscribers.findByTopic(ns, "")
	for key, change := range changeEvent.Changes {
		keyLevel := lis.subscribers.findByTopic(ns, key)
		// 2. notice
		for _, s := range groupLevel {
			lis.notify(s, key, change)
		}
		for _, s := range keyLevel {
			lis.notify(s, key, change)
		}
	}
}

func (lis *changeListener) OnNewestChange(event *storage.FullChangeEvent) {
}

func (lis *changeListener) notify(s *subscriber, keyWithLabel string, change *storage.ConfigChange) {
	if s == nil || s.respChan == nil || change == nil {
		return
	}
	// 1 recover panic caused when interacting with the chan
	defer func() {
		if r := recover(); r != nil {
			log.DefaultLogger.Errorf("panic when notify subscriber. %v", r)
			// make sure unused chan are all deleted
			if lis != nil && lis.subscribers != nil {
				lis.subscribers.remove(s)
			}
		}
	}()
	// 2 prepare response
	res := &configstores.SubscribeResp{StoreName: lis.store.GetStoreName(), AppId: lis.store.GetAppId()}
	item := &configstores.ConfigurationItem{}
	item.Group = s.group
	item.Key, item.Label = lis.store.splitKey(keyWithLabel)
	// TODO add a removed flag in response struct.
	if change.ChangeType != storage.DELETED {
		item.Content = change.NewValue.(string)
		tags, err := lis.store.getAllTags(s.group, keyWithLabel)
		if err != nil {
			//	log and ignore
			log.DefaultLogger.Errorf("Error when querying tags in change_listener: %v", err)
		} else {
			item.Tags = tags
		}
	}
	res.Items = append(res.Items, item)

	select {
	// 3 write
	case s.respChan <- res:
		return
	// 4 close chan if timeout
	case <-time.After(lis.timeout):
		// remove for gc
		lis.subscribers.remove(s)
		close(s.respChan)
		return
	}
}

func (lis *changeListener) addByTopic(namespace string, keyWithLabel string, respChan chan *configstores.SubscribeResp) error {
	return lis.subscribers.addByTopic(namespace, keyWithLabel, respChan)
}

func (lis *changeListener) reset() {
	lis.subscribers.reset()
}
