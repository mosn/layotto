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
	"encoding/json"
	"fmt"
	"github.com/dapr/components-contrib/secretstores"
	moke_secret "mosn.io/layotto/pkg/mock/components/secret"

	"github.com/dapr/components-contrib/bindings"
	"github.com/dapr/components-contrib/pubsub"
	"github.com/dapr/components-contrib/state"
	"github.com/golang/mock/gomock"
	"github.com/phayes/freeport"
	"github.com/stretchr/testify/assert"
	rawGRPC "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/test/bufconn"
	l8grpc "mosn.io/layotto/pkg/grpc"
	"net"
	"testing"

	"errors"
	jsoniter "github.com/json-iterator/go"
	"google.golang.org/grpc/status"
	"io"
	"mosn.io/layotto/components/configstores"
	"mosn.io/layotto/components/hello"
	"mosn.io/layotto/components/lock"
	"mosn.io/layotto/components/rpc"
	"mosn.io/layotto/components/sequencer"
	"mosn.io/layotto/pkg/mock"
	mock_invoker "mosn.io/layotto/pkg/mock/components/invoker"
	mock_lock "mosn.io/layotto/pkg/mock/components/lock"
	mock_pubsub "mosn.io/layotto/pkg/mock/components/pubsub"
	mock_sequencer "mosn.io/layotto/pkg/mock/components/sequencer"
	mock_state "mosn.io/layotto/pkg/mock/components/state"
	mock_appcallback "mosn.io/layotto/pkg/mock/runtime/appcallback"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"

	"time"

	"mosn.io/layotto/components/file"

	"github.com/golang/protobuf/ptypes/any"
	tmock "github.com/stretchr/testify/mock"
	"google.golang.org/grpc"

	mosninvoker "mosn.io/layotto/components/rpc/invoker/mosn"
)

const (
	maxGRPCServerUptime = 100 * time.Millisecond
	testGRPCServerPort  = 19887
)

type MockGrpcServer struct {
	err error
	req *runtimev1pb.SubscribeConfigurationRequest
	grpc.ServerStream
}

func (m *MockGrpcServer) Send(res *runtimev1pb.SubscribeConfigurationResponse) error {
	return nil
}

func (m *MockGrpcServer) Recv() (*runtimev1pb.SubscribeConfigurationRequest, error) {
	return m.req, m.err
}

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

func TestMosnRuntime_publishMessageGRPC(t *testing.T) {
	t.Run("publish success", func(t *testing.T) {
		subResp := &runtimev1pb.TopicEventResponse{
			Status: runtimev1pb.TopicEventResponse_SUCCESS,
		}
		// init grpc server
		mockAppCallbackServer := mock_appcallback.NewMockAppCallbackServer(gomock.NewController(t))
		mockAppCallbackServer.EXPECT().OnTopicEvent(gomock.Any(), gomock.Any()).Return(subResp, nil)

		lis := bufconn.Listen(1024 * 1024)
		s := grpc.NewServer()
		runtimev1pb.RegisterAppCallbackServer(s, mockAppCallbackServer)
		go func() {
			s.Serve(lis)
		}()

		// init callback client
		callbackClient, err := grpc.DialContext(context.Background(), "bufnet", rawGRPC.WithInsecure(), rawGRPC.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
			return lis.Dial()
		}))
		assert.Nil(t, err)

		cloudEvent := map[string]interface{}{
			pubsub.IDField:              "id",
			pubsub.SourceField:          "source",
			pubsub.DataContentTypeField: "content-type",
			pubsub.TypeField:            "type",
			pubsub.SpecVersionField:     "v1.0.0",
			pubsub.DataBase64Field:      "bGF5b3R0bw==",
		}

		data, err := json.Marshal(cloudEvent)
		assert.Nil(t, err)

		msg := &pubsub.NewMessage{
			Data:     data,
			Topic:    "layotto",
			Metadata: make(map[string]string),
		}
		a := NewAPI("", nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)

		var apiForTest = a.(*api)
		//apiForTest.errInt = func(err error, format string, args ...interface{}) {
		//	log.DefaultLogger.Errorf("[runtime] occurs an error: "+err.Error()+", "+format, args...)
		//}
		apiForTest.AppCallbackConn = callbackClient
		apiForTest.json = jsoniter.ConfigFastest
		err = apiForTest.publishMessageGRPC(context.Background(), msg)
		assert.Nil(t, err)
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

func TestGetConfiguration(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockConfigStore := mock.NewMockStore(ctrl)
	api := NewAPI("", nil, map[string]configstores.Store{"mock": mockConfigStore}, nil, nil, nil, nil, nil, nil, nil, nil)
	mockConfigStore.EXPECT().Get(gomock.Any(), gomock.Any()).Return([]*configstores.ConfigurationItem{
		{Key: "sofa", Content: "sofa1"},
	}, nil).Times(1)
	res, err := api.GetConfiguration(context.Background(), &runtimev1pb.GetConfigurationRequest{StoreName: "mock", AppId: "mosn", Keys: []string{"sofa"}})
	assert.Nil(t, err)
	assert.Equal(t, res.Items[0].Key, "sofa")
	assert.Equal(t, res.Items[0].Content, "sofa1")
	_, err = api.GetConfiguration(context.Background(), &runtimev1pb.GetConfigurationRequest{StoreName: "etcd", AppId: "mosn", Keys: []string{"sofa"}})
	assert.Equal(t, err.Error(), "configure store [etcd] don't support now")

}

func TestSaveConfiguration(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockConfigStore := mock.NewMockStore(ctrl)
		mockConfigStore.EXPECT().Set(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, req *configstores.SetRequest) error {
			assert.Equal(t, "appid", req.AppId)
			assert.Equal(t, "mock", req.StoreName)
			assert.Equal(t, 1, len(req.Items))
			return nil
		})
		req := &runtimev1pb.SaveConfigurationRequest{
			StoreName: "mock",
			AppId:     "appid",
			Items: []*runtimev1pb.ConfigurationItem{
				{
					Key:      "key",
					Content:  "value",
					Group:    "  ",
					Label:    "  ",
					Tags:     nil,
					Metadata: nil,
				},
			},
			Metadata: nil,
		}
		api := NewAPI("", nil, map[string]configstores.Store{"mock": mockConfigStore}, nil, nil, nil, nil, nil, nil, nil, nil)
		_, err := api.SaveConfiguration(context.Background(), req)
		assert.Nil(t, err)
	})

	t.Run("unsupport configstore", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockConfigStore := mock.NewMockStore(ctrl)
		api := NewAPI("", nil, map[string]configstores.Store{"mock": mockConfigStore}, nil, nil, nil, nil, nil, nil, nil, nil)
		_, err := api.SaveConfiguration(context.Background(), &runtimev1pb.SaveConfigurationRequest{StoreName: "etcd"})
		assert.Equal(t, err.Error(), "configure store [etcd] don't support now")
	})

}

