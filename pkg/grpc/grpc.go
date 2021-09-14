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

package grpc

import (
	"google.golang.org/grpc"
	"mosn.io/layotto/diagnostics"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
	mgrpc "mosn.io/mosn/pkg/filter/network/grpc"
)

func NewGrpcServer(opts ...Option) mgrpc.RegisteredServer {
	var o grpcOptions
	for _, opt := range opts {
		opt(&o)
	}
	srvMaker := NewDefaultServer
	o.options = append(o.options, grpc.ChainUnaryInterceptor(diagnostics.UnaryInterceptorFilter))
	if o.maker != nil {
		srvMaker = o.maker
	}
	return srvMaker(o.api, o.options...)
}

func NewDefaultServer(api API, opts ...grpc.ServerOption) mgrpc.RegisteredServer {
	s := grpc.NewServer(opts...)
	runtimev1pb.RegisterRuntimeServer(s, api)
	return s
}
