package pubsub

import (
	"context"
	runtimev1pb "github.com/layotto/layotto/spec/proto/runtime/v1"
	"google.golang.org/protobuf/types/known/emptypb"
	"mosn.io/pkg/log"
)

func ListTopicSubscriptions(client runtimev1pb.AppCallbackClient, log log.ErrorLogger) []*runtimev1pb.TopicSubscription {
	resp, err := client.ListTopicSubscriptions(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Errorf("[runtime][ListTopicSubscriptions]error after callback: %s", err)
		return make([]*runtimev1pb.TopicSubscription, 0)
	}
	if resp != nil && len(resp.Subscriptions) > 0 {
		return resp.Subscriptions
	}
	return make([]*runtimev1pb.TopicSubscription, 0)
}