func TestDeleteConfiguration(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockConfigStore := mock.NewMockStore(ctrl)
		mockConfigStore.EXPECT().Delete(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, req *configstores.DeleteRequest) error {
			assert.Equal(t, "appid", req.AppId)
			assert.Equal(t, 1, len(req.Keys))
			assert.Equal(t, "key", req.Keys[0])
			return nil
		})
		req := &runtimev1pb.DeleteConfigurationRequest{
			StoreName: "mock",
			AppId:     "appid",
			Keys:      []string{"key"},
			Metadata:  nil,
		}
		api := NewAPI("", nil, map[string]configstores.Store{"mock": mockConfigStore}, nil, nil, nil, nil, nil, nil, nil, nil)
		_, err := api.DeleteConfiguration(context.Background(), req)
		assert.Nil(t, err)
	})

	t.Run("unsupport configstore", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockConfigStore := mock.NewMockStore(ctrl)
		api := NewAPI("", nil, map[string]configstores.Store{"mock": mockConfigStore}, nil, nil, nil, nil, nil, nil, nil, nil)
		_, err := api.DeleteConfiguration(context.Background(), &runtimev1pb.DeleteConfigurationRequest{StoreName: "etcd"})
		assert.Equal(t, err.Error(), "configure store [etcd] don't support now")
	})

}

func TestSubscribeConfiguration(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConfigStore := mock.NewMockStore(ctrl)
	api := NewAPI("", nil, map[string]configstores.Store{"mock": mockConfigStore}, nil, nil, nil, nil, nil, nil, nil, nil)

	//test not support store type
	grpcServer := &MockGrpcServer{req: &runtimev1pb.SubscribeConfigurationRequest{}, err: nil}
	err := api.SubscribeConfiguration(grpcServer)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "configure store [] don't support now")

	//test
	grpcServer2 := &MockGrpcServer{req: &runtimev1pb.SubscribeConfigurationRequest{}, err: errors.New("exit")}
	err = api.SubscribeConfiguration(grpcServer2)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "exit")
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

		a := NewAPI("", nil, nil,
			map[string]rpc.Invoker{
				mosninvoker.Name: mockInvoker,
			}, nil, nil, nil, nil, nil, nil, nil)

		_, err := a.InvokeService(context.Background(), in)
		assert.Nil(t, err)
	})
}

func TestPublishEvent(t *testing.T) {
	t.Run("invalid pubsub name", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockPubSub := mock_pubsub.NewMockPubSub(ctrl)
		api := NewAPI("", nil, nil, nil, map[string]pubsub.PubSub{"mock": mockPubSub}, nil, nil, nil, nil, nil, nil)
		_, err := api.PublishEvent(context.Background(), &runtimev1pb.PublishEventRequest{})
		assert.Equal(t, "rpc error: code = InvalidArgument desc = pubsub name is empty", err.Error())
	})

	t.Run("invalid topic", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockPubSub := mock_pubsub.NewMockPubSub(ctrl)
		api := NewAPI("", nil, nil, nil, map[string]pubsub.PubSub{"mock": mockPubSub}, nil, nil, nil, nil, nil, nil)
		req := &runtimev1pb.PublishEventRequest{
			PubsubName: "abc",
		}
		_, err := api.PublishEvent(context.Background(), req)
		assert.Equal(t, "rpc error: code = InvalidArgument desc = topic is empty in pubsub abc", err.Error())
	})

	t.Run("component not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockPubSub := mock_pubsub.NewMockPubSub(ctrl)
		api := NewAPI("", nil, nil, nil, map[string]pubsub.PubSub{"mock": mockPubSub}, nil, nil, nil, nil, nil, nil)
		req := &runtimev1pb.PublishEventRequest{
			PubsubName: "abc",
			Topic:      "abc",
		}
		_, err := api.PublishEvent(context.Background(), req)
		assert.Equal(t, "rpc error: code = InvalidArgument desc = pubsub abc not found", err.Error())
	})

	t.Run("publish success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockPubSub := mock_pubsub.NewMockPubSub(ctrl)
		mockPubSub.EXPECT().Publish(gomock.Any()).Return(nil)
		mockPubSub.EXPECT().Features().Return(nil)
		api := NewAPI("", nil, nil, nil, map[string]pubsub.PubSub{"mock": mockPubSub}, nil, nil, nil, nil, nil, nil)
		req := &runtimev1pb.PublishEventRequest{
			PubsubName: "mock",
			Topic:      "abc",
		}
		_, err := api.PublishEvent(context.Background(), req)
		assert.Nil(t, err)
	})

	t.Run("publish net error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockPubSub := mock_pubsub.NewMockPubSub(ctrl)
		mockPubSub.EXPECT().Publish(gomock.Any()).Return(fmt.Errorf("net error"))
		mockPubSub.EXPECT().Features().Return(nil)
		api := NewAPI("", nil, nil, nil, map[string]pubsub.PubSub{"mock": mockPubSub}, nil, nil, nil, nil, nil, nil)
		req := &runtimev1pb.PublishEventRequest{
			PubsubName: "mock",
			Topic:      "abc",
		}
		_, err := api.PublishEvent(context.Background(), req)
		assert.NotNil(t, err)
		assert.Equal(t, "rpc error: code = Internal desc = error when publish to topic abc in pubsub mock: net error", err.Error())
	})
}

