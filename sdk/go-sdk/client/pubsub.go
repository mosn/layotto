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

package client

import (
	"context"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
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
