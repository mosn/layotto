// Copyright 2021 Layotto Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hello

import (
	"context"
	"errors"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/atomic"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"mosn.io/layotto/components/pluggable"
	helloproto "mosn.io/layotto/spec/proto/pluggable/v1/hello"
)

var _ helloproto.HelloServer = (*mockServer)(nil)

type mockServer struct {
	helloproto.UnimplementedHelloServer

	initCalled   atomic.Int32
	onInitCalled func(config *helloproto.HelloConfig)
	initError    error

	sayHelloCalled   atomic.Int32
	onSayHelloCalled func(request *helloproto.HelloRequest)
	sayHelloResponse *helloproto.HelloResponse
	sayHelloError    error
}

func (m *mockServer) Init(ctx context.Context, config *helloproto.HelloConfig) (*emptypb.Empty, error) {
	m.initCalled.Add(1)
	if m.onInitCalled != nil {
		m.onInitCalled(config)
	}
	return &emptypb.Empty{}, m.initError
}

func (m *mockServer) SayHello(ctx context.Context, request *helloproto.HelloRequest) (*helloproto.HelloResponse, error) {
	m.sayHelloCalled.Add(1)
	if m.onSayHelloCalled != nil {
		m.onSayHelloCalled(request)
	}
	return m.sayHelloResponse, m.sayHelloError
}

func TestGRPCHelloComponent(t *testing.T) {
	serverFor := pluggable.TestServerFor(helloproto.RegisterHelloServer, func(cc grpc.ClientConnInterface) *grpcHello {
		client := helloproto.NewHelloClient(cc)
		hello := &grpcHello{}
		hello.client = client
		return hello
	})

	socketServerFor := pluggable.TestSocketServerFor(helloproto.RegisterHelloServer, func(dialer pluggable.GRPCConnectionDialer) HelloService {
		return NewGRPCHello(dialer)
	})

	t.Run("test init should call grpc init and the pluggable component should get the config params", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			return
		}

		const mockType = "mock"
		srv := &mockServer{
			onInitCalled: func(config *helloproto.HelloConfig) {
				assert.Equal(t, mockType, config.Type)
			},
		}
		client, cleanup, err := socketServerFor(srv)
		require.NoError(t, err)
		defer cleanup()
		config := &HelloConfig{
			Type: mockType,
		}
		err = client.Init(config)
		assert.Nil(t, err)
		assert.Equal(t, int32(1), srv.initCalled.Load())
	})

	t.Run("init should return an err when dail returns it", func(t *testing.T) {
		serverFor1 := pluggable.TestServerFor(helloproto.RegisterHelloServer, func(cc grpc.ClientConnInterface) *grpcHello {
			client := helloproto.NewHelloClient(cc)
			hello := &grpcHello{}
			hello.dialer = func(ctx context.Context, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
				return nil, errors.New("dial failed")
			}
			hello.client = client
			return hello
		})

		server := &mockServer{
			initError: errors.New("init error"),
		}

		client, cleanup, err := serverFor1(server)
		require.NoError(t, err)
		defer cleanup()
		err = client.Init(&HelloConfig{})
		assert.NotNil(t, err)
		assert.Equal(t, int32(0), server.initCalled.Load())
	})

	t.Run("init should return an err when grpc method returns it", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			return
		}

		srv := &mockServer{
			initError: errors.New("init error"),
		}
		client, cleanup, err := socketServerFor(srv)
		require.NoError(t, err)
		defer cleanup()
		err = client.Init(&HelloConfig{})
		assert.NotNil(t, err)
		assert.Equal(t, int32(1), srv.initCalled.Load())
	})

	t.Run("Hello should return hello response when response is returned from the grpc call", func(t *testing.T) {
		const fakeName = "fakeName"
		const helloString = "hello"
		server := &mockServer{
			onSayHelloCalled: func(request *helloproto.HelloRequest) {
				assert.Equal(t, fakeName, request.GetName())
			},
			sayHelloResponse: &helloproto.HelloResponse{
				HelloString: helloString,
			},
		}

		client, cleanup, err := serverFor(server)
		require.NoError(t, err)
		defer cleanup()
		response, err := client.Hello(context.TODO(), &HelloRequest{Name: fakeName})
		assert.Equal(t, int32(1), server.sayHelloCalled.Load())
		assert.NoError(t, err)
		assert.Equal(t, response, &HelloResponse{HelloString: helloString})
	})

	t.Run("Hello should return an err when grpc method returns it", func(t *testing.T) {
		const fakeName = "fakeName"
		server := &mockServer{
			onSayHelloCalled: func(request *helloproto.HelloRequest) {
				assert.Equal(t, fakeName, request.GetName())
			},
			sayHelloError: errors.New("say hello error"),
		}

		client, cleanup, err := serverFor(server)
		require.NoError(t, err)
		defer cleanup()
		response, err := client.Hello(context.TODO(), &HelloRequest{Name: fakeName})
		assert.Equal(t, int32(1), server.sayHelloCalled.Load())
		assert.NotNil(t, err)
		assert.Equal(t, response, (*HelloResponse)(nil))
	})
}
