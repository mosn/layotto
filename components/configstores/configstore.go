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

package configstores

import (
	"context"
)

const (
	ServiceName = "configStore"
	All         = "*"
)

const (
	AppId = iota
	Group
	Label
	Key
	Tag
)

// Store is an interface to perform operations on config store
type Store interface {
	//Init init the configuration store.
	Init(config *StoreConfig) error

	// GetSpecificKeysValue get specific key value.
	Get(context.Context, *GetRequest) ([]*ConfigurationItem, error)

	// Set saves configuration into configuration store.
	Set(context.Context, *SetRequest) error

	// Delete deletes configuration from configuration store.
	Delete(context.Context, *DeleteRequest) error

	// Subscribe subscribe the configurations updates.
	Subscribe(*SubscribeReq, chan *SubscribeResp) error

	//StopSubscribe stop subs
	StopSubscribe()

	// GetDefaultGroup returns default group.This method will be invoked if a request doesn't specify the group field
	GetDefaultGroup() string

	// GetDefaultLabel returns default label
	GetDefaultLabel() string
}
