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
	client "mosn.io/layotto/sdk/go-sdk/client"
)

const topicName = "topic1"

func main() {
	// 1. construct client
	cli, err := client.NewClient()
	if err != nil {
		panic(err)
	}
	// 2. publish a new event
	testPublish(cli)
	cli.Close()
}

func testPublish(cli client.Client) error {
	data := []byte("value1")
	err := cli.PublishEvent(context.Background(), "redis",topicName, data)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Published a new event.Topic: %s ,Data: %s \n", topicName, data)
	return err
}
