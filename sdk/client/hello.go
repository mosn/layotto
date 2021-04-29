package client

import (
	"context"
	"gitlab.alipay-inc.com/ant-mesh/runtime/proto/runtime/v1"
)

type SayHelloRequest struct {
	ServiceName string
}

type SayHelloResp struct {
	Hello string
}

func (c *GRPCClient) SayHello(ctx context.Context, in *SayHelloRequest) (*SayHelloResp, error) {
	req := &runtime.SayHelloRequest{ServiceName: in.ServiceName}
	resp, err := c.protoClient.SayHello(ctx, req)
	if err != nil {
		return nil, err
	}
	return &SayHelloResp{Hello: resp.Hello}, nil
}
