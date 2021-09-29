# Layotto (L8):To be the next layer of OSI layer 7

[![codecov](https://codecov.io/gh/mosn/layotto/branch/main/graph/badge.svg?token=10RxwSV6Sz)](https://codecov.io/gh/mosn/layotto)
[![Average time to resolve an issue](http://isitmaintained.com/badge/resolution/mosn/layotto.svg)](http://isitmaintained.com/project/mosn/layotto "Average time to resolve an issue")

<img src="https://raw.githubusercontent.com/mosn/layotto/main/docs/img/logo/grey2-1.svg" height="120px">

Layotto 是一款使用 Golang 开发的应用运行时, 旨在帮助开发人员快速构建云原生应用，帮助应用和基础设施解耦。它为应用提供了各种分布式能力，比如状态管理，配置管理，事件发布订阅等能力，以简化应用的开发。

Layotto 以开源的 [MOSN](https://github.com/mosn/mosn) 为底座，在提供分布式能力以外，提供了 Service Mesh 对于流量的管控能力。

## 诞生背景

Layotto希望可以把 [Multi-Runtime](https://www.infoq.com/articles/multi-runtime-microservice-architecture/) 跟Service Mesh两者的能力结合起来，无论你是使用MOSN还是Envoy或者其他产品作为Service Mesh的数据面，都可以在不增加新的sidecar的前提下，使用Layotto为这些数据面追加Runtime的能力。

例如，通过为MOSN添加Runtime能力，一个Layotto进程可以[既作为istio的数据面](zh/start/istio/start.md) 又提供各种Runtime API（例如Configuration API,Pub/Sub API等）

此外，随着探索实践，我们发现sidecar能做的事情远不止于此。 通过引入[WebAssembly](https://en.wikipedia.org/wiki/WebAssembly) ,我们正在尝试将Layotto做成 FaaS (Function as a service) 和 [reloadable sdk](https://github.com/mosn/layotto/issues/166) 的运行时容器 。

如果您对诞生背景感兴趣，可以看下[这篇演讲](https://mosn.io/layotto/#/zh/blog/mosn-subproject-layotto-opening-a-new-chapter-in-service-grid-application-runtime/index) 。

## 功能

- 服务通信 
- 服务治理，例如流量的劫持和观测，服务限流等
- [作为 istio 的数据面](zh/start/istio/start.md)  
- 配置管理
- 状态管理
- 事件发布订阅
- 健康检查、查询运行时元数据
- 基于WASM的多语言编程

## 工程架构

如下图架构图所示，Layotto 以开源 MOSN 作为底座，在提供了网络层管理能力的同时提供了分布式能力，业务可以通过轻量级的 SDK 直接与 Layotto 进行交互，而无需关注后端的具体的基础设施。

Layotto 提供了多种语言版本的 SDK，SDK 通过 gRPC 与 Layotto 进行交互，应用开发者只需要通过 Layotto 提供的配置文件[配置文件](https://github.com/mosn/layotto/blob/main/configs/runtime_config.json)
来指定自己基础设施类型，而不需要进行任何编码的更改，大大提高了程序的可移植性。

![系统架构图](https://raw.githubusercontent.com/mosn/layotto/main/docs/img/runtime-architecture.png)

## 快速开始

### 使用配置中心API

[通过 Layotto 调用 apollo 配置中心](zh/start/configuration/start-apollo.md)

### 使用Pub/Sub API实现发布/订阅模式

[通过Layotto调用redis，进行消息发布/订阅](zh/start/pubsub/start.md)

### 使用State API进行状态数据的增删改查

[通过Layotto调用redis，进行状态管理](zh/start/state/start.md)

### 使用分布式锁API

[基于redis使用分布式锁](zh/start/lock/start.md)

### 在四层网络进行流量干预

[Dump TCP 流量](zh/start/network_filter/tcpcopy.md)

### 在七层网络进行流量干预

[方法级别限流](zh/start/stream_filter/flow_control.md)

### 进行RPC调用

[Hello World](zh/start/rpc/helloworld.md)

[Dubbo JSON RPC](zh/start/rpc/dubbo_json_rpc.md)

### 跟istio集成

[作为 istio 的数据面](zh/start/istio/start.md)

### 健康检查、查询运行时元数据

[使用 Layotto Actuator 进行健康检查和元数据查询](zh/start/actuator/start.md)

### 使用WASM进行多语言编程

[WASM概述](zh/start/wasm/start.md)

## 设计文档

[Actuator设计文档](zh/design/actuator/actuator-design-doc.md)

[pubsub api以及与dapr component的兼容性](zh/design/pubsub/pubsub-api-and-compability-with-dapr-component.md)

[configuration-api-with-apollo(英文)](en/design/configuration/configuration-api-with-apollo.md)

[rpc设计文档](zh/design/rpc/rpc设计文档.md)

[分布式锁api设计文档](zh/design/lock/lock-api-design.md)


## 社区

| 平台  | 联系方式        |
|:----------|:------------|
| 💬 [钉钉](https://www.dingtalk.com/zh) (推荐) | 群号: 31912621 或者扫描下方二维码 <br> <img src="https://raw.githubusercontent.com/mosn/layotto/main/docs/img/ding-talk-group-1.png" height="200px">

[comment]: <> (| 💬 [微信]&#40;https://www.wechat.com/&#41;  | 扫描下方二维码添加好友，她会邀请您加入微信群 <br> <img src="../img/wechat-group.jpg" height="200px">)

## 如何贡献代码

[组件开发指南](zh/development/developing-component.md)

[Layotto贡献者指南](zh/development/CONTRIBUTING.md)

## FAQ

### 跟dapr有什么差异？

dapr是一款优秀的Runtime产品，但它本身缺失了Service Mesh的能力，而这部分能力对于实际在生产环境落地是至关重要的，因此我们希望把Runtime
跟Service Mesh两种能力结合在一起，满足更复杂的生产落地需求。