func TestGetBulkState(t *testing.T) {
	t.Run("state store not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mock_state.NewMockStore(ctrl)
		mockStore.EXPECT().Features().Return(nil)
		api := NewAPI("", nil, nil, nil, nil, map[string]state.Store{"mock": mockStore}, nil, nil, nil, nil, nil)
		req := &runtimev1pb.GetBulkStateRequest{
			StoreName: "abc",
		}
		_, err := api.GetBulkState(context.Background(), req)
		assert.Equal(t, "rpc error: code = InvalidArgument desc = state store abc is not found", err.Error())
	})

	t.Run("get state error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mock_state.NewMockStore(ctrl)
		mockStore.EXPECT().Features().Return(nil)
		mockStore.EXPECT().BulkGet(gomock.Any()).Return(false, nil, fmt.Errorf("net error"))
		api := NewAPI("", nil, nil, nil, nil, map[string]state.Store{"mock": mockStore}, nil, nil, nil, nil, nil)
		req := &runtimev1pb.GetBulkStateRequest{
			StoreName: "mock",
			Keys:      []string{"mykey"},
		}
		_, err := api.GetBulkState(context.Background(), req)
		assert.Equal(t, "net error", err.Error())
	})

	t.Run("support bulk get", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mock_state.NewMockStore(ctrl)
		mockStore.EXPECT().Features().Return(nil)

		compResp := []state.BulkGetResponse{
			{
				Data:     []byte("mock data"),
				Metadata: nil,
			},
		}
		mockStore.EXPECT().BulkGet(gomock.Any()).Return(true, compResp, nil)
		api := NewAPI("", nil, nil, nil, nil, map[string]state.Store{"mock": mockStore}, nil, nil, nil, nil, nil)
		req := &runtimev1pb.GetBulkStateRequest{
			StoreName: "mock",
			Keys:      []string{"mykey"},
		}
		rsp, err := api.GetBulkState(context.Background(), req)
		assert.Nil(t, err)
		assert.Equal(t, []byte("mock data"), rsp.GetItems()[0].GetData())
	})

	t.Run("don't support bulk get", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mock_state.NewMockStore(ctrl)
		mockStore.EXPECT().Features().Return(nil)

		resp1 := &state.GetResponse{
			Data:     []byte("mock data"),
			Metadata: nil,
		}

		resp2 := &state.GetResponse{
			Data:     []byte("mock data2"),
			Metadata: nil,
		}
		mockStore.EXPECT().BulkGet(gomock.Any()).Return(false, nil, nil)
		mockStore.EXPECT().Get(gomock.Any()).Return(resp1, nil)
		mockStore.EXPECT().Get(gomock.Any()).Return(resp2, nil)
		api := NewAPI("", nil, nil, nil, nil, map[string]state.Store{"mock": mockStore}, nil, nil, nil, nil, nil)
		req := &runtimev1pb.GetBulkStateRequest{
			StoreName: "mock",
			Keys:      []string{"mykey", "mykey2"},
		}
		rsp, err := api.GetBulkState(context.Background(), req)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(rsp.GetItems()))
	})

}

func TestGetState(t *testing.T) {
	t.Run("state store not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mock_state.NewMockStore(ctrl)
		mockStore.EXPECT().Features().Return(nil)
		api := NewAPI("", nil, nil, nil, nil, map[string]state.Store{"mock": mockStore}, nil, nil, nil, nil, nil)
		req := &runtimev1pb.GetStateRequest{
			StoreName: "abc",
		}
		_, err := api.GetState(context.Background(), req)
		assert.Equal(t, "rpc error: code = InvalidArgument desc = state store abc is not found", err.Error())
	})

	t.Run("state store not configured", func(t *testing.T) {
		api := NewAPI("", nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
		req := &runtimev1pb.GetStateRequest{
			StoreName: "abc",
		}
		_, err := api.GetState(context.Background(), req)
		assert.Equal(t, "rpc error: code = FailedPrecondition desc = state store is not configured", err.Error())
	})

	t.Run("get modified state key error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mock_state.NewMockStore(ctrl)
		mockStore.EXPECT().Features().Return(nil)
		api := NewAPI("", nil, nil, nil, nil, map[string]state.Store{"mock": mockStore}, nil, nil, nil, nil, nil)
		req := &runtimev1pb.GetStateRequest{
			StoreName: "mock",
			Key:       "mykey||abc",
		}
		_, err := api.GetState(context.Background(), req)
		assert.Equal(t, "input key/keyPrefix 'mykey||abc' can't contain '||'", err.Error())
	})

	t.Run("get state error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mock_state.NewMockStore(ctrl)
		mockStore.EXPECT().Features().Return(nil)
		mockStore.EXPECT().Get(gomock.Any()).Return(nil, fmt.Errorf("net error"))
		api := NewAPI("", nil, nil, nil, nil, map[string]state.Store{"mock": mockStore}, nil, nil, nil, nil, nil)
		req := &runtimev1pb.GetStateRequest{
			StoreName: "mock",
			Key:       "mykey",
		}
		_, err := api.GetState(context.Background(), req)
		assert.Equal(t, "rpc error: code = Internal desc = fail to get mykey from state store mock: net error", err.Error())
	})

	t.Run("normal", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mock_state.NewMockStore(ctrl)
		mockStore.EXPECT().Features().Return(nil)

		compResp := &state.GetResponse{
			Data:     []byte("mock data"),
			Metadata: nil,
		}
		mockStore.EXPECT().Get(gomock.Any()).Return(compResp, nil)
		api := NewAPI("", nil, nil, nil, nil, map[string]state.Store{"mock": mockStore}, nil, nil, nil, nil, nil)
		req := &runtimev1pb.GetStateRequest{
			StoreName: "mock",
			Key:       "mykey",
		}
		rsp, err := api.GetState(context.Background(), req)
		assert.Nil(t, err)
		assert.Equal(t, []byte("mock data"), rsp.GetData())
	})

}

