# Layotto Pub/Sub, Compatible Dapr Package Scheme

# Needs analysis

1. Support Pub/Sub API
2. Packages that can be reused as much as possible on the architecture

# Overview design

## Whether the corporate architecture：reuses Dapr's sdk and proto

In order to be able to develop a set of API specs with the Dapr and Envoy communities in the future, the Dapr API is being kept as close as possible.

Dapr component libraries can be reused directly; below discuss whether sdk and proto are repeated and how to use them again.

### Problems encountered

1. dapr sdk wrote the pack name in which the call interface was called, with dapr
   ![img.png](/img/mq/design/img.png)
2. We will have differentiated needs, such as new fields, new APIs, if directly using dapr.proto will not be flexible

### Programmes

#### No longer using sdk and proto; detach proto file, neutral path

![img_1.png](/img/mq/design/img_1.png)

We first define an api-spec.proto, a superset of dapr API, with a neutral path name without layotto, based on this proto itself develop a neutral RuntimeAPI sdk.

Then try promoting proto into a community-sanctioned api-spec, or working with other communities to rebuild a path neutral api-spec.proto.

If proto changes in the promotion process does not matter anything, Layotto draws a layer of API below proto to prevent proto;

If not, we can write a dapr adapter first in a neutral sdk, using our sdk to both adjust the dapr and layotto：

![img_2.png](/img/mq/design/img_2.png)

Advantages：

1. Clean.If you want to revert to Dapr's sdk and proto, there is an unavoidable problem：when the API and dapr are different, you need to encapsulate a layer of logic for yourself, which will bring complexity, hacky, hill sense and increase the code reading threshold
2. Extensions between API and Dapr

： Disadvantages

1. Follow Dapr client or proto modified. We may not know, causing inconsistency

## API Design

### Design Principle：To deal with new fields for Dapr API

We want to revert to the Dapr API, but there is certainly a need for customization in the long term.When our API and dapr are different (e.g. just want to give a new API field to Dapr), whether to open a new method name, or to add a field to the old method?
to add a field to the original method, this may cause a field conflict.

Some of the following ideas：

#### New method name as long as the API changes different from Dapr

![img_3.png](/img/mq/design/img_3.png)

New and old methods are supported when new methods are opened.**For example, version v1 is Dapr API and version v2 is extended**

： Disadvantages

1. To support two APIs

#### New fields also use old method names, but jumps in numbers, leave white

： Disadvantages

1. Make no more sense?If Dapr later added this field but the numbers are different, we would find it difficult to do so (e.g. we define it as 10, dapr then as 5, and we have one field that takes up both 5 and 10?)
2. If Dapr adds a similar but nuanced field, we find it difficult to add：to the field?

#### C. Fields are directly added to the conflict without allowing for a conflict (will, of course, try to raise issues for the Dapr community)

In the future, when you really sit together to reach a consensus and make api-spec, you will start a new path to proto, and you will not worry about the current conflict.

#### Conclusion

Discussion decides to follow C line

### Between APP and Layotto

Use the grpc API like Dapr

```protobuf
service AppCallback {
  // Lists all topics subscribed by this app.
  rpc ListTopicSubscriptions(google.protobuf.Empty) returns (ListTopicSubscriptionsResponse) {}

  // Subscribes events from Pubsub
  rpc OnTopicEvent(TopicEventRequest) returns (TopicEventResponse) {}

}
```

```protobuf
Service Dapr LO
  // Publishes events to the specific topic.
  rpc PublishEvent(PublishEventRequest) returns (google.protect.Empt) {}


```

### Between Layotto and Component

In the same way as Dapr;
PublishRequest.Data-and NewMessage.Data-fit for CloudEvent 1.0 (deserializable into map[string]interface{}

```go
// PubSub is the interface for message buses
type PubSub interface {
	Init(metadata Metadata) error
	Features() []Feature
	Publish(req *PublishRequest) error
	Subscribe(req SubscribeRequest, handler func(msg *NewMessage) error) error
	Close() error
}

// PublishRequest is the request to publish a message
type PublishRequest struct {
	Data       []byte            `json:"data"`
	PubsubName string            `json:"pubsubname"`
	Topic      string            `json:"topic"`
	Metadata   map[string]string `json:"metadata"`
}


// NewMessage is an event arriving from a message bus instance
type NewMessage struct {
	Data     []byte            `json:"data"`
	Topic    string            `json:"topic"`
	Metadata map[string]string `json:"metadata"`
}

```

### sidecar knows which port of the callback

Reference Dapr, configure callback port on startup.The cost is that sidecar can only serve one process.

Select this option for the current period

### How to maintain the antecedents of the subscription list

Ask an app on sidecar startup to get subscription at once.There is therefore a requirement for the order of start and start the app first.

Follow up the app to optimize it into scheduled poll app

### Subscription support declaration configuration

First issue only supports the format of the interface callback and then optimized with the declaration configuration

## Config Design

![img.png](/img/mq/design/config.png)

The relevant configuration of the app is in the loaded, and the code you want to reconfigure the configuration API, etc. (see below).

**Q: How to pass configuration data to Dapr and Layotto components**

A: pass metadata's data through the Init interface to components

# Future Work

## A Bigger Control Plane

The Control Plane of Service Mesh serves only RPC, but in Runtime API the configuration of components also needs to be distributed in clusters; components also need to be discovered, routed, and so also have their own control Plane.

It would be better to have a Bigger Control Plane, which integrates RPC's and all middleware configuration data.

may require extension of the ODS protocol, such as runtime Discovery Service

## Subscription Support Configuration

Subscription is now obtained by callback app, which can be added to get subscription by configuration

## appcallback supports tls

## Detach component configuration and personality configuration (callback port, app-id)

Component configuration and app profile (callback port,app-id) are placed together with the following question：

1. Not good at making cluster configuration leader
2. Component control configuration is not available (e.g. Dapr can limit app-id1 access only to topic_id1)

![img_4.png](/img/mq/design/img_4.png)

Need to reconstruct the original component logic

## Tracing
