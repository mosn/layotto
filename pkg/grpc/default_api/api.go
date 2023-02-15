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
	"errors"
	"sync"

	"github.com/dapr/components-contrib/secretstores"

	"github.com/dapr/components-contrib/bindings"
	"github.com/dapr/components-contrib/pubsub"
	"github.com/dapr/components-contrib/state"
	jsoniter "github.com/json-iterator/go"
	"google.golang.org/grpc"
	"mosn.io/pkg/log"

	"mosn.io/layotto/components/configstores"
	"mosn.io/layotto/components/file"
	"mosn.io/layotto/components/hello"
	"mosn.io/layotto/components/lock"
	"mosn.io/layotto/components/rpc"
	"mosn.io/layotto/components/sequencer"
	grpc_api "mosn.io/layotto/pkg/grpc"
	"mosn.io/layotto/pkg/grpc/dapr"
	dapr_common_v1pb "mosn.io/layotto/pkg/grpc/dapr/proto/common/v1"
	dapr_v1pb "mosn.io/layotto/pkg/grpc/dapr/proto/runtime/v1"
	"mosn.io/layotto/spec/proto/runtime/v1"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
)

const (
	Metadata_key_pubsubName = "pubsubName"
)

var (
	ErrNoInstance = errors.New("no instance found")
	bytesPool     = sync.Pool{
		New: func() interface{} {
			// set size to 100kb
			return new([]byte)
		},
	}
	// FIXME I put it here for compatibility.Don't write singleton like this !
	// LayottoAPISingleton should be refactored and deleted.
	LayottoAPISingleton API
)

type API interface {
	//Layotto Service methods
	runtime.RuntimeServer
	// GrpcAPI related
	grpc_api.GrpcAPI
}

// api is a default implementation for MosnRuntimeServer.
type api struct {
	daprAPI                  dapr.DaprGrpcAPI
	appId                    string
	hellos                   map[string]hello.HelloService
	configStores             map[string]configstores.Store
	rpcs                     map[string]rpc.Invoker
	pubSubs                  map[string]pubsub.PubSub
	stateStores              map[string]state.Store
	transactionalStateStores map[string]state.TransactionalStore
	fileOps                  map[string]file.File
	lockStores               map[string]lock.LockStore
	sequencers               map[string]sequencer.Store
	sendToOutputBindingFn    func(name string, req *bindings.InvokeRequest) (*bindings.InvokeResponse, error)
	secretStores             map[string]secretstores.SecretStore
	// app callback
	AppCallbackConn   *grpc.ClientConn
	topicPerComponent map[string]TopicSubscriptions
	// json
	json jsoniter.API
}

func (a *api) Init(conn *grpc.ClientConn) error {
	// 1. set connection
	a.AppCallbackConn = conn
	return a.startSubscribing()
}

func (a *api) Register(rawGrpcServer *grpc.Server) error {
	LayottoAPISingleton = a
	runtimev1pb.RegisterRuntimeServer(rawGrpcServer, a)
	return nil
}

func NewGrpcAPI(ac *grpc_api.ApplicationContext) grpc_api.GrpcAPI {
	return NewAPI(ac.AppId,
		ac.Hellos, ac.ConfigStores, ac.Rpcs, ac.PubSubs, ac.StateStores, ac.Files, ac.LockStores, ac.Sequencers,
		ac.SendToOutputBindingFn, ac.SecretStores)
}

func NewAPI(
	appId string,
	hellos map[string]hello.HelloService,
	configStores map[string]configstores.Store,
	rpcs map[string]rpc.Invoker,
	pubSubs map[string]pubsub.PubSub,
	stateStores map[string]state.Store,
	files map[string]file.File,
	lockStores map[string]lock.LockStore,
	sequencers map[string]sequencer.Store,
	sendToOutputBindingFn func(name string, req *bindings.InvokeRequest) (*bindings.InvokeResponse, error),
	secretStores map[string]secretstores.SecretStore,
) API {
	// filter out transactionalStateStores
	transactionalStateStores := map[string]state.TransactionalStore{}
	for key, store := range stateStores {
		if state.FeatureTransactional.IsPresent(store.Features()) {
			transactionalStateStores[key] = store.(state.TransactionalStore)
		}
	}
	dAPI := dapr.NewDaprServer(appId, hellos, configStores, rpcs, pubSubs,
		stateStores, transactionalStateStores,
		files, lockStores, sequencers, sendToOutputBindingFn, secretStores)
	// construct
	return &api{
		daprAPI:                  dAPI,
		appId:                    appId,
		hellos:                   hellos,
		configStores:             configStores,
		rpcs:                     rpcs,
		pubSubs:                  pubSubs,
		stateStores:              stateStores,
		transactionalStateStores: transactionalStateStores,
		fileOps:                  files,
		lockStores:               lockStores,
		sequencers:               sequencers,
		sendToOutputBindingFn:    sendToOutputBindingFn,
		secretStores:             secretStores,
		json:                     jsoniter.ConfigFastest,
	}
}

func (a *api) SayHello(ctx context.Context, in *runtimev1pb.SayHelloRequest) (*runtimev1pb.SayHelloResponse, error) {
	h, err := a.getHello(in.ServiceName)
	if err != nil {
		log.DefaultLogger.Errorf("[runtime] [grpc.say_hello] get hello error: %v", err)
		return nil, err
	}
	// create hello request based on pb.go struct
	req := &hello.HelloRequest{
		Name: in.Name,
	}
	resp, err := h.Hello(ctx, req)
	if err != nil {
		log.DefaultLogger.Errorf("[runtime] [grpc.say_hello] request hello error: %v", err)
		return nil, err
	}
	// create response base on hello.Response
	return &runtimev1pb.SayHelloResponse{
		Hello: resp.HelloString,
		Data:  in.Data,
	}, nil

}

func (a *api) getHello(name string) (hello.HelloService, error) {
	if len(a.hellos) == 0 {
		return nil, ErrNoInstance
	}
	h, ok := a.hellos[name]
	if !ok {
		return nil, ErrNoInstance
	}
	return h, nil
}

func (a *api) InvokeService(ctx context.Context, in *runtimev1pb.InvokeServiceRequest) (*runtimev1pb.InvokeResponse, error) {
	// convert request
	var msg *dapr_common_v1pb.InvokeRequest
	if in != nil && in.Message != nil {
		msg = &dapr_common_v1pb.InvokeRequest{
			Method:      in.Message.Method,
			Data:        in.Message.Data,
			ContentType: in.Message.ContentType,
		}
		if in.Message.HttpExtension != nil {
			msg.HttpExtension = &dapr_common_v1pb.HTTPExtension{
				Verb:        dapr_common_v1pb.HTTPExtension_Verb(in.Message.HttpExtension.Verb),
				Querystring: in.Message.HttpExtension.Querystring,
			}
		}
	}
	// delegate to dapr api implementation
	daprResp, err := a.daprAPI.InvokeService(ctx, &dapr_v1pb.InvokeServiceRequest{
		Id:      in.Id,
		Message: msg,
	})
	// handle error
	if err != nil {
		return nil, err
	}

	// convert resp
	return &runtimev1pb.InvokeResponse{
		Data:        daprResp.Data,
		ContentType: daprResp.ContentType,
	}, nil
}
