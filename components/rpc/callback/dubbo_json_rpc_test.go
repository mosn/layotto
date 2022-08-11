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
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"mosn.io/layotto/components/rpc"
)

func Test_beforeFactory_Create(t *testing.T) {
	b := beforeFactory{}
	f := b.Create()
	req := &rpc.RPCRequest{
		Ctx:     context.Background(),
		Id:      "1",
		Timeout: 300,
		Method:  "Hello",
		Header:  make(map[string][]string),
	}
	newReq, err := f(req)
	assert.Nil(t, err)
	assert.Equal(t, "1", newReq.Method)
	assert.Equal(t, "1", newReq.Header.Get("x-services"))
	assert.Equal(t, "Hello", newReq.Header.Get("x-method"))
	assert.Equal(t, "application/json", newReq.Header.Get("content-type"))
	assert.Equal(t, "application/json", newReq.Header.Get("accept"))
}

func Test_beforeFactory_Init(t *testing.T) {
	b := &beforeFactory{}
	err := b.Init(nil)
	assert.Nil(t, err)
}

func Test_beforeFactory_Name(t *testing.T) {
	b := &beforeFactory{}
	assert.Equal(t, "dubbo_json_rpc", b.Name())
}
