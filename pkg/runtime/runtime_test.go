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
	"errors"
	"fmt"
	"net"
	"testing"

	aws2 "mosn.io/layotto/components/oss/aws"

	s3ext "mosn.io/layotto/pkg/grpc/extension/s3"

	"mosn.io/layotto/components/oss"

	"github.com/dapr/components-contrib/bindings"
	"google.golang.org/grpc/test/bufconn"

	"mosn.io/layotto/components/custom"
	"mosn.io/layotto/components/hello/helloworld"
	"mosn.io/layotto/components/sequencer"
	sequencer_etcd "mosn.io/layotto/components/sequencer/etcd"
	sequencer_redis "mosn.io/layotto/components/sequencer/redis"
	sequencer_zookeeper "mosn.io/layotto/components/sequencer/zookeeper"
	"mosn.io/layotto/pkg/grpc"
	"mosn.io/layotto/pkg/grpc/default_api"
	mock_appcallback "mosn.io/layotto/pkg/mock/runtime/appcallback"
	mbindings "mosn.io/layotto/pkg/runtime/bindings"
	runtime_sequencer "mosn.io/layotto/pkg/runtime/sequencer"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"

	"github.com/dapr/components-contrib/pubsub"
	"github.com/dapr/components-contrib/state"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	rawGRPC "google.golang.org/grpc"
	"mosn.io/pkg/log"

	"mosn.io/layotto/components/configstores"
	"mosn.io/layotto/components/hello"
	"mosn.io/layotto/components/lock"
	mock_component "mosn.io/layotto/components/pkg/mock"
	"mosn.io/layotto/components/rpc"
	"mosn.io/layotto/pkg/mock"
	mock_invoker "mosn.io/layotto/pkg/mock/components/invoker"
	mock_lock "mosn.io/layotto/pkg/mock/components/lock"
	mock_pubsub "mosn.io/layotto/pkg/mock/components/pubsub"
	mock_sequencer "mosn.io/layotto/pkg/mock/components/sequencer"
	mock_state "mosn.io/layotto/pkg/mock/components/state"
	mlock "mosn.io/layotto/pkg/runtime/lock"
	mpubsub "mosn.io/layotto/pkg/runtime/pubsub"
	mstate "mosn.io/layotto/pkg/runtime/state"
)

func TestNewMosnRuntime(t *testing.T) {
	runtimeConfig := &MosnRuntimeConfig{}
	rt := NewMosnRuntime(runtimeConfig)
	assert.NotNil(t, rt)
	rt.Stop()
}

func TestMosnRuntime_GetInfo(t *testing.T) {
	runtimeConfig := &MosnRuntimeConfig{}
	rt := NewMosnRuntime(runtimeConfig)
	runtimeInfo := rt.GetInfo()
	assert.NotNil(t, runtimeInfo)
	rt.Stop()
}

type superPubsub interface {
	custom.Component
	sayGoodBye() string
}

type superPubsubImpl struct {
	custom.Component
}

func (s *superPubsubImpl) sayGoodBye() string {
	return "good bye!"
}

func newSuperPubsub() custom.Component {
	return &superPubsubImpl{mock_component.NewCustomComponentMock()}
}

type mockGrpcAPI struct {
	comp superPubsub
}

func (m *mockGrpcAPI) Init(conn *rawGRPC.ClientConn) error {
	return nil
}

func (m mockGrpcAPI) Register(rawGrpcServer *rawGRPC.Server) error {
	return nil
}

func (m *mockGrpcAPI) sayGoodBye() string {
	return m.comp.sayGoodBye()
}

