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
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/dapr/components-contrib/pubsub"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"mosn.io/pkg/log"

	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
)

type streamer struct {
	subscribers map[string]*conn
	lock        sync.RWMutex
}

type conn struct {
	lock             sync.RWMutex
	streamLock       sync.Mutex
	stream           runtimev1pb.Runtime_SubscribeTopicEventsServer
	publishResponses map[string]chan *runtimev1pb.SubscribeTopicEventsRequestProcessed
}

// SubscribeTopicEvents is called by the layotto runtime to ad hoc stream
// subscribe to topics. If gRPC API server closes, returns func early with nil
// to close stream.
func (a *api) SubscribeTopicEvents(stream runtimev1pb.Runtime_SubscribeTopicEventsServer) error {
	errCh := make(chan error, 2)
	subDone := make(chan struct{})

	go func() {
		<-subDone
	}()

	go func() {
		errCh <- a.streamSubscribe(stream, subDone)
	}()

	return <-errCh
}

func (a *api) streamSubscribe(stream runtimev1pb.Runtime_SubscribeTopicEventsServer, subDone chan struct{}) error {
	defer close(subDone)

	subscribeTopicEventsRequest, err := stream.Recv()
	if err != nil {
		return err
	}

	initialRequest := subscribeTopicEventsRequest.GetInitialRequest()

	if initialRequest == nil {
		return errors.New("initial request is required")
	}

	if len(initialRequest.GetPubsubName()) == 0 {
		return errors.New("pubsubName is required")
	}

	if len(initialRequest.GetTopic()) == 0 {
		return errors.New("topic is required")
	}

	if a.topicPerComponent == nil {
		a.topicPerComponent = make(map[string]TopicSubscriptions)
	}

	if a.streamer == nil {
		a.streamer = &streamer{
			subscribers: make(map[string]*conn),
		}
	}

	if _, ok := a.topicPerComponent[initialRequest.PubsubName]; !ok {
		a.topicPerComponent[initialRequest.PubsubName] = TopicSubscriptions{topic2Details: make(map[string]Details)}
	}

	a.topicPerComponent[initialRequest.PubsubName].topic2Details[initialRequest.Topic] = Details{
		metadata: initialRequest.Metadata,
	}

	if len(a.topicPerComponent) > 0 {
		for pubsubName, v := range a.topicPerComponent {
			topics := []string{}
			for topic := range v.topic2Details {
				topics = append(topics, topic)
			}
			log.DefaultLogger.Infof("[runtime][streamSubscribe]app is subscribed to the following topics: %v through pubsub=%s", topics, pubsubName)
		}
	}

	if a.pubSubs[initialRequest.PubsubName] == nil {
		return errors.New("pubsub " + initialRequest.PubsubName + " is not initialized.")
	}

	if err = stream.Send(&runtimev1pb.SubscribeTopicEventsResponse{
		SubscribeTopicEventsResponseType: &runtimev1pb.SubscribeTopicEventsResponse_InitialResponse{
			InitialResponse: new(runtimev1pb.SubscribeTopicEventsResponseInitial),
		},
	}); err != nil {
		return err
	}

	if err = a.pubSubs[initialRequest.PubsubName].Subscribe(pubsub.SubscribeRequest{
		Topic:    initialRequest.Topic,
		Metadata: a.topicPerComponent[initialRequest.PubsubName].topic2Details[initialRequest.Topic].metadata,
	}, func(ctx context.Context, msg *pubsub.NewMessage) error {
		if msg.Metadata == nil {
			msg.Metadata = make(map[string]string, 1)
		}
		msg.Metadata[Metadata_key_pubsubName] = initialRequest.PubsubName
		return a.publishMessageForStream(ctx, msg, initialRequest.PubsubName)
	}); err != nil {
		log.DefaultLogger.Warnf("[runtime][beginPubSub]failed to subscribe to topic %s: %s", initialRequest.Topic, err)
		return err
	}

	return a.streamer.Subscribe(stream, initialRequest)
}

