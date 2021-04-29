package grpc

import (
	"google.golang.org/grpc"
	mgrpc "mosn.io/mosn/pkg/filter/network/grpc"
)

type grpcOptions struct {
	api     API
	maker   NewServer
	options []grpc.ServerOption
}

type Option func(o *grpcOptions)

// add an api for grpc server.
// api CANNOT be nil.
func WithAPI(api API) Option {
	return func(o *grpcOptions) {
		o.api = api
	}
}

type NewServer func(api API, opts ...grpc.ServerOption) mgrpc.RegisteredServer

func WithNewServer(f NewServer) Option {
	return func(o *grpcOptions) {
		o.maker = f
	}
}

func WithGrpcOptions(options ...grpc.ServerOption) Option {
	return func(o *grpcOptions) {
		o.options = append(o.options, options...)
	}
}
