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

package client

import (
	"context"

	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
)

type SayHelloRequest struct {
	ServiceName string
}

type SayHelloResp struct {
	Hello string
}

func (c *GRPCClient) SayHello(ctx context.Context, in *SayHelloRequest) (*SayHelloResp, error) {
	req := &runtimev1pb.SayHelloRequest{
		ServiceName: in.ServiceName,
	}
	resp, err := c.protoClient.SayHello(ctx, req)
	if err != nil {
		return nil, err
	}
	return &SayHelloResp{Hello: resp.Hello}, nil
}
