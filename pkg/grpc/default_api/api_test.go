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

package default_api

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	l8grpc "mosn.io/layotto/pkg/grpc"

	"mosn.io/layotto/components/hello"
	"mosn.io/layotto/components/rpc"
	"mosn.io/layotto/pkg/mock"
	mock_invoker "mosn.io/layotto/pkg/mock/components/invoker"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"

	"time"

	"github.com/golang/protobuf/ptypes/any"
	tmock "github.com/stretchr/testify/mock"
	"google.golang.org/grpc"

	mosninvoker "mosn.io/layotto/components/rpc/invoker/mosn"
)

const (
	maxGRPCServerUptime = 100 * time.Millisecond
	testGRPCServerPort  = 19887
)

type mockGRPCAPI struct {
	API
}

func (m *mockGRPCAPI) SayHello(ctx context.Context, in *runtimev1pb.SayHelloRequest) (*runtimev1pb.SayHelloResponse, error) {
	return &runtimev1pb.SayHelloResponse{}, nil
}

func TestStartServerAPI(t *testing.T) {
	port := testGRPCServerPort
	server := startTestRuntimeAPIServer(port, &mockGRPCAPI{})
	defer server.Stop()
}

func TestSayHello(t *testing.T) {
	t.Run("request ok", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockHello := mock.NewMockHelloService(ctrl)
		api := &api{hellos: map[string]hello.HelloService{
			"mock": mockHello,
		}}
		mockHello.EXPECT().Hello(gomock.Any()).Return(&hello.HelloReponse{
			HelloString: "mock hello",
		}, nil).Times(1)
		resp, err := api.SayHello(context.Background(), &runtimev1pb.SayHelloRequest{
			ServiceName: "mock",
		})
		if err != nil {
			t.Fatalf("say hello request failed: %v", err)
		}
		if resp.Hello != "mock hello" {
			t.Fatalf("say hello response is not expected: %v", resp)
		}
	})

	t.Run("no hello stored", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockHello := mock.NewMockHelloService(ctrl)
		api := &api{hellos: map[string]hello.HelloService{
			"mock": mockHello,
		}}
		_, err := api.SayHello(context.Background(), &runtimev1pb.SayHelloRequest{
			ServiceName: "no register",
		})
		if err != ErrNoInstance {
			t.Fatalf("expected got a no instance error, but got %v", err)
		}
	})

	t.Run("empty say hello", func(t *testing.T) {
		api := &api{hellos: map[string]hello.HelloService{}}
		_, err := api.SayHello(context.Background(), &runtimev1pb.SayHelloRequest{
			ServiceName: "mock",
		})
		if err != ErrNoInstance {
			t.Fatalf("expected got a no instance error, but got %v", err)
		}
	})
}

func startTestRuntimeAPIServer(port int, testAPIServer API) *grpc.Server {
	lis, _ := net.Listen("tcp", fmt.Sprintf(":%d", port))
	opts := []grpc.ServerOption{grpc.WriteBufferSize(1)}

	server := grpc.NewServer(opts...)
	go func() {
		runtimev1pb.RegisterRuntimeServer(server, testAPIServer)
		if err := server.Serve(lis); err != nil {
			panic(err)
		}
	}()

	time.Sleep(maxGRPCServerUptime)

	return server
}

type MockInvoker struct {
	tmock.Mock
}

func (m *MockInvoker) Init(config rpc.RpcConfig) error {
	args := m.Called(config)
	return args.Error(0)
}

func (m *MockInvoker) Invoke(ctx context.Context, req *rpc.RPCRequest) (*rpc.RPCResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*rpc.RPCResponse), args.Error(1)
}

func TestInvokeService(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		resp := &rpc.RPCResponse{
			Header: rpc.RPCHeader{
				"header1": []string{"value1"},
			},
			ContentType: "application/json",
			Data:        []byte("resp data"),
		}

		mockInvoker := mock_invoker.NewMockInvoker(gomock.NewController(t))
		mockInvoker.EXPECT().Invoke(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, req *rpc.RPCRequest) (*rpc.RPCResponse, error) {
				assert.Equal(t, "id1", req.Id)
				assert.Equal(t, "POST", req.Method)
				assert.Equal(t, "application/json", req.ContentType)
				return resp, nil
			})
		httpMethod := int32(runtimev1pb.HTTPExtension_POST)
		in := &runtimev1pb.InvokeServiceRequest{
			Id: "id1",
			Message: &runtimev1pb.CommonInvokeRequest{
				Method:      "POST",
				Data:        &any.Any{},
				ContentType: "application/json",
				HttpExtension: &runtimev1pb.HTTPExtension{
					Verb:        runtimev1pb.HTTPExtension_Verb(httpMethod),
					Querystring: "",
				},
			},
		}

		a := NewAPI(
			"",
			nil,
			nil,
			map[string]rpc.Invoker{
				mosninvoker.Name: mockInvoker,
			},
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
		)

		_, err := a.InvokeService(context.Background(), in)
		assert.Nil(t, err)
	})
}

func createTestClient(port int) *grpc.ClientConn {
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	return conn
}

func TestNewGrpcServer(t *testing.T) {
	apiInterface := &api{}
	_, err := l8grpc.NewGrpcServer(l8grpc.WithGrpcAPIs([]l8grpc.GrpcAPI{apiInterface}), l8grpc.WithNewServer(l8grpc.NewDefaultServer), l8grpc.WithGrpcOptions())
	if err != nil {
		t.Error()
		return
	}
}

func startTestServerAPI(port int, srv runtimev1pb.RuntimeServer) *grpc.Server {
	lis, _ := net.Listen("tcp", fmt.Sprintf(":%d", port))

	server := grpc.NewServer()
	go func() {
		runtimev1pb.RegisterRuntimeServer(server, srv)
		if err := server.Serve(lis); err != nil {
			panic(err)
		}
	}()

	// wait until server starts
	time.Sleep(maxGRPCServerUptime)

	return server
}
