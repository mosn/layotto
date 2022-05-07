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

	"github.com/dapr/components-contrib/state"
	"github.com/gammazero/workerpool"
	"github.com/golang/protobuf/ptypes/empty"
	jsoniter "github.com/json-iterator/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"mosn.io/pkg/log"

	"mosn.io/layotto/pkg/common"
	dapr_common_v1pb "mosn.io/layotto/pkg/grpc/dapr/proto/common/v1"
	dapr_v1pb "mosn.io/layotto/pkg/grpc/dapr/proto/runtime/v1"
	"mosn.io/layotto/pkg/messages"
	state2 "mosn.io/layotto/pkg/runtime/state"
)

func (d *daprGrpcAPI) SaveState(ctx context.Context, in *dapr_v1pb.SaveStateRequest) (*emptypb.Empty, error) {
	// 1. get store
	store, err := d.getStateStore(in.StoreName)
	if err != nil {
		log.DefaultLogger.Errorf("[runtime] [grpc.SaveState] error: %v", err)
		return &emptypb.Empty{}, err
	}
	// 2. convert requests
	reqs := []state.SetRequest{}
	for _, s := range in.States {
		key, err := state2.GetModifiedStateKey(s.Key, in.StoreName, d.appId)
		if err != nil {
			return &emptypb.Empty{}, err
		}
		reqs = append(reqs, *StateItem2SetRequest(s, key))
	}
	// 3. query
	err = store.BulkSet(reqs)
	// 4. check result
	if err != nil {
		err = d.wrapDaprComponentError(err, messages.ErrStateSave, in.StoreName, err.Error())
		log.DefaultLogger.Errorf("[runtime] [grpc.SaveState] error: %v", err)
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}

// GetState obtains the state for a specific key.
func (d *daprGrpcAPI) GetState(ctx context.Context, request *dapr_v1pb.GetStateRequest) (*dapr_v1pb.GetStateResponse, error) {
	// 1. get store
	store, err := d.getStateStore(request.StoreName)
	if err != nil {
		log.DefaultLogger.Errorf("[runtime] [grpc.GetState] error: %v", err)
		return nil, err
	}
	// 2. generate the actual key
	key, err := state2.GetModifiedStateKey(request.Key, request.StoreName, d.appId)
	if err != nil {
		return &dapr_v1pb.GetStateResponse{}, err
	}
	req := &state.GetRequest{
		Key:      key,
		Metadata: request.GetMetadata(),
		Options: state.GetStateOption{
			Consistency: StateConsistencyToString(request.Consistency),
		},
	}
	// 3. query
	compResp, err := store.Get(req)
	// 4. check result
	if err != nil {
		err = status.Errorf(codes.Internal, messages.ErrStateGet, request.Key, request.StoreName, err.Error())
		log.DefaultLogger.Errorf("[runtime] [grpc.GetState] %v", err)
		return &dapr_v1pb.GetStateResponse{}, err
	}
	return GetResponse2GetStateResponse(compResp), nil
}

func (d *daprGrpcAPI) GetBulkState(ctx context.Context, request *dapr_v1pb.GetBulkStateRequest) (*dapr_v1pb.GetBulkStateResponse, error) {
	// 1. get store
	store, err := d.getStateStore(request.StoreName)
	if err != nil {
		log.DefaultLogger.Errorf("[runtime] [grpc.GetBulkState] error: %v", err)
		return &dapr_v1pb.GetBulkStateResponse{}, err
	}

	bulkResp := &dapr_v1pb.GetBulkStateResponse{}
	if len(request.Keys) == 0 {
		return bulkResp, nil
	}

	// 2. store.BulkGet
	// 2.1. convert reqs
	reqs := make([]state.GetRequest, len(request.Keys))
	for i, k := range request.Keys {
		key, err := state2.GetModifiedStateKey(k, request.StoreName, d.appId)
		if err != nil {
			return &dapr_v1pb.GetBulkStateResponse{}, err
		}
		r := state.GetRequest{
			Key:      key,
			Metadata: request.GetMetadata(),
		}
		reqs[i] = r
	}
	// 2.2. query
	support, responses, err := store.BulkGet(reqs)
	if err != nil {
		return bulkResp, err
	}
	// 2.3. parse and return result if store supports this method
	if support {
		for i := 0; i < len(responses); i++ {
			bulkResp.Items = append(bulkResp.Items, BulkGetResponse2BulkStateItem(&responses[i]))
		}
		return bulkResp, nil
	}

	// 3. Simulate the method if the store doesn't support it
	n := len(reqs)
	pool := workerpool.New(int(request.Parallelism))
	resultCh := make(chan *dapr_v1pb.BulkStateItem, n)
	for i := 0; i < n; i++ {
		pool.Submit(generateGetStateTask(store, &reqs[i], resultCh))
	}
	pool.StopWait()
	for {
		select {
		case item, ok := <-resultCh:
			if !ok {
				return bulkResp, nil
			}
			bulkResp.Items = append(bulkResp.Items, item)
		default:
			return bulkResp, nil
		}
	}
}

func (d *daprGrpcAPI) QueryStateAlpha1(ctx context.Context, request *dapr_v1pb.QueryStateRequest) (*dapr_v1pb.QueryStateResponse, error) {
	ret := &dapr_v1pb.QueryStateResponse{}

	// 1. get state store component
	store, err := d.getStateStore(request.StoreName)
	if err != nil {
		log.DefaultLogger.Errorf("[runtime] [grpc.QueryStateAlpha1] error: %v", err)
		return ret, err
	}

	// 2. check if this store has the query feature
	querier, ok := store.(state.Querier)
	if !ok {
		err = status.Errorf(codes.Unimplemented, messages.ErrNotFound, "Query")
		log.DefaultLogger.Errorf("[runtime] [grpc.QueryStateAlpha1] error: %v", err)
		return ret, err
	}

	// 3. Unmarshal query dsl
	var req state.QueryRequest
	if err = jsoniter.Unmarshal([]byte(request.GetQuery()), &req.Query); err != nil {
		err = status.Errorf(codes.InvalidArgument, messages.ErrMalformedRequest, err.Error())
		log.DefaultLogger.Errorf("[runtime] [grpc.QueryStateAlpha1] error: %v", err)
		return ret, err
	}
	req.Metadata = request.GetMetadata()

	// 4. delegate to the store
	resp, err := querier.Query(&req)
	// 5. convert response
	if err != nil {
		err = status.Errorf(codes.Internal, messages.ErrStateQuery, request.GetStoreName(), err.Error())
		log.DefaultLogger.Errorf("[runtime] [grpc.QueryStateAlpha1] error: %v", err)
		return ret, err
	}
	if resp == nil || len(resp.Results) == 0 {
		return ret, nil
	}

	ret.Results = make([]*dapr_v1pb.QueryStateItem, len(resp.Results))
	ret.Token = resp.Token
	ret.Metadata = resp.Metadata

	for i := range resp.Results {
		ret.Results[i] = &dapr_v1pb.QueryStateItem{
			Key:  state2.GetOriginalStateKey(resp.Results[i].Key),
			Data: resp.Results[i].Data,
		}
	}
	return ret, nil
}

func (d *daprGrpcAPI) DeleteState(ctx context.Context, request *dapr_v1pb.DeleteStateRequest) (*emptypb.Empty, error) {
	// 1. get store
	store, err := d.getStateStore(request.StoreName)
	if err != nil {
		log.DefaultLogger.Errorf("[runtime] [grpc.DeleteState] error: %v", err)
		return &emptypb.Empty{}, err
	}
	// 2. generate the actual key
	key, err := state2.GetModifiedStateKey(request.Key, request.StoreName, d.appId)
	if err != nil {
		return &empty.Empty{}, err
	}
	// 3. convert and send request
	err = store.Delete(DeleteStateRequest2DeleteRequest(request, key))
	// 4. check result
	if err != nil {
		err = d.wrapDaprComponentError(err, messages.ErrStateDelete, request.Key, err.Error())
		log.DefaultLogger.Errorf("[runtime] [grpc.DeleteState] error: %v", err)
		return &empty.Empty{}, err
	}
	return &empty.Empty{}, nil
}

func (d *daprGrpcAPI) DeleteBulkState(ctx context.Context, request *dapr_v1pb.DeleteBulkStateRequest) (*emptypb.Empty, error) {
	// 1. get store
	store, err := d.getStateStore(request.StoreName)
	if err != nil {
		log.DefaultLogger.Errorf("[runtime] [grpc.DeleteBulkState] error: %v", err)
		return &empty.Empty{}, err
	}
	// 2. convert request
	reqs := make([]state.DeleteRequest, 0, len(request.States))
	for _, item := range request.States {
		key, err := state2.GetModifiedStateKey(item.Key, request.StoreName, d.appId)
		if err != nil {
			return &empty.Empty{}, err
		}
		reqs = append(reqs, *StateItem2DeleteRequest(item, key))
	}
	// 3. send request
	err = store.BulkDelete(reqs)
	// 4. check result
	if err != nil {
		log.DefaultLogger.Errorf("[runtime] [grpc.DeleteBulkState] error: %v", err)
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}

func (d *daprGrpcAPI) ExecuteStateTransaction(ctx context.Context, request *dapr_v1pb.ExecuteStateTransactionRequest) (*emptypb.Empty, error) {
	// 1. check params
	if d.stateStores == nil || len(d.stateStores) == 0 {
		err := status.Error(codes.FailedPrecondition, messages.ErrStateStoresNotConfigured)
		log.DefaultLogger.Errorf("[runtime] [grpc.ExecuteStateTransaction] error: %v", err)
		return &emptypb.Empty{}, err
	}
	storeName := request.StoreName
	if d.stateStores[storeName] == nil {
		err := status.Errorf(codes.InvalidArgument, messages.ErrStateStoreNotFound, storeName)
		log.DefaultLogger.Errorf("[runtime] [grpc.ExecuteStateTransaction] error: %v", err)
		return &emptypb.Empty{}, err
	}
	// 2. find store
	store, ok := d.transactionalStateStores[storeName]
	if !ok {
		err := status.Errorf(codes.Unimplemented, messages.ErrStateStoreNotSupported, storeName)
		log.DefaultLogger.Errorf("[runtime] [grpc.ExecuteStateTransaction] error: %v", err)
		return &emptypb.Empty{}, err
	}
	// 3. convert request
	operations := []state.TransactionalStateOperation{}
	for _, op := range request.Operations {
		// 3.1. extract and validate fields
		var operation state.TransactionalStateOperation
		var req = op.Request
		// tolerant npe
		if req == nil {
			log.DefaultLogger.Warnf("[runtime] [grpc.ExecuteStateTransaction] one of TransactionalStateOperation.Request is nil")
			continue
		}
		key, err := state2.GetModifiedStateKey(req.Key, request.StoreName, d.appId)
		if err != nil {
			return &emptypb.Empty{}, err
		}
		// 3.2. prepare TransactionalStateOperation struct according to the operation type
		switch state.OperationType(op.OperationType) {
		case state.Upsert:
			operation = state.TransactionalStateOperation{
				Operation: state.Upsert,
				Request:   *StateItem2SetRequest(req, key),
			}
		case state.Delete:
			operation = state.TransactionalStateOperation{
				Operation: state.Delete,
				Request:   *StateItem2DeleteRequest(req, key),
			}
		default:
			err := status.Errorf(codes.Unimplemented, messages.ErrNotSupportedStateOperation, op.OperationType)
			log.DefaultLogger.Errorf("[runtime] [grpc.ExecuteStateTransaction] error: %v", err)
			return &emptypb.Empty{}, err
		}
		operations = append(operations, operation)
	}
	// 4. submit transactional request
	err := store.Multi(&state.TransactionalStateRequest{
		Operations: operations,
		Metadata:   request.Metadata,
	})
	// 5. check result
	if err != nil {
		err = status.Errorf(codes.Internal, messages.ErrStateTransaction, err.Error())
		log.DefaultLogger.Errorf("[runtime] [grpc.ExecuteStateTransaction] error: %v", err)
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}

func (d *daprGrpcAPI) getStateStore(name string) (state.Store, error) {
	// check if the stateStores exists
	if d.stateStores == nil || len(d.stateStores) == 0 {
		return nil, status.Error(codes.FailedPrecondition, messages.ErrStateStoresNotConfigured)
	}
	// check name
	if d.stateStores[name] == nil {
		return nil, status.Errorf(codes.InvalidArgument, messages.ErrStateStoreNotFound, name)
	}
	return d.stateStores[name], nil
}

func StateItem2SetRequest(grpcReq *dapr_common_v1pb.StateItem, key string) *state.SetRequest {
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

func GetResponse2GetStateResponse(compResp *state.GetResponse) *dapr_v1pb.GetStateResponse {
	// Initialize an element of type GetStateResponse
	resp := &dapr_v1pb.GetStateResponse{}
	// check if the compResp exists
	if compResp != nil {
		resp.Etag = common.PointerToString(compResp.ETag)
		resp.Data = compResp.Data
		resp.Metadata = compResp.Metadata
	}
	return resp
}

func StateConsistencyToString(c dapr_common_v1pb.StateOptions_StateConsistency) string {
	// check
	switch c {
	case dapr_common_v1pb.StateOptions_CONSISTENCY_EVENTUAL:
		return "eventual"
	case dapr_common_v1pb.StateOptions_CONSISTENCY_STRONG:
		return "strong"
	}
	return ""
}

func StateConcurrencyToString(c dapr_common_v1pb.StateOptions_StateConcurrency) string {
	// check the StateOptions of StateOptions_StateConcurrency
	switch c {
	case dapr_common_v1pb.StateOptions_CONCURRENCY_FIRST_WRITE:
		return "first-write"
	case dapr_common_v1pb.StateOptions_CONCURRENCY_LAST_WRITE:
		return "last-write"
	}

	return ""
}

// wrapDaprComponentError parse and wrap error from dapr component
func (d *daprGrpcAPI) wrapDaprComponentError(err error, format string, args ...interface{}) error {
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

func generateGetStateTask(store state.Store, req *state.GetRequest, resultCh chan *dapr_v1pb.BulkStateItem) func() {
	return func() {
		// get
		r, err := store.Get(req)
		// convert
		var item *dapr_v1pb.BulkStateItem
		if err != nil {
			item = &dapr_v1pb.BulkStateItem{
				Key:   state2.GetOriginalStateKey(req.Key),
				Error: err.Error(),
			}
		} else {
			item = GetResponse2BulkStateItem(r, state2.GetOriginalStateKey(req.Key))
		}
		// collect result
		select {
		case resultCh <- item:
		default:
			//never happen
			log.DefaultLogger.Errorf("[api.generateGetStateTask] can not push result to the resultCh. item: %+v", item)
		}
	}
}

// converting from BulkGetResponse to BulkStateItem
func BulkGetResponse2BulkStateItem(compResp *state.BulkGetResponse) *dapr_v1pb.BulkStateItem {
	if compResp == nil {
		return &dapr_v1pb.BulkStateItem{}
	}
	return &dapr_v1pb.BulkStateItem{
		Key:      state2.GetOriginalStateKey(compResp.Key),
		Data:     compResp.Data,
		Etag:     common.PointerToString(compResp.ETag),
		Metadata: compResp.Metadata,
		Error:    compResp.Error,
	}
}

// converting from GetResponse to BulkStateItem
func GetResponse2BulkStateItem(compResp *state.GetResponse, key string) *dapr_v1pb.BulkStateItem {
	// convert
	resp := &dapr_v1pb.BulkStateItem{}
	resp.Key = key
	if compResp != nil {
		resp.Data = compResp.Data
		resp.Etag = common.PointerToString(compResp.ETag)
		resp.Metadata = compResp.Metadata
	}
	return resp
}

// converting from DeleteStateRequest to DeleteRequest
func DeleteStateRequest2DeleteRequest(grpcReq *dapr_v1pb.DeleteStateRequest, key string) *state.DeleteRequest {
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
func StateItem2DeleteRequest(grpcReq *dapr_common_v1pb.StateItem, key string) *state.DeleteRequest {
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
