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

package converter

import (
	"github.com/dapr/components-contrib/state"
	"mosn.io/layotto/pkg/common"
	runtime_state "mosn.io/layotto/pkg/runtime/state"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
)

func GetResponse2GetStateResponse(compResp *state.GetResponse) *runtimev1pb.GetStateResponse {
	resp := &runtimev1pb.GetStateResponse{}
	if compResp != nil {
		resp.Etag = common.PointerToString(compResp.ETag)
		resp.Data = compResp.Data
		resp.Metadata = compResp.Metadata
	}
	return resp
}

func GetResponse2BulkStateItem(compResp *state.GetResponse, key string) *runtimev1pb.BulkStateItem {
	resp := &runtimev1pb.BulkStateItem{}
	resp.Key = key
	if compResp != nil {
		resp.Data = compResp.Data
		resp.Etag = common.PointerToString(compResp.ETag)
		resp.Metadata = compResp.Metadata
	}
	return resp
}

func BulkGetResponse2BulkStateItem(compResp *state.BulkGetResponse) *runtimev1pb.BulkStateItem {
	if compResp == nil {
		return &runtimev1pb.BulkStateItem{}
	}
	return &runtimev1pb.BulkStateItem{
		Key:      runtime_state.GetOriginalStateKey(compResp.Key),
		Data:     compResp.Data,
		Etag:     common.PointerToString(compResp.ETag),
		Metadata: compResp.Metadata,
		Error:    compResp.Error,
	}
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
			Consistency: runtime_state.StateConsistencyToString(grpcReq.Options.Consistency),
			Concurrency: runtime_state.StateConcurrencyToString(grpcReq.Options.Concurrency),
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
			Concurrency: runtime_state.StateConcurrencyToString(grpcReq.Options.Concurrency),
			Consistency: runtime_state.StateConsistencyToString(grpcReq.Options.Consistency),
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
			Concurrency: runtime_state.StateConcurrencyToString(grpcReq.Options.Concurrency),
			Consistency: runtime_state.StateConsistencyToString(grpcReq.Options.Consistency),
		}
	}
	return req
}
