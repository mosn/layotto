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

package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"testing"

	"github.com/dapr/components-contrib/pubsub"
	"github.com/dapr/components-contrib/state"
	"github.com/golang/mock/gomock"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	rawGRPC "google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	"mosn.io/pkg/log"

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
	mlock "mosn.io/layotto/pkg/runtime/lock"
	mpubsub "mosn.io/layotto/pkg/runtime/pubsub"
	msequencer "mosn.io/layotto/pkg/runtime/sequencer"
	mstate "mosn.io/layotto/pkg/runtime/state"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
)

func TestNewMosnRuntime(t *testing.T) {
	runtimeConfig := &MosnRuntimeConfig{}
	rt := NewMosnRuntime(runtimeConfig)
	assert.NotNil(t, rt)
}

func TestMosnRuntime_GetInfo(t *testing.T) {
	runtimeConfig := &MosnRuntimeConfig{}
	rt := NewMosnRuntime(runtimeConfig)
	runtimeInfo := rt.GetInfo()
	assert.NotNil(t, runtimeInfo)
}

func TestMosnRuntime_Run(t *testing.T) {
	t.Run("run succ", func(t *testing.T) {
		runtimeConfig := &MosnRuntimeConfig{}
		rt := NewMosnRuntime(runtimeConfig)
		server, err := rt.Run()
		assert.Nil(t, err)
		assert.NotNil(t, server)
	})

	t.Run("no runtime config", func(t *testing.T) {
		rt := NewMosnRuntime(nil)
		_, err := rt.Run()
		assert.NotNil(t, err)
		assert.Equal(t, "[runtime] init error:no runtimeConfig", err.Error())
	})
}

func TestMosnRuntime_initAppCallbackConnection(t *testing.T) {
	t.Run("init success", func(t *testing.T) {
		// prepare app callback grpc server
		port := 8888
		listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%v", port))
		assert.Nil(t, err)
		defer listener.Close()
		svr := rawGRPC.NewServer()
		go func() {
			svr.Serve(listener)
		}()
		cfg := &MosnRuntimeConfig{
			AppManagement: AppConfig{
				AppId:            "",
				GrpcCallbackPort: port,
			},
		}
		// construct MosnRuntime
		m := NewMosnRuntime(cfg)
		// test initAppCallbackConnection
		err = m.initAppCallbackConnection()
		// assert
		assert.Nil(t, err)
	})
}

func TestMosnRuntime_initPubSubs(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		// mock callback response
		subResp := &runtimev1pb.ListTopicSubscriptionsResponse{
			Subscriptions: []*runtimev1pb.TopicSubscription{
				{
					PubsubName: "mock",
					Topic:      "layotto",
					Metadata:   nil,
				},
			},
		}
		// init grpc server
		mockAppCallbackServer := mock_appcallback.NewMockAppCallbackServer(gomock.NewController(t))
		mockAppCallbackServer.EXPECT().ListTopicSubscriptions(gomock.Any(), gomock.Any()).Return(subResp, nil)

		lis := bufconn.Listen(1024 * 1024)
		s := rawGRPC.NewServer()
		runtimev1pb.RegisterAppCallbackServer(s, mockAppCallbackServer)
		go func() {
			s.Serve(lis)
		}()

		// init callback client
		callbackClient, err := rawGRPC.DialContext(context.Background(), "bufnet", rawGRPC.WithInsecure(), rawGRPC.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
			return lis.Dial()
		}))
		assert.Nil(t, err)

		// mock pubsub component
		mockPubSub := mock_pubsub.NewMockPubSub(gomock.NewController(t))
		mockPubSub.EXPECT().Init(gomock.Any()).Return(nil)
		mockPubSub.EXPECT().Subscribe(gomock.Any(), gomock.Any()).Return(nil)
		f := func() pubsub.PubSub {
			return mockPubSub
		}

		cfg := &MosnRuntimeConfig{
			PubSubManagement: map[string]mpubsub.Config{
				"mock": {
					Metadata: map[string]string{
						"target": "layotto",
					},
				},
			},
		}
		// construct MosnRuntime
		m := NewMosnRuntime(cfg)
		m.AppCallbackConn = callbackClient
		m.errInt = func(err error, format string, args ...interface{}) {
			log.DefaultLogger.Errorf("[runtime] occurs an error: "+err.Error()+", "+format, args...)
		}
		// test initPubSubs
		err = m.initPubSubs(mpubsub.NewFactory("mock", f))
		// assert result
		assert.Nil(t, err)
	})
}

func TestMosnRuntime_initStates(t *testing.T) {
	t.Run("init success", func(t *testing.T) {
		// prepare mock
		mockStateStore := mock_state.NewMockStore(gomock.NewController(t))
		mockStateStore.EXPECT().Init(gomock.Any()).Return(nil)
		f := func() state.Store {
			return mockStateStore
		}

		cfg := &MosnRuntimeConfig{
			StateManagement: map[string]mstate.Config{
				"mock": {
					Metadata: map[string]string{
						"target": "layotto",
					},
				},
			},
		}
		// construct MosnRuntime
		m := NewMosnRuntime(cfg)
		m.errInt = func(err error, format string, args ...interface{}) {
			log.DefaultLogger.Errorf("[runtime] occurs an error: "+err.Error()+", "+format, args...)
		}
		// test initStates
		err := m.initStates(mstate.NewFactory("mock", f))
		// assert result
		assert.Nil(t, err)
	})
}

