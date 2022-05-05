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
	"testing"

	"github.com/stretchr/testify/assert"

	"mosn.io/layotto/components/pkg/info"
)

type fakeInvoker struct {
}

func (f *fakeInvoker) Init(config RpcConfig) error {
	return nil
}

func (f *fakeInvoker) Invoke(ctx context.Context, req *RPCRequest) (*RPCResponse, error) {
	return nil, nil
}

func TestNewRegistry(t *testing.T) {
	runtimeInfo := info.NewRuntimeInfo()
	registry := NewRegistry(runtimeInfo)
	assert.NotNil(t, registry)
}

func TestNewRpcFactory(t *testing.T) {
	name := "fake"
	fm := func() Invoker {
		return nil
	}
	factory := NewRpcFactory(name, fm)
	assert.NotNil(t, factory)
	assert.Equal(t, "fake", factory.Name)
}

func Test_rpcRegistry_Create(t *testing.T) {
	t.Run("exist", func(t *testing.T) {
		name := "fake"
		fm := func() Invoker {
			return &fakeInvoker{}
		}
		factory := NewRpcFactory(name, fm)
		r := NewRegistry(info.NewRuntimeInfo())
		r.Register(factory)
		invoker, err := r.Create("fake")
		assert.Nil(t, err)
		assert.NotNil(t, invoker)
	})

	t.Run("not exist", func(t *testing.T) {
		r := NewRegistry(info.NewRuntimeInfo())
		invoker, err := r.Create("fake")
		assert.Equal(t, "service component fake is not registered", err.Error())
		assert.Nil(t, invoker)
	})
}

func Test_rpcRegistry_Register(t *testing.T) {
	name := "fake"
	fm := func() Invoker {
		return &fakeInvoker{}
	}
	factory := NewRpcFactory(name, fm)
	runtimeInfo := info.NewRuntimeInfo()
	r := NewRegistry(runtimeInfo)
	r.Register(factory)
	invoker, err := r.Create("fake")
	assert.Nil(t, err)
	assert.NotNil(t, invoker)
}