func TestMosnRuntime_Run(t *testing.T) {
	t.Run("run succ", func(t *testing.T) {
		runtimeConfig := &MosnRuntimeConfig{}
		rt := NewMosnRuntime(runtimeConfig)
		server, err := rt.Run(
			// register your grpc API here
			WithGrpcAPI(
				default_api.NewGrpcAPI,
			),
		)
		assert.Nil(t, err)
		assert.NotNil(t, server)
		rt.Stop()
	})
	t.Run("run succesfully with initRuntimeStage", func(t *testing.T) {
		runtimeConfig := &MosnRuntimeConfig{}
		rt := NewMosnRuntime(runtimeConfig)
		etcdCustomComponent := mock_component.NewCustomComponentMock()
		compType := "xxx_store"
		compName := "etcd"
		rt.AppendInitRuntimeStage(func(o *runtimeOptions, m *MosnRuntime) error {
			m.SetCustomComponent(compType, compName, etcdCustomComponent)
			return nil
		})
		expect := false
		server, err := rt.Run(
			// register your grpc API here
			WithGrpcAPI(
				default_api.NewGrpcAPI,
				func(ac *grpc.ApplicationContext) grpc.GrpcAPI {
					if ac.CustomComponent[compType][compName] == etcdCustomComponent {
						expect = true
					}
					return &mockGrpcAPI{}
				},
			),
		)
		assert.True(t, expect)
		assert.Nil(t, err)
		assert.NotNil(t, server)
		rt.Stop()
	})
	t.Run("run with initRuntimeStage error", func(t *testing.T) {
		runtimeConfig := &MosnRuntimeConfig{}
		rt := NewMosnRuntime(runtimeConfig)
		rt.AppendInitRuntimeStage(nil)
		var expectErr error = errors.New("expected")
		rt.AppendInitRuntimeStage(func(o *runtimeOptions, m *MosnRuntime) error {
			return expectErr
		})
		_, err := rt.Run(
			// register your grpc API here
			WithGrpcAPI(
				default_api.NewGrpcAPI,
			),
		)
		assert.Equal(t, err, expectErr)
	})

	t.Run("no runtime config", func(t *testing.T) {
		rt := NewMosnRuntime(nil)
		_, err := rt.Run(
			// register your grpc API here
			WithGrpcAPI(
				default_api.NewGrpcAPI,
			),
		)
		assert.NotNil(t, err)
		assert.Equal(t, "[runtime] init error:no runtimeConfig", err.Error())
		rt.Stop()
	})
	t.Run("component init error", func(t *testing.T) {
		// mock pubsub component
		mockPubSub := mock_pubsub.NewMockPubSub(gomock.NewController(t))
		errExpected := errors.New("init error")
		mockPubSub.EXPECT().Init(gomock.Any()).Return(errExpected)
		//mockPubSub.EXPECT().Subscribe(gomock.Any(), gomock.Any()).Return(nil)
		f := func() pubsub.PubSub {
			return mockPubSub
		}

		// 2. construct runtime
		cfg := &MosnRuntimeConfig{
			PubSubManagement: map[string]mpubsub.Config{
				"demo": {
					Type: "mock",
					Metadata: map[string]string{
						"target": "layotto",
					},
				},
			},
		}
		rt := NewMosnRuntime(cfg)

		// 3. Run
		_, err := rt.Run(
			// Hello
			WithHelloFactory(
				hello.NewHelloFactory("helloworld", helloworld.NewHelloWorld),
			),
			// register your grpc API here
			WithGrpcAPI(
				default_api.NewGrpcAPI,
			),
			// PubSub
			WithPubSubFactory(
				mpubsub.NewFactory("mock", f),
			),
			// Sequencer
			WithSequencerFactory(
				runtime_sequencer.NewFactory("etcd", func() sequencer.Store {
					return sequencer_etcd.NewEtcdSequencer(log.DefaultLogger)
				}),
				runtime_sequencer.NewFactory("redis", func() sequencer.Store {
					return sequencer_redis.NewStandaloneRedisSequencer(log.DefaultLogger)
				}),
				runtime_sequencer.NewFactory("zookeeper", func() sequencer.Store {
					return sequencer_zookeeper.NewZookeeperSequencer(log.DefaultLogger)
				}),
			),
		)
		// 4. assert
		assert.NotNil(t, err)
		assert.True(t, err == errExpected)

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
		// mock pubsub component
		mockPubSub := mock_pubsub.NewMockPubSub(gomock.NewController(t))
		mockPubSub.EXPECT().Init(gomock.Any()).Return(nil)
		f := func() pubsub.PubSub {
			return mockPubSub
		}

		cfg := &MosnRuntimeConfig{
			PubSubManagement: map[string]mpubsub.Config{
				"demo": {
					Type: "mock",
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
		// test initPubSubs
		err := m.initPubSubs(mpubsub.NewFactory("mock", f))
		// assert result
		assert.Nil(t, err)
	})
}

func TestMosnRuntime_initPubSubsNotExistMetadata(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		// mock pubsub component
		mockPubSub := mock_pubsub.NewMockPubSub(gomock.NewController(t))
		mockPubSub.EXPECT().Init(gomock.Any()).Return(nil)
		f := func() pubsub.PubSub {
			return mockPubSub
		}

		cfg := &MosnRuntimeConfig{
			PubSubManagement: map[string]mpubsub.Config{
				"demo": {
					Type: "mock",
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
		// test initPubSubs
		err := m.initPubSubs(mpubsub.NewFactory("mock", f))
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
					Type: "status",
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
		err := m.initStates(mstate.NewFactory("status", f))
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
				"rpc": {},
			},
		}
		// construct MosnRuntime
		m := NewMosnRuntime(cfg)
		m.errInt = func(err error, format string, args ...interface{}) {
			log.DefaultLogger.Errorf("[runtime] occurs an error: "+err.Error()+", "+format, args...)
		}
		// test initRpcs method
		err := m.initRpcs(rpc.NewRpcFactory("rpc", f))
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
				"mock": {
					Type: "store_config",
				},
			},
		}
		m := NewMosnRuntime(cfg)
		m.errInt = func(err error, format string, args ...interface{}) {
			log.DefaultLogger.Errorf("[runtime] occurs an error: "+err.Error()+", "+format, args...)
		}
		err := m.initConfigStores(configstores.NewStoreFactory("store_config", f))
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
				"mock": {
					Type: "hello",
				},
			},
		}
		m := NewMosnRuntime(cfg)
		m.errInt = func(err error, format string, args ...interface{}) {
			log.DefaultLogger.Errorf("[runtime] occurs an error: "+err.Error()+", "+format, args...)
		}
		err := m.initHellos(hello.NewHelloFactory("hello", f))
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
				"mock": {
					Type: "sequencers",
				},
			},
		}
		m := NewMosnRuntime(cfg)
		m.errInt = func(err error, format string, args ...interface{}) {
			log.DefaultLogger.Errorf("[runtime] occurs an error: "+err.Error()+", "+format, args...)
		}
		err := m.initSequencers(runtime_sequencer.NewFactory("sequencers", f))
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
				"mock": {
					Type: "lock",
				},
			},
		}
		m := NewMosnRuntime(cfg)
		m.errInt = func(err error, format string, args ...interface{}) {
			log.DefaultLogger.Errorf("[runtime] occurs an error: "+err.Error()+", "+format, args...)
		}
		err := m.initLocks(mlock.NewFactory("lock", f))
		assert.Nil(t, err)
	})
}

