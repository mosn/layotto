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

func (a *api) GetSecret(ctx context.Context, in *runtimev1pb.GetSecretRequest) (*runtimev1pb.GetSecretResponse, error) {
	daprResp, err := a.daprAPI.GetSecret(ctx, &dapr_v1pb.GetSecretRequest{
		StoreName: in.StoreName,
		Key:       in.Key,
		Metadata:  in.Metadata,
	})
	if err != nil {
		return &runtimev1pb.GetSecretResponse{}, err
	}
	return &runtimev1pb.GetSecretResponse{Data: daprResp.Data}, nil
}

func (a *api) GetBulkSecret(ctx context.Context, in *runtimev1pb.GetBulkSecretRequest) (*runtimev1pb.GetBulkSecretResponse, error) {
	daprResp, err := a.daprAPI.GetBulkSecret(ctx, &dapr_v1pb.GetBulkSecretRequest{
		StoreName: in.StoreName,
		Metadata:  in.Metadata,
	})
	if err != nil {
		return &runtimev1pb.GetBulkSecretResponse{}, err
	}
	return &runtimev1pb.GetBulkSecretResponse{
		Data: convertSecretResponseMap(daprResp.Data),
	}, nil
}

func convertSecretResponseMap(data map[string]*dapr_v1pb.SecretResponse) map[string]*runtimev1pb.SecretResponse {
	if data == nil {
		return nil
	}
	result := make(map[string]*runtimev1pb.SecretResponse)
	for k, v := range data {
		var converted *runtimev1pb.SecretResponse
		if v != nil {
			converted = &runtimev1pb.SecretResponse{
				Secrets: v.Secrets,
			}
		}
		result[k] = converted
	}
	return result
}
