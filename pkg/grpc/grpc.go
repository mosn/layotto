package grpc

import (
	runtimev1pb "github.com/layotto/layotto/proto/runtime/v1"
	"google.golang.org/grpc"
	mgrpc "mosn.io/mosn/pkg/filter/network/grpc"
)

func NewGrpcServer(opts ...Option) mgrpc.RegisteredServer {
	var o grpcOptions
	for _, opt := range opts {
		opt(&o)
	}
	srvMaker := NewDefaultServer
	if o.maker != nil {
		srvMaker = o.maker
	}
	return srvMaker(o.api, o.options...)
}

func NewDefaultServer(api API, opts ...grpc.ServerOption) mgrpc.RegisteredServer {
	s := grpc.NewServer(opts...)
	runtimev1pb.RegisterMosnRuntimeServer(s, api)
	return s
}