type MockBindings struct {
}

func (m *MockBindings) Init(metadata bindings.Metadata) error {
	//do nothing
	return nil
}

func (m *MockBindings) Invoke(req *bindings.InvokeRequest) (*bindings.InvokeResponse, error) {
	return nil, nil
}

func (m *MockBindings) Operations() []bindings.OperationKind {
	return nil
}

func TestMosnRuntime_initOutputBinding(t *testing.T) {
	cfg := &MosnRuntimeConfig{}
	m := NewMosnRuntime(cfg)
	assert.Nil(t, m.outputBindings["mockOutbindings"])

	registry := mbindings.NewOutputBindingFactory("mock_outbindings", func() bindings.OutputBinding {
		return &MockBindings{}
	})
	mdata := make(map[string]string)
	m.RuntimeConfig().Bindings = make(map[string]mbindings.Metadata)
	m.runtimeConfig.Bindings["mockOutbindings"] = mbindings.Metadata{
		Type:     "mock_outbindings",
		Metadata: mdata,
	}
	err := m.initOutputBinding(registry)
	assert.Nil(t, err)
	assert.NotNil(t, m.outputBindings["mockOutbindings"])
}

func TestMosnRuntime_runWithCustomComponentAndAPI(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		kind := "super_pubsub"
		compName := "demo"
		compType := "etcd"
		// 1. construct config
		cfg := &MosnRuntimeConfig{
			CustomComponent: map[string]map[string]custom.Config{
				kind: {
					compName: custom.Config{
						Type:     compType,
						Version:  "",
						Metadata: nil,
					},
				},
			},
		}
		// 2. construct runtime
		rt := NewMosnRuntime(cfg)
		var customAPI *mockGrpcAPI
		// 3. Run
		server, err := rt.Run(
			WithErrInterceptor(func(err error, format string, args ...interface{}) {
				panic(err)
			}),
			// register your grpc API here
			WithGrpcAPI(
				default_api.NewGrpcAPI,
				func(ac *grpc.ApplicationContext) grpc.GrpcAPI {
					comp := ac.CustomComponent[kind][compName].(superPubsub)
					customAPI = &mockGrpcAPI{comp: comp}
					return customAPI
				},
			),
			// Custom components
			WithCustomComponentFactory(kind,
				custom.NewComponentFactory(compType, newSuperPubsub),
			),
			// Hello
			WithHelloFactory(
				hello.NewHelloFactory("helloworld", helloworld.NewHelloWorld),
			),
			// Sequencer
			WithSequencerFactory(
				runtime_sequencer.NewFactory(compType, func() sequencer.Store {
					return sequencer_etcd.NewEtcdSequencer(log.DefaultLogger)
				}),
				runtime_sequencer.NewFactory("redis", func() sequencer.Store {
					return sequencer_redis.NewStandaloneRedisSequencer(log.DefaultLogger)
				}),
				runtime_sequencer.NewFactory("zookeeper", func() sequencer.Store {
					return sequencer_zookeeper.NewZookeeperSequencer(log.DefaultLogger)
				}),
			),
		)
		// 4. assert
		assert.Nil(t, err)
		assert.NotNil(t, server)
		// 5. invoke customAPI
		bye := customAPI.sayGoodBye()
		assert.Equal(t, bye, "good bye!")
		// 6. stop
		rt.Stop()
	})
}

