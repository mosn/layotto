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

package dapr

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/dapr/components-contrib/bindings"
	"github.com/dapr/components-contrib/state"
	"github.com/golang/mock/gomock"
	"github.com/phayes/freeport"
	"github.com/stretchr/testify/assert"

	grpc_api "mosn.io/layotto/pkg/grpc"
	mock_state "mosn.io/layotto/pkg/mock/components/state"

	"errors"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	dapr_v1pb "mosn.io/layotto/pkg/grpc/dapr/proto/runtime/v1"
)

const (
	maxGRPCServerUptime = 100 * time.Millisecond
)

type MockTxStore struct {
	state.Store
	state.TransactionalStore
}

func (m *MockTxStore) Init(metadata state.Metadata) error {
	return m.Store.Init(metadata)
}

func TestNewDaprAPI_Alpha(t *testing.T) {
	port, _ := freeport.GetFreePort()
	ctrl := gomock.NewController(t)
	mockStore := mock_state.NewMockStore(ctrl)
	mockStore.EXPECT().Features().Return([]state.Feature{state.FeatureTransactional})

	mockTxStore := mock_state.NewMockTransactionalStore(gomock.NewController(t))

	store := &MockTxStore{
		mockStore,
		mockTxStore,
	}
	// construct API
	grpcAPI := NewDaprAPI_Alpha(&grpc_api.ApplicationContext{
		StateStores: map[string]state.Store{"mock": store},
		SendToOutputBindingFn: func(name string, req *bindings.InvokeRequest) (*bindings.InvokeResponse, error) {
			if name == "error-binding" {
				return nil, errors.New("error when invoke binding")
			}
			return &bindings.InvokeResponse{Data: []byte("ok")}, nil
		}})
	err := grpcAPI.Init(nil)
	if err != nil {
		t.Errorf("grpcAPI.Init error")
		return
	}
	// test type assertion
	_, ok := grpcAPI.(dapr_v1pb.DaprServer)
	if !ok {
		t.Errorf("Can not cast grpcAPI to DaprServer")
		return
	}
	srv, ok := grpcAPI.(DaprGrpcAPI)
	if !ok {
		t.Errorf("Can not cast grpcAPI to DaprServer")
		return
	}
	// test invokeBinding
	server := startDaprServerForTest(port, srv)
	defer server.Stop()

	clientConn := createTestClient(port)
	defer clientConn.Close()

	client := dapr_v1pb.NewDaprClient(clientConn)
	_, err = client.InvokeBinding(context.Background(), &dapr_v1pb.InvokeBindingRequest{})
	assert.Nil(t, err)
	_, err = client.InvokeBinding(context.Background(), &dapr_v1pb.InvokeBindingRequest{Name: "error-binding"})
	assert.Equal(t, codes.Internal, status.Code(err))
}

func startDaprServerForTest(port int, srv DaprGrpcAPI) *grpc.Server {
	lis, _ := net.Listen("tcp", fmt.Sprintf(":%d", port))

	server := grpc.NewServer()
	go func() {
		srv.Register(server)
		if err := server.Serve(lis); err != nil {
			panic(err)
		}
	}()

	// wait until server starts
	time.Sleep(maxGRPCServerUptime)

	return server
}

func createTestClient(port int) *grpc.ClientConn {
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	return conn
}
