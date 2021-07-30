// CODE ATTRIBUTION: https://github.com/dapr/go-sdk
// Modified the import package to use layotto's pb
// We use same sdk code with Dapr's for state API because we want to keep compatible with Dapr state API
package client

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	pb "mosn.io/layotto/spec/proto/runtime/v1"
)

// PublishEvent publishes data onto specific pubsub topic.
func (c *GRPCClient) PublishEvent(ctx context.Context, pubsubName, topicName string, data []byte) error {
	if pubsubName == "" {
		return errors.New("pubsubName name required")
	}
	if topicName == "" {
		return errors.New("topic name required")
	}

	envelop := &pb.PublishEventRequest{
		PubsubName: pubsubName,
		Topic:      topicName,
		Data:       data,
	}

	_, err := c.protoClient.PublishEvent(ctx, envelop)
	if err != nil {
		return errors.Wrapf(err, "error publishing event unto %s topic", topicName)
	}

	return nil
}

// PublishEventfromCustomContent serializes an struct and publishes its contents as data (JSON) onto topic in specific pubsub component.
func (c *GRPCClient) PublishEventfromCustomContent(ctx context.Context, pubsubName, topicName string, data interface{}) error {
	if pubsubName == "" {
		return errors.New("pubsubName name required")
	}
	if topicName == "" {
		return errors.New("topic name required")
	}

	bytes, err := json.Marshal(data)

	if err != nil {
		return errors.WithMessage(err, "error serializing input struct")
	}

	envelop := &pb.PublishEventRequest{
		PubsubName:      pubsubName,
		Topic:           topicName,
		Data:            bytes,
		DataContentType: "application/json",
	}

	_, err = c.protoClient.PublishEvent(ctx, envelop)

	if err != nil {
		return errors.Wrapf(err, "error publishing event unto %s topic", topicName)
	}

	return nil
}
