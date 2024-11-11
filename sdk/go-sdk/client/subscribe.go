package client

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
)

const (
	// SubscriptionResponseStatusSuccess means message is processed successfully.
	SubscriptionResponseStatusSuccess SubscriptionResponseStatus = "SUCCESS"
	// SubscriptionResponseStatusRetry means message to be retried by Dapr.
	SubscriptionResponseStatusRetry SubscriptionResponseStatus = "RETRY"
	// SubscriptionResponseStatusDrop means warning is logged and message is dropped.
	SubscriptionResponseStatusDrop SubscriptionResponseStatus = "DROP"
)

type SubscriptionResponseStatus string

type SubscriptionResponse struct {
	Status SubscriptionResponseStatus `json:"status"`
}

type SubscriptionHandleFunction func(request *runtimev1pb.TopicEventRequest) SubscriptionResponseStatus

type SubscriptionRequest struct {
	PubsubName      string
	Topic           string
	DeadLetterTopic *string
	Metadata        map[string]string
}

type Subscription struct {
	stream runtimev1pb.Runtime_SubscribeTopicEventsClient
	// lock locks concurrent writes to subscription stream.
	lock   sync.Mutex
	closed atomic.Bool
}

type SubscriptionMessage struct {
	*runtimev1pb.TopicEventRequest
	sub *Subscription
}

func (c *GRPCClient) Subscribe(ctx context.Context, request SubscriptionRequest) (*Subscription, error) {
	stream, err := c.subscribeInitialRequest(ctx, request)
	if err != nil {
		return nil, err
	}

	s := &Subscription{
		stream: stream,
	}

	return s, nil
}

func (c *GRPCClient) SubscribeWithHandler(ctx context.Context, request SubscriptionRequest, handler SubscriptionHandleFunction) (func() error, error) {
	s, err := c.Subscribe(ctx, request)
	if err != nil {
		return nil, err
	}

	go func() {
		defer s.Close()

		for {
			msg, err := s.Receive()
			if err != nil {
				if !s.closed.Load() {
					logger.Printf("Error receiving messages from subscription pubsub=%s topic=%s, closing subscription: %s",
						request.PubsubName, request.Topic, err)
				}
				return
			}

			go func() {
				if err := msg.respondStatus(handler(msg.TopicEventRequest)); err != nil {
					logger.Printf("Error responding to topic with event status pubsub=%s topic=%s message_id=%s: %s",
						request.PubsubName, request.Topic, msg.Id, err)
				}
			}()
		}
	}()

	return s.Close, nil
}

func (s *Subscription) Close() error {
	if !s.closed.CompareAndSwap(false, true) {
		return errors.New("subscription already closed")
	}

	return s.stream.CloseSend()
}

func (s *Subscription) Receive() (*SubscriptionMessage, error) {
	resp, err := s.stream.Recv()
	if err != nil {
		return nil, err
	}
	event := resp.GetEventMessage()

	eventRequest := &runtimev1pb.TopicEventRequest{
		Id:              event.GetId(),
		Source:          event.GetSource(),
		Type:            event.GetType(),
		SpecVersion:     event.GetSpecVersion(),
		DataContentType: event.GetDataContentType(),
		Data:            event.GetData(),
		Topic:           event.GetTopic(),
		PubsubName:      event.GetPubsubName(),
	}

	return &SubscriptionMessage{
		sub:               s,
		TopicEventRequest: eventRequest,
	}, nil
}

func (s *SubscriptionMessage) Success() error {
	return s.respond(runtimev1pb.TopicEventResponse_SUCCESS)
}

func (s *SubscriptionMessage) Retry() error {
	return s.respond(runtimev1pb.TopicEventResponse_RETRY)
}

func (s *SubscriptionMessage) Drop() error {
	return s.respond(runtimev1pb.TopicEventResponse_DROP)
}

func (s *SubscriptionMessage) respondStatus(status SubscriptionResponseStatus) error {
	var statuspb runtimev1pb.TopicEventResponse_TopicEventResponseStatus
	switch status {
	case SubscriptionResponseStatusSuccess:
		statuspb = runtimev1pb.TopicEventResponse_SUCCESS
	case SubscriptionResponseStatusRetry:
		statuspb = runtimev1pb.TopicEventResponse_RETRY
	case SubscriptionResponseStatusDrop:
		statuspb = runtimev1pb.TopicEventResponse_DROP
	default:
		return fmt.Errorf("unknown status, expected one of %s, %s, %s: %s",
			SubscriptionResponseStatusSuccess, SubscriptionResponseStatusRetry,
			SubscriptionResponseStatusDrop, status)
	}

	return s.respond(statuspb)
}

func (s *SubscriptionMessage) respond(status runtimev1pb.TopicEventResponse_TopicEventResponseStatus) error {
	s.sub.lock.Lock()
	defer s.sub.lock.Unlock()

	return s.sub.stream.Send(&runtimev1pb.SubscribeTopicEventsRequest{
		SubscribeTopicEventsRequestType: &runtimev1pb.SubscribeTopicEventsRequest_EventProcessed{
			EventProcessed: &runtimev1pb.SubscribeTopicEventsRequestProcessed{
				Id:     s.Id,
				Status: &runtimev1pb.TopicEventResponse{Status: status},
			},
		},
	})
}

func (c *GRPCClient) subscribeInitialRequest(ctx context.Context, request SubscriptionRequest) (runtimev1pb.Runtime_SubscribeTopicEventsClient, error) {
	if len(request.PubsubName) == 0 {
		return nil, errors.New("pubsub name required")
	}

	if len(request.Topic) == 0 {
		return nil, errors.New("topic required")
	}

	stream, err := c.protoClient.SubscribeTopicEvents(ctx)
	if err != nil {
		return nil, err
	}

	err = stream.Send(&runtimev1pb.SubscribeTopicEventsRequest{
		SubscribeTopicEventsRequestType: &runtimev1pb.SubscribeTopicEventsRequest_InitialRequest{
			InitialRequest: &runtimev1pb.SubscribeTopicEventsRequestInitial{
				PubsubName: request.PubsubName, Topic: request.Topic,
				Metadata: request.Metadata, DeadLetterTopic: request.DeadLetterTopic,
			},
		},
	})
	if err != nil {
		return nil, errors.Join(err, stream.CloseSend())
	}

	resp, err := stream.Recv()
	if err != nil {
		return nil, errors.Join(err, stream.CloseSend())
	}

	switch resp.GetSubscribeTopicEventsResponseType().(type) {
	case *runtimev1pb.SubscribeTopicEventsResponse_InitialResponse:
	default:
		return nil, fmt.Errorf("unexpected initial response from server : %v", resp)
	}

	return stream, nil
}
