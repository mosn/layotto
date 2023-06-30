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

// StoreConfig wraps configuration for a store implementation
type StoreConfig struct {
	Type      string            `json:"type"`
	StoreName string            `json:"store_name"`
	AppId     string            `json:"app_id"`
	Address   []string          `json:"address"`
	TimeOut   string            `json:"timeout"`
	Metadata  map[string]string `json:"metadata"`
}

// GetRequest is the object describing a get configuration request
type GetRequest struct {
	AppId    string
	Group    string
	Label    string
	Keys     []string
	Metadata map[string]string
}

// SetRequest is the object describing a save configuration request
type SetRequest struct {
	StoreName string
	AppId     string
	Items     []*ConfigurationItem
}

// ConfigurationItem represents a configuration item with key, content and other information.
type ConfigurationItem struct {
	Key      string
	Content  string
	Group    string
	Label    string
	Tags     map[string]string
	Metadata map[string]string
}

// DeleteRequest is the object describing a delete configuration request
type DeleteRequest struct {
	AppId    string
	Group    string
	Label    string
	Keys     []string
	Metadata map[string]string
}

// SubscribeReq is the object describing a subscription request
type SubscribeReq struct {
	AppId    string
	Group    string
	Label    string
	Keys     []string
	Metadata map[string]string
}

// SubscribeResp is the object describing a response for subscription
type SubscribeResp struct {
	StoreName string
	AppId     string
	Items     []*ConfigurationItem
}