func TestSaveState(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mock_state.NewMockStore(ctrl)
		mockStore.EXPECT().Features().Return(nil)
		mockStore.EXPECT().BulkSet(gomock.Any()).DoAndReturn(func(reqs []state.SetRequest) error {
			assert.Equal(t, 1, len(reqs))
			assert.Equal(t, "abc", reqs[0].Key)
			assert.Equal(t, []byte("mock data"), reqs[0].Value)
			return nil
		})
		api := NewAPI("", nil, nil, nil, nil, map[string]state.Store{"mock": mockStore}, nil, nil, nil, nil, nil)
		req := &runtimev1pb.SaveStateRequest{
			StoreName: "mock",
			States: []*runtimev1pb.StateItem{
				{
					Key:   "abc",
					Value: []byte("mock data"),
				},
			},
		}
		_, err := api.SaveState(context.Background(), req)
		assert.Nil(t, err)
	})

	t.Run("save error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mock_state.NewMockStore(ctrl)
		mockStore.EXPECT().Features().Return(nil)
		mockStore.EXPECT().BulkSet(gomock.Any()).Return(fmt.Errorf("net error"))
		api := NewAPI("", nil, nil, nil, nil, map[string]state.Store{"mock": mockStore}, nil, nil, nil, nil, nil)
		req := &runtimev1pb.SaveStateRequest{
			StoreName: "mock",
			States: []*runtimev1pb.StateItem{
				{
					Key:   "abc",
					Value: []byte("mock data"),
				},
			},
		}
		_, err := api.SaveState(context.Background(), req)
		assert.NotNil(t, err)
		assert.Equal(t, "rpc error: code = Internal desc = failed saving state in state store mock: net error", err.Error())
	})
}

func TestDeleteState(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mock_state.NewMockStore(ctrl)
		mockStore.EXPECT().Features().Return(nil)
		mockStore.EXPECT().Delete(gomock.Any()).DoAndReturn(func(req *state.DeleteRequest) error {
			assert.Equal(t, "abc", req.Key)
			return nil
		})
		api := NewAPI("", nil, nil, nil, nil, map[string]state.Store{"mock": mockStore}, nil, nil, nil, nil, nil)
		req := &runtimev1pb.DeleteStateRequest{
			StoreName: "mock",
			Key:       "abc",
		}
		_, err := api.DeleteState(context.Background(), req)
		assert.Nil(t, err)
	})

	t.Run("net error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mock_state.NewMockStore(ctrl)
		mockStore.EXPECT().Features().Return(nil)
		mockStore.EXPECT().Delete(gomock.Any()).Return(fmt.Errorf("net error"))
		api := NewAPI("", nil, nil, nil, nil, map[string]state.Store{"mock": mockStore}, nil, nil, nil, nil, nil)
		req := &runtimev1pb.DeleteStateRequest{
			StoreName: "mock",
			Key:       "abc",
		}
		_, err := api.DeleteState(context.Background(), req)
		assert.NotNil(t, err)
		assert.Equal(t, "rpc error: code = Internal desc = failed deleting state with key abc: net error", err.Error())
	})
}

func TestDeleteBulkState(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mock_state.NewMockStore(ctrl)
		mockStore.EXPECT().Features().Return(nil)
		mockStore.EXPECT().BulkDelete(gomock.Any()).DoAndReturn(func(reqs []state.DeleteRequest) error {
			assert.Equal(t, "abc", reqs[0].Key)
			return nil
		})
		api := NewAPI("", nil, nil, nil, nil, map[string]state.Store{"mock": mockStore}, nil, nil, nil, nil, nil)
		req := &runtimev1pb.DeleteBulkStateRequest{
			StoreName: "mock",
			States: []*runtimev1pb.StateItem{
				{
					Key: "abc",
				},
			},
		}
		_, err := api.DeleteBulkState(context.Background(), req)
		assert.Nil(t, err)
	})

	t.Run("net error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mock_state.NewMockStore(ctrl)
		mockStore.EXPECT().Features().Return(nil)
		mockStore.EXPECT().BulkDelete(gomock.Any()).Return(fmt.Errorf("net error"))
		api := NewAPI("", nil, nil, nil, nil, map[string]state.Store{"mock": mockStore}, nil, nil, nil, nil, nil)
		req := &runtimev1pb.DeleteBulkStateRequest{
			StoreName: "mock",
			States: []*runtimev1pb.StateItem{
				{
					Key: "abc",
				},
			},
		}
		_, err := api.DeleteBulkState(context.Background(), req)
		assert.NotNil(t, err)
		assert.Equal(t, "net error", err.Error())
	})
}

type MockTxStore struct {
	state.Store
	state.TransactionalStore
}

func (m *MockTxStore) Init(metadata state.Metadata) error {
	return m.Store.Init(metadata)
}

