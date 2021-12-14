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
	"errors"
	"github.com/dapr/components-contrib/bindings"
	"github.com/dapr/components-contrib/pubsub"
	"github.com/dapr/components-contrib/state"
	jsoniter "github.com/json-iterator/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/emptypb"
	"mosn.io/layotto/components/configstores"
	"mosn.io/layotto/components/file"
	"mosn.io/layotto/components/hello"
	"mosn.io/layotto/components/lock"
	runtime_common "mosn.io/layotto/components/pkg/common"
	"mosn.io/layotto/components/rpc"
	mosninvoker "mosn.io/layotto/components/rpc/invoker/mosn"
	"mosn.io/layotto/components/sequencer"
	grpc_api "mosn.io/layotto/pkg/grpc"
	dapr_common_v1pb "mosn.io/layotto/pkg/grpc/dapr/proto/common/v1"
	"mosn.io/layotto/pkg/grpc/dapr/proto/runtime/v1"
	dapr_v1pb "mosn.io/layotto/pkg/grpc/dapr/proto/runtime/v1"
	"mosn.io/layotto/pkg/messages"
	mgrpc "mosn.io/mosn/pkg/filter/network/grpc"
	"mosn.io/pkg/log"
	"strings"
)

type DaprGrpcAPI interface {
	dapr_v1pb.DaprServer
	grpc_api.GrpcAPI
}

type daprGrpcAPI struct {
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
	// app callback
	AppCallbackConn *grpc.ClientConn
	// json
	json jsoniter.API
}

func (d *daprGrpcAPI) Init(conn *grpc.ClientConn) error {
	// 1. set connection
	d.AppCallbackConn = conn
	return d.startSubscribing()
}

func (d *daprGrpcAPI) startSubscribing() error {
	// TODO
	return nil
}

func (d *daprGrpcAPI) Register(s *grpc.Server, registeredServer mgrpc.RegisteredServer) (mgrpc.RegisteredServer, error) {
	dapr_v1pb.RegisterDaprServer(s, d)
	return registeredServer, nil
}

func (d *daprGrpcAPI) InvokeService(ctx context.Context, in *runtime.InvokeServiceRequest) (*dapr_common_v1pb.InvokeResponse, error) {
	msg := in.GetMessage()
	req := &rpc.RPCRequest{
		Ctx:         ctx,
		Id:          in.GetId(),
		Method:      msg.GetMethod(),
		ContentType: msg.GetContentType(),
		Data:        msg.GetData().GetValue(),
	}
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		req.Header = rpc.RPCHeader(md)
	} else {
		req.Header = rpc.RPCHeader(map[string][]string{})
	}
	if ext := msg.GetHttpExtension(); ext != nil {
		req.Header["verb"] = []string{ext.Verb.String()}
		req.Header["query_string"] = []string{ext.GetQuerystring()}
	}

	invoker, ok := d.rpcs[mosninvoker.Name]
	if !ok {
		return nil, errors.New("invoker not init")
	}

	resp, err := invoker.Invoke(ctx, req)
	if err != nil {
		return nil, runtime_common.ToGrpcError(err)
	}

	if resp.Header != nil {
		header := metadata.Pairs()
		for k, values := range resp.Header {
			// fix https://github.com/mosn/layotto/issues/285
			if strings.EqualFold("content-length", k) {
				continue
			}
			header.Set(k, values...)
		}
		grpc.SetHeader(ctx, header)
	}
	return &dapr_common_v1pb.InvokeResponse{
		ContentType: resp.ContentType,
		Data:        &anypb.Any{Value: resp.Data},
	}, nil
}

func (d *daprGrpcAPI) GetState(ctx context.Context, request *runtime.GetStateRequest) (*runtime.GetStateResponse, error) {
	panic("implement me")
}

func (d *daprGrpcAPI) GetBulkState(ctx context.Context, request *runtime.GetBulkStateRequest) (*runtime.GetBulkStateResponse, error) {
	panic("implement me")
}

func (d *daprGrpcAPI) SaveState(ctx context.Context, request *runtime.SaveStateRequest) (*emptypb.Empty, error) {
	panic("implement me")
}

func (d *daprGrpcAPI) QueryStateAlpha1(ctx context.Context, request *runtime.QueryStateRequest) (*runtime.QueryStateResponse, error) {
	panic("implement me")
}

func (d *daprGrpcAPI) DeleteState(ctx context.Context, request *runtime.DeleteStateRequest) (*emptypb.Empty, error) {
	panic("implement me")
}

func (d *daprGrpcAPI) DeleteBulkState(ctx context.Context, request *runtime.DeleteBulkStateRequest) (*emptypb.Empty, error) {
	panic("implement me")
}

func (d *daprGrpcAPI) ExecuteStateTransaction(ctx context.Context, request *runtime.ExecuteStateTransactionRequest) (*emptypb.Empty, error) {
	panic("implement me")
}

func (d *daprGrpcAPI) PublishEvent(ctx context.Context, request *runtime.PublishEventRequest) (*emptypb.Empty, error) {
	panic("implement me")
}

