package main

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/layotto/layotto/spec/proto/runtime/v1"
	"google.golang.org/grpc"
	"net"
	"strconv"
)

const topicName = "topic1"

func main() {
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
	runtime.RegisterAppCallbackServer(grpcServer, &AppCallbackServerImpl{})
	fmt.Printf("Start listening on port %v ...... \n", port)
	err = grpcServer.Serve(lis)
	if err != nil {
		panic(err)
	}
}

type AppCallbackServerImpl struct {
}

func (a *AppCallbackServerImpl) ListTopicSubscriptions(ctx context.Context, empty *empty.Empty) (*runtime.ListTopicSubscriptionsResponse, error) {
	result := &runtime.ListTopicSubscriptionsResponse{}
	ts := &runtime.TopicSubscription{
		PubsubName: "redis",
		Topic:      topicName,
		Metadata:   nil,
	}
	result.Subscriptions = append(result.Subscriptions, ts)
	return result, nil
}

func (a *AppCallbackServerImpl) OnTopicEvent(ctx context.Context, request *runtime.TopicEventRequest) (*runtime.TopicEventResponse, error) {
	fmt.Printf("Received a new event.Topic: %s , Data:%s \n", request.Topic, request.Data)
	return &runtime.TopicEventResponse{}, nil
}
