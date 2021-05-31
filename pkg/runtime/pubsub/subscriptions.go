// REF: https://github.com/dapr/dapr/blob/master/pkg/runtime/pubsub/subscriptions.go
package pubsub

import (
	"context"
	runtimev1pb "github.com/layotto/layotto/spec/proto/runtime/v1"
	"google.golang.org/protobuf/types/known/emptypb"
	"mosn.io/pkg/log"
)

const (
	getTopicsError         = "error getting topic list from app: %s"
	deserializeTopicsError = "error getting topics from app: %s"
	noSubscriptionsError   = "user app did not subscribe to any topic"
	subscriptionKind       = "Subscription"
)

func GetSubscriptionsGRPC(channel runtimev1pb.AppCallbackClient, log log.ErrorLogger) []Subscription {
	var subscriptions []Subscription

	resp, err := channel.ListTopicSubscriptions(context.Background(), &emptypb.Empty{})
	if err != nil {
		// Unexpected response: both GRPC and HTTP have to log the same level.
		log.Errorf(getTopicsError, err)
	} else {
		if resp == nil || resp.Subscriptions == nil || len(resp.Subscriptions) == 0 {
			log.Debugf(noSubscriptionsError)
		} else {
			for _, s := range resp.Subscriptions {
				subscriptions = append(subscriptions, Subscription{
					PubsubName: s.PubsubName,
					Topic:      s.GetTopic(),
					Metadata:   s.GetMetadata(),
				})
			}
		}
	}
	return subscriptions
}
