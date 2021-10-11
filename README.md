# Layotto (L8):To be the next layer of OSI layer 7

[![codecov](https://codecov.io/gh/mosn/layotto/branch/main/graph/badge.svg?token=10RxwSV6Sz)](https://codecov.io/gh/mosn/layotto)
[![Average time to resolve an issue](http://isitmaintained.com/badge/resolution/mosn/layotto.svg)](http://isitmaintained.com/project/mosn/layotto "Average time to resolve an issue")

<img src="https://raw.githubusercontent.com/mosn/layotto/main/docs/img/logo/grey2-1.svg" height="120px">

[查看中文版本](https://mosn.io/layotto/#/zh/README)

Layotto is an application runtime developed using Golang, which provides various distributed capabilities for applications, such as state management, configuration management, and event pub/sub capabilities to simplify application development.

Layotto is built on the open source data plane [MOSN](https://github.com/mosn/mosn) .In addition to providing distributed building blocks, Layotto can also serve as the data plane of Service Mesh and has the ability to control traffic.

## Motivation

Layotto aims to combine [Multi-Runtime](https://www.infoq.com/articles/multi-runtime-microservice-architecture/) with Service Mesh into one sidecar. No matter which product you are using as the Service Mesh data plane (e.g. MOSN,Envoy or any other product), you can always attach Layotto to it and add Multi-Runtime capabilities without adding new sidecars. 

For example, by adding Runtime capabilities to MOSN, a Layotto process can both [serve as the data plane of istio](https://mosn.io/layotto/#/en/start/istio/start.md) and provide various Runtime APIs (such as Configuration API, Pub/Sub API, etc.)

In addition, we were surprised to find that a sidecar can do much more than that. We are trying to make Layotto even the runtime container of FaaS (Function as a service) and [reloadable sdk](https://github.com/mosn/layotto/issues/166) with the magic power of [WebAssembly](https://en.wikipedia.org/wiki/WebAssembly) .

## Features

- Service Communication
- Service Governance.Such as traffic hijacking and observation, service rate limiting, etc
- [As the data plane of istio](https://mosn.io/layotto/#/en/start/istio/start)
- Configuration management
- State management
- Event publish and subscribe
- Health check, query runtime metadata
- [FaaS model based on WASM and Runtime](docs/en/start/faas/start.md)

## Project Architecture

As shown in the architecture diagram below, Layotto uses the open source MOSN as the base to provide network layer management capabilities while providing distributed capabilities. The business logic can directly interact with Layotto through a lightweight SDK without paying attention to the specific back-end infrastructure.

Layotto provides sdks in various languages. The sdk interacts with Layotto through grpc. Application developers only need to specify their own infrastructure type through the configuration file [configure file](./configs/runtime_config.json) provided by Layotto. No coding changes are required, which greatly improves the portability of the program.

![Architecture](https://raw.githubusercontent.com/mosn/layotto/main/docs/img/runtime-architecture.png)

## Quickstarts and Samples

### Get started with Layotto

See the quick start guide [configuration demo with apollo](https://mosn.io/layotto/#/en/start/configuration/start-apollo) that can help you get started with Layotto.

### Use Pub/Sub API

[Implementing Pub/Sub Pattern using Layotto and Redis](https://mosn.io/layotto/#/en/start/pubsub/start)

### Use State API to manage state

[State management demo with redis](https://mosn.io/layotto/#/en/start/state/start)

### Use Distributed Lock API

[Distributed Lock API demo with redis](https://mosn.io/layotto/#/en/start/lock/start)

### Traffic intervention on the 4th layer network

[Dump TCP Traffic](https://mosn.io/layotto/#/en/start/network_filter/tcpcopy)

### Flow Control on the 7th layer network

[Method Level Flow Control](https://mosn.io/layotto/#/en/start/stream_filter/flow_control)

### Health check and metadata query

[Use Layotto Actuator for health check and metadata query](https://mosn.io/layotto/#/en/start/actuator/start)

### Service Invocation

[Hello World](https://mosn.io/layotto/#/en/start/rpc/helloworld)

[Dubbo JSON RPC](https://mosn.io/layotto/#/en/start/rpc/dubbo_json_rpc)

### Integrate with istio

[As the data plane of istio](https://mosn.io/layotto/#/en/start/istio/start)

### FaaS model based on WASM and Runtime

[FaaS on Layotto](docs/en/start/faas/start.md)

## Design Documents

[actuator-design-doc](https://mosn.io/layotto/#/en/design/actuator/actuator-design-doc)

[configuration-api-with-apollo](https://mosn.io/layotto/#/en/design/configuration/configuration-api-with-apollo)

[pubsub-api-and-compability-with-dapr-component](https://mosn.io/layotto/#/en/design/pubsub/pubsub-api-and-compability-with-dapr-component)

[rpc-design-doc](https://mosn.io/layotto/#/en/design/rpc/rpc-design-doc)

[distributed-lock-api-design](https://mosn.io/layotto/#/en/design/lock/lock-api-design)

[FaaS design](https://mosn.io/layotto/#/en/design/faas/faas-poc-design.md)

## Community

### Contact Us

| Platform  | Link        |
|:----------|:------------|
| 💬 [DingTalk](https://www.dingtalk.com/en) (preferred) | Search the group number: 31912621 or scan the QR code below <br> <img src="https://raw.githubusercontent.com/mosn/layotto/main/docs/img/ding-talk-group-1.png" height="200px">

[comment]: <> (| 💬 [Wechat]&#40;https://www.wechat.com/en/&#41;  | Scan the QR code below and she will invite you into the wechat group <br> <img src="/img/wechat-group.jpg" height="200px">)

## Contributing to Layotto

See the [Development Guide](https://mosn.io/layotto/#/en/development/CONTRIBUTING) to get started with building and developing.

## FAQ

### Difference with dapr?

dapr is an excellent Runtime product, but it lacks the ability of Service Mesh, which is necessary for the Runtime 
product used in production environment, so we hope to combine Runtime with Service Mesh into one sidecar to meet 
more complex production requirements.