func TestMosnRuntime_runWithPubsub(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		// mock pubsub component
		mockPubSub := mock_pubsub.NewMockPubSub(gomock.NewController(t))
		mockPubSub.EXPECT().Init(gomock.Any()).Return(nil)
		mockPubSub.EXPECT().Subscribe(gomock.Any(), gomock.Any()).Return(nil)
		f := func() pubsub.PubSub {
			return mockPubSub
		}

		// 2. construct runtime
		rt, _ := runtimeWithCallbackConnection(t)

		// 3. Run
		server, err := rt.Run(
			WithErrInterceptor(func(err error, format string, args ...interface{}) {
				panic(err)
			}),
			// Hello
			WithHelloFactory(
				hello.NewHelloFactory("helloworld", helloworld.NewHelloWorld),
			),
			// register your grpc API here
			WithGrpcAPI(
				default_api.NewGrpcAPI,
			),
			// PubSub
			WithPubSubFactory(
				mpubsub.NewFactory("mock", f),
			),
			// Sequencer
			WithSequencerFactory(
				runtime_sequencer.NewFactory("etcd", func() sequencer.Store {
					return sequencer_etcd.NewEtcdSequencer(log.DefaultLogger)
				}),
				runtime_sequencer.NewFactory("redis", func() sequencer.Store {
					return sequencer_redis.NewStandaloneRedisSequencer(log.DefaultLogger)
				}),
				runtime_sequencer.NewFactory("zookeeper", func() sequencer.Store {
					return sequencer_zookeeper.NewZookeeperSequencer(log.DefaultLogger)
				}),
			),
		)
		// 4. assert
		assert.Nil(t, err)
		assert.NotNil(t, server)

		// 5. stop
		rt.Stop()
	})

	t.Run("init_with_callback", func(t *testing.T) {
		cloudEvent := constructCloudEvent()
		data, err := json.Marshal(cloudEvent)
		assert.Nil(t, err)
		// mock pubsub component
		mockPubSub := mock_pubsub.NewMockPubSub(gomock.NewController(t))
		mockPubSub.EXPECT().Init(gomock.Any()).Return(nil)
		mockPubSub.EXPECT().Subscribe(gomock.Any(), gomock.Any()).DoAndReturn(func(req pubsub.SubscribeRequest, handler pubsub.Handler) error {
			if req.Topic == "layotto" {
				return handler(context.Background(), &pubsub.NewMessage{
					Data:     data,
					Topic:    "layotto",
					Metadata: nil,
				})
			}
			return nil
		})
		f := func() pubsub.PubSub {
			return mockPubSub
		}

		// 2. construct runtime
		rt, mockAppCallbackServer := runtimeWithCallbackConnection(t)

		topicResp := &runtimev1pb.TopicEventResponse{Status: runtimev1pb.TopicEventResponse_SUCCESS}
		mockAppCallbackServer.EXPECT().OnTopicEvent(gomock.Any(), gomock.Any()).Return(topicResp, nil)
		// 3. Run
		server, err := rt.Run(
			// register your grpc API here
			WithGrpcAPI(
				default_api.NewGrpcAPI,
			),
			// PubSub
			WithPubSubFactory(
				mpubsub.NewFactory("mock", f),
			),
		)
		// 4. assert
		assert.Nil(t, err)
		assert.NotNil(t, server)

		// 5. stop
		rt.Stop()
	})

	t.Run("callback_fail_then_retry", func(t *testing.T) {
		cloudEvent := constructCloudEvent()
		data, err := json.Marshal(cloudEvent)
		assert.Nil(t, err)
		// mock pubsub component
		mockPubSub := mock_pubsub.NewMockPubSub(gomock.NewController(t))
		mockPubSub.EXPECT().Init(gomock.Any()).Return(nil)
		mockPubSub.EXPECT().Subscribe(gomock.Any(), gomock.Any()).DoAndReturn(func(req pubsub.SubscribeRequest, handler pubsub.Handler) error {
			if req.Topic == "layotto" {
				err := handler(context.Background(), &pubsub.NewMessage{
					Data:     data,
					Topic:    "layotto",
					Metadata: nil,
				})
				assert.NotNil(t, err)
				return nil
			}
			return nil
		})
		f := func() pubsub.PubSub {
			return mockPubSub
		}

		// 2. construct runtime
		rt, mockAppCallbackServer := runtimeWithCallbackConnection(t)

		topicResp := &runtimev1pb.TopicEventResponse{Status: runtimev1pb.TopicEventResponse_RETRY}
		mockAppCallbackServer.EXPECT().OnTopicEvent(gomock.Any(), gomock.Any()).Return(topicResp, nil)
		// 3. Run
		server, err := rt.Run(
			// register your grpc API here
			WithGrpcAPI(
				default_api.NewGrpcAPI,
			),
			// PubSub
			WithPubSubFactory(
				mpubsub.NewFactory("mock", f),
			),
		)
		// 4. assert
		assert.Nil(t, err)
		assert.NotNil(t, server)

		// 5. stop
		rt.Stop()
	})

	t.Run("callback_drop", func(t *testing.T) {
		cloudEvent := constructCloudEvent()
		data, err := json.Marshal(cloudEvent)
		assert.Nil(t, err)
		// mock pubsub component
		mockPubSub := mock_pubsub.NewMockPubSub(gomock.NewController(t))
		mockPubSub.EXPECT().Init(gomock.Any()).Return(nil)
		mockPubSub.EXPECT().Subscribe(gomock.Any(), gomock.Any()).DoAndReturn(func(req pubsub.SubscribeRequest, handler pubsub.Handler) error {
			if req.Topic == "layotto" {
				err := handler(context.Background(), &pubsub.NewMessage{
					Data:     data,
					Topic:    "layotto",
					Metadata: nil,
				})
				assert.Nil(t, err)
				return nil
			}
			return nil
		})
		f := func() pubsub.PubSub {
			return mockPubSub
		}

		// 2. construct runtime
		rt, mockAppCallbackServer := runtimeWithCallbackConnection(t)

		topicResp := &runtimev1pb.TopicEventResponse{Status: runtimev1pb.TopicEventResponse_DROP}
		mockAppCallbackServer.EXPECT().OnTopicEvent(gomock.Any(), gomock.Any()).Return(topicResp, nil)
		// 3. Run
		server, err := rt.Run(
			// register your grpc API here
			WithGrpcAPI(
				default_api.NewGrpcAPI,
			),
			// PubSub
			WithPubSubFactory(
				mpubsub.NewFactory("mock", f),
			),
		)
		// 4. assert
		assert.Nil(t, err)
		assert.NotNil(t, server)

		// 5. stop
		rt.Stop()
	})
}

