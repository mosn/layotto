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
