# Layotto (L8):To be the next layer of OSI layer 7

[![codecov](https://codecov.io/gh/mosn/layotto/branch/main/graph/badge.svg?token=10RxwSV6Sz)](https://codecov.io/gh/mosn/layotto)
[![Average time to resolve an issue](http://isitmaintained.com/badge/resolution/mosn/layotto.svg)](http://isitmaintained.com/project/mosn/layotto "Average time to resolve an issue")

<img src="https://raw.githubusercontent.com/mosn/layotto/main/docs/img/logo/grey2-1.svg" height="120px">

Layotto(/leɪˈɒtəʊ/) is an application runtime developed using Golang, which provides various distributed capabilities for applications, such as state management, configuration management, and event pub/sub capabilities to simplify application development.

Layotto is built on the open source data plane [MOSN](https://github.com/mosn/mosn) .In addition to providing distributed building blocks, Layotto can also serve as the data plane of Service Mesh and has the ability to control traffic.

## Motivation
Layotto aims to combine [Multi-Runtime](https://www.infoq.com/articles/multi-runtime-microservice-architecture/) with Service Mesh into one sidecar. No matter which product you are using as the Service Mesh data plane (e.g. MOSN,Envoy or any other product), you can always attach Layotto to it and add Multi-Runtime capabilities without adding new sidecars.

For example, by adding Runtime capabilities to MOSN, a Layotto process can both [serve as the data plane of istio](https://mosn.io/layotto/#/en/start/istio/start.md) and provide various Runtime APIs (such as Configuration API, Pub/Sub API, etc.)

In addition, we were surprised to find that a sidecar can do much more than that. We are trying to make Layotto even the runtime container of FaaS (Function as a service) and [reloadable sdk](https://github.com/mosn/layotto/issues/166) with the magic power of [WebAssembly](https://en.wikipedia.org/wiki/WebAssembly) .

## Features

- Service Communication
- Service Governance.Such as traffic hijacking and observation, service rate limiting, etc
- [As the data plane of istio](en/start/istio/start)
- Configuration management
- State management
- Event publish and subscribe
- Health check, query runtime metadata
- Multilingual programming based on WASM

## Project Architecture

As shown in the architecture diagram below, Layotto uses the open source MOSN as the base to provide network layer management capabilities while providing distributed capabilities. The business logic can directly interact with Layotto through a lightweight SDK without paying attention to the specific back-end infrastructure.

Layotto provides sdk in various languages. The sdk interacts with Layotto through grpc. Application developers only need to specify their own infrastructure type through the [configuration file](https://github.com/mosn/layotto/blob/main/configs/runtime_config.json) provided by Layotto. No coding changes are required, which greatly improves the portability of the program.

![Architecture](https://raw.githubusercontent.com/mosn/layotto/main/docs/img/runtime-architecture.png)

## Quickstarts

### Get started with Layotto

You can try the [configuration demo with apollo](en/start/configuration/start-apollo.md) to get started with Layotto.

For other features,see the demos below:

### API

|  API            | status |                               quick start                             |                                components                                 | desc |
|  -------------  | :----: | :--------------------------------------------------------------------:|:-------------------------------------------------------------------------:|---- |
| State           | ✅     | [demo](https://mosn.io/layotto/#/en/start/state/start)                | [list](https://mosn.io/layotto/#/en/component_specs/state/common)         | Write/Query the data of the Key/Value model |
| Pub/Sub         | ✅     | [demo](https://mosn.io/layotto/#/en/start/pubsub/start)               | [list](https://mosn.io/layotto/#/en/component_specs/pubsub/redis)         | Publish/Subscribe message through various Message Queue |
| Service Invoke  | ✅     | [demo](https://mosn.io/layotto/#/en/start/rpc/helloworld)             | [list](https://mosn.io/layotto/#/en/start/rpc/helloworld)                 | Call Service through MOSN (another istio data plane)|
| Config          | ✅     | [demo](https://mosn.io/layotto/#/en/start/configuration/start-apollo) | [list](https://mosn.io/layotto/#/en/component_specs/configuration/apollo) | Write/Query/Subscribe the config through various Config Center|
| Lock            | ✅     | [demo](https://mosn.io/layotto/#/en/start/lock/start)                 | [list](https://mosn.io/layotto/#/en/component_specs/lock/common)          | Distribute lock implementation|
| Sequencer       | ✅     | [demo](https://mosn.io/layotto/#/en/start/sequencer/start)            | [list](https://mosn.io/layotto/#/en/component_specs/sequencer/common)     | Distribube auto increment ID generator |
| File            | ✅     | TODO                                                                  | [list](https://mosn.io/layotto/#/en/component_specs/file/oss)             | File API implementation|
| Binding         | ✅     | TODO                                                                  | TODO                                                                      | Transparent data transmission API |

### Actuator

|  feature       | status |                         quick start                       |                         desc                         |
|  ------------- | :----: | :--------------------------------------------------------:|------------------------------------------------------|
| Health Check   | ✅     | [demo](https://mosn.io/layotto/#/en/start/actuator/start) | Query health state of app and components in Layotto  |
| Metadata Query | ✅     | [demo](https://mosn.io/layotto/#/en/start/actuator/start) | Query metadata in Layotto/app                        |

### Traffic Control

|  feature      | status |                              quick start                              |                               desc                              |
|  -----------  | :----: | :--------------------------------------------------------------------:|-----------------------------------------------------------------|
| TCP Copy      | ✅     | [demo](https://mosn.io/layotto/#/en/start/network_filter/tcpcopy)     | Dump the tcp traffic received by Layotto into local file system |
| Flow Control  | ✅     | [demo](https://mosn.io/layotto/#/en/start/stream_filter/flow_control) | limit access to the APIs provided by Layotto                    |

### Layotto + WebAssembly

|  feature       | status |                       quick start                      |                               desc                                  |
|  ------------- | :----: | :-----------------------------------------------------:|---------------------------------------------------------------------|
| Go (TinyGo)    | ✅     | [demo](https://mosn.io/layotto/#/en/start/wasm/start)  | Compile Code written by TinyGo to *.wasm and run in Layotto         |
| Rust           | TODO     | TODO | Compile Code written by Rust to *.wasm and run in Layotto           |
| AssemblyScript | TODO     | TODO | Compile Code written by AssemblyScript to *.wasm and run in Layotto |

### FaaS (Layotto + WebAssembly + k8s)

|  feature       | status |                       quick start                      |                               desc                                  |
|  ------------- | :----: | :-----------------------------------------------------:|---------------------------------------------------------------------|
| Go (TinyGo)    | ✅     | [demo](https://mosn.io/layotto/#/en/start/faas/start)  | Compile Code written by TinyGo to *.wasm and run in Layotto And Scheduled by k8s. |
| Rust           | TODO     | TODO  | Compile Code written by Rust to *.wasm and run in Layotto And Scheduled by k8s.            |
| AssemblyScript | TODO     | TODO  | Compile Code written by AssemblyScript to *.wasm and run in Layotto And Scheduled by k8s.  |

### Service Mesh
| feature | status |                       quick start                      |            desc            |
| ------- | :----: | :-----------------------------------------------------:|----------------------------|
| istio   | ✅     | [demo](https://mosn.io/layotto/#/en/start/istio/start) | As the data plane of istio |


## Design Documents

[actuator-design-doc](en/design/actuator/actuator-design-doc.md)

[configuration-api-with-apollo](en/design/configuration/configuration-api-with-apollo.md)

[pubsub-api-and-compability-with-dapr-component](en/design/pubsub/pubsub-api-and-compability-with-dapr-component.md)

[rpc-design-doc](en/design/rpc/rpc-design-doc.md)

[distributed-lock-api-design](en/design/lock/lock-api-design.md)

[FaaS design](en/design/faas/faas-poc-design.md)

## Presentations
  * [Layotto - A new chapter of Service Mesh and Application Runtime](https://www.youtube.com/watch?v=5v8gTrFUDk8)
  * [WebAssembly + Application Runtime = A New Era of FaaS?](https://www.youtube.com/watch?v=g01CJ4S9Qao)

## Community

### Contact Us

| Platform  | Link        |
|:----------|:------------|
| 💬 [DingTalk](https://www.dingtalk.com/en) (preferred) | Search the group number: 31912621 or scan the QR code below <br> <img src="https://raw.githubusercontent.com/mosn/layotto/main/docs/img/ding-talk-group-1.png?raw=true" height="200px">

[comment]: <> (| 💬 [Wechat]&#40;https://www.wechat.com/en/&#41;  | Scan the QR code below and she will invite you into the wechat group <br> <img src="img/wechat-group.jpg" height="200px">)

## Contributing to Layotto

See the Development Guide [contributing](CONTRIBUTING.md) to get started with building and developing.

## FAQ

### Difference with dapr?

dapr is an excellent Runtime product, but it lacks the ability of Service Mesh, which is necessary for the Runtime 
product used in production environment, so we hope to combine Runtime with Service Mesh into one sidecar to meet 
more complex production requirements.
