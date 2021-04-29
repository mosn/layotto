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

const (
	ADDED ConfigChangeType = iota
	MODIFIED
	DELETED
)

//ChangeListener 监听器
type ChangeListener interface {
	//OnChange 增加变更监控
	OnChange(event *ChangeEvent)

	//OnNewestChange 监控最新变更
	OnNewestChange(event *FullChangeEvent)
}

//config change type
type ConfigChangeType int

//config change event
type baseChangeEvent struct {
	Namespace      string
	NotificationID int64
}

//config change event
type ChangeEvent struct {
	baseChangeEvent
	Changes map[string]*ConfigChange
}

type ConfigChange struct {
	OldValue   interface{}
	NewValue   interface{}
	ChangeType ConfigChangeType
}

// all config change event
type FullChangeEvent struct {
	baseChangeEvent
	Changes map[string]interface{}
}

//create modify config change
func createModifyConfigChange(oldValue interface{}, newValue interface{}) *ConfigChange {
	return &ConfigChange{
		OldValue:   oldValue,
		NewValue:   newValue,
		ChangeType: MODIFIED,
	}
}

//create add config change
func createAddConfigChange(newValue interface{}) *ConfigChange {
	return &ConfigChange{
		NewValue:   newValue,
		ChangeType: ADDED,
	}
}

//create delete config change
func createDeletedConfigChange(oldValue interface{}) *ConfigChange {
	return &ConfigChange{
		OldValue:   oldValue,
		ChangeType: DELETED,
	}
}

//base on changeList create Change event
func createConfigChangeEvent(changes map[string]*ConfigChange, nameSpace string, notificationID int64) *ChangeEvent {
	c := &ChangeEvent{
		Changes: changes,
	}
	c.Namespace = nameSpace
	c.NotificationID = notificationID
	return c
}
