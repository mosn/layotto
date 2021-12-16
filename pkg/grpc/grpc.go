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
	mgrpc "mosn.io/mosn/pkg/filter/network/grpc"
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
	s := grpc.NewServer(opts...)
	// create registeredServer to manage lifecycle of the grpc server
	var registeredServer mgrpc.RegisteredServer = s
	var err error = nil
	// loop registering grpc api
	for _, grpcAPI := range apis {
		registeredServer, err = grpcAPI.Register(s, registeredServer)
		if err != nil {
			return registeredServer, err
		}
	}
	return registeredServer, nil
}
