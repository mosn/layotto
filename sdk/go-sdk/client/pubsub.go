package client

import (
	"context"
	runtimev1pb "github.com/layotto/layotto/spec/proto/runtime/v1"
)

type PublishEventRequest struct {
	// The name of the pubsub component
	PubsubName string `json:"pubsub_name,omitempty"`
	// The pubsub topic
	Topic string `json:"topic,omitempty"`
	// The data which will be published to topic.
	Data []byte `json:"data,omitempty"`
	// The content type for the data (optional).
	DataContentType string `json:"data_content_type,omitempty"`
	// The metadata passing to pub components
	//
	// metadata property:
	// - key : the key of the message.
	Metadata map[string]string `json:"metadata,omitempty"`
}

func (c *GRPCClient) PublishEvent(ctx context.Context, in *PublishEventRequest) error {
	req := &runtimev1pb.PublishEventRequest{
		PubsubName:      in.PubsubName,
		Topic:           in.Topic,
		Data:            in.Data,
		DataContentType: in.DataContentType,
		Metadata:        in.Metadata,
	}
	_, err := c.protoClient.PublishEvent(ctx, req)
	return err
}
