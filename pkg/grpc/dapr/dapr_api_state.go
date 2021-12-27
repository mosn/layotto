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
	if in == nil {
		return &emptypb.Empty{}, status.Error(codes.InvalidArgument, "SaveStateRequest is nil")
	}
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
