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

package mosn

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"mosn.io/layotto/components/rpc"
	"mosn.io/layotto/components/rpc/invoker/mosn/channel"
)

func Test_mosnInvoker_Init(t *testing.T) {
	t.Run("invalid config", func(t *testing.T) {
		invoker := NewMosnInvoker()
		conf := rpc.RpcConfig{
			Config: []byte("invoker"),
		}
		err := invoker.Init(conf)
		assert.NotNil(t, err)
		assert.Equal(t, "invalid character 'i' looking for beginning of value", err.Error())
	})

	t.Run("missing channel config", func(t *testing.T) {
		invoker := NewMosnInvoker()
		conf := rpc.RpcConfig{
			Config: []byte(`{"channel": []}`),
		}
		err := invoker.Init(conf)
		assert.NotNil(t, err)
		assert.Equal(t, "missing channel config", err.Error())
	})

	t.Run("channel not register", func(t *testing.T) {
		invoker := NewMosnInvoker()
		conf := rpc.RpcConfig{
			Config: []byte(`{"channel": [{"protocol":"fake"}]}`),
		}
		err := invoker.Init(conf)
		assert.NotNil(t, err)
		assert.Equal(t, "channel fake not found", err.Error())
	})

	t.Run("success", func(t *testing.T) {
		channel.RegistChannel("fake", func(config channel.ChannelConfig) (rpc.Channel, error) {
			return nil, nil
		})
		invoker := NewMosnInvoker()
		conf := rpc.RpcConfig{
			Config: []byte(`{"channel": [{"protocol":"fake"}]}`),
		}
		err := invoker.Init(conf)
		assert.Nil(t, err)
	})

}

func Test_mosnInvoker_Invoke(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		channel.RegistChannel("fake", func(config channel.ChannelConfig) (rpc.Channel, error) {
			return &fakeChannel{}, nil
		})

		invoker := NewMosnInvoker()
		conf := rpc.RpcConfig{
			Config: []byte(`{"channel": [{"protocol":"fake", "size": 2, "listener": "mosn"}]}`),
		}
		err := invoker.Init(conf)
		assert.Nil(t, err)

		req := &rpc.RPCRequest{
			Ctx:     context.Background(),
			Id:      "1",
			Timeout: 100,
			Method:  "Hello",
			Data:    []byte("hello"),
			Header:  map[string][]string{},
		}
		rsp, err := invoker.Invoke(context.Background(), req)
		assert.Nil(t, err)
		assert.Equal(t, "hello world!", string(rsp.Data))

		req.Header[rpc.RequestTimeoutMs] = []string{"0"}
		req.Timeout = 0
		rsp, err = invoker.Invoke(context.Background(), req)
		assert.Nil(t, err)
		assert.Equal(t, "hello world!", string(rsp.Data))

		assert.Equal(t, int32(3000), req.Timeout)

		req.Header[rpc.RequestTimeoutMs] = []string{"100000"}
		req.Timeout = 0
		rsp, err = invoker.Invoke(context.Background(), req)
		assert.Nil(t, err)
		assert.Equal(t, "hello world!", string(rsp.Data))

		assert.Equal(t, int32(100000), req.Timeout)
	})

	t.Run("panic", func(t *testing.T) {
		invoker := NewMosnInvoker()

		// miss call Init(), invoker.ch will be nil
		req := &rpc.RPCRequest{
			Ctx:     context.Background(),
			Id:      "1",
			Timeout: 100,
			Method:  "Hello",
			Data:    []byte("hello"),
		}
		_, err := invoker.Invoke(context.Background(), req)
		assert.NotNil(t, err)
		assert.Equal(t, "[runtime][rpc]mosn invoker panic: runtime error: invalid memory address or nil pointer dereference", err.Error())
	})
}

type fakeChannel struct {
}

func (c *fakeChannel) Do(req *rpc.RPCRequest) (*rpc.RPCResponse, error) {
	rsp := &rpc.RPCResponse{
		Ctx:         context.Background(),
		Header:      req.Header,
		ContentType: "application/json",
		Data:        append(req.Data, []byte(" world!")...),
	}

	return rsp, nil
}
