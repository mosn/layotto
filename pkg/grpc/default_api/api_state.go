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

	"github.com/dapr/components-contrib/state"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"mosn.io/pkg/log"

	"mosn.io/layotto/pkg/common"
	"mosn.io/layotto/pkg/messages"
	state2 "mosn.io/layotto/pkg/runtime/state"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
)

// GetState obtains the state for a specific key.
func (a *api) GetState(ctx context.Context, in *runtimev1pb.GetStateRequest) (*runtimev1pb.GetStateResponse, error) {
	// check if the StateRequest is exists
	if in == nil {
		return &runtimev1pb.GetStateResponse{}, status.Error(codes.InvalidArgument, "GetStateRequest is nil")
	}

	// get store
	store, err := a.getStateStore(in.GetStoreName())
	if err != nil {
		log.DefaultLogger.Errorf("[runtime] [grpc.GetState] error: %v", err)
		return nil, err
	}

	// generate the actual key
	key, err := state2.GetModifiedStateKey(in.GetKey(), in.GetStoreName(), a.appId)
	if err != nil {
		return &runtimev1pb.GetStateResponse{}, err
	}
	req := &state.GetRequest{
		Key:      key,
		Metadata: in.GetMetadata(),
		Options: state.GetStateOption{
			Consistency: StateConsistencyToString(in.GetConsistency()),
		},
	}

	// query
	compResp, err := store.Get(ctx, req)

	// check result
	if err != nil {
		err = status.Errorf(codes.Internal, messages.ErrStateGet, in.GetKey(), in.GetStoreName(), err.Error())
		log.DefaultLogger.Errorf("[runtime] [grpc.GetState] %v", err)
		return &runtimev1pb.GetStateResponse{}, err
	}

	return GetResponse2GetStateResponse(compResp), nil
}

