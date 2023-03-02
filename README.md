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

[查看中文版本](https://mosn.io/layotto/#/zh/README)

Layotto(/leɪˈɒtəʊ/) is an application runtime developed using Golang, which provides various distributed capabilities for applications, such as state management, configuration management, and event pub/sub capabilities to simplify application development.

Layotto is built on the open source data plane [MOSN](https://github.com/mosn/mosn) .In addition to providing distributed building blocks, Layotto can also serve as the data plane of Service Mesh and has the ability to control traffic.

## Motivation

Layotto aims to combine [Multi-Runtime](https://www.infoq.com/articles/multi-runtime-microservice-architecture/) with Service Mesh into one sidecar. No matter which product you are using as the Service Mesh data plane (e.g. MOSN,Envoy or any other product), you can always attach Layotto to it and add Multi-Runtime capabilities without adding new sidecars.

For examples, by adding Runtime capabilities to MOSN, a Layotto process can both [serve as the data plane of istio](https://mosn.io/layotto/#/en/start/istio/) and provide various Runtime APIs (such as Configuration API, Pub/Sub API, etc.)

In addition, we were surprised to find that a sidecar can do much more than that. We are trying to make Layotto even the runtime container of FaaS (Function as a service) with the magic power of [WebAssembly](https://en.wikipedia.org/wiki/WebAssembly) .

## Features

- Service Communication
- Service Governance.Such as traffic hijacking and observation, service rate limiting, etc
- [As the data plane of istio](https://mosn.io/layotto/#/en/start/istio/)
- Configuration management
- State management
- Event publish and subscribe
- Health check, query runtime metadata
- [FaaS model based on WASM and Runtime](https://mosn.io/layotto/#/en/start/faas/start.md)

## Project Architecture

As shown in the architecture diagram below, Layotto uses the open source MOSN as the base to provide network layer management capabilities while providing distributed capabilities. The business logic can directly interact with Layotto through a lightweight SDK without paying attention to the specific back-end infrastructure.

Layotto provides SDKs in various languages. The sdk interacts with Layotto through grpc. Application developers only need to specify their own infrastructure type through the configuration file [configure file](./configs/runtime_config.json) provided by Layotto. No coding changes are required, which greatly improves the portability of the program.

![Architecture](https://raw.githubusercontent.com/mosn/layotto/main/docs/img/runtime-architecture.png)

## Quickstarts

### Get started with Layotto

You can try the quickstart demos below to get started with Layotto. In addition, you can experience the [online laboratory](https://mosn.io/layotto/#/en/start/lab)

### API

| API            | status |                              quick start                              |                          desc                                                           |
| -------------- | :----: | :-------------------------------------------------------------------: | -------------------------------------------------------------- |
| State          |   ✅   |        [demo](https://mosn.io/layotto/#/en/start/state/start)         |     Write/Query the data of the Key/Value model                    |
| Pub/Sub        |   ✅   |        [demo](https://mosn.io/layotto/#/en/start/pubsub/start)        |     Publish/Subscribe message through various Message Queue        |
| Service Invoke |   ✅   |       [demo](https://mosn.io/layotto/#/en/start/rpc/helloworld)       |         Call Service through MOSN (another istio data plane)           |
| Config         |   ✅   | [demo](https://mosn.io/layotto/#/en/start/configuration/start-apollo) | Write/Query/Subscribe the config through various Config Center |
| Lock           |   ✅   |         [demo](https://mosn.io/layotto/#/en/start/lock/start)         |     Distributed lock API                                 |
| Sequencer      |   ✅   |      [demo](https://mosn.io/layotto/#/en/start/sequencer/start)       |   Generate distributed unique and incremental ID               |
| File           |   ✅   |                                 TODO                                  |      File API implementation                                        |
| Binding        |   ✅   |                                 TODO                                  |      Transparent data transmission API                              |


### Service Mesh

| feature | status |                      quick start                       | desc                       |
| ------- | :----: | :----------------------------------------------------: | -------------------------- |
| istio   |   ✅   | [demo](https://mosn.io/layotto/#/en/start/istio/) | As the data plane of istio |

### Extendability

| feature    | status |                           quick start                            | desc                       |
| ---------- | :----: | :--------------------------------------------------------------: | -------------------------- |
| API plugin |   ✅   | [demo](https://mosn.io/layotto/#/en/start/api_plugin/helloworld) | You can add your own API ! |

### Actuator

| feature        | status |                        quick start                        | desc                                                |
| -------------- | :----: | :-------------------------------------------------------: | --------------------------------------------------- |
| Health Check   |   ✅   | [demo](https://mosn.io/layotto/#/en/start/actuator/start) | Query health state of app and components in Layotto |
| Metadata Query |   ✅   | [demo](https://mosn.io/layotto/#/en/start/actuator/start) | Query metadata in Layotto/app                       |

### Traffic Control

| feature      | status |                              quick start                              | desc                                                            |
| ------------ | :----: | :-------------------------------------------------------------------: | --------------------------------------------------------------- |
| TCP Copy     |   ✅   |   [demo](https://mosn.io/layotto/#/en/start/network_filter/tcpcopy)   | Dump the tcp traffic received by Layotto into local file system |
| Flow Control |   ✅   | [demo](https://mosn.io/layotto/#/en/start/stream_filter/flow_control) | limit access to the APIs provided by Layotto                    |

### Write your bussiness logic using WASM

| feature        | status |                      quick start                      | desc                                                                 |
| -------------- | :----: | :---------------------------------------------------: | -------------------------------------------------------------------- |
| Go (TinyGo)    |   ✅   | [demo](https://mosn.io/layotto/#/en/start/wasm/start) | Compile Code written by TinyGo to \*.wasm and run in Layotto         |
| Rust           |   ✅   | [demo](https://mosn.io/layotto/#/en/start/wasm/start) | Compile Code written by Rust to \*.wasm and run in Layotto           |
| AssemblyScript |   ✅   | [demo](https://mosn.io/layotto/#/en/start/wasm/start) | Compile Code written by AssemblyScript to \*.wasm and run in Layotto |

### As a FaaS(Serverless) runtime (Layotto + WebAssembly + k8s)

| feature        | status |                      quick start                      | desc                                                                                       |
| -------------- | :----: | :---------------------------------------------------: | ------------------------------------------------------------------------------------------ |
| Go (TinyGo)    |   ✅   | [demo](https://mosn.io/layotto/#/en/start/faas/start) | Compile Code written by TinyGo to \*.wasm and run in Layotto And Scheduled by k8s.         |
| Rust           |   ✅   | [demo](https://mosn.io/layotto/#/en/start/faas/start) | Compile Code written by Rust to \*.wasm and run in Layotto And Scheduled by k8s.           |
| AssemblyScript |   ✅   | [demo](https://mosn.io/layotto/#/en/start/faas/start) | Compile Code written by AssemblyScript to \*.wasm and run in Layotto And Scheduled by k8s. |

## Presentations

- [Layotto - A new chapter of Service Mesh and Application Runtime](https://www.youtube.com/watch?v=5v8gTrFUDk8)
- [WebAssembly + Application Runtime = A New Era of FaaS?](https://www.youtube.com/watch?v=g01CJ4S9Qao)

## Landscapes

<p align="center">
<img src="https://landscape.cncf.io/images/left-logo.svg" width="150"/>&nbsp;&nbsp;<img src="https://landscape.cncf.io/images/right-logo.svg" width="200"/>
<br/><br/>
Layotto enriches the <a href="https://landscape.cncf.io/serverless">CNCF CLOUD NATIVE Landscape.</a>
</p>

## Community

### Contact Us

| Platform                                               | Link                                                                                                                                                                           |
| :----------------------------------------------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 💬 [DingTalk](https://www.dingtalk.com/en) (preferred) | Search the group number: 31912621 or scan the QR code below <br> <img src="https://gw.alipayobjects.com/mdn/rms_5891a1/afts/img/A*HUKsSrFthtUAAAAAAAAAAAAAARQnAQ" height="200px"> |

[comment]: <> (| 💬 [Wechat]&#40;https://www.wechat.com/en/&#41; | Scan the QR code below and she will invite you into the wechat group <br> <img src="/img/wechat-group.jpg" height="200px">)

## How to contribute

[Where to start? Check "Community tasks" list!](https://github.com/mosn/layotto/issues/108)

As a programming enthusiast , have you ever felt that you want to participate in the development of an open source project, but don't know where to start?
In order to help everyone better participate in open source projects, our community will regularly publish community tasks to help everyone learn by doing!

[Document Contribution Guide](docs/en/development/contributing-doc.md)

[Component Development Guide](docs/en/development/developing-component.md)

[Layotto Github Workflows](docs/en/development/github-workflows.md)

[Layotto Commands Guide](docs/en/development/commands.md)

[Layotto Contributor Guide](docs/en/development/CONTRIBUTING.md)

## Contributors

Thank y'all!

<a href="https://github.com/mosn/layotto/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=mosn/layotto" />
</a>

## Design Documents

[Actuator Design Doc](docs/en/design/actuator/actuator-design-doc.md)

[Configuration API with Apollo](docs/en/design/configuration/configuration-api-with-apollo.md)

[Pubsub API and Compability with Dapr Component](docs/en/design/pubsub/pubsub-api-and-compability-with-dapr-component.md)

[RPC Design Doc](docs/en/design/rpc/rpc-design-doc.md)

[Distributed Lock API Design](docs/en/design/lock/lock-api-design.md)

[FaaS Design](docs/en/design/faas/faas-poc-design.md)

## FAQ

### Difference with dapr?

dapr is an excellent Runtime product, but it lacks the ability of Service Mesh, which is necessary for the Runtime
product used in production environment, so we hope to combine Runtime with Service Mesh into one sidecar to meet
more complex production requirements.
