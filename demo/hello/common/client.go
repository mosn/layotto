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

package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"mosn.io/layotto/sdk/go-sdk/client"
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
		log.Fatal(err)
	}

	defer cli.Close()

	res, err := cli.SayHello(context.Background(), &client.SayHelloRequest{
		ServiceName: storeName,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res.Hello)
}
