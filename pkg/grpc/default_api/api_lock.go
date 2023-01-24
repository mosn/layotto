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

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"mosn.io/pkg/log"

	"mosn.io/layotto/components/lock"
	"mosn.io/layotto/pkg/messages"
	runtime_lock "mosn.io/layotto/pkg/runtime/lock"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
)

func (a *api) TryLock(ctx context.Context, req *runtimev1pb.TryLockRequest) (*runtimev1pb.TryLockResponse, error) {
	// 1. validate
	if a.lockStores == nil || len(a.lockStores) == 0 {
		err := status.Error(codes.FailedPrecondition, messages.ErrLockStoresNotConfigured)
		log.DefaultLogger.Errorf("[runtime] [grpc.TryLock] error: %v", err)
		return &runtimev1pb.TryLockResponse{}, err
	}
	if req.ResourceId == "" {
		err := status.Errorf(codes.InvalidArgument, messages.ErrResourceIdEmpty, req.StoreName)
		return &runtimev1pb.TryLockResponse{}, err
	}
	if req.LockOwner == "" {
		err := status.Errorf(codes.InvalidArgument, messages.ErrLockOwnerEmpty, req.StoreName)
		return &runtimev1pb.TryLockResponse{}, err
	}
	if req.Expire <= 0 {
		err := status.Errorf(codes.InvalidArgument, messages.ErrExpireNotPositive, req.StoreName)
		return &runtimev1pb.TryLockResponse{}, err
	}
	// 2. find store component
	store, ok := a.lockStores[req.StoreName]
	if !ok {
		return &runtimev1pb.TryLockResponse{}, status.Errorf(codes.InvalidArgument, messages.ErrLockStoreNotFound, req.StoreName)
	}
	// 3. convert request
	compReq := TryLockRequest2ComponentRequest(req)
	// modify key
	var err error
	compReq.ResourceId, err = runtime_lock.GetModifiedLockKey(compReq.ResourceId, req.StoreName, a.appId)
	if err != nil {
		log.DefaultLogger.Errorf("[runtime] [grpc.TryLock] error: %v", err)
		return &runtimev1pb.TryLockResponse{}, err
	}
	// 4. delegate to the component
	compResp, err := store.TryLock(compReq)
	if err != nil {
		log.DefaultLogger.Errorf("[runtime] [grpc.TryLock] error: %v", err)
		return &runtimev1pb.TryLockResponse{}, err
	}
	// 5. convert response
	resp := TryLockResponse2GrpcResponse(compResp)
	return resp, nil
}

func (a *api) Unlock(ctx context.Context, req *runtimev1pb.UnlockRequest) (*runtimev1pb.UnlockResponse, error) {
	// 1. validate
	if a.lockStores == nil || len(a.lockStores) == 0 {
		err := status.Error(codes.FailedPrecondition, messages.ErrLockStoresNotConfigured)
		log.DefaultLogger.Errorf("[runtime] [grpc.Unlock] error: %v", err)
		return newInternalErrorUnlockResponse(), err
	}
	if req.ResourceId == "" {
		err := status.Errorf(codes.InvalidArgument, messages.ErrResourceIdEmpty, req.StoreName)
		return newInternalErrorUnlockResponse(), err
	}
	if req.LockOwner == "" {
		err := status.Errorf(codes.InvalidArgument, messages.ErrLockOwnerEmpty, req.StoreName)
		return newInternalErrorUnlockResponse(), err
	}
	// 2. find store component
	store, ok := a.lockStores[req.StoreName]
	if !ok {
		return newInternalErrorUnlockResponse(), status.Errorf(codes.InvalidArgument, messages.ErrLockStoreNotFound, req.StoreName)
	}
	// 3. convert request
	compReq := UnlockGrpc2ComponentRequest(req)
	// modify key
	var err error
	compReq.ResourceId, err = runtime_lock.GetModifiedLockKey(compReq.ResourceId, req.StoreName, a.appId)
	if err != nil {
		log.DefaultLogger.Errorf("[runtime] [grpc.TryLock] error: %v", err)
		return newInternalErrorUnlockResponse(), err
	}
	// 4. delegate to the component
	compResp, err := store.Unlock(compReq)
	if err != nil {
		log.DefaultLogger.Errorf("[runtime] [grpc.Unlock] error: %v", err)
		return newInternalErrorUnlockResponse(), err
	}
	// 5. convert response
	resp := UnlockComp2GrpcResponse(compResp)
	return resp, nil
}

func (a *api) LockKeepAlive(ctx context.Context, request *runtimev1pb.LockKeepAliveRequest) (*runtimev1pb.LockKeepAliveResponse, error) {
	return nil, nil
}

func newInternalErrorUnlockResponse() *runtimev1pb.UnlockResponse {
	return &runtimev1pb.UnlockResponse{
		Status: runtimev1pb.UnlockResponse_INTERNAL_ERROR,
	}
}

func TryLockRequest2ComponentRequest(req *runtimev1pb.TryLockRequest) *lock.TryLockRequest {
	result := &lock.TryLockRequest{}
	if req == nil {
		return result
	}
	result.ResourceId = req.ResourceId
	result.LockOwner = req.LockOwner
	result.Expire = req.Expire
	return result
}

func TryLockResponse2GrpcResponse(compResponse *lock.TryLockResponse) *runtimev1pb.TryLockResponse {
	result := &runtimev1pb.TryLockResponse{}
	if compResponse == nil {
		return result
	}
	result.Success = compResponse.Success
	return result
}

func UnlockGrpc2ComponentRequest(req *runtimev1pb.UnlockRequest) *lock.UnlockRequest {
	result := &lock.UnlockRequest{}
	if req == nil {
		return result
	}
	result.ResourceId = req.ResourceId
	result.LockOwner = req.LockOwner
	return result
}

func UnlockComp2GrpcResponse(compResp *lock.UnlockResponse) *runtimev1pb.UnlockResponse {
	result := &runtimev1pb.UnlockResponse{}
	if compResp == nil {
		return result
	}
	result.Status = runtimev1pb.UnlockResponse_Status(compResp.Status)
	return result
}
