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

package pluggable

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCConnectionDialer func(ctx context.Context, opts ...grpc.DialOption) (*grpc.ClientConn, error)

// WithOptions returns a new connection dialer that adds the new options to it.
func (g GRPCConnectionDialer) WithOptions(newOpts ...grpc.DialOption) GRPCConnectionDialer {
	return func(ctx context.Context, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
		return g(ctx, append(opts, newOpts...)...)
	}
}

// SocketDialer creates a dialer for the given socket.
func SocketDialer(socket string, additionalOpts ...grpc.DialOption) GRPCConnectionDialer {
	return func(ctx context.Context, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
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
