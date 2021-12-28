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
	_ "net/http/pprof"

	"github.com/dapr/components-contrib/state"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/emptypb"
	dapr_common_v1pb "mosn.io/layotto/pkg/grpc/dapr/proto/common/v1"
	dapr_v1pb "mosn.io/layotto/pkg/grpc/dapr/proto/runtime/v1"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"

	"google.golang.org/grpc/status"
)

// GetState obtains the state for a specific key.
func (a *api) GetState(ctx context.Context, in *runtimev1pb.GetStateRequest) (*runtimev1pb.GetStateResponse, error) {
	if in == nil {
		return &runtimev1pb.GetStateResponse{}, status.Error(codes.InvalidArgument, "GetStateRequest is nil")
	}
	daprReq := &dapr_v1pb.GetStateRequest{
		StoreName:   in.GetStoreName(),
		Key:         in.GetKey(),
		Consistency: dapr_common_v1pb.StateOptions_StateConsistency(in.GetConsistency()),
		Metadata:    in.GetMetadata(),
	}
	resp, err := a.daprAPI.GetState(ctx, daprReq)
	if err != nil {
		return &runtimev1pb.GetStateResponse{}, err
	}
	return convertGetStateResponseToRuntimeResponse(resp), nil
}

func (a *api) SaveState(ctx context.Context, in *runtimev1pb.SaveStateRequest) (*emptypb.Empty, error) {
	if in == nil {
		return &emptypb.Empty{}, status.Error(codes.InvalidArgument, "SaveStateRequest is nil")
	}
	// convert request
	daprReq := &dapr_v1pb.SaveStateRequest{
		StoreName: in.StoreName,
		States:    convertStatesToDaprPB(in.States),
	}
	// delegate to dapr api implementation
	return a.daprAPI.SaveState(ctx, daprReq)
}

func (a *api) GetBulkState(ctx context.Context, in *runtimev1pb.GetBulkStateRequest) (*runtimev1pb.GetBulkStateResponse, error) {
	if in == nil {
		return &runtimev1pb.GetBulkStateResponse{}, status.Error(codes.InvalidArgument, "GetBulkStateRequest is nil")
	}
	daprReq := &dapr_v1pb.GetBulkStateRequest{
		StoreName:   in.GetStoreName(),
		Keys:        in.GetKeys(),
		Parallelism: in.GetParallelism(),
		Metadata:    in.GetMetadata(),
	}
	resp, err := a.daprAPI.GetBulkState(ctx, daprReq)
	if err != nil {
		return &runtimev1pb.GetBulkStateResponse{}, err
	}
	return convertGetBulkStateResponseToRuntimeResponse(resp), nil
}

func convertStatesToDaprPB(states []*runtimev1pb.StateItem) []*dapr_common_v1pb.StateItem {
	dStates := make([]*dapr_common_v1pb.StateItem, 0)
	if states == nil {
		return dStates
	}
	for _, s := range states {
		ds := &dapr_common_v1pb.StateItem{
			Key:      s.Key,
			Value:    s.Value,
			Metadata: s.Metadata,
		}
		if s.Etag != nil {
			ds.Etag = &dapr_common_v1pb.Etag{Value: s.Etag.Value}
		}
		if s.Options != nil {
			ds.Options = &dapr_common_v1pb.StateOptions{
				Concurrency: dapr_common_v1pb.StateOptions_StateConcurrency(s.Options.Concurrency),
				Consistency: dapr_common_v1pb.StateOptions_StateConsistency(s.Options.Consistency),
			}
		}
		dStates = append(dStates, ds)
	}
	return dStates
}

func convertGetStateResponseToRuntimeResponse(response *dapr_v1pb.GetStateResponse) *runtimev1pb.GetStateResponse {
	res := &runtimev1pb.GetStateResponse{
		Data:     response.GetData(),
		Etag:     response.GetEtag(),
		Metadata: response.GetMetadata(),
	}
	return res
}

func convertGetBulkStateResponseToRuntimeResponse(resp *dapr_v1pb.GetBulkStateResponse) *runtimev1pb.GetBulkStateResponse {
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
	return ret
}

func StateItem2SetRequest(grpcReq *runtimev1pb.StateItem, key string) *state.SetRequest {
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

func DeleteStateRequest2DeleteRequest(grpcReq *runtimev1pb.DeleteStateRequest, key string) *state.DeleteRequest {
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

func StateItem2DeleteRequest(grpcReq *runtimev1pb.StateItem, key string) *state.DeleteRequest {
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

func StateConsistencyToString(c runtimev1pb.StateOptions_StateConsistency) string {
	switch c {
	case runtimev1pb.StateOptions_CONSISTENCY_EVENTUAL:
		return "eventual"
	case runtimev1pb.StateOptions_CONSISTENCY_STRONG:
		return "strong"
	}

	return ""
}

func StateConcurrencyToString(c runtimev1pb.StateOptions_StateConcurrency) string {
	switch c {
	case runtimev1pb.StateOptions_CONCURRENCY_FIRST_WRITE:
		return "first-write"
	case runtimev1pb.StateOptions_CONCURRENCY_LAST_WRITE:
		return "last-write"
	}

	return ""
}
