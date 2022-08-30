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
)

// GrpcAPI is the interface of API plugin. It has lifecycle related methods
type GrpcAPI interface {
	// init this API before binding it to the grpc server.
	// For example,you can call app to query their subscriptions.
	Init(conn *grpc.ClientConn) error

	// Bind this API to the grpc server
	Register(rawGrpcServer *grpc.Server) error
}

// NewGrpcAPI is the constructor of GrpcAPI
type NewGrpcAPI func(applicationContext *ApplicationContext) GrpcAPI
