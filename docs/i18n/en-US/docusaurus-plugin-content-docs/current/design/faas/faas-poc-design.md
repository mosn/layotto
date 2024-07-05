# FaaS Design Document

## Structure design

![img.png](/img/faas/faas-design.jpg)

In the package FaaS, the following two issues were resolved mainly by the following question：

1. What is the relationship between compiled waste files and mirrors?
   1. The target waste file is eventually built into a normal mirror, pushed to Dockerhub. The pulse pulls the image as well as the original image. However, when running, the target waste file will be extracted from the image to load it separately.
2. How can k8s be managed to deploy wasm documents?
   1. Based on k8s' life-cycle management and movement strategy, the Containerd-shim-layotto-v2 plugins are customized in connection with Containerd's v2 interfaces and converted to Layotto Runtime when the container is running, for example, the k8s's creation of a container is turned into a function that loads the waste shape and runs.
   2. A sandbox isolation environment with a good Web Assembly and Layotto as a base can load functions that run multiple wasm forms, all of which are running in a process but not influential, and this nanoprocity thinks more than a docker can make full use of resources.

### Core components

#### A, [WebAssembly (wasm)](https://webassembly.org/)

Wasm1,wasm2 in the corresponding architectural charts, which codify and run the developed function as `*.wasm` as the form of function existing, using the sandbox isolation provided by WebAssembly technology for purposes that do not affect each function.

#### B,[Layotto](https://github.com/mosn/layotto)

Positions are designed to provide services, resources and security for functions.As a base from which functions will operate, provide access to the entrance to the infrastructure including WebAssembly running, functions can use maximum resource limits, functions for system call permission validation, etc.

#### C,[Containerd](https://containerd.io/)

When officially supported containers are running, docker is one of the most scenic implementations currently used, and security containers such as kata and gvisor are also using the technology and Layotto builds on their thinking and integrates the function loading process into concrete implementation when the container is running.

#### D,[Containerd-shim-layotto-v2](https://github.com/layotto/containerd-wasm)

Based on Containerd's V2 interface definition, the logic of the running of the container is customized, such as creating the container to perform the operation to allow Layotto load and run the wasm function.

#### E,[Kubernetes](https://kubernetes.io/)

The factual standard of the current container schedule, life-cycle management and movement strategy are excellent, and the containerd-based solution is designed to combine function movement with k8s ecologically perfect.

### Runtime ABI

#### A. [proxy-wasm-go-sdk](https://github.com/layotto/proxy-wasm-go-sdk)

The interface of function access to system resources and infrastructure services is defined and implemented on a community-based basis [proxy-waste/spec] (https://github.com/proxy-wasm/spec) that brings together the [Runtime API](https://github.com/mosn/layotto/blob/main/spec/proto/runtime/v1/runtime.proto) and adds ABI to infrastructure visits.

#### B. [proxy-wasm-go-host](https://github.com/layotto/proxy-wasm-go-host)

Concrete logic for Runtime ABI in Layotto.
