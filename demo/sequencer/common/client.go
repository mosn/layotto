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

	client "mosn.io/layotto/sdk/go-sdk/client"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
)

const (
	key = "key666"
)

var storeName string

func init() {
	flag.StringVar(&storeName, "s", "", "set `storeName`")
}

func main() {
	flag.Parse()
	if storeName == "" {
		panic("storeName is empty.")
	}

	cli, err := client.NewClient()
	if err != nil {
		panic(err)
	}
	defer cli.Close()
	ctx := context.Background()
	fmt.Printf("Try to get next id.Key:%s \n", key)
	for i := 0; i < 10; i++ {
		id, err := cli.GetNextId(ctx, &runtimev1pb.GetNextIdRequest{
			StoreName: storeName,
			Key:       key,
			Options:   nil,
			Metadata:  nil,
		})
		if err != nil {
			panic(err)
		}
		fmt.Printf("Next id:%v \n", id)
	}
	fmt.Println("Demo success!")
}
