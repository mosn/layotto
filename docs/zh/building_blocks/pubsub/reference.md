# 发布/订阅 API
## 什么是Pub/Sub API
Pub/Sub API用于实现发布/订阅模式。发布/订阅模式允许微服务使用消息相互通信。 **生产者或发布者**将消息发送至指定Topic，并且不知道接收消息的应用程序。而**消费者**将订阅该主题并收到它的消息，并且不知道什么应用程序生产了这些消息。消息队列（message broker）作为中间人，负责将每条消息的转发。 当您需要将微服务解偶时，此模式特别有用。

Pub/Sub API 提供至少一次（at-least-once）的保证，并与各种消息代理和队列系统集成。 您的应用程序可以使用同一套Pub/Sub API操作不同的消息队列。

## 何时使用Pub/Sub API，好处是什么？
如果您的应用需要访问消息队列(message queue)进行事件发布订阅，那么使用Pub/Sub API是一个不错的选择，它有以下好处：

- 多（云）环境部署：同一套业务代码部署在不同环境

中立的API可以帮助您的应用和MQ供应商、云厂商解耦，能够不改代码部署在不同的云上。

- 多语言复用中间件：同一套消息中间件能支持不同语言的应用

如果您的公司内部有不同语言开发的应用（例如同时有java和python应用），那么传统做法是为每种语言开发一套sdk。

使用Pub/Sub API可以帮助您免去维护多语言sdk的烦恼，不同语言的应用可以用同一套grpc API和Layotto交互。

## 如何使用Pub/Sub API
您可以通过grpc调用Pub/Sub API，接口定义在[runtime.proto](https://github.com/mosn/layotto/blob/main/spec/proto/runtime/v1/runtime.proto) 中。

使用前需要先对组件进行配置，详细的配置说明见[发布/订阅组件文档](zh/component_specs/pubsub/common.md)

### 使用示例
Layotto client sdk封装了grpc调用的逻辑，使用sdk调用Pub/Sub API的示例可以参考[快速开始：使用Pub/Sub API](zh/start/pubsub/start.md)

### PublishEvent
用于发布事件到指定topic

```protobuf
// Publishes events to the specific topic.
rpc PublishEvent(PublishEventRequest) returns (google.protobuf.Empty) {}
```

为避免文档和代码不一致，详细入参和返回值请参考[runtime.proto](https://github.com/mosn/layotto/blob/main/spec/proto/runtime/v1/runtime.proto)

### 订阅事件
订阅事件需要应用实现两个grpc接口，供Layotto回调：


```protobuf
  // Lists all topics subscribed by this app.
  rpc ListTopicSubscriptions(google.protobuf.Empty) returns (ListTopicSubscriptionsResponse) {}

  // Subscribes events from Pubsub
  rpc OnTopicEvent(TopicEventRequest) returns (TopicEventResponse) {}

```

为避免文档和代码不一致，详细入参和返回值请参考[appcallback.proto](https://github.com/mosn/layotto/blob/main/spec/proto/runtime/v1/appcallback.proto)