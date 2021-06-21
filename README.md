# Layotto (L8):To be the next layer of OSI layer 7
<img src="/img/logo/grey2-1.svg" height="120px">

[查看中文版本](README-ZH.md)

Layotto is an application runtime developed using Golang, which provides various distributed capabilities for applications, such as state management, configuration management, and event pub/sub capabilities to simplify application development.

Layotto uses the open source [MOSN](https://github.com/mosn/mosn) as the base, in addition to providing distributed capabilities, it also provides Service Mesh's ability to control traffic.

## Features

- Service Communication
- Service Governance.Such as traffic hijacking and observation, service rate limiting, etc
- Configuration management
- State management
- Event publish and subscribe
- Health check, query runtime metadata
- Multilingual programming based on WASM

## Project Architecture

As shown in the architecture diagram below, Layotto uses the open source MOSN as the base to provide network layer management capabilities while providing distributed capabilities. The business logic can directly interact with Layotto through a lightweight SDK without paying attention to the specific back-end infrastructure.

Layotto provides sdk in various languages. The sdk interacts with Layotto through grpc. Application developers only need to specify their own infrastructure type through the configuration file [configure file](./configs/runtime_config.json) provided by Layotto. No coding changes are required, which greatly improves the portability of the program.

![Architecture](img/runtime-architecture.png)

## Quickstarts and Samples

### Get started with Layotto

See the quick start guide [configuration demo with apollo](docs/en/start/configuration/start-apollo.md) that can help you get started with Layotto.

### Use Pub/Sub API

[Implementing Pub/Sub Pattern using Layotto and Redis](docs/en/start/pubsub/start.md)

### Use State API to manage state

[State management demo with redis](docs/en/start/state/start.md)

### Traffic intervention on the 4th layer network

[Dump TCP Traffic](docs/en/start/network_filter/tcpcopy.md)

### Flow Control on the 7th layer network

[Method Level Flow Control](docs/en/start/stream_filter/flow_control.md)

### Health check and metadata query

[Use Layotto Actuator for health check and metadata query](docs/en/start/actuator/start.md)

### Service Invocation

[Hello World](docs/en/start/rpc/helloworld.md)

[Dubbo JSON RPC](docs/en/start/rpc/dubbo_json_rpc.md)

### Multilingual programming based on WASM

[WASM on Layotto](docs/en/start/wasm/start.md)

## Design Documents

[configuration-api-with-apollo](docs/en/design/configuration/configuration-api-with-apollo.md)

[pubsub-api-and-compability-with-dapr-component](docs/en/design/pubsub/pubsub-api-and-compability-with-dapr-component.md)

[rpc-design-doc](docs/en/design/rpc/rpc-design-doc.md)

[actuator-design-doc](docs/en/design/actuator/actuator-design-doc.md)

## Community

### Contact Us

| Platform  | Link        |
|:----------|:------------|
| 💬 [DingTalk](https://www.dingtalk.com/en) (preferred) | Search the group number: 31912621 or scan the QR code below <br> <img src="/img/ding-talk-group-1.png" height="200px">

[comment]: <> (| 💬 [Wechat]&#40;https://www.wechat.com/en/&#41;  | Scan the QR code below and she will invite you into the wechat group <br> <img src="/img/wechat-group.jpg" height="200px">)

## Contributing to Layotto

See the Development Guide [contributing](CONTRIBUTING.md) to get started with building and developing.
