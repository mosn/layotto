// Code generated by github.com/seeflood/protoc-gen-p6 .

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

package client

import (
	context "context"

	grpc "google.golang.org/grpc"

	email "mosn.io/layotto/spec/proto/extension/v1/email"
	phone "mosn.io/layotto/spec/proto/extension/v1/phone"
	s3 "mosn.io/layotto/spec/proto/extension/v1/s3"
	v1 "mosn.io/layotto/spec/proto/runtime/v1"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context

// Client is the interface for runtime client implementation.
type Client interface {
	runtimeAPI

	s3.ObjectStorageServiceClient

	// "mosn.io/layotto/spec/proto/extension/v1/email"
	email.EmailServiceClient

	// "mosn.io/layotto/spec/proto/extension/v1/phone"
	phone.PhoneCallServiceClient
}

// NewClientWithConnection instantiates runtime client using specific connection.
func NewClientWithConnection(conn *grpc.ClientConn) Client {
	return &GRPCClient{
		connection:                 conn,
		protoClient:                v1.NewRuntimeClient(conn),
		ObjectStorageServiceClient: s3.NewObjectStorageServiceClient(conn),
		// "mosn.io/layotto/spec/proto/extension/v1/email"
		EmailServiceClient: email.NewEmailServiceClient(conn),

		// "mosn.io/layotto/spec/proto/extension/v1/phone"
		PhoneCallServiceClient: phone.NewPhoneCallServiceClient(conn),
	}
}

// GRPCClient is the gRPC implementation of runtime client.
type GRPCClient struct {
	connection  *grpc.ClientConn
	protoClient v1.RuntimeClient
	s3.ObjectStorageServiceClient
	// "mosn.io/layotto/spec/proto/extension/v1/email"
	email.EmailServiceClient
	// "mosn.io/layotto/spec/proto/extension/v1/phone"
	phone.PhoneCallServiceClient
}
