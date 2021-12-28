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

	"github.com/gammazero/workerpool"

	"mosn.io/layotto/pkg/common"

	"github.com/dapr/components-contrib/state"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	dapr_common_v1pb "mosn.io/layotto/pkg/grpc/dapr/proto/common/v1"
	dapr_v1pb "mosn.io/layotto/pkg/grpc/dapr/proto/runtime/v1"
	"mosn.io/layotto/pkg/messages"
	state2 "mosn.io/layotto/pkg/runtime/state"
	"mosn.io/pkg/log"
)

func (d *daprGrpcAPI) SaveState(ctx context.Context, in *dapr_v1pb.SaveStateRequest) (*emptypb.Empty, error) {
	// 1. get store
	store, err := d.getStateStore(in.StoreName)
	if err != nil {
		log.DefaultLogger.Errorf("[runtime] [grpc.SaveState] error: %v", err)
		return &emptypb.Empty{}, err
	}
	// 2. convert requests
	var reqs []state.SetRequest
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

func (d *daprGrpcAPI) GetState(ctx context.Context, request *dapr_v1pb.GetStateRequest) (*dapr_v1pb.GetStateResponse, error) {
	store, err := d.getStateStore(request.StoreName)
	if err != nil {
		log.DefaultLogger.Errorf("[runtime] [grpc.GetState] error: %v", err)
		return nil, err
	}
	key, err := state2.GetModifiedStateKey(request.Key, request.StoreName, d.appId)
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
	panic("implement me")
}

func (d *daprGrpcAPI) DeleteState(ctx context.Context, request *dapr_v1pb.DeleteStateRequest) (*emptypb.Empty, error) {
	panic("implement me")
}

func (d *daprGrpcAPI) DeleteBulkState(ctx context.Context, request *dapr_v1pb.DeleteBulkStateRequest) (*emptypb.Empty, error) {
	panic("implement me")
}

func (d *daprGrpcAPI) ExecuteStateTransaction(ctx context.Context, request *dapr_v1pb.ExecuteStateTransactionRequest) (*emptypb.Empty, error) {
	panic("implement me")
}

func (d *daprGrpcAPI) getStateStore(name string) (state.Store, error) {
	if d.stateStores == nil || len(d.stateStores) == 0 {
		return nil, status.Error(codes.FailedPrecondition, messages.ErrStateStoresNotConfigured)
	}

	if d.stateStores[name] == nil {
		return nil, status.Errorf(codes.InvalidArgument, messages.ErrStateStoreNotFound, name)
	}
	return d.stateStores[name], nil
}

func StateItem2SetRequest(grpcReq *dapr_common_v1pb.StateItem, key string) *state.SetRequest {
	req := &state.SetRequest{
		Key: key,
	}
	if grpcReq == nil {
		return req
	}
	req.Metadata = grpcReq.Metadata
	req.Value = grpcReq.Value
	if grpcReq.Etag != nil {
		req.ETag = &grpcReq.Etag.Value
	}
	if grpcReq.Options != nil {
		req.Options = state.SetStateOption{
			Consistency: StateConsistencyToString(grpcReq.Options.Consistency),
			Concurrency: StateConcurrencyToString(grpcReq.Options.Concurrency),
		}
	}
	return req
}

func GetResponse2GetStateResponse(compResp *state.GetResponse) *dapr_v1pb.GetStateResponse {
	resp := &dapr_v1pb.GetStateResponse{}
	if compResp != nil {
		resp.Etag = common.PointerToString(compResp.ETag)
		resp.Data = compResp.Data
		resp.Metadata = compResp.Metadata
	}
	return resp
}

func StateConsistencyToString(c dapr_common_v1pb.StateOptions_StateConsistency) string {
	switch c {
	case dapr_common_v1pb.StateOptions_CONSISTENCY_EVENTUAL:
		return "eventual"
	case dapr_common_v1pb.StateOptions_CONSISTENCY_STRONG:
		return "strong"
	}
	return ""
}

func StateConcurrencyToString(c dapr_common_v1pb.StateOptions_StateConcurrency) string {
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

func GetResponse2BulkStateItem(compResp *state.GetResponse, key string) *dapr_v1pb.BulkStateItem {
	resp := &dapr_v1pb.BulkStateItem{}
	resp.Key = key
	if compResp != nil {
		resp.Data = compResp.Data
		resp.Etag = common.PointerToString(compResp.ETag)
		resp.Metadata = compResp.Metadata
	}
	return resp
}
