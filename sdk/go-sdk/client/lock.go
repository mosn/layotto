package client

import (
	"context"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
)

func (c *GRPCClient) TryLock(ctx context.Context, req *runtimev1pb.TryLockRequest) (*runtimev1pb.TryLockResponse, error) {
	return c.protoClient.TryLock(ctx, req)
}

func (c *GRPCClient) Unlock(ctx context.Context, req *runtimev1pb.UnlockRequest) (*runtimev1pb.UnlockResponse, error) {
	return c.protoClient.Unlock(ctx, req)
}
