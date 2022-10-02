package servicebus

import (
	"context"
	azservicebus "github.com/dapr/components-contrib/pubsub/azure/servicebus"
	delay_queue "mosn.io/layotto/components/delay_queue"

	"github.com/dapr/components-contrib/pubsub"
	"github.com/dapr/kit/logger"
)

type azureServiceBus struct {
	pubsub.PubSub
}

// NewAzureServiceBus returns a new Azure ServiceBus pub-sub implementation.
func NewAzureServiceBus(logger logger.Logger) pubsub.PubSub {
	return &azureServiceBus{
		PubSub: azservicebus.NewAzureServiceBus(logger),
	}
}

func (a *azureServiceBus) PublishDelayMessage(ctx context.Context, request *delay_queue.DelayMessageRequest) (*delay_queue.DelayMessageResponse, error) {
	// convert ScheduledEnqueueTimeUtc

	req := &pubsub.PublishRequest{
		Data:       request.Data,
		PubsubName: request.ComponentName,
		Topic:      request.Topic,
		Metadata:   request.Metadata,
	}
	err := a.Publish(req)
	return nil, err
}
