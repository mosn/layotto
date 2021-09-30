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

package callback

import (
	"encoding/json"

	"mosn.io/layotto/components/rpc"
)

func init() {
	RegisterBeforeInvoke(&beforeFactory{})
}

// beforeFactory is BeforeFactory implement
type beforeFactory struct {
}

func (b *beforeFactory) Name() string {
	return "dubbo_json_rpc"
}

func (b *beforeFactory) Init(json.RawMessage) error {
	return nil
}

// Create is set some header before handle RPCRequest
func (b *beforeFactory) Create() func(*rpc.RPCRequest) (*rpc.RPCRequest, error) {
	return func(request *rpc.RPCRequest) (*rpc.RPCRequest, error) {
		request.Header["x-services"] = []string{request.Id}
		request.Header["x-method"] = []string{request.Method}
		request.Header["content-type"] = []string{"application/json"}
		request.Header["accept"] = []string{"application/json"}

		request.Method = request.Id
		return request, nil
	}
}
