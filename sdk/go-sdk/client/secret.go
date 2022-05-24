// Copyright 2021 Layotto Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"context"

	"google.golang.org/grpc"

	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
)

func (c *GRPCClient) GetSecret(ctx context.Context, in *runtimev1pb.GetSecretRequest, opts ...grpc.CallOption) (*runtimev1pb.GetSecretResponse, error) {

	return c.protoClient.GetSecret(ctx, in)
}
func (c *GRPCClient) GetBulkSecret(ctx context.Context, in *runtimev1pb.GetBulkSecretRequest, opts ...grpc.CallOption) (*runtimev1pb.GetBulkSecretResponse, error) {

	return c.protoClient.GetBulkSecret(ctx, in)
}
