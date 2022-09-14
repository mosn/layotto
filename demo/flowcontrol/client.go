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
	"time"

	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:34904", grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(2*time.Second))
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	c := runtimev1pb.NewRuntimeClient(conn)

	for i := 0; i < 10; i++ {
		r, err := c.SayHello(context.Background(), &runtimev1pb.SayHelloRequest{
			ServiceName: "helloworld",
		})
		if err != nil {
			fmt.Println("get an error: ", err)
		} else {
			fmt.Println("get a message: ", r.GetHello())
		}
	}

}
