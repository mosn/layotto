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
	mgrpc "mosn.io/mosn/pkg/filter/network/grpc"

	"mosn.io/layotto/diagnostics"
)

func NewGrpcServer(opts ...Option) (mgrpc.RegisteredServer, error) {
	var o grpcOptions
	for _, opt := range opts {
		opt(&o)
	}
	srvMaker := NewDefaultServer
	o.options = append(o.options, grpc.ChainUnaryInterceptor(diagnostics.UnaryInterceptorFilter))
	o.options = append(o.options, grpc.ChainStreamInterceptor(diagnostics.StreamInterceptorFilter))
	if o.maker != nil {
		srvMaker = o.maker
	}
	return srvMaker(o.apis, o.options...)
}

func NewDefaultServer(apis []GrpcAPI, opts ...grpc.ServerOption) (mgrpc.RegisteredServer, error) {
	return NewRawGrpcServer(apis, opts...)
}

func NewRawGrpcServer(apis []GrpcAPI, opts ...grpc.ServerOption) (*grpc.Server, error) {
	s := grpc.NewServer(opts...)
	var err error
	// loop registering grpc api
	for _, grpcAPI := range apis {
		err = grpcAPI.Register(s)
		if err != nil {
			return s, err
		}
	}
	return s, nil
}
