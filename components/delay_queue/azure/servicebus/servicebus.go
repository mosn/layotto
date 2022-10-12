// Copyright 2021 Layotto Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package servicebus

import (
	"context"
	"net/http"
	"time"

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
	nowUtc := time.Now().UTC()
	enqueueTime := nowUtc.Add(time.Second * time.Duration(request.DelayInSeconds))
	request.Metadata["metadata.ScheduledEnqueueTimeUtc"] = enqueueTime.Format(http.TimeFormat)

	req := &pubsub.PublishRequest{
		Data:       request.Data,
		PubsubName: request.ComponentName,
		Topic:      request.Topic,
		Metadata:   request.Metadata,
	}
	err := a.Publish(req)
	return nil, err
}
