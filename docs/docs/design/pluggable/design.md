# Pluggable Component 设计文档

## 背景

当前 Layotto 的 component 都是实现在 Layotto 的工程里面的，这要求用户要想使用新的 component，必须使用 golang 语言开发，同时必须在 Layotto 工程中实现，然后统一编译。
对于多语言用户来说非常不友好，因此 Layotto 需要提供pluggable components 的能力，允许用户可以通过任何语言实现自己的component，Layotto 通过 grpc 协议和外部的 component 进行通信。

## 方案

- 基于 uds（unix domain socket）实现本地跨语言组件服务发现，降低通信开销。
- 基于 proto 实现组件跨语言实现能力。

## 数据流架构

![](/img/pluggable/layotto_datatflow.png)

这是当前用户调用 sdk 开始的数据流向。虚线部分是与 pluggable component 主要参与的数据流。

### 组件发现

![](/img/pluggable/layotto.png)

如上图所示，用户自定义组件启动 socket 服务，并将 socket 文件放到指定目录中。 layotto 启动时，会读取该目录中的所有 socket 文件（跳过文件夹），并建立 socket 连接。

目前，layotto 向 dapr 对齐，不负责用户组件的生命周期，服务期间若用户组件下线，不会进行重连，该组件服务无法使用。
后面根据社区使用情况，决定 layotto 是否需要支持进程管理模块，或是使用一个单独的服务来管理。

由于 windows 对于 uds 的支持还不是很完善，且 layotto 本身取消了对 windows 的兼容，所以新特性采用的 uds 发现模式未对 windows 系统做兼容。

## 组件注册

如上面的数据流架构图所示，用户注册的组件需要实现 pluggable proto 定义的 grpc 服务。 layotto 会根据 grpc 接口，实现 go interface 接口，这里
对应于数据流图中的 wrap component。wrap component 与 build-in component 对 layotto runtime 来说没有任何区别，对于用户来说也没有特殊的感知。
 
layotto 通过 grpc reflect 库，获取到用户提供服务实现了哪些组件，注册到全局的组件注册中心，供用户使用。