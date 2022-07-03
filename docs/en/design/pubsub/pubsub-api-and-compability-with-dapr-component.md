# Pub/Sub API and compatibility with Dapr's packages
# 1. Requirements
1. Support Pub/Sub API
2. The architecture can reuse Dapr's packages as much as possible

# 2. High-level design
## 2.1. Architecture: whether to reuse Dapr's sdk and proto
In order to develop a set of universally accepted API specs with the Dapr and Envoy communities in the future, we try to be consistent with the Dapr API at present.

Dapr's component library can be reused directly; the following discusses whether sdk and proto are reused and how to reuse.

### The problems we are facing

1. The SDK of dapr is hard-coded to call the package name of the API, and there is 'dapr' word in the name
   ![img.png](../../../img/mq/design/img.png)
2. We will have differentiated requirements in the future, such as new fields and new APIs. If we use dapr.proto directly, it will be inflexible
### Solution

#### Do not reuse sdk and proto; move the proto file to a neutral path
![img_1.png](../../../img/mq/design/img_1.png)

We first define an api-spec.proto, this proto is a superset of dapr API, and the path name is neutral without the word 'layotto' or 'dapr'.Based on this proto, we can develop a neutral RuntimeAPI sdk.

Later, try to promote the proto into an api-spec accepted by the runtime community, or rebuild a path-neutral api-spec.proto with other communities.

It does not matter if the proto changes during the promotion process. Layotto internally extracts an API layer under the proto to prevent proto changes;

If it is not easy to push, we can write a dapr adapter in the neutral SDK in the short term, and use our SDK to adjust dapr and layotto:
![img_2.png](../../../img/mq/design/img_2.png)

Advantages:

1. Neat and tidy. If you want to reuse Dapr's sdk and proto, there is an inevitable problem: when the API and dapr are different, you need to encapsulate a layer of your own logic, which will bring complexity, hacky, a sense of copycat, and raise the difficulty of code reading.
1. Extendible when the APIs are different from Dapr's

Disadvantages:

1. Subsequent Dapr client or proto changes, we may not know, resulting in inconsistencies


## 2.2. API Design
### 2.2.1. Design principle: How to add fields to Dapr's API
We want to reuse Dapr API, but in the long run, there will definitely be customization requirements. When our API and dapr's are different (for example, we just want to add a new field to a certain API of Dapr), should we create a new method name or just add a field to the original method?

If we add a field to the original method, it may cause field conflicts.

After serveral discussion,we finally decide to add fields directly in that situation.Conflicts of API are inevitable (of course,we will try to raise pull requests to the Dapr community to avoid conflicts)

In the future, when everyone really sits together to reach a consensus and build api-spec, a new proto with a new path will be created. Anyway, there will be a new proto at that time, so don't worry about the current conflict.

### 2.2.2. Between APP and Layotto
Use the same grpc API as Dapr

```protobuf
service AppCallback {
  // Lists all topics subscribed by this app.
  rpc ListTopicSubscriptions(google.protobuf.Empty) returns (ListTopicSubscriptionsResponse) {}

  // Subscribes events from Pubsub
  rpc OnTopicEvent(TopicEventRequest) returns (TopicEventResponse) {}

}
```

```protobuf
service Dapr {
  // Publishes events to the specific topic.
  rpc PublishEvent(PublishEventRequest) returns (google.protobuf.Empty) {}
}

```

### 2.2.3. Between Layotto and Component
Use the same as Dapr;
PublishRequest.Data and NewMessage.Data put json data conforming to CloudEvent 1.0 specification (can be deserialized and put into map[string]interface{})

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

### 2.2.4. How does the sidecar know which port to call back

Configure the callback port at startup. The price is that the sidecar can only serve one process.

Temporarily choose this plan in this issue

### 2.2.5. How to keep the subscription list real-time

The app is called when the sidecar starts, and the subscription relationship is obtained at that time. Therefore, there are requirements for the startup sequence. Start the app first.

It can be optimized into a timed polling mechanism in the future

### 2.2.6. Does the subscription relationship support declarative configuration?

In the first phase, only the way of opening an API for callback is supported, and the subsequent optimization will be declarative configuration or dynamic configuration.

## 2.3. Config Design
![img.png](../../../img/mq/design/config.png)

# 3. Future Work
## A Bigger Control Plane

The Control Plane in the Service Mesh era only serves RPC, but in the Runtime API era, component configuration also needs to be distributed by the cluster; components also need service discovery and routing, so components also need their own Control Plane.

It is convenient to have a Bigger Control Plane that integrates RPC and all middleware configuration data

Maybe we have to extend the xDS protocol, like 'runtime Discovery Service'.

## Subscription relationship support configuration

The subscription relationship is now obtained by the callback mechanism.We want the subscription relationship be obtained through configuration.

## appcallback support TLS


## Separate component configuration and personality configuration (callback port, app-id)
The current component configuration and app personality configuration (callback port, app-id) are put together, and there are some problems:

1. It's not easy to distribute the configuration to the whole cluster
1. Can't do component access control (for example, Dapr can restrict app-id1 to only access topic_id1)
![img_4.png](../../../img/mq/design/img_4.png)

Need to refactor the original component logic.

## Tracing
