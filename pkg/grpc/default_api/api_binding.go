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

	dapr_v1pb "mosn.io/layotto/pkg/grpc/dapr/proto/runtime/v1"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
)

func (a *api) InvokeBinding(ctx context.Context, in *runtimev1pb.InvokeBindingRequest) (*runtimev1pb.InvokeBindingResponse, error) {
	daprResp, err := a.daprAPI.InvokeBinding(ctx, &dapr_v1pb.InvokeBindingRequest{
		Name:      in.Name,
		Data:      in.Data,
		Metadata:  in.Metadata,
		Operation: in.Operation,
	})
	if err != nil {
		return &runtimev1pb.InvokeBindingResponse{}, err
	}
	return &runtimev1pb.InvokeBindingResponse{
		Data:     daprResp.Data,
		Metadata: daprResp.Metadata,
	}, nil
}
