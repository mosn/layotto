# Layotto (L8):To be the next layer of OSI layer 7

[![codecov](https://codecov.io/gh/mosn/layotto/branch/main/graph/badge.svg?token=10RxwSV6Sz)](https://codecov.io/gh/mosn/layotto)

<img src="https://raw.githubusercontent.com/mosn/layotto/main/docs/img/logo/grey2-1.svg" height="120px">

[查看中文版本](https://mosn.io/layotto/#/zh/README)

Layotto is an application runtime developed using Golang, which provides various distributed capabilities for applications, such as state management, configuration management, and event pub/sub capabilities to simplify application development.

Layotto uses the open source [MOSN](https://github.com/mosn/mosn) as the base, in addition to providing distributed capabilities, it also provides Service Mesh's ability to control traffic.

## Motivation

Layotto aims to combine Runtime with Service Mesh into one sidecar. No matter which product you are using as the Service Mesh data plane (e.g. Mosn,Envoy or any other product), you can always attach Layotto to it and add Multi-Runtime capabilities without adding new sidecars. 

For example, the Layotto process can both [serve as the data plane of istio](https://mosn.io/layotto/#/en/start/istio/start.md) and provide various Runtime APIs (such as Configuration API, Pub/Sub API, etc.)

## Features

- Service Communication
- Service Governance.Such as traffic hijacking and observation, service rate limiting, etc
- [As the data plane of istio](https://mosn.io/layotto/#/en/start/istio/start)
- Configuration management
- State management
- Event publish and subscribe
- Health check, query runtime metadata
- Multilingual programming based on WASM

## Project Architecture

As shown in the architecture diagram below, Layotto uses the open source MOSN as the base to provide network layer management capabilities while providing distributed capabilities. The business logic can directly interact with Layotto through a lightweight SDK without paying attention to the specific back-end infrastructure.

Layotto provides sdk in various languages. The sdk interacts with Layotto through grpc. Application developers only need to specify their own infrastructure type through the configuration file [configure file](./configs/runtime_config.json) provided by Layotto. No coding changes are required, which greatly improves the portability of the program.

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

[As a data plane in istio](https://mosn.io/layotto/#/en/start/istio/start)

### Multilingual programming based on WASM

[WASM on Layotto](https://mosn.io/layotto/#/en/start/wasm/start)

## Design Documents

[actuator-design-doc](https://mosn.io/layotto/#/en/design/actuator/actuator-design-doc)

[configuration-api-with-apollo](https://mosn.io/layotto/#/en/design/configuration/configuration-api-with-apollo)

[pubsub-api-and-compability-with-dapr-component](https://mosn.io/layotto/#/en/design/pubsub/pubsub-api-and-compability-with-dapr-component)

[rpc-design-doc](https://mosn.io/layotto/#/en/design/rpc/rpc-design-doc)

[distributed-lock-api-design](https://mosn.io/layotto/#/en/design/lock/lock-api-design)

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