func TestMosnRuntime_initRpc(t *testing.T) {
	t.Run("init success", func(t *testing.T) {
		// prepare mock
		mockInvoker := mock_invoker.NewMockInvoker(gomock.NewController(t))
		mockInvoker.EXPECT().Init(gomock.Any()).Return(nil)
		f := func() rpc.Invoker {
			return mockInvoker
		}

		cfg := &MosnRuntimeConfig{
			RpcManagement: map[string]rpc.RpcConfig{
				"mock": {},
			},
		}
		// construct MosnRuntime
		m := NewMosnRuntime(cfg)
		m.errInt = func(err error, format string, args ...interface{}) {
			log.DefaultLogger.Errorf("[runtime] occurs an error: "+err.Error()+", "+format, args...)
		}
		// test initRpcs method
		err := m.initRpcs(rpc.NewRpcFactory("mock", f))
		// assert
		assert.Nil(t, err)
	})
}

func TestMosnRuntime_initConfigStores(t *testing.T) {
	t.Run("init success", func(t *testing.T) {
		mockStore := mock.NewMockStore(gomock.NewController(t))
		mockStore.EXPECT().Init(gomock.Any()).Return(nil)
		f := func() configstores.Store {
			return mockStore
		}

		cfg := &MosnRuntimeConfig{
			ConfigStoreManagement: map[string]configstores.StoreConfig{
				"mock": {},
			},
		}
		m := NewMosnRuntime(cfg)
		m.errInt = func(err error, format string, args ...interface{}) {
			log.DefaultLogger.Errorf("[runtime] occurs an error: "+err.Error()+", "+format, args...)
		}
		err := m.initConfigStores(configstores.NewStoreFactory("mock", f))
		assert.Nil(t, err)
	})
}

func TestMosnRuntime_initHellos(t *testing.T) {
	t.Run("init success", func(t *testing.T) {
		mockHello := mock.NewMockHelloService(gomock.NewController(t))
		mockHello.EXPECT().Init(gomock.Any()).Return(nil)
		f := func() hello.HelloService {
			return mockHello
		}

		cfg := &MosnRuntimeConfig{
			HelloServiceManagement: map[string]hello.HelloConfig{
				"mock": {},
			},
		}
		m := NewMosnRuntime(cfg)
		m.errInt = func(err error, format string, args ...interface{}) {
			log.DefaultLogger.Errorf("[runtime] occurs an error: "+err.Error()+", "+format, args...)
		}
		err := m.initHellos(hello.NewHelloFactory("mock", f))
		assert.Nil(t, err)
	})
}

func TestMosnRuntime_initSequencers(t *testing.T) {
	t.Run("init success", func(t *testing.T) {
		mockStore := mock_sequencer.NewMockStore(gomock.NewController(t))
		mockStore.EXPECT().Init(gomock.Any()).Return(nil)
		f := func() sequencer.Store {
			return mockStore
		}

		cfg := &MosnRuntimeConfig{
			SequencerManagement: map[string]sequencer.Config{
				"mock": {},
			},
		}
		m := NewMosnRuntime(cfg)
		m.errInt = func(err error, format string, args ...interface{}) {
			log.DefaultLogger.Errorf("[runtime] occurs an error: "+err.Error()+", "+format, args...)
		}
		err := m.initSequencers(msequencer.NewFactory("mock", f))
		assert.Nil(t, err)
	})
}

func TestMosnRuntime_initLocks(t *testing.T) {
	t.Run("init success", func(t *testing.T) {
		mockLockStore := mock_lock.NewMockLockStore(gomock.NewController(t))
		mockLockStore.EXPECT().Init(gomock.Any()).Return(nil)
		f := func() lock.LockStore {
			return mockLockStore
		}

		cfg := &MosnRuntimeConfig{
			LockManagement: map[string]lock.Config{
				"mock": {},
			},
		}
		m := NewMosnRuntime(cfg)
		m.errInt = func(err error, format string, args ...interface{}) {
			log.DefaultLogger.Errorf("[runtime] occurs an error: "+err.Error()+", "+format, args...)
		}
		err := m.initLocks(mlock.NewFactory("mock", f))
		assert.Nil(t, err)
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
		s := rawGRPC.NewServer()
		runtimev1pb.RegisterAppCallbackServer(s, mockAppCallbackServer)
		go func() {
			s.Serve(lis)
		}()

		// init callback client
		callbackClient, err := rawGRPC.DialContext(context.Background(), "bufnet", rawGRPC.WithInsecure(), rawGRPC.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
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

		cfg := &MosnRuntimeConfig{}
		m := NewMosnRuntime(cfg)
		m.errInt = func(err error, format string, args ...interface{}) {
			log.DefaultLogger.Errorf("[runtime] occurs an error: "+err.Error()+", "+format, args...)
		}
		m.AppCallbackConn = callbackClient
		m.json = jsoniter.ConfigFastest
		err = m.publishMessageGRPC(context.Background(), msg)
		assert.Nil(t, err)
	})
}
