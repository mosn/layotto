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
	"strings"

	"github.com/dapr/components-contrib/bindings"
	"github.com/dapr/components-contrib/pubsub"
	"github.com/dapr/components-contrib/secretstores"
	"github.com/dapr/components-contrib/state"
	jsoniter "github.com/json-iterator/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
	"mosn.io/pkg/log"

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
	dapr_v1pb "mosn.io/layotto/pkg/grpc/dapr/proto/runtime/v1"
	"mosn.io/layotto/pkg/messages"
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
	secretStores             map[string]secretstores.SecretStore
	// app callback
	AppCallbackConn   *grpc.ClientConn
	topicPerComponent map[string]TopicSubscriptions
	// json
	json jsoniter.API
}

func (d *daprGrpcAPI) Init(conn *grpc.ClientConn) error {
	// 1. set connection
	d.AppCallbackConn = conn
	return d.startSubscribing()
}

func (d *daprGrpcAPI) Register(rawGrpcServer *grpc.Server) error {
	dapr_v1pb.RegisterDaprServer(rawGrpcServer, d)
	return nil
}

func (d *daprGrpcAPI) InvokeService(ctx context.Context, in *dapr_v1pb.InvokeServiceRequest) (*dapr_common_v1pb.InvokeResponse, error) {
	// 1. convert request to RPCRequest,which is the parameter for RPC components
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

	// 2. route to the specific rpc.Invoker component.
	// Only support mosn component now.
	invoker, ok := d.rpcs[mosninvoker.Name]
	if !ok {
		return nil, errors.New("invoker not init")
	}

	// 3. delegate to the rpc.Invoker component
	resp, err := invoker.Invoke(ctx, req)

	// 4. convert result
	if err != nil {
		return nil, runtime_common.ToGrpcError(err)
	}
	// 5. convert result
	if !resp.Success && resp.Error != nil {
		return nil, runtime_common.ToGrpcError(resp.Error)
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

func (d *daprGrpcAPI) InvokeBinding(ctx context.Context, in *dapr_v1pb.InvokeBindingRequest) (*dapr_v1pb.InvokeBindingResponse, error) {
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

func (d *daprGrpcAPI) isSecretAllowed(storeName string, key string) bool {
	// TODO: add permission control
	return true
}

// NewDaprAPI_Alpha construct a grpc_api.GrpcAPI which implements DaprServer.
// Currently it only support Dapr's InvokeService and InvokeBinding API.
// Note: this feature is still in Alpha state and we don't recommend that you use it in your production environment.
func NewDaprAPI_Alpha(ac *grpc_api.ApplicationContext) grpc_api.GrpcAPI {
	// filter out transactionalStateStores
	transactionalStateStores := map[string]state.TransactionalStore{}
	for key, store := range ac.StateStores {
		if state.FeatureTransactional.IsPresent(store.Features()) {
			transactionalStateStores[key] = store.(state.TransactionalStore)
		}
	}
	return NewDaprServer(ac.AppId,
		ac.Hellos, ac.ConfigStores, ac.Rpcs, ac.PubSubs, ac.StateStores, transactionalStateStores,
		ac.Files, ac.LockStores, ac.Sequencers,
		ac.SendToOutputBindingFn, ac.SecretStores)
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
	secretStores map[string]secretstores.SecretStore,
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
		secretStores:             secretStores,
	}
}
