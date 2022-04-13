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
	"flag"
	"fmt"
	"log"

	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/anypb"
)

func main() {
	data := flag.String("d", `{"jsonrpc":"2.0","method":"GetUser","params":["A003"],"id":9527}`, "-d")
	flag.Parse()

	conn, err := grpc.Dial("localhost:34904", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	cli := runtimev1pb.NewRuntimeClient(conn)
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	resp, err := cli.InvokeService(
		ctx,
		&runtimev1pb.InvokeServiceRequest{
			Id: "org.apache.dubbo.samples.UserProvider",
			Message: &runtimev1pb.CommonInvokeRequest{
				Method:        "GetUser",
				ContentType:   "",
				Data:          &anypb.Any{Value: []byte(*data)},
				HttpExtension: &runtimev1pb.HTTPExtension{Verb: runtimev1pb.HTTPExtension_POST},
			},
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(resp.Data.GetValue()))
}
