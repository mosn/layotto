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

package client

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	layotto_grpc "mosn.io/layotto/pkg/grpc"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
)

// localRuntimeClient implements RuntimeClient interface
// and do in-process message passing
type localRuntimeClient struct {
	api layotto_grpc.API
}

func newLocalRuntimeClient() *localRuntimeClient {
	api := layotto_grpc.NewAPI(appId, hellos, configStores, rpcs, pubSubs, stateStores, files, lockStores, sequencers, sendToOutputBindingFn)
	return &localRuntimeClient{
		api: api,
	}
}

func (l *localRuntimeClient) SayHello(ctx context.Context, in *runtimev1pb.SayHelloRequest, opts ...grpc.CallOption) (*runtimev1pb.SayHelloResponse, error) {
	return l.api.SayHello(ctx, in)
}

func (l *localRuntimeClient) InvokeService(ctx context.Context, in *runtimev1pb.InvokeServiceRequest, opts ...grpc.CallOption) (*runtimev1pb.InvokeResponse, error) {
	return l.api.InvokeService(ctx, in)
}

func (l *localRuntimeClient) GetConfiguration(ctx context.Context, in *runtimev1pb.GetConfigurationRequest, opts ...grpc.CallOption) (*runtimev1pb.GetConfigurationResponse, error) {
	return l.api.GetConfiguration(ctx, in)
}

func (l *localRuntimeClient) SaveConfiguration(ctx context.Context, in *runtimev1pb.SaveConfigurationRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return l.api.SaveConfiguration(ctx, in)
}

func (l *localRuntimeClient) DeleteConfiguration(ctx context.Context, in *runtimev1pb.DeleteConfigurationRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return l.api.DeleteConfiguration(ctx, in)
}

func (l *localRuntimeClient) SubscribeConfiguration(ctx context.Context, opts ...grpc.CallOption) (runtimev1pb.Runtime_SubscribeConfigurationClient, error) {
	panic("implement me")
}

func (l *localRuntimeClient) TryLock(ctx context.Context, in *runtimev1pb.TryLockRequest, opts ...grpc.CallOption) (*runtimev1pb.TryLockResponse, error) {
	return l.api.TryLock(ctx, in)
}

func (l *localRuntimeClient) Unlock(ctx context.Context, in *runtimev1pb.UnlockRequest, opts ...grpc.CallOption) (*runtimev1pb.UnlockResponse, error) {
	return l.api.Unlock(ctx, in)
}

func (l *localRuntimeClient) GetNextId(ctx context.Context, in *runtimev1pb.GetNextIdRequest, opts ...grpc.CallOption) (*runtimev1pb.GetNextIdResponse, error) {
	return l.api.GetNextId(ctx, in)
}

func (l *localRuntimeClient) GetState(ctx context.Context, in *runtimev1pb.GetStateRequest, opts ...grpc.CallOption) (*runtimev1pb.GetStateResponse, error) {
	return l.api.GetState(ctx, in)
}

func (l *localRuntimeClient) GetBulkState(ctx context.Context, in *runtimev1pb.GetBulkStateRequest, opts ...grpc.CallOption) (*runtimev1pb.GetBulkStateResponse, error) {
	return l.api.GetBulkState(ctx, in)
}

func (l *localRuntimeClient) SaveState(ctx context.Context, in *runtimev1pb.SaveStateRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return l.api.SaveState(ctx, in)
}

func (l *localRuntimeClient) DeleteState(ctx context.Context, in *runtimev1pb.DeleteStateRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return l.api.DeleteState(ctx, in)
}

func (l *localRuntimeClient) DeleteBulkState(ctx context.Context, in *runtimev1pb.DeleteBulkStateRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return l.api.DeleteBulkState(ctx, in)
}

func (l *localRuntimeClient) ExecuteStateTransaction(ctx context.Context, in *runtimev1pb.ExecuteStateTransactionRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return l.api.ExecuteStateTransaction(ctx, in)
}

func (l *localRuntimeClient) PublishEvent(ctx context.Context, in *runtimev1pb.PublishEventRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return l.api.PublishEvent(ctx, in)
}

func (l *localRuntimeClient) GetFile(ctx context.Context, in *runtimev1pb.GetFileRequest, opts ...grpc.CallOption) (runtimev1pb.Runtime_GetFileClient, error) {
	panic("implement me")
}

func (l *localRuntimeClient) PutFile(ctx context.Context, opts ...grpc.CallOption) (runtimev1pb.Runtime_PutFileClient, error) {
	panic("implement me")
}

func (l *localRuntimeClient) ListFile(ctx context.Context, in *runtimev1pb.ListFileRequest, opts ...grpc.CallOption) (*runtimev1pb.ListFileResp, error) {
	panic("implement me")
}

func (l *localRuntimeClient) DelFile(ctx context.Context, in *runtimev1pb.DelFileRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return l.api.DelFile(ctx, in)
}

func (l *localRuntimeClient) InvokeBinding(ctx context.Context, in *runtimev1pb.InvokeBindingRequest, opts ...grpc.CallOption) (*runtimev1pb.InvokeBindingResponse, error) {
	return l.api.InvokeBinding(ctx, in)
}