func (d *daprGrpcAPI) InvokeBinding(ctx context.Context, in *runtime.InvokeBindingRequest) (*runtime.InvokeBindingResponse, error) {
	req := &bindings.InvokeRequest{
		Metadata:  in.Metadata,
		Operation: bindings.OperationKind(in.Operation),
	}
	if in.Data != nil {
		req.Data = in.Data
	}

	r := &dapr_v1pb.InvokeBindingResponse{}
	resp, err := d.sendToOutputBindingFn(in.Name, req)
	if err != nil {
		err = status.Errorf(codes.Internal, messages.ErrInvokeOutputBinding, in.Name, err.Error())
		log.DefaultLogger.Errorf("call out binding fail, err:%+v", err)
		return r, err
	}

	if resp != nil {
		r.Data = resp.Data
		r.Metadata = resp.Metadata
	}
	return r, nil
}

func (d *daprGrpcAPI) GetSecret(ctx context.Context, request *runtime.GetSecretRequest) (*runtime.GetSecretResponse, error) {
	panic("implement me")
}

func (d *daprGrpcAPI) GetBulkSecret(ctx context.Context, request *runtime.GetBulkSecretRequest) (*runtime.GetBulkSecretResponse, error) {
	panic("implement me")
}

func (d *daprGrpcAPI) RegisterActorTimer(ctx context.Context, request *runtime.RegisterActorTimerRequest) (*emptypb.Empty, error) {
	panic("implement me")
}

func (d *daprGrpcAPI) UnregisterActorTimer(ctx context.Context, request *runtime.UnregisterActorTimerRequest) (*emptypb.Empty, error) {
	panic("implement me")
}

func (d *daprGrpcAPI) RegisterActorReminder(ctx context.Context, request *runtime.RegisterActorReminderRequest) (*emptypb.Empty, error) {
	panic("implement me")
}

func (d *daprGrpcAPI) UnregisterActorReminder(ctx context.Context, request *runtime.UnregisterActorReminderRequest) (*emptypb.Empty, error) {
	panic("implement me")
}

func (d *daprGrpcAPI) GetActorState(ctx context.Context, request *runtime.GetActorStateRequest) (*runtime.GetActorStateResponse, error) {
	panic("implement me")
}

func (d *daprGrpcAPI) ExecuteActorStateTransaction(ctx context.Context, request *runtime.ExecuteActorStateTransactionRequest) (*emptypb.Empty, error) {
	panic("implement me")
}

func (d *daprGrpcAPI) InvokeActor(ctx context.Context, request *runtime.InvokeActorRequest) (*runtime.InvokeActorResponse, error) {
	panic("implement me")
}

func (d *daprGrpcAPI) GetConfigurationAlpha1(ctx context.Context, request *runtime.GetConfigurationRequest) (*runtime.GetConfigurationResponse, error) {
	panic("implement me")
}

func (d *daprGrpcAPI) SubscribeConfigurationAlpha1(request *runtime.SubscribeConfigurationRequest, server runtime.Dapr_SubscribeConfigurationAlpha1Server) error {
	panic("implement me")
}

func (d *daprGrpcAPI) GetMetadata(ctx context.Context, empty *emptypb.Empty) (*runtime.GetMetadataResponse, error) {
	panic("implement me")
}

func (d *daprGrpcAPI) SetMetadata(ctx context.Context, request *runtime.SetMetadataRequest) (*emptypb.Empty, error) {
	panic("implement me")
}

func (d *daprGrpcAPI) Shutdown(ctx context.Context, empty *emptypb.Empty) (*emptypb.Empty, error) {
	panic("implement me")
}

// NewDaprAPI_Alpha construct a grpc_api.GrpcAPI which implements DaprServer.
// Currently it only support Dapr's InvokeService and InvokeBinding API.
// Note: this feature is still in Alpha state and we don't recommend that you use it in your production environment.
func NewDaprAPI_Alpha(
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
) grpc_api.GrpcAPI {
	// filter out transactionalStateStores
	transactionalStateStores := map[string]state.TransactionalStore{}
	for key, store := range stateStores {
		if state.FeatureTransactional.IsPresent(store.Features()) {
			transactionalStateStores[key] = store.(state.TransactionalStore)
		}
	}
	return NewDaprServer(appId, hellos, configStores, rpcs, pubSubs,
		stateStores, transactionalStateStores,
		files, lockStores, sequencers, sendToOutputBindingFn)
}

func NewDaprServer(
	appId string,
	hellos map[string]hello.HelloService,
	configStores map[string]configstores.Store,
	rpcs map[string]rpc.Invoker,
	pubSubs map[string]pubsub.PubSub,
	stateStores map[string]state.Store,
	transactionalStateStores map[string]state.TransactionalStore,
	files map[string]file.File,
	lockStores map[string]lock.LockStore,
	sequencers map[string]sequencer.Store,
	sendToOutputBindingFn func(name string, req *bindings.InvokeRequest) (*bindings.InvokeResponse, error),
) DaprGrpcAPI {
	// construct
	return &daprGrpcAPI{
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
		json:                     jsoniter.ConfigFastest,
	}
}
