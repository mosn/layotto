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
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
)

var factory ClientFactory = &ClientFactoryImpl{}

func RegisterClientFactory(cf ClientFactory) {
	factory = cf
}

type ClientFactory interface {
	NewClientWithAddress(address string) (*grpc.ClientConn, runtimev1pb.RuntimeClient, error)
}

type ClientFactoryImpl struct {
}

// NewClientWithAddress instantiates runtime using specific address (including port).
func (cf *ClientFactoryImpl) NewClientWithAddress(address string) (*grpc.ClientConn, runtimev1pb.RuntimeClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, nil, errors.Wrapf(err, "error creating connection to '%s': %v", address, err)
	}

	return conn, runtimev1pb.NewRuntimeClient(conn), nil
}