func constructCloudEvent() map[string]interface{} {
	cloudEvent := make(map[string]interface{})
	cloudEvent[pubsub.IDField] = "1"
	cloudEvent[pubsub.SpecVersionField] = "1"
	cloudEvent[pubsub.SourceField] = "adsdafdas"
	cloudEvent[pubsub.DataContentTypeField] = "application/json"
	cloudEvent[pubsub.TypeField] = "adsdafdas"
	return cloudEvent
}

func runtimeWithCallbackConnection(t *testing.T) (*MosnRuntime, *mock_appcallback.MockAppCallbackServer) {
	compName := "demo"
	compType := "mock"
	// 1. prepare callback
	// mock callback response
	subResp := &runtimev1pb.ListTopicSubscriptionsResponse{
		Subscriptions: []*runtimev1pb.TopicSubscription{
			{
				PubsubName: compName,
				Topic:      "layotto",
				Metadata:   nil,
			},
		},
	}
	// init grpc server for callback
	mockAppCallbackServer := mock_appcallback.NewMockAppCallbackServer(gomock.NewController(t))
	mockAppCallbackServer.EXPECT().ListTopicSubscriptions(gomock.Any(), gomock.Any()).Return(subResp, nil)

	lis := bufconn.Listen(1024 * 1024)
	s := rawGRPC.NewServer()
	runtimev1pb.RegisterAppCallbackServer(s, mockAppCallbackServer)
	go func() {
		s.Serve(lis)
	}()

	// 2. construct those necessary fields for mosn runtime
	// init callback client
	callbackClient, err := rawGRPC.DialContext(context.Background(), "bufnet", rawGRPC.WithInsecure(), rawGRPC.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
		return lis.Dial()
	}))
	assert.Nil(t, err)

	// 3. construct mosn runtime
	cfg := &MosnRuntimeConfig{
		PubSubManagement: map[string]mpubsub.Config{
			compName: {
				Type: compType,
				Metadata: map[string]string{
					"target": "layotto",
				},
			},
		},
	}
	rt := NewMosnRuntime(cfg)
	rt.AppCallbackConn = callbackClient
	return rt, mockAppCallbackServer
}

func TestMosnRuntimeWithOssConfig(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		// 1. construct config
		cfg := &MosnRuntimeConfig{
			Oss: map[string]oss.Config{
				"awsdemo": {Type: "aws.oss"},
			},
		}
		// 2. construct runtime
		rt := NewMosnRuntime(cfg)
		// 3. Run
		server, err := rt.Run(
			// register your grpc API here
			WithGrpcAPI(
				default_api.NewGrpcAPI,
				s3ext.NewS3Server,
			),
			WithOssFactory(
				oss.NewFactory("aws.oss", aws2.NewAwsOss),
			),
		)
		// 4. assert
		assert.Equal(t, "invalid argument", err.Error())
		assert.Nil(t, server)
		// 5. stop
		rt.Stop()
	})
}
