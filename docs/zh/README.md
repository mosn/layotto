<div align="center">
  <h1>Layotto (L8): To be the next layer of OSI layer 7</h1>
  <img src="https://raw.githubusercontent.com/mosn/layotto/main/docs/img/logo/grey2-1.svg" height="120px">

[![Layotto Env Pipeline 🌊](https://github.com/mosn/layotto/actions/workflows/quickstart-checker.yml/badge.svg)](https://github.com/mosn/layotto/actions/workflows/quickstart-checker.yml)
[![Layotto Dev Pipeline 🌊](https://github.com/mosn/layotto/actions/workflows/layotto-ci.yml/badge.svg)](https://github.com/mosn/layotto/actions/workflows/layotto-ci.yml)

[![GoDoc](https://godoc.org/mosn.io/layotto?status.svg)](https://godoc.org/mosn.io/layotto)
[![Go Report Card](https://goreportcard.com/badge/github.com/mosn/layotto)](https://goreportcard.com/report/mosn.io/layotto)
[![codecov](https://codecov.io/gh/mosn/layotto/branch/main/graph/badge.svg?token=10RxwSV6Sz)](https://codecov.io/gh/mosn/layotto)
[![Average time to resolve an issue](http://isitmaintained.com/badge/resolution/mosn/layotto.svg)](http://isitmaintained.com/project/mosn/layotto "Average time to resolve an issue")

</div>

Layotto(/leɪˈɒtəʊ/) 是一款使用 Golang 开发的应用运行时, 旨在帮助开发人员快速构建云原生应用，帮助应用和基础设施解耦。它为应用提供了各种分布式能力，比如状态管理，配置管理，事件发布订阅等能力，以简化应用的开发。

Layotto 以开源的 [MOSN](https://github.com/mosn/mosn) 为底座，在提供分布式能力以外，提供了 Service Mesh 对于流量的管控能力。

## 诞生背景

Layotto 希望可以把 [Multi-Runtime](https://www.infoq.com/articles/multi-runtime-microservice-architecture/) 跟 Service
Mesh 两者的能力结合起来，无论你是使用 MOSN 还是 Envoy 或者其他产品作为 Service Mesh 的数据面，都可以在不增加新的 sidecar 的前提下，使用 Layotto 为这些数据面追加 Runtime 的能力。

例如，通过为 MOSN 添加 Runtime 能力，一个 Layotto 进程可以[既作为 istio 的数据面](zh/start/istio/start.md) 又提供各种 Runtime API（例如 Configuration API,Pub/Sub API 等）

此外，随着探索实践，我们发现 sidecar 能做的事情远不止于此。 通过引入[WebAssembly](https://en.wikipedia.org/wiki/WebAssembly) ,我们正在尝试将 Layotto 做成 FaaS (Function as a service)
和 [reloadable sdk](https://github.com/mosn/layotto/issues/166) 的运行时容器 。

如果您对诞生背景感兴趣，可以看下[这篇演讲](https://mosn.io/layotto/#/zh/blog/mosn-subproject-layotto-opening-a-new-chapter-in-service-grid-application-runtime/index)
。

## 功能

- 服务通信
- 服务治理，例如流量的劫持和观测，服务限流等
- [作为 istio 的数据面](zh/start/istio/start.md)
- 配置管理
- 状态管理
- 事件发布订阅
- 健康检查、查询运行时元数据
- 基于 WASM 的多语言编程

## 工程架构

如下图架构图所示，Layotto 以开源 MOSN 作为底座，在提供了网络层管理能力的同时提供了分布式能力，业务可以通过轻量级的 SDK 直接与 Layotto 进行交互，而无需关注后端的具体的基础设施。

Layotto 提供了多种语言版本的 SDK，SDK 通过 gRPC 与 Layotto 进行交互。

如果您想把应用部署到不同的云平台（例如将阿里云上的应用部署到 AWS），您只需要在 Layotto 提供的 [配置文件](https://github.com/mosn/layotto/blob/main/configs/runtime_config.json)
里修改配置、指定自己想用的基础设施类型，不需要修改应用的代码就能让应用拥有"跨云部署"能力，大大提高了程序的可移植性。

![系统架构图](https://raw.githubusercontent.com/mosn/layotto/main/docs/img/runtime-architecture.png)

## 快速开始

### Get started with Layotto

您可以尝试 demo [通过 Layotto 调用 apollo 配置中心](zh/start/configuration/start-apollo.md) 来体验 Layotto

其他功能的 demo 见下.

### API

| API            | status |                              quick start                              |                                components                                 | desc                             |
| -------------- | :----: | :-------------------------------------------------------------------: | :-----------------------------------------------------------------------: | -------------------------------- |
| State          |   ✅   |        [demo](https://mosn.io/layotto/#/zh/start/state/start)         |     [list](https://mosn.io/layotto/#/zh/component_specs/state/common)     | 提供读写 KV 模型存储的数据的能力 |
| Pub/Sub        |   ✅   |        [demo](https://mosn.io/layotto/#/zh/start/pubsub/start)        |     [list](https://mosn.io/layotto/#/zh/component_specs/pubsub/redis)     | 提供消息的发布/订阅能力          |
| Service Invoke |   ✅   |       [demo](https://mosn.io/layotto/#/zh/start/rpc/helloworld)       |         [list](https://mosn.io/layotto/#/zh/start/rpc/helloworld)         | 通过 MOSN 进行服务调用           |
| Config         |   ✅   | [demo](https://mosn.io/layotto/#/zh/start/configuration/start-apollo) | [list](https://mosn.io/layotto/#/zh/component_specs/configuration/apollo) | 提供配置增删改查及订阅的能力     |
| Lock           |   ✅   |         [demo](https://mosn.io/layotto/#/zh/start/lock/start)         |     [list](https://mosn.io/layotto/#/zh/component_specs/lock/common)      | 提供 lock/unlock 分布式锁的实现  |
| Sequencer      |   ✅   |      [demo](https://mosn.io/layotto/#/zh/start/sequencer/start)       |   [list](https://mosn.io/layotto/#/zh/component_specs/sequencer/common)   | 提供获取分布式自增 ID 的能力     |
| File           |   ✅   |         [demo](https://mosn.io/layotto/#/zh/start/file/start)         |     [list](https://mosn.io/layotto/#/zh/component_specs/file/common)      | 提供访问文件的能力               |
| Binding        |   ✅   |                                 TODO                                  |                                   TODO                                    | 提供透传数据的能力               |

### 可扩展性

| feature  | status |                           quick start                            | desc                        |
| -------- | :----: | :--------------------------------------------------------------: | --------------------------- |
| API 插件 |   ✅   | [demo](https://mosn.io/layotto/#/zh/start/api_plugin/helloworld) | 为 Layotto 添加您自己的 API |

### 可观测性


| feature    | status |                         quick start                         | desc                  |
|------------| :----: |:-----------------------------------------------------------:|-----------------------|
| Skywalking |   ✅   | [demo](https://mosn.io/layotto/#/zh/start/trace/skywalking) | Layotto 接入 Skywalking |


### Actuator

| feature        | status |                        quick start                        | desc                                  |
| -------------- | :----: | :-------------------------------------------------------: | ------------------------------------- |
| Health Check   |   ✅   | [demo](https://mosn.io/layotto/#/zh/start/actuator/start) | 查询 Layotto 依赖的各种组件的健康状态 |
| Metadata Query |   ✅   | [demo](https://mosn.io/layotto/#/zh/start/actuator/start) | 查询 Layotto 或应用对外暴露的元信息   |

### 流量控制

| feature      | status |                              quick start                              | desc                                       |
| ------------ | :----: | :-------------------------------------------------------------------: | ------------------------------------------ |
| TCP Copy     |   ✅   |   [demo](https://mosn.io/layotto/#/zh/start/network_filter/tcpcopy)   | 把 Layotto 收到的 TCP 数据 dump 到本地文件 |
| Flow Control |   ✅   | [demo](https://mosn.io/layotto/#/zh/start/stream_filter/flow_control) | 限制访问 Layotto 对外提供的 API            |

### 在 Sidecar 中用 WebAssembly (WASM) 写业务逻辑

| feature        | status |                      quick start                      | desc                                                             |
| -------------- | :----: | :---------------------------------------------------: | ---------------------------------------------------------------- |
| Go (TinyGo)    |   ✅   | [demo](https://mosn.io/layotto/#/zh/start/wasm/start) | 把用 TinyGo 开发的代码编译成 \*.wasm 文件跑在 Layotto 上         |
| Rust           | 待开发 |                        待开发                         | 把用 Rust 开发的代码编译成 \*.wasm 文件跑在 Layotto 上           |
| AssemblyScript | 待开发 |                        待开发                         | 把用 AssemblyScript 开发的代码编译成 \*.wasm 文件跑在 Layotto 上 |

### 作为 Serverless 的运行时，通过 WebAssembly (WASM) 写 FaaS

| feature        | status |                      quick start                      | desc                                                                                      |
| -------------- | :----: | :---------------------------------------------------: | ----------------------------------------------------------------------------------------- |
| Go (TinyGo)    |   ✅   | [demo](https://mosn.io/layotto/#/zh/start/faas/start) | 把用 TinyGo 开发的代码编译成 \*.wasm 文件跑在 Layotto 上，并且使用 k8s 进行调度。         |
| Rust           | 待开发 |                        待开发                         | 把用 Rust 开发的代码编译成 \*.wasm 文件跑在 Layotto 上，并且使用 k8s 进行调度。           |
| AssemblyScript | 待开发 |                        待开发                         | 把用 AssemblyScript 开发的代码编译成 \*.wasm 文件跑在 Layotto 上，并且使用 k8s 进行调度。 |

### Service Mesh

| feature | status |                      quick start                       | desc                          |
| ------- | :----: | :----------------------------------------------------: | ----------------------------- |
| istio   |   ✅   | [demo](https://mosn.io/layotto/#/zh/start/istio/start) | 跟 istio 集成，作为它的数据面 |

## Landscapes

<p align="center">
<img src="https://landscape.cncf.io/images/left-logo.svg" width="150"/>&nbsp;&nbsp;<img src="https://landscape.cncf.io/images/right-logo.svg" width="200"/>
<br/><br/>
Layotto enriches the <a href="https://landscape.cncf.io/serverless">CNCF CLOUD NATIVE Landscape.</a>
</p>

## 社区

| 平台                                          | 联系方式                                                                                                                                             |
| :-------------------------------------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------- |
| 💬 [钉钉](https://www.dingtalk.com/zh) (用户群) | 群号: 31912621 或者扫描下方二维码 <br> <img src="https://raw.githubusercontent.com/mosn/layotto/main/docs/img/ding-talk-group-1.png" height="200px"> <br> |
| 💬 [钉钉](https://www.dingtalk.com/zh) (社区会议群) | 群号：41585216 <br> [Layotto 在每周五晚 8 点进行社区会议，欢迎所有人](zh/community/meeting.md) |

[comment]: <> (| 💬 [微信]&#40;https://www.wechat.com/&#41; | 扫描下方二维码添加好友，她会邀请您加入微信群 <br> <img src="../img/wechat-group.jpg" height="200px">)

## 如何贡献

[新手攻略：从零开始成为 Layotto 贡献者](zh/development/start-from-zero.md)

[从哪下手？看看"新手任务"列表](https://github.com/mosn/layotto/issues/108#issuecomment-872779356)

作为技术同学，你是否有过“想参与某个开源项目的开发、但是不知道从何下手”的感觉？
为了帮助大家更好的参与开源项目，社区会定期发布适合新手的新手开发任务，帮助大家 learning by doing!

[文档贡献指南](zh/development/contributing-doc.md)

[组件开发指南](zh/development/developing-component.md)

[Layotto Github Workflow 指南](zh/development/github-workflows.md)

[Layotto 命令行指南](zh/development/commands.md)

[Layotto 贡献者指南](zh/development/CONTRIBUTING.md)

## 设计文档

[Actuator 设计文档](zh/design/actuator/actuator-design-doc.md)

[Pubsub API 与 Dapr Component 的兼容性](zh/design/pubsub/pubsub-api-and-compability-with-dapr-component.md)

[Configuration API with Apollo(英文)](en/design/configuration/configuration-api-with-apollo.md)

[RPC 设计文档](zh/design/rpc/rpc设计文档.md)

[分布式锁 API 设计文档](zh/design/lock/lock-api-design.md)

[FaaS 设计文档](zh/design/faas/faas-poc-design.md)

## FAQ

### 跟 dapr 有什么差异？

dapr 是一款优秀的 Runtime 产品，但它本身缺失了 Service Mesh 的能力，而这部分能力对于实际在生产环境落地是至关重要的，因此我们希望把 Runtime 跟 Service Mesh 两种能力结合在一起，满足更复杂的生产落地需求。
