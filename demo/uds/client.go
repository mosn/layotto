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

package main

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
)

const (
	protocol = "unix"
	sockAddr = "/tmp/client-proxy.sock"
)

func main() {
	var (
		credentials = insecure.NewCredentials() // No SSL/TLS
		dialer      = func(ctx context.Context, addr string) (net.Conn, error) {
			var d net.Dialer
			return d.DialContext(ctx, protocol, addr)
		}
		grpcoptions = []grpc.DialOption{
			grpc.WithTransportCredentials(credentials),
			grpc.WithContextDialer(dialer),
		}
	)
	conn, err := grpc.Dial(sockAddr, grpcoptions...)
	if err != nil {
		panic(err)
	}
	client := runtimev1pb.NewRuntimeClient(conn)
	hello, err := client.SayHello(context.Background(), &runtimev1pb.SayHelloRequest{
		ServiceName: "helloworld",
		Name:        "eva",
	})
	if err != nil {
		panic(err)
	}
	if hello.Hello != "greeting, eva" {
		panic(fmt.Errorf("assertion failed! result is %v", hello.Hello))
	}
	fmt.Println(hello.Hello)
}