func TestExecuteStateTransaction(t *testing.T) {
	t.Run("state store not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mock_state.NewMockStore(ctrl)
		mockStore.EXPECT().Features().Return(nil)
		api := NewAPI("", nil, nil, nil, nil, map[string]state.Store{"mock": mockStore}, nil, nil, nil, nil, nil)
		req := &runtimev1pb.ExecuteStateTransactionRequest{
			StoreName: "abc",
		}
		_, err := api.ExecuteStateTransaction(context.Background(), req)
		assert.Equal(t, "rpc error: code = InvalidArgument desc = state store abc is not found", err.Error())
	})

	t.Run("state store not configured", func(t *testing.T) {
		api := NewAPI("", nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
		req := &runtimev1pb.ExecuteStateTransactionRequest{
			StoreName: "abc",
		}
		_, err := api.ExecuteStateTransaction(context.Background(), req)
		assert.Equal(t, "rpc error: code = FailedPrecondition desc = state store is not configured", err.Error())
	})

	t.Run("normal", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mock_state.NewMockStore(ctrl)
		mockStore.EXPECT().Features().Return([]state.Feature{state.FeatureTransactional})

		mockTxStore := mock_state.NewMockTransactionalStore(gomock.NewController(t))
		mockTxStore.EXPECT().Multi(gomock.Any()).DoAndReturn(func(req *state.TransactionalStateRequest) error {
			assert.Equal(t, 2, len(req.Operations))
			assert.Equal(t, "mosn", req.Metadata["runtime"])
			assert.Equal(t, state.Upsert, req.Operations[0].Operation)
			assert.Equal(t, state.Delete, req.Operations[1].Operation)
			return nil
		})

		store := &MockTxStore{
			mockStore,
			mockTxStore,
		}

		api := NewAPI("", nil, nil, nil, nil, map[string]state.Store{"mock": store}, nil, nil, nil, nil, nil)
		req := &runtimev1pb.ExecuteStateTransactionRequest{
			StoreName: "mock",
			Operations: []*runtimev1pb.TransactionalStateOperation{
				{
					OperationType: string(state.Upsert),
					Request: &runtimev1pb.StateItem{
						Key:   "upsert",
						Value: []byte("mock data"),
					},
				},
				{
					OperationType: string(state.Delete),
					Request: &runtimev1pb.StateItem{
						Key: "delete_abc",
					},
				},
				{
					OperationType: string(state.Delete),
				},
			},
			Metadata: map[string]string{
				"runtime": "mosn",
			},
		}
		_, err := api.ExecuteStateTransaction(context.Background(), req)
		assert.Nil(t, err)
	})

	t.Run("net error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockStore := mock_state.NewMockStore(ctrl)
		mockStore.EXPECT().Features().Return([]state.Feature{state.FeatureTransactional})

		mockTxStore := mock_state.NewMockTransactionalStore(gomock.NewController(t))
		mockTxStore.EXPECT().Multi(gomock.Any()).Return(fmt.Errorf("net error"))

		store := &MockTxStore{
			mockStore,
			mockTxStore,
		}
		api := NewAPI("", nil, nil, nil, nil, map[string]state.Store{"mock": store}, nil, nil, nil, nil, nil)
		req := &runtimev1pb.ExecuteStateTransactionRequest{
			StoreName: "mock",
			Operations: []*runtimev1pb.TransactionalStateOperation{
				{
					OperationType: string(state.Upsert),
					Request: &runtimev1pb.StateItem{
						Key:   "upsert",
						Value: []byte("mock data"),
					},
				},
				{
					OperationType: string(state.Delete),
					Request: &runtimev1pb.StateItem{
						Key: "delete_abc",
					},
				},
				{
					OperationType: string(state.Delete),
				},
			},
			Metadata: map[string]string{
				"runtime": "mosn",
			},
		}
		_, err := api.ExecuteStateTransaction(context.Background(), req)
		assert.NotNil(t, err)
		assert.Equal(t, "rpc error: code = Internal desc = error while executing state transaction: net error", err.Error())
	})
}

