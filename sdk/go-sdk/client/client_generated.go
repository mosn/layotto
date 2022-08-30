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

package client

import (
	"google.golang.org/grpc"

	v1 "mosn.io/layotto/spec/proto/runtime/v1"

	"mosn.io/layotto/spec/proto/extension/v1/s3"
)

// Client is the interface for runtime client implementation.
type Client interface {
	runtimeAPI

	s3.ObjectStorageServiceClient
}

// NewClientWithConnection instantiates runtime client using specific connection.
func NewClientWithConnection(conn *grpc.ClientConn) Client {
	return &GRPCClient{
		connection:                 conn,
		protoClient:                v1.NewRuntimeClient(conn),
		ObjectStorageServiceClient: s3.NewObjectStorageServiceClient(conn),
	}
}

// GRPCClient is the gRPC implementation of runtime client.
type GRPCClient struct {
	connection  *grpc.ClientConn
	protoClient v1.RuntimeClient
	s3.ObjectStorageServiceClient
}
