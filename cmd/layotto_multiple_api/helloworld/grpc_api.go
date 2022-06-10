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

package helloworld

import (
	"context"
	"fmt"

	rawGRPC "google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"mosn.io/pkg/log"

	"mosn.io/layotto/cmd/layotto_multiple_api/helloworld/component"
	"mosn.io/layotto/components/lock"
	"mosn.io/layotto/pkg/grpc"
	grpc_api "mosn.io/layotto/pkg/grpc"
)

const kind = "helloworld"

// This demo will always use this component name.
const componentName = "demo"

func NewHelloWorldAPI(ac *grpc_api.ApplicationContext) grpc.GrpcAPI {
	// 1. convert custom components
	name2component := make(map[string]component.HelloWorld)
	if len(ac.CustomComponent) != 0 {
		// we only care about those components of type "helloworld"
		name2comp, ok := ac.CustomComponent[kind]
		if ok && len(name2comp) > 0 {
			for name, v := range name2comp {
				// convert them using type assertion
				comp, ok := v.(component.HelloWorld)
				if !ok {
					errMsg := fmt.Sprintf("custom component %s does not implement HelloWorld interface", name)
					log.DefaultLogger.Errorf(errMsg)
				}
				name2component[name] = comp
			}
		}
	}
	// 2. construct your API implementation
	return &server{
		appId: ac.AppId,
		// Your API plugin can store and use all the components.
		// For example,this demo set all the LockStore components here.
		name2LockStore: ac.LockStores,
		// Custom components of type "helloworld"
		name2component: name2component,
	}
}

// server is used to implement helloworld.GreeterServer.
type server struct {
	appId string
	// custom components which implements the `HelloWorld` interface
	name2component map[string]component.HelloWorld
	// LockStore components. They are not used in this demo, we put them here as a demo.
	name2LockStore map[string]lock.LockStore
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer.SayHello
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	if _, ok := s.name2component[componentName]; !ok {
		return &pb.HelloReply{Message: "We don't want to talk with you!"}, nil
	}
	message, err := s.name2component[componentName].SayHello(in.GetName())
	if err != nil {
		return nil, err
	}
	return &pb.HelloReply{Message: message}, nil
}

func (s *server) Init(conn *rawGRPC.ClientConn) error {
	return nil
}

func (s *server) Register(rawGrpcServer *rawGRPC.Server) error {
	pb.RegisterGreeterServer(rawGrpcServer, s)
	return nil
}