func TestTryLock(t *testing.T) {
	t.Run("lock store not configured", func(t *testing.T) {
		api := NewAPI("", nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
		req := &runtimev1pb.TryLockRequest{
			StoreName: "abc",
		}
		_, err := api.TryLock(context.Background(), req)
		assert.Equal(t, "rpc error: code = FailedPrecondition desc = lock store is not configured", err.Error())
	})

	t.Run("resourceid empty", func(t *testing.T) {
		mockLockStore := mock_lock.NewMockLockStore(gomock.NewController(t))
		api := NewAPI("", nil, nil, nil, nil, nil, nil, map[string]lock.LockStore{"mock": mockLockStore}, nil, nil, nil)
		req := &runtimev1pb.TryLockRequest{
			StoreName: "abc",
		}
		_, err := api.TryLock(context.Background(), req)
		assert.Equal(t, "rpc error: code = InvalidArgument desc = ResourceId is empty in lock store abc", err.Error())
	})

	t.Run("lock owner empty", func(t *testing.T) {
		mockLockStore := mock_lock.NewMockLockStore(gomock.NewController(t))
		api := NewAPI("", nil, nil, nil, nil, nil, nil, map[string]lock.LockStore{"mock": mockLockStore}, nil, nil, nil)
		req := &runtimev1pb.TryLockRequest{
			StoreName:  "abc",
			ResourceId: "resource",
		}
		_, err := api.TryLock(context.Background(), req)
		assert.Equal(t, "rpc error: code = InvalidArgument desc = LockOwner is empty in lock store abc", err.Error())
	})

	t.Run("lock expire is not positive", func(t *testing.T) {
		mockLockStore := mock_lock.NewMockLockStore(gomock.NewController(t))
		api := NewAPI("", nil, nil, nil, nil, nil, nil, map[string]lock.LockStore{"mock": mockLockStore}, nil, nil, nil)
		req := &runtimev1pb.TryLockRequest{
			StoreName:  "abc",
			ResourceId: "resource",
			LockOwner:  "owner",
		}
		_, err := api.TryLock(context.Background(), req)
		assert.Equal(t, "rpc error: code = InvalidArgument desc = Expire is not positive in lock store abc", err.Error())
	})

	t.Run("lock store not found", func(t *testing.T) {
		mockLockStore := mock_lock.NewMockLockStore(gomock.NewController(t))
		api := NewAPI("", nil, nil, nil, nil, nil, nil, map[string]lock.LockStore{"mock": mockLockStore}, nil, nil, nil)
		req := &runtimev1pb.TryLockRequest{
			StoreName:  "abc",
			ResourceId: "resource",
			LockOwner:  "owner",
			Expire:     1,
		}
		_, err := api.TryLock(context.Background(), req)
		assert.Equal(t, "rpc error: code = InvalidArgument desc = lock store abc not found", err.Error())
	})

	t.Run("normal", func(t *testing.T) {
		mockLockStore := mock_lock.NewMockLockStore(gomock.NewController(t))
		mockLockStore.EXPECT().TryLock(gomock.Any()).DoAndReturn(func(req *lock.TryLockRequest) (*lock.TryLockResponse, error) {
			assert.Equal(t, "lock|||resource", req.ResourceId)
			assert.Equal(t, "owner", req.LockOwner)
			assert.Equal(t, int32(1), req.Expire)
			return &lock.TryLockResponse{
				Success: true,
			}, nil
		})
		api := NewAPI("", nil, nil, nil, nil, nil, nil, map[string]lock.LockStore{"mock": mockLockStore}, nil, nil, nil)
		req := &runtimev1pb.TryLockRequest{
			StoreName:  "mock",
			ResourceId: "resource",
			LockOwner:  "owner",
			Expire:     1,
		}
		resp, err := api.TryLock(context.Background(), req)
		assert.Nil(t, err)
		assert.Equal(t, true, resp.Success)
	})

}

func TestUnlock(t *testing.T) {
	t.Run("lock store not configured", func(t *testing.T) {
		api := NewAPI("", nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
		req := &runtimev1pb.UnlockRequest{
			StoreName: "abc",
		}
		_, err := api.Unlock(context.Background(), req)
		assert.Equal(t, "rpc error: code = FailedPrecondition desc = lock store is not configured", err.Error())
	})

	t.Run("resourceid empty", func(t *testing.T) {
		mockLockStore := mock_lock.NewMockLockStore(gomock.NewController(t))
		api := NewAPI("", nil, nil, nil, nil, nil, nil, map[string]lock.LockStore{"mock": mockLockStore}, nil, nil, nil)
		req := &runtimev1pb.UnlockRequest{
			StoreName: "abc",
		}
		_, err := api.Unlock(context.Background(), req)
		assert.Equal(t, "rpc error: code = InvalidArgument desc = ResourceId is empty in lock store abc", err.Error())
	})

	t.Run("lock owner empty", func(t *testing.T) {
		mockLockStore := mock_lock.NewMockLockStore(gomock.NewController(t))
		api := NewAPI("", nil, nil, nil, nil, nil, nil, map[string]lock.LockStore{"mock": mockLockStore}, nil, nil, nil)
		req := &runtimev1pb.UnlockRequest{
			StoreName:  "abc",
			ResourceId: "resource",
		}
		_, err := api.Unlock(context.Background(), req)
		assert.Equal(t, "rpc error: code = InvalidArgument desc = LockOwner is empty in lock store abc", err.Error())
	})

	t.Run("lock store not found", func(t *testing.T) {
		mockLockStore := mock_lock.NewMockLockStore(gomock.NewController(t))
		api := NewAPI("", nil, nil, nil, nil, nil, nil, map[string]lock.LockStore{"mock": mockLockStore}, nil, nil, nil)
		req := &runtimev1pb.UnlockRequest{
			StoreName:  "abc",
			ResourceId: "resource",
			LockOwner:  "owner",
		}
		_, err := api.Unlock(context.Background(), req)
		assert.Equal(t, "rpc error: code = InvalidArgument desc = lock store abc not found", err.Error())
	})

	t.Run("normal", func(t *testing.T) {
		mockLockStore := mock_lock.NewMockLockStore(gomock.NewController(t))
		mockLockStore.EXPECT().Unlock(gomock.Any()).DoAndReturn(func(req *lock.UnlockRequest) (*lock.UnlockResponse, error) {
			assert.Equal(t, "lock|||resource", req.ResourceId)
			assert.Equal(t, "owner", req.LockOwner)
			return &lock.UnlockResponse{
				Status: lock.SUCCESS,
			}, nil
		})
		api := NewAPI("", nil, nil, nil, nil, nil, nil, map[string]lock.LockStore{"mock": mockLockStore}, nil, nil, nil)
		req := &runtimev1pb.UnlockRequest{
			StoreName:  "mock",
			ResourceId: "resource",
			LockOwner:  "owner",
		}
		resp, err := api.Unlock(context.Background(), req)
		assert.Nil(t, err)
		assert.Equal(t, runtimev1pb.UnlockResponse_SUCCESS, resp.Status)
	})
}

func TestGetNextId(t *testing.T) {
	t.Run("sequencers not configured", func(t *testing.T) {
		api := NewAPI("", nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
		req := &runtimev1pb.GetNextIdRequest{
			StoreName: "abc",
		}
		_, err := api.GetNextId(context.Background(), req)
		assert.Equal(t, "rpc error: code = FailedPrecondition desc = Sequencer store is not configured", err.Error())
	})

	t.Run("seq key empty", func(t *testing.T) {
		mockSequencerStore := mock_sequencer.NewMockStore(gomock.NewController(t))
		api := NewAPI("", nil, nil, nil, nil, nil, nil, nil, map[string]sequencer.Store{"mock": mockSequencerStore}, nil, nil)
		req := &runtimev1pb.GetNextIdRequest{
			StoreName: "abc",
		}
		_, err := api.GetNextId(context.Background(), req)
		assert.Equal(t, "rpc error: code = InvalidArgument desc = Key is empty in sequencer store abc", err.Error())
	})

	t.Run("sequencer store not found", func(t *testing.T) {
		mockSequencerStore := mock_sequencer.NewMockStore(gomock.NewController(t))
		api := NewAPI("", nil, nil, nil, nil, nil, nil, nil, map[string]sequencer.Store{"mock": mockSequencerStore}, nil, nil)
		req := &runtimev1pb.GetNextIdRequest{
			StoreName: "abc",
			Key:       "next key",
		}
		_, err := api.GetNextId(context.Background(), req)
		assert.Equal(t, "rpc error: code = InvalidArgument desc = Sequencer store abc not found", err.Error())
	})

	t.Run("auto increment is strong", func(t *testing.T) {
		mockSequencerStore := mock_sequencer.NewMockStore(gomock.NewController(t))
		mockSequencerStore.EXPECT().GetNextId(gomock.Any()).
			DoAndReturn(func(req *sequencer.GetNextIdRequest) (*sequencer.GetNextIdResponse, error) {
				assert.Equal(t, "sequencer|||next key", req.Key)
				assert.Equal(t, sequencer.STRONG, req.Options.AutoIncrement)
				return &sequencer.GetNextIdResponse{
					NextId: 10,
				}, nil
			})
		api := NewAPI("", nil, nil, nil, nil, nil, nil, nil, map[string]sequencer.Store{"mock": mockSequencerStore}, nil, nil)
		req := &runtimev1pb.GetNextIdRequest{
			StoreName: "mock",
			Key:       "next key",
			Options: &runtimev1pb.SequencerOptions{
				Increment: runtimev1pb.SequencerOptions_STRONG,
			},
		}
		rsp, err := api.GetNextId(context.Background(), req)
		assert.Nil(t, err)
		assert.Equal(t, int64(10), rsp.NextId)
	})

	t.Run("net error", func(t *testing.T) {
		mockSequencerStore := mock_sequencer.NewMockStore(gomock.NewController(t))
		mockSequencerStore.EXPECT().GetNextId(gomock.Any()).Return(nil, fmt.Errorf("net error"))
		api := NewAPI("", nil, nil, nil, nil, nil, nil, nil, map[string]sequencer.Store{"mock": mockSequencerStore}, nil, nil)
		req := &runtimev1pb.GetNextIdRequest{
			StoreName: "mock",
			Key:       "next key",
			Options: &runtimev1pb.SequencerOptions{
				Increment: runtimev1pb.SequencerOptions_STRONG,
			},
		}
		_, err := api.GetNextId(context.Background(), req)
		assert.NotNil(t, err)
		assert.Equal(t, "net error", err.Error())
	})
}

func SendData(w net.Conn) {
	w.Write([]byte("testFile"))
	w.Close()
}
func TestGetFile(t *testing.T) {
	r, w := net.Pipe()
	ctrl := gomock.NewController(t)
	mockFile := mock.NewMockFile(ctrl)
	mockStream := mock.NewMockRuntime_GetFileServer(ctrl)
	api := NewAPI("", nil, nil, nil, nil, nil, map[string]file.File{"mock": mockFile}, nil, nil, nil, nil)
	err := api.GetFile(&runtimev1pb.GetFileRequest{StoreName: "mock1"}, mockStream)
	assert.Equal(t, err, status.Errorf(codes.InvalidArgument, "not supported store type: mock1"))
	metadata := make(map[string]string)
	mockFile.EXPECT().Get(context.Background(), &file.GetFileStu{FileName: "", Metadata: metadata}).Return(r, nil).Times(1)
	mockStream.EXPECT().Send(&runtimev1pb.GetFileResponse{Data: []byte("testFile")}).Times(1)
	mockStream.EXPECT().Context().Return(context.Background())
	go SendData(w)
	api.GetFile(&runtimev1pb.GetFileRequest{StoreName: "mock"}, mockStream)
}

func TestPutFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockFile := mock.NewMockFile(ctrl)
	mockStream := mock.NewMockRuntime_PutFileServer(ctrl)
	api := NewAPI("", nil, nil, nil, nil, nil, map[string]file.File{"mock": mockFile}, nil, nil, nil, nil)

	mockStream.EXPECT().Recv().Return(nil, io.EOF).Times(1)
	err := api.PutFile(mockStream)
	assert.Nil(t, err)

	mockStream.EXPECT().Recv().Return(&runtimev1pb.PutFileRequest{StoreName: "mock1"}, nil).Times(1)
	err = api.PutFile(mockStream)
	assert.Equal(t, err, status.Errorf(codes.InvalidArgument, "not support store type: mock1"))

	mockStream.EXPECT().Recv().Return(&runtimev1pb.PutFileRequest{StoreName: "mock"}, nil).Times(1)
	stream := newPutObjectStreamReader(nil, mockStream)
	Metadata := make(map[string]string)
	mockStream.EXPECT().Context().Return(context.Background())
	mockFile.EXPECT().Put(context.Background(), &file.PutFileStu{DataStream: stream, FileName: "", Metadata: Metadata}).Return(errors.New("err occur")).Times(1)
	err = api.PutFile(mockStream)
	s, _ := status.FromError(err)
	assert.Equal(t, s.Message(), "err occur")
}

func TestListFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockFile := mock.NewMockFile(ctrl)
	api := NewAPI("", nil, nil, nil, nil, nil, map[string]file.File{"mock": mockFile}, nil, nil, nil, nil)
	request := &runtimev1pb.FileRequest{StoreName: "mock1"}
	request.Metadata = make(map[string]string)
	resp, err := api.ListFile(context.Background(), &runtimev1pb.ListFileRequest{Request: request})
	assert.Nil(t, resp)
	assert.Equal(t, err, status.Errorf(codes.InvalidArgument, "not support store type: mock1"))
	request = &runtimev1pb.FileRequest{StoreName: "mock", Name: "test"}
	request.Metadata = make(map[string]string)
	mockFile.EXPECT().List(context.Background(), &file.ListRequest{DirectoryName: request.Name, Metadata: request.Metadata}).Return(&file.ListResp{Files: nil, Marker: "hello", IsTruncated: true}, nil).Times(1)
	resp, err = api.ListFile(context.Background(), &runtimev1pb.ListFileRequest{Request: request})
	assert.Nil(t, err)
	assert.Equal(t, resp.Marker, "hello")
	assert.Equal(t, resp.IsTruncated, true)
	mockFile.EXPECT().List(context.Background(), &file.ListRequest{DirectoryName: request.Name, Metadata: request.Metadata}).Return(&file.ListResp{}, errors.New("test fail")).Times(1)
	resp, err = api.ListFile(context.Background(), &runtimev1pb.ListFileRequest{Request: request})
	assert.NotNil(t, err)
	info := &file.FilesInfo{FileName: "hello", Size: 10, LastModified: "2021.11.12"}
	files := make([]*file.FilesInfo, 0)
	files = append(files, info)
	mockFile.EXPECT().List(context.Background(), &file.ListRequest{DirectoryName: request.Name, Metadata: request.Metadata}).Return(&file.ListResp{Files: files}, nil).Times(1)
	resp, err = api.ListFile(context.Background(), &runtimev1pb.ListFileRequest{Request: request})
	assert.Equal(t, len(resp.Files), 1)
	assert.Equal(t, resp.Files[0].FileName, "hello")
}

func TestDelFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockFile := mock.NewMockFile(ctrl)
	api := NewAPI("", nil, nil, nil, nil, nil, map[string]file.File{"mock": mockFile}, nil, nil, nil, nil)
	request := &runtimev1pb.FileRequest{StoreName: "mock1"}
	request.Metadata = make(map[string]string)
	resp, err := api.DelFile(context.Background(), &runtimev1pb.DelFileRequest{Request: request})
	assert.Nil(t, resp)
	assert.Equal(t, err, status.Errorf(codes.InvalidArgument, "not support store type: mock1"))
	request = &runtimev1pb.FileRequest{StoreName: "mock", Name: "test"}
	request.Metadata = make(map[string]string)
	mockFile.EXPECT().Del(context.Background(), &file.DelRequest{FileName: request.Name, Metadata: request.Metadata}).Return(nil).Times(1)
	_, err = api.DelFile(context.Background(), &runtimev1pb.DelFileRequest{Request: request})
	assert.Nil(t, err)
	mockFile.EXPECT().Del(context.Background(), &file.DelRequest{FileName: request.Name, Metadata: request.Metadata}).Return(errors.New("test fail")).Times(1)
	_, err = api.DelFile(context.Background(), &runtimev1pb.DelFileRequest{Request: request})
	assert.NotNil(t, err)
}

