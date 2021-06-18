# Layotto (L8):To be the next layer of OSI layer 7
[查看中文版本](README-ZH.md)

Layotto is an application runtime developed using Golang, which provides various distributed capabilities for applications, such as state management, configuration management, and event-driven capabilities to simplify application development.

Layotto uses the open source [MOSN](https://github.com/mosn/mosn) as the base, in addition to providing distributed capabilities, it also provides Service Mesh's ability to control traffic.

## Features

- Hijacking and observation of data traffic
- The current limiting capability of the service
- Configuration center read and write monitoring capabilities

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

## Contributing to Layotto

See the Development Guide [contributing](CONTRIBUTING.md) to get started with building and developing.

## Community

### Contact Us

Use [DingTalk](https://www.dingtalk.com/en) to scan the QR code below to join the Layotto user exchange group.

![Ding Talk Group QR Code](img/ding-talk-group-1.jpg)

Or through Ding Talk search group number 31912621, join the user exchange group.