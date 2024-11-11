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
	"errors"
	"flag"
	"log"
	"time"

	"mosn.io/layotto/sdk/go-sdk/client"

	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
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
	testDynamicSub()
}

func testDynamicSub() {
	cli, err := client.NewClient()

	// Use SubscribeWithHandler API to subscribe to a topic.
	stop, err := cli.SubscribeWithHandler(context.Background(), client.SubscriptionRequest{
		PubsubName: storeName,
		Topic:      "hello",
		Metadata:   nil,
	}, eventHandler)

	//Use Subscribe API to subscribe to a topic.
	sub, err := cli.Subscribe(context.Background(), client.SubscriptionRequest{
		PubsubName: storeName,
		Topic:      "topic1",
		Metadata:   nil,
	})

	if err != nil {
		log.Fatalf("failed to subscribe to topic: %v", err)
	}

	msg, err := sub.Receive()
	if err != nil {
		log.Fatalf("Error receiving message: %v", err)
	}
	log.Printf(">>[Subscribe API]Received message\n")
	log.Printf("event - PubsubName: %s, Topic: %s, ID: %s, Data: %s\n", msg.PubsubName, msg.Topic, msg.Id, msg.Data)

	// Use _MUST_ always signal the result of processing the message, else the
	// message will not be considered as processed and will be redelivered or
	// dead lettered.
	if err := msg.Success(); err != nil {
		log.Fatalf("error sending message success: %v", err)
	}

	time.Sleep(time.Second * 10)

	if err := errors.Join(stop(), sub.Close()); err != nil {
		log.Fatal(err)
	}

	if err != nil {
		panic(err)
	}
	cli.Close()
}

func eventHandler(request *runtimev1pb.TopicEventRequest) client.SubscriptionResponseStatus {
	log.Printf(">>[SubscribeWithHandler API] Received message\n")
	log.Printf("event - PubsubName: %s, Topic: %s, ID: %s, Data: %s\n", request.PubsubName, request.Topic, request.Id, request.Data)
	return client.SubscriptionResponseStatusSuccess
}
