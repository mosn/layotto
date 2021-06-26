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
	"testing"

	"github.com/stretchr/testify/assert"

	"mosn.io/layotto/components/rpc"
)

type bf struct{}

func (b *bf) Name() string {
	return "before"
}

func (b *bf) Init(message json.RawMessage) error {
	return nil
}

func (b *bf) Create() func(*rpc.RPCRequest) (*rpc.RPCRequest, error) {
	return func(request *rpc.RPCRequest) (*rpc.RPCRequest, error) {
		request.Data = []byte("before")
		return request, nil
	}
}

type af struct{}

func (a *af) Name() string {
	return "after"
}

func (a *af) Init(message json.RawMessage) error {
	return nil
}

func (a *af) Create() func(*rpc.RPCResponse) (*rpc.RPCResponse, error) {
	return func(response *rpc.RPCResponse) (*rpc.RPCResponse, error) {
		response.Data = []byte("after")
		return response, nil
	}
}

func init() {
	RegisterBeforeInvoke(&bf{})
	RegisterAfterInvoke(&af{})
}

func TestCallback(t *testing.T) {
	cb := NewCallback()

	cb.AddBeforeInvoke(rpc.CallbackFunc{Name: "before"})
	req := &rpc.RPCRequest{}
	cb.BeforeInvoke(req)
	assert.Equal(t, "before", string(req.Data))

	cb.AddAfterInvoke(rpc.CallbackFunc{Name: "after"})
	resp := &rpc.RPCResponse{}
	cb.AfterInvoke(resp)
	assert.Equal(t, "after", string(resp.Data))
}
