package client

import (
	"context"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
)

func (c *GRPCClient) GetNextId(ctx context.Context, req *runtimev1pb.GetNextIdRequest) (*runtimev1pb.GetNextIdResponse, error) {
	return c.protoClient.GetNextId(ctx, req)
}
