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

package rpc

import (
	"context"
	"encoding/json"
	"strings"
)

const (
	TargetAddress    = "rpc_target_address"
	RequestTimeoutMs = "rpc_request_timeout"
)

const (
	DefaultRequestTimeoutMs = 3000
)

// RPCHeader is storage header info
type RPCHeader map[string][]string

// Range is handle RPCHeader info
func (r RPCHeader) Range(f func(key string, value string) bool) {
	if len(r) == 0 {
		return
	}

	for k, values := range r {
		if ok := f(k, strings.Join(values, ",")); !ok {
			break
		}
	}
}

// Get is get RPCHeader info
func (r RPCHeader) Get(key string) string {
	if r == nil {
		return ""
	}
	values, ok := r[key]
	if !ok {
		return ""
	}
	return strings.Join(values, ",")
}

// RPCRequest is request info
type RPCRequest struct {
	// context
	Ctx context.Context
	// request id
	Id          string
	Timeout     int32
	Method      string
	ContentType string
	Header      RPCHeader
	Data        []byte
}

// RPCResponse is response info
type RPCResponse struct {
	Ctx         context.Context
	Header      RPCHeader
	ContentType string
	Data        []byte
	Success     bool
	Error       error
}

type RpcConfig struct {
	Config json.RawMessage
}

// Invoker is interface for init rpc config or invoke rpc request
type Invoker interface {
	Init(config RpcConfig) error
	Invoke(ctx context.Context, req *RPCRequest) (*RPCResponse, error)
}

// Callback is interface for before invoke or after invoke
type Callback interface {
	// AddBeforeInvoke is add BeforeInvoke func
	AddBeforeInvoke(CallbackFunc)
	// AddAfterInvoke is add AfterInvoke func
	AddAfterInvoke(CallbackFunc)

	// BeforeInvoke is used to invoke beforeInvoke callbacks
	BeforeInvoke(*RPCRequest) (*RPCRequest, error)
	// AfterInvoke is used to invoke afterInvoke callbacks
	AfterInvoke(*RPCResponse) (*RPCResponse, error)
}

// CallbackFunc is Callback implement
type CallbackFunc struct {
	Name   string          `json:"name"`
	Config json.RawMessage `json:"config"`
}

// Channel is handle RPCRequest to RPCResponse
type Channel interface {
	Do(*RPCRequest) (*RPCResponse, error)
}
