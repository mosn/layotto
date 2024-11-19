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

package default_api

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/dapr/components-contrib/pubsub"
	"github.com/golang/mock/gomock"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"

	mock "mosn.io/layotto/pkg/mock/components/pubsub"
	"mosn.io/layotto/pkg/mock/runtime"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
)

func TestSubscribeTopicEvents(t *testing.T) {
	t.Run("SubscribeTopicEvents duplicate initial request", func(t *testing.T) {
		pubsubName := "test"
		topic := "layotto"
		ctrl := gomock.NewController(t)
		stream := runtime.NewMockRuntime_SubscribeTopicEventsServer(ctrl)
		stream.EXPECT().Send(gomock.Any()).Return(nil).Times(1)
		stream.EXPECT().Context().Return(context.Background()).AnyTimes()
		stream.EXPECT().Recv().Return(&runtimev1pb.SubscribeTopicEventsRequest{
			SubscribeTopicEventsRequestType: &runtimev1pb.SubscribeTopicEventsRequest_InitialRequest{
				InitialRequest: &runtimev1pb.SubscribeTopicEventsRequestInitial{
					PubsubName: pubsubName, Topic: topic,
					Metadata: make(map[string]string),
				},
			}}, nil).AnyTimes()

		a := NewAPI("", nil, nil, nil, make(map[string]pubsub.PubSub), nil, nil, nil, nil, nil, nil)

		var apiForTest = a.(*api)

		m := mock.NewMockPubSub(ctrl)
		m.EXPECT().Subscribe(gomock.Any(), gomock.Any()).Return(nil).Times(1)
		apiForTest.pubSubs["test"] = m

		apiForTest.streamer = &streamer{
			subscribers: make(map[string]*conn),
		}

		apiForTest.json = jsoniter.ConfigFastest

		err := apiForTest.SubscribeTopicEvents(stream)
		assert.Error(t, err, "Expected error")
	})
}

func TestPublishMessageForStream(t *testing.T) {
	t.Run("publish success", func(t *testing.T) {
		pubsubName := "test"
		topic := "layotto"
		ctrl := gomock.NewController(t)
		stream := runtime.NewMockRuntime_SubscribeTopicEventsServer(ctrl)
		stream.EXPECT().Send(gomock.Any()).Return(nil).Times(1)
		stream.EXPECT().Recv().Return(&runtimev1pb.SubscribeTopicEventsRequest{
			SubscribeTopicEventsRequestType: &runtimev1pb.SubscribeTopicEventsRequest_InitialRequest{
				InitialRequest: &runtimev1pb.SubscribeTopicEventsRequestInitial{
					PubsubName: pubsubName, Topic: topic,
					Metadata: make(map[string]string),
				},
			}}, nil).Times(1)

		cloudEvent := map[string]interface{}{
			pubsub.IDField:              "1",
			pubsub.SourceField:          "source",
			pubsub.DataContentTypeField: "content-type",
			pubsub.TypeField:            "type",
			pubsub.SpecVersionField:     "v1.0.0",
			pubsub.DataBase64Field:      "bGF5b3R0bw==",
		}

		data, err := json.Marshal(cloudEvent)
		assert.Nil(t, err)

		msg := &pubsub.NewMessage{
			Data:     data,
			Topic:    topic,
			Metadata: make(map[string]string),
		}
		a := NewAPI("", nil, nil, nil, make(map[string]pubsub.PubSub), nil, nil, nil, nil, nil, nil)

		var apiForTest = a.(*api)

		apiForTest.streamer = &streamer{
			subscribers: make(map[string]*conn),
		}

		apiForTest.streamer.subscribers["___test||layotto"] = &conn{
			stream:           stream,
			publishResponses: make(map[string]chan *runtimev1pb.SubscribeTopicEventsRequestProcessed),
		}

		_ = apiForTest.SubscribeTopicEvents(stream)
		apiForTest.json = jsoniter.ConfigFastest

		go func() {
			time.Sleep(1 * time.Second)
			ch := apiForTest.streamer.subscribers["___test||layotto"].publishResponses["1"]
			ch <- &runtimev1pb.SubscribeTopicEventsRequestProcessed{
				Id: "1",
				Status: &runtimev1pb.TopicEventResponse{
					Status: runtimev1pb.TopicEventResponse_SUCCESS,
				},
			}
		}()

		err = apiForTest.publishMessageForStream(context.Background(), msg, pubsubName)
		assert.Nil(t, err)
	})
}