func TestGetFileMeta(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockFile := mock.NewMockFile(ctrl)
	api := NewAPI("", nil, nil, nil, nil, nil, map[string]file.File{"mock": mockFile}, nil, nil, nil, nil)
	request := &runtimev1pb.GetFileMetaRequest{Request: nil}
	resp, err := api.GetFileMeta(context.Background(), request)
	assert.Nil(t, resp)
	st, _ := status.FromError(err)
	assert.Equal(t, st.Message(), "request can't be nil")
	request.Request = &runtimev1pb.FileRequest{StoreName: "mock", Name: "test"}
	meta := make(map[string]string)
	re := &file.FileMetaResp{
		Size:         10,
		LastModified: "123",
		Metadata: map[string][]string{
			"test": []string{},
		},
	}
	mockFile.EXPECT().Stat(context.Background(), &file.FileMetaRequest{FileName: request.Request.Name, Metadata: meta}).Return(re, nil).Times(1)
	resp, err = api.GetFileMeta(context.Background(), request)
	assert.Equal(t, resp.LastModified, "123")
	assert.Equal(t, int(resp.Size), 10)
}

func createTestClient(port int) *grpc.ClientConn {
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	return conn
}

func TestGetSecret(t *testing.T) {
	fakeStore := moke_secret.FakeSecretStore{}
	fakeStores := map[string]secretstores.SecretStore{
		"store1": fakeStore,
		"store2": fakeStore,
	}
	testCases := []struct {
		testName         string
		storeName        string
		key              string
		errorExcepted    bool
		expectedResponse string
		expectedError    codes.Code
	}{
		{
			testName:         "Good Key from store",
			storeName:        "store1",
			key:              "good-key",
			errorExcepted:    false,
			expectedResponse: "life is good",
			expectedError:    codes.OK,
		},
		{
			testName:         "error occur with error-key",
			storeName:        "store2",
			key:              "error-key",
			errorExcepted:    true,
			expectedResponse: "null",
			expectedError:    codes.Internal,
		},
	}
	// Setup API server
	fakeAPI := NewAPI("", nil, nil, nil, nil, nil, nil, nil, nil,
		nil, fakeStores)

	// Run test server
	port, _ := freeport.GetFreePort()
	server := startTestRuntimeAPIServer(port, fakeAPI)
	defer server.Stop()

	// Create gRPC test client
	clientConn := createTestClient(port)
	defer clientConn.Close()
	// act
	client := runtimev1pb.NewRuntimeClient(clientConn)
	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			req := &runtimev1pb.GetSecretRequest{
				StoreName: tt.storeName,
				Key:       tt.key,
			}
			resp, err := client.GetSecret(context.Background(), req)

			if !tt.errorExcepted {
				assert.NoError(t, err, "Expected no error")
				assert.Equal(t, resp.Data[tt.key], tt.expectedResponse, "Expected responses to be same")
			} else {
				assert.Error(t, err, "Expected error")
				assert.Equal(t, tt.expectedError, status.Code(err))
			}
		})
	}
}

