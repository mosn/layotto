# Layotto (L8):To be the next layer of OSI layer 7
<img src="/img/logo/grey2-1.svg" height="120px">

[查看中文版本](README-ZH.md)

Layotto is a network proxy developed using Golang, which can be used as the sidecar of Service Mesh, the [multi-runtime of microservices and cloud native applications](https://www.infoq.com/articles/multi-runtime-microservice-architecture/) and also the runtime of Serverless.

As the multi-runtime of cloud native applications,Layotto provides various distributed building blocks for applications, such as state management, configuration management, and event-driven capabilities to simplify application development.At the same time, Layotto gives application portability, enabling cross-cloud deployment and unbinding with cloud vendors.

As the sidecar of Service Mesh,Layotto is based on the open source network proxy [MOSN](https://github.com/mosn/mosn) .In addition to providing service communication capabilities, it also provides Service Mesh's ability of service governance.

So,what is Layotto?

Layotto = Service Mesh + Multi-Runtime + Serverless Runtime = Layer8

## Features

- Service Mesh
    - Service communication
    - Service governance. Such as the hijacking and observation of data traffic, service rate limiting, etc.
- [Multi-Runtime](https://www.infoq.com/articles/multi-runtime-microservice-architecture/)
    - State management
    - Configuration management
    - Event pub/sub
    - Distributed lock
    - ...
- Serverless runtime based on WASM
- Observability
    - health examination
    - Query runtime metadata
    - ...
- Others: all the cool features of [MOSN](https://github.com/mosn/mosn) , such as smooth hot upgrade

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

### Service Invocation

[Hello World](docs/en/start/rpc/helloworld.md)

[Dubbo JSON RPC](docs/en/start/rpc/dubbo_json_rpc.md)

### Traffic intervention on the 4th layer network

[Dump TCP Traffic](docs/en/start/network_filter/tcpcopy.md)

### Flow Control on the 7th layer network

[Method Level Flow Control](docs/en/start/stream_filter/flow_control.md)

### Health check and metadata query

[Use Layotto Actuator for health check and metadata query](docs/en/start/actuator/start.md)

## Community

### Contact Us

| Platform  | Link        |
|:----------|:------------|
| 💬 [DingTalk](https://www.dingtalk.com/en) (preferred) | Search the group number: 31912621 or scan the QR code below <br> <img src="/img/ding-talk-group-1.png" height="200px">
| 💬 [Wechat](https://www.wechat.com/en/)  | Scan the QR code below and she will invite you into the wechat group <br> <img src="/img/wechat-group.jpg" height="200px">

## Contributing to Layotto

See the Development Guide [contributing](CONTRIBUTING.md) to get started with building and developing.
