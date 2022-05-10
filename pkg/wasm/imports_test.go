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

package wasm

import (
	"context"
	"testing"

	"github.com/dapr/components-contrib/state"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	proxywasm "mosn.io/proxy-wasm-go-host/proxywasm/v1"

	"mosn.io/layotto/components/rpc"
	mosninvoker "mosn.io/layotto/components/rpc/invoker/mosn"
	"mosn.io/layotto/pkg/grpc/default_api"
	mock_invoker "mosn.io/layotto/pkg/mock/components/invoker"
	mock_state "mosn.io/layotto/pkg/mock/components/state"
)

func TestImportsHandler(t *testing.T) {
	d := &LayottoHandler{}
	assert.Equal(t, d.Log(proxywasm.LogLevelCritical, "msg"), proxywasm.WasmResultOk)
	assert.Equal(t, d.Log(proxywasm.LogLevelError, "msg"), proxywasm.WasmResultOk)
	assert.Equal(t, d.Log(proxywasm.LogLevelWarn, "msg"), proxywasm.WasmResultOk)
	assert.Equal(t, d.Log(proxywasm.LogLevelInfo, "msg"), proxywasm.WasmResultOk)
	assert.Equal(t, d.Log(proxywasm.LogLevelDebug, "msg"), proxywasm.WasmResultOk)
	assert.Equal(t, d.Log(proxywasm.LogLevelTrace, "msg"), proxywasm.WasmResultOk)
}

func TestGetState(t *testing.T) {
	d := &LayottoHandler{}
	t.Run("normal", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mock_state.NewMockStore(ctrl)
		mockStore.EXPECT().Features().Return(nil)

		compResp := &state.GetResponse{
			Data:     []byte("mock data"),
			Metadata: nil,
		}
		mockStore.EXPECT().Get(gomock.Any()).Return(compResp, nil)
		default_api.LayottoAPISingleton = default_api.NewAPI("", nil, nil, nil, nil, map[string]state.Store{"mock": mockStore}, nil, nil, nil, nil, nil)
		value, ok := d.GetState("mock", "mykey")
		assert.Equal(t, proxywasm.WasmResultOk, ok)
		assert.Equal(t, "mock data", value)
	})
}

func TestInvokeService(t *testing.T) {
	d := &LayottoHandler{}
	t.Run("normal", func(t *testing.T) {
		resp := &rpc.RPCResponse{
			ContentType: "application/json",
			Data:        []byte("100"),
		}

		mockInvoker := mock_invoker.NewMockInvoker(gomock.NewController(t))
		mockInvoker.EXPECT().Invoke(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, req *rpc.RPCRequest) (*rpc.RPCResponse, error) {
				assert.Equal(t, "id_2", req.Id)
				return resp, nil
			})
		default_api.LayottoAPISingleton = default_api.NewAPI("", nil, nil, map[string]rpc.Invoker{mosninvoker.Name: mockInvoker}, nil, nil, nil, nil, nil, nil, nil)
		result, ok := d.InvokeService("id_2", "", "book1")
		assert.Equal(t, proxywasm.WasmResultOk, ok)
		assert.Equal(t, "100", result)
	})
}