func TestGetBulkSecret(t *testing.T) {

	fakeStore := moke_secret.FakeSecretStore{}
	fakeStores := map[string]secretstores.SecretStore{
		"store1": fakeStore,
	}
	expectedResponse := "life is good"

	testCases := []struct {
		testName         string
		storeName        string
		key              string
		errorExcepted    bool
		expectedResponse string
		expectedError    codes.Code
	}{
		{
			testName:         "Good Key from unrestricted store",
			storeName:        "store1",
			key:              "good-key",
			errorExcepted:    false,
			expectedResponse: expectedResponse,
		},
	}
	// Setup API server
	fakeAPI := NewAPI("", nil, nil, nil, nil, nil, nil, nil, nil,
		nil, fakeStores)

	// Run test server
	port, _ := freeport.GetFreePort()
	server := startTestRuntimeAPIServer(port, fakeAPI)
	defer server.Stop()

	// Create gRPC test client
	clientConn := createTestClient(port)
	defer clientConn.Close()

	// act
	client := runtimev1pb.NewRuntimeClient(clientConn)

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			req := &runtimev1pb.GetBulkSecretRequest{
				StoreName: tt.storeName,
			}
			resp, err := client.GetBulkSecret(context.Background(), req)

			if !tt.errorExcepted {
				assert.NoError(t, err, "Expected no error")
				assert.Equal(t, resp.Data[tt.key].Secrets[tt.key], tt.expectedResponse, "Expected responses to be same")
			} else {
				assert.Error(t, err, "Expected error")
				assert.Equal(t, tt.expectedError, status.Code(err))
			}
		})
	}

}

func TestNewGrpcServer(t *testing.T) {
	apiInterface := &api{}
	_, err := l8grpc.NewGrpcServer(l8grpc.WithGrpcAPIs([]l8grpc.GrpcAPI{apiInterface}), l8grpc.WithNewServer(l8grpc.NewDefaultServer), l8grpc.WithGrpcOptions())
	if err != nil {
		t.Error()
		return
	}
}

func TestInvokeBinding(t *testing.T) {
	port, _ := freeport.GetFreePort()
	srv := NewAPI("", nil, nil, nil, nil, nil, nil, nil, nil,
		func(name string, req *bindings.InvokeRequest) (*bindings.InvokeResponse, error) {
			if name == "error-binding" {
				return nil, errors.New("error when invoke binding")
			}
			return &bindings.InvokeResponse{Data: []byte("ok")}, nil
		}, nil)
	server := startTestServerAPI(port, srv)
	defer server.Stop()

	clientConn := createTestClient(port)
	defer clientConn.Close()

	client := runtimev1pb.NewRuntimeClient(clientConn)
	_, err := client.InvokeBinding(context.Background(), &runtimev1pb.InvokeBindingRequest{})
	assert.Nil(t, err)
	_, err = client.InvokeBinding(context.Background(), &runtimev1pb.InvokeBindingRequest{Name: "error-binding"})
	assert.Equal(t, codes.Internal, status.Code(err))
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
