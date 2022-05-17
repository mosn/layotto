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

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	dapr_common_v1pb "mosn.io/layotto/pkg/grpc/dapr/proto/common/v1"
	dapr_v1pb "mosn.io/layotto/pkg/grpc/dapr/proto/runtime/v1"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
)

// GetState obtains the state for a specific key.
func (a *api) GetState(ctx context.Context, in *runtimev1pb.GetStateRequest) (*runtimev1pb.GetStateResponse, error) {
	// Check if the StateRequest is exists
	if in == nil {
		return &runtimev1pb.GetStateResponse{}, status.Error(codes.InvalidArgument, "GetStateRequest is nil")
	}
	// convert request
	daprReq := &dapr_v1pb.GetStateRequest{
		StoreName:   in.GetStoreName(),
		Key:         in.GetKey(),
		Consistency: dapr_common_v1pb.StateOptions_StateConsistency(in.GetConsistency()),
		Metadata:    in.GetMetadata(),
	}
	// Generate response by request
	resp, err := a.daprAPI.GetState(ctx, daprReq)
	if err != nil {
		return &runtimev1pb.GetStateResponse{}, err
	}
	return &runtimev1pb.GetStateResponse{
		Data:     resp.GetData(),
		Etag:     resp.GetEtag(),
		Metadata: resp.GetMetadata(),
	}, nil
}

func (a *api) SaveState(ctx context.Context, in *runtimev1pb.SaveStateRequest) (*emptypb.Empty, error) {
	// Check if the request is nil
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

// GetBulkState gets a batch of state data
func (a *api) GetBulkState(ctx context.Context, in *runtimev1pb.GetBulkStateRequest) (*runtimev1pb.GetBulkStateResponse, error) {
	if in == nil {
		return &runtimev1pb.GetBulkStateResponse{}, status.Error(codes.InvalidArgument, "GetBulkStateRequest is nil")
	}
	// convert request
	daprReq := &dapr_v1pb.GetBulkStateRequest{
		StoreName:   in.GetStoreName(),
		Keys:        in.GetKeys(),
		Parallelism: in.GetParallelism(),
		Metadata:    in.GetMetadata(),
	}
	// Generate response by request
	resp, err := a.daprAPI.GetBulkState(ctx, daprReq)
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
	if in == nil {
		return &emptypb.Empty{}, status.Error(codes.InvalidArgument, "DeleteStateRequest is nil")
	}
	// convert request
	daprReq := &dapr_v1pb.DeleteStateRequest{
		StoreName: in.GetStoreName(),
		Key:       in.GetKey(),
		Etag:      convertEtagToDaprPB(in.Etag),
		Options:   convertOptionsToDaprPB(in.Options),
		Metadata:  in.GetMetadata(),
	}
	return a.daprAPI.DeleteState(ctx, daprReq)
}

func (a *api) DeleteBulkState(ctx context.Context, in *runtimev1pb.DeleteBulkStateRequest) (*empty.Empty, error) {
	if in == nil {
		return &emptypb.Empty{}, status.Error(codes.InvalidArgument, "DeleteBulkStateRequest is nil")
	}
	// convert request
	daprReq := &dapr_v1pb.DeleteBulkStateRequest{
		StoreName: in.GetStoreName(),
		States:    convertStatesToDaprPB(in.States),
	}
	return a.daprAPI.DeleteBulkState(ctx, daprReq)
}

func (a *api) ExecuteStateTransaction(ctx context.Context, in *runtimev1pb.ExecuteStateTransactionRequest) (*emptypb.Empty, error) {
	if in == nil {
		return &emptypb.Empty{}, status.Error(codes.InvalidArgument, "ExecuteStateTransactionRequest is nil")
	}
	// convert request
	daprReq := &dapr_v1pb.ExecuteStateTransactionRequest{
		StoreName:  in.GetStoreName(),
		Operations: convertTransactionalStateOperationToDaprPB(in.Operations),
		Metadata:   in.GetMetadata(),
	}
	return a.daprAPI.ExecuteStateTransaction(ctx, daprReq)
}

// some code for converting from runtimev1pb to dapr_common_v1pb

func convertEtagToDaprPB(etag *runtimev1pb.Etag) *dapr_common_v1pb.Etag {
	if etag == nil {
		return &dapr_common_v1pb.Etag{}
	}
	return &dapr_common_v1pb.Etag{Value: etag.GetValue()}
}
func convertOptionsToDaprPB(op *runtimev1pb.StateOptions) *dapr_common_v1pb.StateOptions {
	if op == nil {
		return &dapr_common_v1pb.StateOptions{}
	}
	return &dapr_common_v1pb.StateOptions{
		Concurrency: dapr_common_v1pb.StateOptions_StateConcurrency(op.Concurrency),
		Consistency: dapr_common_v1pb.StateOptions_StateConsistency(op.Consistency),
	}
}

func convertStatesToDaprPB(states []*runtimev1pb.StateItem) []*dapr_common_v1pb.StateItem {
	dStates := make([]*dapr_common_v1pb.StateItem, 0, len(states))
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
			ds.Etag = convertEtagToDaprPB(s.Etag)
		}
		if s.Options != nil {
			ds.Options = convertOptionsToDaprPB(s.Options)
		}
		dStates = append(dStates, ds)
	}
	return dStates
}

func convertTransactionalStateOperationToDaprPB(ops []*runtimev1pb.TransactionalStateOperation) []*dapr_v1pb.TransactionalStateOperation {
	ret := make([]*dapr_v1pb.TransactionalStateOperation, 0, len(ops))
	for i := 0; i < len(ops); i++ {
		op := ops[i]
		var req *dapr_common_v1pb.StateItem
		if op.Request != nil {
			req = &dapr_common_v1pb.StateItem{
				Key:      op.GetRequest().GetKey(),
				Value:    op.GetRequest().GetValue(),
				Etag:     convertEtagToDaprPB(op.GetRequest().GetEtag()),
				Metadata: op.GetRequest().GetMetadata(),
				Options:  convertOptionsToDaprPB(op.GetRequest().GetOptions()),
			}
		}
		ret = append(ret, &dapr_v1pb.TransactionalStateOperation{
			OperationType: op.OperationType,
			Request:       req,
		})
	}
	return ret
}
