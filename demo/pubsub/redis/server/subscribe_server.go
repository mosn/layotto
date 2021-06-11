package main

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
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
		PubsubName: "redis",
		Topic:      topicName,
		Metadata:   nil,
	}
	result.Subscriptions = append(result.Subscriptions, ts)
	return result, nil
}

func (a *AppCallbackServerImpl) OnTopicEvent(ctx context.Context, request *runtimev1pb.TopicEventRequest) (*runtimev1pb.TopicEventResponse, error) {
	fmt.Printf("Received a new event.Topic: %s , Data:%s \n", request.Topic, request.Data)
	return &runtimev1pb.TopicEventResponse{}, nil
}
