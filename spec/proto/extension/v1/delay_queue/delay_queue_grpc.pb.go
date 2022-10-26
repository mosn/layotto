// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.17.3
// source: spec/proto/extension/v1/delay_queue/delay_queue.proto

package delay_queue

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// DelayQueueClient is the client API for DelayQueue service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DelayQueueClient interface {
	// Publish a delay message
	PublishDelayMessage(ctx context.Context, in *DelayMessageRequest, opts ...grpc.CallOption) (*DelayMessageResponse, error)
}

type delayQueueClient struct {
	cc grpc.ClientConnInterface
}

func NewDelayQueueClient(cc grpc.ClientConnInterface) DelayQueueClient {
	return &delayQueueClient{cc}
}

func (c *delayQueueClient) PublishDelayMessage(ctx context.Context, in *DelayMessageRequest, opts ...grpc.CallOption) (*DelayMessageResponse, error) {
	out := new(DelayMessageResponse)
	err := c.cc.Invoke(ctx, "/spec.proto.extension.v1.delay_queue.DelayQueue/PublishDelayMessage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DelayQueueServer is the server API for DelayQueue service.
// All implementations should embed UnimplementedDelayQueueServer
// for forward compatibility
type DelayQueueServer interface {
	// Publish a delay message
	PublishDelayMessage(context.Context, *DelayMessageRequest) (*DelayMessageResponse, error)
}

// UnimplementedDelayQueueServer should be embedded to have forward compatible implementations.
type UnimplementedDelayQueueServer struct {
}

func (UnimplementedDelayQueueServer) PublishDelayMessage(context.Context, *DelayMessageRequest) (*DelayMessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PublishDelayMessage not implemented")
}

// UnsafeDelayQueueServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DelayQueueServer will
// result in compilation errors.
type UnsafeDelayQueueServer interface {
	mustEmbedUnimplementedDelayQueueServer()
}

func RegisterDelayQueueServer(s grpc.ServiceRegistrar, srv DelayQueueServer) {
	s.RegisterService(&DelayQueue_ServiceDesc, srv)
}

func _DelayQueue_PublishDelayMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DelayMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DelayQueueServer).PublishDelayMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/spec.proto.extension.v1.delay_queue.DelayQueue/PublishDelayMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DelayQueueServer).PublishDelayMessage(ctx, req.(*DelayMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// DelayQueue_ServiceDesc is the grpc.ServiceDesc for DelayQueue service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DelayQueue_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "spec.proto.extension.v1.delay_queue.DelayQueue",
	HandlerType: (*DelayQueueServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PublishDelayMessage",
			Handler:    _DelayQueue_PublishDelayMessage_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "spec/proto/extension/v1/delay_queue/delay_queue.proto",
}