func (a *api) SaveState(ctx context.Context, in *runtimev1pb.SaveStateRequest) (*emptypb.Empty, error) {
	// Check if the request is nil
	if in == nil {
		return &emptypb.Empty{}, status.Error(codes.InvalidArgument, "SaveStateRequest is nil")
	}

	// get store
	store, err := a.getStateStore(in.StoreName)
	if err != nil {
		log.DefaultLogger.Errorf("[runtime] [grpc.SaveState] error: %v", err)
		return &emptypb.Empty{}, err
	}

	// convert requests
	reqs := []state.SetRequest{}
	for _, s := range in.States {
		key, err := state2.GetModifiedStateKey(s.Key, in.StoreName, a.appId)
		if err != nil {
			return &emptypb.Empty{}, err
		}
		reqs = append(reqs, *StateItem2SetRequest(s, key))
	}

	// query
	err = store.BulkSet(ctx, reqs, state.BulkStoreOpts{})

	// check result
	if err != nil {
		err = a.wrapDaprComponentError(err, messages.ErrStateSave, in.StoreName, err.Error())
		log.DefaultLogger.Errorf("[runtime] [grpc.SaveState] error: %v", err)
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, nil
}

// GetBulkState gets a batch of state data
func (a *api) GetBulkState(ctx context.Context, in *runtimev1pb.GetBulkStateRequest) (*runtimev1pb.GetBulkStateResponse, error) {
	// Check if the request is nil
	if in == nil {
		return &runtimev1pb.GetBulkStateResponse{}, status.Error(codes.InvalidArgument, "GetBulkStateRequest is nil")
	}

	// Generate response by request
	resp, err := a.getBulkState(ctx, in)
	if err != nil {
		return &runtimev1pb.GetBulkStateResponse{}, err
	}

	ret := &runtimev1pb.GetBulkStateResponse{Items: make([]*runtimev1pb.BulkStateItem, 0)}
	for _, item := range resp.Items {
		ret.Items = append(ret.Items, &runtimev1pb.BulkStateItem{
			Key:      item.GetKey(),
			Data:     item.GetData(),
			Etag:     item.GetEtag(),
			Error:    item.GetError(),
			Metadata: item.GetMetadata(),
		})
	}

	return ret, nil
}

func (a *api) DeleteState(ctx context.Context, in *runtimev1pb.DeleteStateRequest) (*emptypb.Empty, error) {
	// Check if the request is nil
	if in == nil {
		return &emptypb.Empty{}, status.Error(codes.InvalidArgument, "DeleteStateRequest is nil")
	}

	// get store
	store, err := a.getStateStore(in.GetStoreName())
	if err != nil {
		log.DefaultLogger.Errorf("[runtime] [grpc.DeleteState] error: %v", err)
		return &emptypb.Empty{}, err
	}

	// generate the actual key
	key, err := state2.GetModifiedStateKey(in.GetKey(), in.GetStoreName(), a.appId)
	if err != nil {
		return &empty.Empty{}, err
	}

	// convert and send request
	err = store.Delete(ctx, DeleteStateRequest2DeleteRequest(in, key))

	// 4. check result
	if err != nil {
		err = a.wrapDaprComponentError(err, messages.ErrStateDelete, in.GetKey(), err.Error())
		log.DefaultLogger.Errorf("[runtime] [grpc.DeleteState] error: %v", err)
		return &empty.Empty{}, err
	}

	return &empty.Empty{}, nil
}

func (a *api) DeleteBulkState(ctx context.Context, in *runtimev1pb.DeleteBulkStateRequest) (*empty.Empty, error) {
	// Check if the request is nil
	if in == nil {
		return &emptypb.Empty{}, status.Error(codes.InvalidArgument, "DeleteBulkStateRequest is nil")
	}

	// get store
	store, err := a.getStateStore(in.GetStoreName())
	if err != nil {
		log.DefaultLogger.Errorf("[runtime] [grpc.DeleteBulkState] error: %v", err)
		return &empty.Empty{}, err
	}

	// convert request
	reqs := make([]state.DeleteRequest, 0, len(in.GetStates()))
	for _, item := range in.States {
		key, err := state2.GetModifiedStateKey(item.Key, in.GetStoreName(), a.appId)
		if err != nil {
			return &empty.Empty{}, err
		}
		reqs = append(reqs, *StateItem2DeleteRequest(item, key))
	}

	// send request
	err = store.BulkDelete(ctx, reqs, state.BulkStoreOpts{})

	// check result
	if err != nil {
		log.DefaultLogger.Errorf("[runtime] [grpc.DeleteBulkState] error: %v", err)
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, nil
}

func (a *api) ExecuteStateTransaction(ctx context.Context, in *runtimev1pb.ExecuteStateTransactionRequest) (*emptypb.Empty, error) {
	// Check if the request is nil
	if in == nil {
		return &emptypb.Empty{}, status.Error(codes.InvalidArgument, "ExecuteStateTransactionRequest is nil")
	}

	// 1. check params
	if a.stateStores == nil || len(a.stateStores) == 0 {
		err := status.Error(codes.FailedPrecondition, messages.ErrStateStoresNotConfigured)
		log.DefaultLogger.Errorf("[runtime] [grpc.ExecuteStateTransaction] error: %v", err)
		return &emptypb.Empty{}, err
	}
	storeName := in.GetStoreName()
	if a.stateStores[storeName] == nil {
		err := status.Errorf(codes.InvalidArgument, messages.ErrStateStoreNotFound, storeName)
		log.DefaultLogger.Errorf("[runtime] [grpc.ExecuteStateTransaction] error: %v", err)
		return &emptypb.Empty{}, err
	}

	// find store
	store, ok := a.transactionalStateStores[storeName]
	if !ok {
		err := status.Errorf(codes.Unimplemented, messages.ErrStateStoreNotSupported, storeName)
		log.DefaultLogger.Errorf("[runtime] [grpc.ExecuteStateTransaction] error: %v", err)
		return &emptypb.Empty{}, err
	}

	// convert request
	operations := []state.TransactionalStateOperation{}
	for _, op := range in.Operations {
		// extract and validate fields
		var operation state.TransactionalStateOperation
		var req = op.Request
		// tolerant npe
		if req == nil {
			log.DefaultLogger.Warnf("[runtime] [grpc.ExecuteStateTransaction] one of TransactionalStateOperation.Request is nil")
			continue
		}
		key, err := state2.GetModifiedStateKey(req.Key, in.GetStoreName(), a.appId)
		if err != nil {
			return &emptypb.Empty{}, err
		}
		// 3.2. prepare TransactionalStateOperation struct according to the operation type
		switch state.OperationType(op.OperationType) {
		case state.OperationUpsert:
			operation = *StateItem2SetRequest(req, key)
		case state.OperationDelete:
			operation = *StateItem2DeleteRequest(req, key)
		default:
			err := status.Errorf(codes.Unimplemented, messages.ErrNotSupportedStateOperation, op.OperationType)
			log.DefaultLogger.Errorf("[runtime] [grpc.ExecuteStateTransaction] error: %v", err)
			return &emptypb.Empty{}, err
		}
		operations = append(operations, operation)
	}

	// submit transactional request
	err := store.Multi(ctx, &state.TransactionalStateRequest{
		Operations: operations,
		Metadata:   in.GetMetadata(),
	})

	// check result
	if err != nil {
		err = status.Errorf(codes.Internal, messages.ErrStateTransaction, err.Error())
		log.DefaultLogger.Errorf("[runtime] [grpc.ExecuteStateTransaction] error: %v", err)
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, nil
}

// the specific processing logic of GetBulkState
func (a *api) getBulkState(ctx context.Context, in *runtimev1pb.GetBulkStateRequest) (*runtimev1pb.GetBulkStateResponse, error) {
	// get store
	store, err := a.getStateStore(in.GetStoreName())
	if err != nil {
		log.DefaultLogger.Errorf("[runtime] [grpc.GetBulkState] error: %v", err)
		return &runtimev1pb.GetBulkStateResponse{}, err
	}

	bulkResp := &runtimev1pb.GetBulkStateResponse{}
	if len(in.GetKeys()) == 0 {
		return bulkResp, nil
	}

	// store.BulkGet
	// convert reqs
	reqs := make([]state.GetRequest, len(in.GetKeys()))
	for i, k := range in.GetKeys() {
		key, err := state2.GetModifiedStateKey(k, in.GetStoreName(), a.appId)
		if err != nil {
			return &runtimev1pb.GetBulkStateResponse{}, err
		}
		r := state.GetRequest{
			Key:      key,
			Metadata: in.GetMetadata(),
		}
		reqs[i] = r
	}

	// query
	responses, err := store.BulkGet(ctx, reqs, state.BulkGetOpts{})
	if err != nil {
		return bulkResp, err
	}

	for i := 0; i < len(responses); i++ {
		bulkResp.Items = append(bulkResp.Items, BulkGetResponse2BulkStateItem(&responses[i]))
	}

	return bulkResp, nil
}

func (a *api) getStateStore(name string) (state.Store, error) {
	// check if the stateStores exists
	if a.stateStores == nil || len(a.stateStores) == 0 {
		return nil, status.Error(codes.FailedPrecondition, messages.ErrStateStoresNotConfigured)
	}
	// check name
	if a.stateStores[name] == nil {
		return nil, status.Errorf(codes.InvalidArgument, messages.ErrStateStoreNotFound, name)
	}
	return a.stateStores[name], nil
}

func StateConsistencyToString(c runtimev1pb.StateOptions_StateConsistency) string {
	// check
	switch c {
	case runtimev1pb.StateOptions_CONSISTENCY_EVENTUAL:
		return "eventual"
	case runtimev1pb.StateOptions_CONSISTENCY_STRONG:
		return "strong"
	}
	return ""
}

func StateConcurrencyToString(c runtimev1pb.StateOptions_StateConcurrency) string {
	// check the StateOptions of StateOptions_StateConcurrency
	switch c {
	case runtimev1pb.StateOptions_CONCURRENCY_FIRST_WRITE:
		return "first-write"
	case runtimev1pb.StateOptions_CONCURRENCY_LAST_WRITE:
		return "last-write"
	}

	return ""
}

func GetResponse2GetStateResponse(compResp *state.GetResponse) *runtimev1pb.GetStateResponse {
	// Initialize an element of type GetStateResponse
	resp := &runtimev1pb.GetStateResponse{}
	// check if the compResp exists
	if compResp != nil {
		resp.Etag = common.PointerToString(compResp.ETag)
		resp.Data = compResp.Data
		resp.Metadata = compResp.Metadata
	}
	return resp
}

func StateItem2SetRequest(grpcReq *runtimev1pb.StateItem, key string) *state.SetRequest {
	// Set the key for the request
	req := &state.SetRequest{
		Key: key,
	}
	// check if the grpcReq exists
	if grpcReq == nil {
		return req
	}
	// Assign the value of grpcReq property to req
	req.Metadata = grpcReq.Metadata
	req.Value = grpcReq.Value
	// Check grpcReq.Etag
	if grpcReq.Etag != nil {
		req.ETag = &grpcReq.Etag.Value
	}
	// Check grpcReq.Options
	if grpcReq.Options != nil {
		req.Options = state.SetStateOption{
			Consistency: StateConsistencyToString(grpcReq.Options.Consistency),
			Concurrency: StateConcurrencyToString(grpcReq.Options.Concurrency),
		}
	}
	return req
}

// wrapDaprComponentError parse and wrap error from dapr component
func (a *api) wrapDaprComponentError(err error, format string, args ...interface{}) error {
	e, ok := err.(*state.ETagError)
	if !ok {
		return status.Errorf(codes.Internal, format, args...)
	}
	// check the Kind of error
	switch e.Kind() {
	case state.ETagMismatch:
		return status.Errorf(codes.Aborted, format, args...)
	case state.ETagInvalid:
		return status.Errorf(codes.InvalidArgument, format, args...)
	}

	return status.Errorf(codes.Internal, format, args...)
}

// converting from BulkGetResponse to BulkStateItem
func BulkGetResponse2BulkStateItem(compResp *state.BulkGetResponse) *runtimev1pb.BulkStateItem {
	if compResp == nil {
		return &runtimev1pb.BulkStateItem{}
	}
	return &runtimev1pb.BulkStateItem{
		Key:      state2.GetOriginalStateKey(compResp.Key),
		Data:     compResp.Data,
		Etag:     common.PointerToString(compResp.ETag),
		Metadata: compResp.Metadata,
		Error:    compResp.Error,
	}
}

// converting from GetResponse to BulkStateItem
func GetResponse2BulkStateItem(compResp *state.GetResponse, key string) *runtimev1pb.BulkStateItem {
	// convert
	resp := &runtimev1pb.BulkStateItem{}
	resp.Key = key
	if compResp != nil {
		resp.Data = compResp.Data
		resp.Etag = common.PointerToString(compResp.ETag)
		resp.Metadata = compResp.Metadata
	}
	return resp
}

// converting from DeleteStateRequest to DeleteRequest
func DeleteStateRequest2DeleteRequest(grpcReq *runtimev1pb.DeleteStateRequest, key string) *state.DeleteRequest {
	// convert
	req := &state.DeleteRequest{
		Key: key,
	}
	if grpcReq == nil {
		return req
	}
	req.Metadata = grpcReq.Metadata
	if grpcReq.Etag != nil {
		req.ETag = &grpcReq.Etag.Value
	}
	if grpcReq.Options != nil {
		req.Options = state.DeleteStateOption{
			Concurrency: StateConcurrencyToString(grpcReq.Options.Concurrency),
			Consistency: StateConsistencyToString(grpcReq.Options.Consistency),
		}
	}
	return req
}

// converting from StateItem to DeleteRequest
func StateItem2DeleteRequest(grpcReq *runtimev1pb.StateItem, key string) *state.DeleteRequest {
	//convert
	req := &state.DeleteRequest{
		Key: key,
	}
	if grpcReq == nil {
		return req
	}
	req.Metadata = grpcReq.Metadata
	if grpcReq.Etag != nil {
		req.ETag = &grpcReq.Etag.Value
	}
	if grpcReq.Options != nil {
		req.Options = state.DeleteStateOption{
			Concurrency: StateConcurrencyToString(grpcReq.Options.Concurrency),
			Consistency: StateConsistencyToString(grpcReq.Options.Consistency),
		}
	}
	return req
}
