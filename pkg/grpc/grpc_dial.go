package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type GRPCConnectionDialer func(ctx context.Context, name string, opts ...grpc.DialOption) (*grpc.ClientConn, error)

// metadataInstanceID is used to differentiate between multiples instance of the same component.
const metadataInstanceID = "x-component-instance"

// instanceIDStreamInterceptor returns a grpc client unary interceptor that adds the instanceID on outgoing metadata.
// instanceID is used for multiplexing connection if the component supports it.
func instanceIDUnaryInterceptor(instanceID string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return invoker(metadata.AppendToOutgoingContext(ctx, metadataInstanceID, instanceID), method, req, reply, cc, opts...)
	}
}

// instanceIDStreamInterceptor returns a grpc client stream interceptor that adds the instanceID on outgoing metadata.
// instanceID is used for multiplexing connection if the component supports it.
func instanceIDStreamInterceptor(instanceID string) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		return streamer(metadata.AppendToOutgoingContext(ctx, metadataInstanceID, instanceID), desc, cc, method, opts...)
	}
}

// SocketDialer creates a dialer for the given socket.
func SocketDialer(socket string, additionalOpts ...grpc.DialOption) GRPCConnectionDialer {
	return func(ctx context.Context, name string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
		additionalOpts = append(additionalOpts, grpc.WithStreamInterceptor(instanceIDStreamInterceptor(name)), grpc.WithUnaryInterceptor(instanceIDUnaryInterceptor(name)))
		return SocketDial(ctx, socket, append(additionalOpts, opts...)...)
	}
}

// SocketDial creates a grpc connection using the given socket.
func SocketDial(ctx context.Context, socket string, additionalOpts ...grpc.DialOption) (*grpc.ClientConn, error) {
	udsSocket := "unix://" + socket
	//log.Debugf("using socket defined at '%s'", udsSocket)
	additionalOpts = append(additionalOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	grpcConn, err := grpc.DialContext(ctx, udsSocket, additionalOpts...)
	if err != nil {
		return nil, fmt.Errorf("unable to open GRPC connection using socket '%s': %w", udsSocket, err)
	}
	return grpcConn, nil
}
