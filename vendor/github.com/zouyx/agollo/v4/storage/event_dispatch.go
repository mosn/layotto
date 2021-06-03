/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package storage

import (
	"errors"
	"fmt"
	"github.com/zouyx/agollo/v4/component/log"
	"regexp"
)

const (
	fmtInvalidKey = "invalid key format for key %s"
)

var (
	//ErrNilListener 为没有找到listener的错误
	ErrNilListener = errors.New("nil listener")
)

// Event generated when any config changes
type Event struct {
	EventType ConfigChangeType
	Key       string
	Value     interface{}
}

// Listener All Listener should implement this Interface
type Listener interface {
	Event(event *Event)
}

//Dispatcher is the observer
type Dispatcher struct {
	listeners map[string][]Listener
}

// UseEventDispatch 用于开启事件分发功能
func UseEventDispatch() *Dispatcher {
	eventDispatch := new(Dispatcher)
	eventDispatch.listeners = make(map[string][]Listener)
	return eventDispatch
}

// RegisterListener 是为某些key注释Listener的方法
func (d *Dispatcher) RegisterListener(listenerObject Listener, keys ...string) error {
	log.Infof("start add  key %v add listener", keys)
	if listenerObject == nil {
		return ErrNilListener
	}

	for _, key := range keys {
		if invalidKey(key) {
			return fmt.Errorf(fmtInvalidKey, key)
		}

		listenerList, ok := d.listeners[key]
		if !ok {
			d.listeners[key] = make([]Listener, 0)
		}

		for _, listener := range listenerList {
			if listener == listenerObject {
				log.Infof("key %s had listener", key)
				return nil
			}
		}
		// append new listener
		listenerList = append(listenerList, listenerObject)
		d.listeners[key] = listenerList
	}
	return nil
}

func invalidKey(key string) bool {
	_, err := regexp.Compile(key)
	if err != nil {
		return true
	}
	return false
}

// UnRegisterListener 用于为某些key注释Listener
func (d *Dispatcher) UnRegisterListener(listenerObj Listener, keys ...string) error {
	if listenerObj == nil {
		return ErrNilListener
	}

	for _, key := range keys {
		listenerList, ok := d.listeners[key]
		if !ok {
			continue
		}

		newListenerList := make([]Listener, 0)
		// remove listener
		for _, listener := range listenerList {
			if listener == listenerObj {
				continue
			}
			newListenerList = append(newListenerList, listener)
		}

		// assign latest listener list
		d.listeners[key] = newListenerList
	}
	return nil
}

//OnChange 实现Apollo的ChangeEvent处理
func (d *Dispatcher) OnChange(changeEvent *ChangeEvent) {
	if changeEvent == nil {
		return
	}
	log.Logger.Infof("get change event for namespace %s", changeEvent.Namespace)
	for key, event := range changeEvent.Changes {
		d.dispatchEvent(key, event)
	}
}

func (d *Dispatcher) OnNewestChange(event *FullChangeEvent) {

}

func (d *Dispatcher) dispatchEvent(eventKey string, event *ConfigChange) {
	for regKey, listenerList := range d.listeners {
		matched, err := regexp.MatchString(regKey, eventKey)
		if err != nil {
			log.Logger.Errorf("regular expression for key %s error %s", eventKey, err)
			continue
		}
		if matched {
			for _, listener := range listenerList {
				log.Logger.Info("event generated for %s key %s", regKey, eventKey)
				go listener.Event(convertToEvent(eventKey, event))
			}
		}
	}
}

func convertToEvent(key string, event *ConfigChange) *Event {
	e := &Event{
		EventType: event.ChangeType,
		Key:       key,
	}
	switch event.ChangeType {
	case ADDED:
		e.Value = event.NewValue
	case MODIFIED:
		e.Value = event.NewValue
	case DELETED:
		e.Value = event.OldValue
	}
	return e
}
