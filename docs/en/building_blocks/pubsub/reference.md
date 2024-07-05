# Pub/Sub API
## What is Pub/Sub API
Pub/Sub API is used to implement the publish/subscribe model. The publish/subscribe model allows microservices to communicate with each other using messages. **The producer or publisher** sends the message to the specified topic and does not know the application receiving the message. The **consumer** will subscribe to the topic and receive its messages, and does not know what application produced these messages. The message broker acts as an intermediary and is responsible for forwarding each message. This mode is especially useful when you need to decouple microservices.

Pub/Sub API provides at-least-once guarantee and integrates with various message brokers and queue systems. Your application can use the same set of Pub/Sub API to operate different message queues.
## When to use Pub/Sub API and what are the benefits?
If your application needs to access the message queue for event publishing and subscription, then using Pub/Sub API is a good choice. It has the following benefits:

- Multi (cloud) environment deployment: the same application code can be deployed in different environments

A neutral API can help your application decouple from storage vendors and cloud vendors, and be able to deploy on different clouds without changing the code.

- Multi-language reuse middleware: the same set of message middleware can support applications in different languages

If your company has applications developed in different languages (for example, there are both java and python applications), then the traditional approach is to develop a set of SDKs for each language.

Using Pub/Sub API can help you avoid the trouble of maintaining multilingual SDKs. Applications in different languages can use the same set of grpc API to interact with Layotto.

## How to use Pub/Sub API
You can call Pub/Sub API through grpc. The API is defined in [runtime.proto](https://github.com/mosn/layotto/blob/main/spec/proto/runtime/v1/runtime.proto).

The component needs to be configured before use. For detailed configuration instructions, see [publish/subscribe component documentation](zh/component_specs/pubsub/common.md)

### Example
Layotto client sdk encapsulates the logic of grpc call. For examples of using sdk to call Pub/Sub API, please refer to [Quick Start: Use Pub/Sub API](en/start/pubsub/start.md)

### PublishEvent
Used to publish events to the specified topic

```protobuf
// Publishes events to the specific topic.
rpc PublishEvent(PublishEventRequest) returns (google.protobuf.Empty) {}
```

To avoid inconsistencies between the documentation and the code, please refer to [runtime.proto](https://github.com/mosn/layotto/blob/main/spec/proto/runtime/v1/runtime.proto) for detailed input parameters and return values

### Subscribe to events
To subscribe to events, the application needs to implement two grpc APIs for Layotto to call back:


```protobuf
  // Lists all topics subscribed by this app.
  rpc ListTopicSubscriptions(google.protobuf.Empty) returns (ListTopicSubscriptionsResponse) {}

  // Subscribes events from Pubsub
  rpc OnTopicEvent(TopicEventRequest) returns (TopicEventResponse) {}

```

To avoid inconsistencies between the documentation and the code, please refer to [appcallback.proto](https://github.com/mosn/layotto/blob/main/spec/proto/runtime/v1/appcallback.proto) for detailed input parameters and return values