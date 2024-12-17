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
	"net"
	"strconv"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"

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
	// start a grpc server for callback
	testSub()
}

func testSub() {
	port := 9999
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic("failed to listen on port " + strconv.Itoa(port))
	}
	grpcServer := grpc.NewServer()
	runtimev1pb.RegisterAppCallbackServer(grpcServer, &AppCallbackServerImpl{})
	fmt.Printf("Start listening on port %v ...... \n", port)
	err = grpcServer.Serve(lis)
	if err != nil {
		panic(err)
	}
}

type AppCallbackServerImpl struct {
}

func (a *AppCallbackServerImpl) ListTopicSubscriptions(ctx context.Context, empty *empty.Empty) (*runtimev1pb.ListTopicSubscriptionsResponse, error) {
	result := &runtimev1pb.ListTopicSubscriptionsResponse{}
	ts := &runtimev1pb.TopicSubscription{
		PubsubName: storeName,
		Topic:      "hello",
		Metadata:   nil,
	}
	result.Subscriptions = append(result.Subscriptions, ts)
	result.Subscriptions = append(result.Subscriptions, &runtimev1pb.TopicSubscription{
		PubsubName: storeName,
		Topic:      "topic1",
		Metadata:   nil,
	})
	return result, nil
}

func (a *AppCallbackServerImpl) OnTopicEvent(ctx context.Context, request *runtimev1pb.TopicEventRequest) (*runtimev1pb.TopicEventResponse, error) {
	fmt.Printf("Received a new event.Topic: %s , Data: %s \n", request.Topic, request.Data)
	return &runtimev1pb.TopicEventResponse{}, nil
}
