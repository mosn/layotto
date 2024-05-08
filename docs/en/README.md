<div align="center">
  <h1>Layotto (L8): To be the next layer of OSI layer 7</h1>
  <img src="https://gw.alipayobjects.com/zos/bmw-prod/65518bfc-8ba5-4234-a5c5-2bc065e3a5f0.svg" height="120px">
[! [Layotto Env Pipeline 🌊] (https://github.com/mosn/layotto/actions/workflows/quickstart-checker.yml/badge.svg)] (https://github.com/mosn/layotto/actions/workflows/quickstart-checker.yml)
[! [Layotto Dev Pipeline 🌊] (https://github.com/mosn/layotto/actions/workflows/layotto-ci.yml/badge.svg)] (https://github.com/mosn/layotto/actions/workflows/layotto-ci.yml)

[! [GoDoc] (https://godoc.org/mosn.io/layotto?status.svg)] (https://godoc.org/mosn.io/layotto)
[! [Go Report Card] (https://goreportcard.com/badge/github.com/mosn/layotto)] (https://goreportcard.com/report/mosn.io/layotto)
[! [codecov] (https://codecov.io/gh/mosn/layotto/branch/main/graph/badge.svg?token=10RxwSV6Sz)] (https://codecov.io/gh/mosn/layotto)
[! [Average time to resolve an issue] (http://isitmaintained.com/badge/resolution/mosn/layotto.svg)] (http://isitmaintained.com/project/mosn/layotto "Average time to resolve an issue")

</div>

Layotto (/leɪˈɒtəʊ/) is an application runtime developed with Golang designed to help developers quickly build cloud-native applications that decouple applications from infrastructure. It provides various distributed capabilities for applications, such as state management, configuration management, and event publishing and subscribing, to simplify application development.

Based on the open source [MOSN](https://github.com/mosn/mosn), Layotto provides Service Mesh with the ability to manage and control traffic in addition to distributed capabilities.

## Birth Background

Layotto wants to be able to use [Multi-Runtime](https://www.infoq.com/articles/multi-runtime-microservice-architecture/) with Service
Whether you're using MOSN or Envoy or something else as the data plane of a Service Mesh, you can use Layotto to append runtimes to those data planes without adding new sidecars.

For example, by adding runtime capabilities to MOSNs, a Layotto process can [act as both the data side of istio] (en/start/istio/) and provide various runtime APIs (e.g., Configuration API, Pub/Sub API, etc.)

In addition, as we explored and practiced, we found that sidecars can do much more than that. With the introduction of WebAssembly(https://en.wikipedia.org/wiki/WebAssembly), we are trying to make Layotto a runtime container for FaaS (Function as a service).

If you are interested in the background of the birth, you can take a look at [this speech] (https://mosn.io/layotto/#/zh/blog/mosn-subproject-layotto-opening-a-new-chapter-in-service-grid-application-runtime/index)
。

## Features

- Service Communications
-wearService governance, such as traffic hijacking and observation, service throttling, etc
- [as istio's data plane](en/start/istio/)
- Configuration management
- State management
- Event publishing and subscribing
- Health checks and runtime metadata queries
- WASM-based multilingual programming

## Engineering Architecture

As shown in the architecture diagram shown in the figure below, Layotto uses open-source MOSN as the base, providing network layer management capabilities and distributed capabilities at the same time, and businesses can directly interact with Layotto through lightweight SDKs without paying attention to the specific infrastructure of the backend.

Layotto provides multiple language versions of the SDK that interact with Layotto via gRPC.

If you want to deploy your application to a different cloud platform (for example, deploy an application on Alibaba Cloud to AWS), you only need to configure the file provided by Layotto https://github.com/mosn/layotto/blob/main/configs/runtime_config.json.
, modify the configuration and specify the type of infrastructure you want to use, and you can make the application have the ability of "cross-cloud deployment" without modifying the application code, which greatly improves the portability of the program.

<img src="https://gw.alipayobjects.com/mdn/rms_5891a1/afts/img/A*oRkFR63JB7cAAAAAAAAAAAAAARQnAQ" />
## Get started quickly

### Get started with Layotto

You can try the following Quickstart demo to experience the features of Layotto; Or try [Online Lab] (en/start/lab.md)

### API

| API            | status |                              quick start                              |                               desc                             |
| -------------- | :----: | :-------------------------------------------------------------------: | -------------------------------- |
| State|   ✅    |        [demo] (https://mosn.io/layotto/#/zh/start/state/start)         |     Provides the ability to read and write data stored by the KV model |
| Pub/Sub        |   ✅    |        [demo] (https://mosn.io/layotto/#/zh/start/pubsub/start)        |     Provides publish/subscribe capabilities for messages |
| Service Invoke |   ✅    |       [demo] (https://mosn.io/layotto/#/zh/start/rpc/helloworld)       |      Service calls via MOSN |
| Config         |   ✅    | [demo] (https://mosn.io/layotto/#/zh/start/configuration/start-apollo) |   Provides the ability to add, delete, modify, query, and subscribe
| Lock           |   ✅    |         [demo] (https://mosn.io/layotto/#/zh/start/lock/start)         |    Provides an implementation of lock/unlock distributed locks |
| Sequencer      |   ✅    |      [demo] (https://mosn.io/layotto/#/zh/start/sequencer/start)       |  Provides the ability to obtain distributed auto-increment IDs |
| File           |   ✅    |         [demo] (https://mosn.io/layotto/#/zh/start/file/start)         |   Provides the ability to access files |
| Binding        |   ✅    |                                 TODO                                  |  Provides the ability to transparently transmit data|

### Service Mesh

| feature | status |                      quick start                       | desc                          |
| ------- | :----: | :----------------------------------------------------: | ----------------------------- |
| Istio   |   ✅    | [demo] (https://mosn.io/layotto/#/zh/start/istio/) | Integrates with Istio as Istio's data plane |

### Scalability

| feature  | status |                           quick start                            | desc                        |
| -------- | :----: | :--------------------------------------------------------------: | --------------------------- |
| API Plugin |   ✅    | [demo] (https://mosn.io/layotto/#/zh/start/api_plugin/helloworld) | Add your own API for Layotto

### Observability

| feature    | status |                         quick start                         | desc                    |
| ---------- | :----: | :---------------------------------------------------------: | ----------------------- |
| Skywalking |   ✅    | [demo] (https://mosn.io/layotto/#/zh/start/trace/skywalking) | Layotto access Skywalking |

### Actuator

| feature        | status |                        quick start                        | desc                                  |
| -------------- | :----: | :-------------------------------------------------------: | ------------------------------------- |
| Health Check   |   ✅    | [demo] (https://mosn.io/layotto/#/zh/start/actuator/start) | Query the health status of the various components that Layotto depends on
| Metadata Query |   ✅    | [demo] (https://mosn.io/layotto/#/zh/start/actuator/start) | Query the meta information exposed by Layotto or the app |

### Flow control

| feature      | status |                              quick start                              | desc                                       |
| ------------ | :----: | :-------------------------------------------------------------------: | ------------------------------------------ |
| TCP Copy     |   ✅    |   [demo] (https://mosn.io/layotto/#/zh/start/network_filter/tcpcopy)   | Dump the TCP data received by Layotto to the local file |
| Flow Control |   ✅| [demo] (https://mosn.io/layotto/#/zh/start/stream_filter/flow_control) | Restrict access to Layotto's external APIs

### Write business logic in WebAssembly (WASM) in Sidecar

| feature        | status |                      quick start                      | desc                                                             |
| -------------- | :----: | :---------------------------------------------------: | ---------------------------------------------------------------- |
| Go (TinyGo)    |   ✅   | [demo] (https://mosn.io/layotto/#/zh/start/wasm/start) | Compile the code developed with TinyGo into a \*.wasm file and run it on Layotto |
| Rust           |   ✅   | [demo] (https://mosn.io/layotto/#/zh/start/wasm/start) | Compile code developed in Rust into a \*.wasm file and run it on Layotto |
| AssemblyScript |   ✅   | [demo] (https://mosn.io/layotto/#/zh/start/wasm/start) | Compile the code developed with AssemblyScript into a \*.wasm file and run it on Layotto |

### As a serverless runtime, write FaaS via WebAssembly (WASM).

| feature        | status |                      quick start| desc                                                                                      |
| -------------- | :----: | :---------------------------------------------------: | ----------------------------------------------------------------------------------------- |
| Go (TinyGo)    |   ✅   | [demo] (https://mosn.io/layotto/#/zh/start/faas/start) | Compile the code developed with TinyGo into a \*.wasm file and run it on Layotto, and use k8s for scheduling. |
| Rust           |   ✅   | [demo] (https://mosn.io/layotto/#/zh/start/faas/start) | Compile the code developed in Rust into a \*.wasm file and run it on Layotto, and use k8s for scheduling. |
| AssemblyScript |   ✅   | [demo] (https://mosn.io/layotto/#/zh/start/faas/start) | Compile the code developed with AssemblyScript into a \*.wasm file and run it on Layotto, and use k8s for scheduling. |

## Landscapes

<p align="center">
<img src="https://landscape.cncf.io/images/left-logo.svg" width="150"/>&nbsp;&nbsp;<img src="https://landscape.cncf.io/images/right-logo.svg" width="200"/>
<br/><br/>
Layotto enriches the <a href="https://landscape.cncf.io/serverless">CNCF CLOUD NATIVE Landscape.</a>
</p>

## Community

| Platform | Contact |
| :------------------------------------------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 💬 [DingTalk](https://www.dingtalk.com/zh) (User Group) | Group ID: 31912621 Or scan the QR code below <br> <img src="https://gw.alipayobjects.com/mdn/rms_5891a1/afts/img/A*--KAT7yyxXoAAAAAAAAAAAAAAAARQnAQ" height="200px"> <br> |
| 💬 [DingTalk](https://www.dingtalk.com/zh) (Community Meeting Group) | Group ID: 41585216 <br> [Layotto has a community meeting every Friday at 8 p.m., everyone is welcome](en/community/meeting.md) |

[comment]: <> (| 💬 [WeChat](https://www.wechat.com/) | Scan the QR code below to add friends, and she will invite you to join the WeChat group <br> <img src="../img/wechat-group.jpg" height="200px">)

## How to contribute

[Novice attack.]Slightly: Become a Layotto contributor from scratch](en/development/start-from-zero.md)

[Where to start?] Take a look at the "Starter Tasks" list] (https://github.com/mosn/layotto/issues/108#issuecomment-872779356)

As a technology student, have you ever felt like you want to get involved in the development of an open source project, but don't know where to start?
In order to help you better participate in open source projects, the community will regularly release novice development tasks suitable for novices to help you learn by doing!

[Documentation Contribution Guide] (zh/development/contributing-doc.md)

[Component Development Guide] (zh/development/developing-component.md)

[Layotto Github Workflow Guide] (zh/development/github-workflows.md)

[Layotto Command Line Guide] (zh/development/commands.md)

[Layotto Contributor Guide] (zh/development/CONTRIBUTING.md)

## Contributors

Thanks to all the contributors!

<a href="https://github.com/mosn/layotto/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=mosn/layotto" />
</a>

## Design Documentation

[Actuator Design Documentation] (zh/design/actuator/actuator-design-doc.md)

[Pubsub API Compatibility with Dapr Component] (zh/design/pubsub/pubsub-api-and-compability-with-dapr-component.md)

[Configuration API with Apollo] (zh/design/configuration/configuration-api-with-apollo.md)

[RPC Design Document] (en/design/rpc/rpc设计document.md)

[Distributed Lock API Design Document] (zh/design/lock/lock-api-design.md)

[FaaS Design Document] (zh/design/faas/faas-poc-design.md)

## FAQ

### What's the difference with DAPR?

dapr is an excellent runtime product, but it lacks the capabilities of Service Mesh, which is essential for actual implementation in the production environment, so we hope to combine the two capabilities of Runtime and Service Mesh to meet more complex production implementation needs.