func (s *streamer) Subscribe(stream runtimev1pb.Runtime_SubscribeTopicEventsServer, req *runtimev1pb.SubscribeTopicEventsRequestInitial) error {
	s.lock.Lock()
	key := s.StreamerKey(req.GetPubsubName(), req.GetTopic())
	if _, ok := s.subscribers[key]; ok {
		s.lock.Unlock()
		return fmt.Errorf("already subscribed to pubsub %q topic %q", req.GetPubsubName(), req.GetTopic())
	}

	conn := &conn{
		stream:           stream,
		publishResponses: make(map[string]chan *runtimev1pb.SubscribeTopicEventsRequestProcessed),
	}
	s.subscribers[key] = conn

	log.DefaultLogger.Infof("Subscribing to pubsub '%s' topic '%s'", req.GetPubsubName(), req.GetTopic())
	s.lock.Unlock()

	defer func() {
		s.lock.Lock()
		delete(s.subscribers, key)
		s.lock.Unlock()
	}()

	for {
		resp, err := stream.Recv()

		s, ok := status.FromError(err)

		if (ok && s.Code() == codes.Canceled) ||
			errors.Is(err, context.Canceled) ||
			errors.Is(err, io.EOF) {
			log.DefaultLogger.Infof("Unsubscribed from pubsub '%s' topic '%s'", req.GetPubsubName(), req.GetTopic())
			return err
		}

		if err != nil {
			log.DefaultLogger.Errorf("error receiving message from client stream: %s", err)
			return err
		}

		eventResp := resp.GetEventProcessed()
		if eventResp == nil {
			return errors.New("duplicate initial request received")
		}
		go func() {
			conn.notifyPublishResponse(stream.Context(), eventResp)
		}()
	}
}

func (a *api) publishMessageForStream(ctx context.Context, msg *pubsub.NewMessage, pubsubName string) error {
	a.streamer.lock.RLock()
	key := a.streamer.StreamerKey(pubsubName, msg.Topic)
	conn, ok := a.streamer.subscribers[key]
	a.streamer.lock.RUnlock()
	if !ok {
		return fmt.Errorf("no streamer subscribed to pubsub %q topic %q", pubsubName, msg.Topic)
	}

	envelope, cloudEvent, _ := a.envelopeFromSubscriptionMessage(ctx, msg)

	ch, defFn := conn.registerPublishResponse(envelope.GetId())
	defer defFn()

	conn.streamLock.Lock()
	err := conn.stream.Send(&runtimev1pb.SubscribeTopicEventsResponse{
		SubscribeTopicEventsResponseType: &runtimev1pb.SubscribeTopicEventsResponse_EventMessage{
			EventMessage: envelope,
		},
	})
	if err != nil {
		log.DefaultLogger.Errorf("error sending message to client stream: %s", err)
		return err
	}
	conn.streamLock.Unlock()

	var resp *runtimev1pb.SubscribeTopicEventsRequestProcessed
	select {
	case <-ctx.Done():
		return ctx.Err()
	case resp = <-ch:
	}

	// 5. Check result
	return retryStrategy(err, resp.Status, cloudEvent)
}

func (c *conn) notifyPublishResponse(ctx context.Context, resp *runtimev1pb.SubscribeTopicEventsRequestProcessed) {
	c.lock.RLock()
	ch, ok := c.publishResponses[resp.GetId()]
	c.lock.RUnlock()

	if !ok {
		log.DefaultLogger.Errorf("no client stream expecting publish response for id %q", resp.GetId())
		return
	}

	select {
	case <-ctx.Done():
	case ch <- resp:
	}
}

func (c *conn) registerPublishResponse(id string) (chan *runtimev1pb.SubscribeTopicEventsRequestProcessed, func()) {
	ch := make(chan *runtimev1pb.SubscribeTopicEventsRequestProcessed)
	c.lock.Lock()
	c.publishResponses[id] = ch
	c.lock.Unlock()
	return ch, func() {
		c.lock.Lock()
		delete(c.publishResponses, id)
		c.lock.Unlock()
	}
}

func (s *streamer) StreamerKey(pubsub, topic string) string {
	return "___" + pubsub + "||" + topic
}
