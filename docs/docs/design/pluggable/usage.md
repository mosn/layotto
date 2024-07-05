# Pluggable Component 使用文档

该示例展示了如何通过 Layotto 提供的可插拔组件能力，用户实现并注册自己的组件。并通过 Layotto sdk 调用，来验证自己组件编写的正确性。

## step1.编写并运行可插拔组件

接下来，运行已经编写好的代码

```shell
cd demo/pluggable/hello
go run .
```

打印如下结果表示服务启动成功

```shell
start grpc server
```

> 1. 以 go 实现 hello 组件为例，在 `layotto/spec/proto/pluggable` 中找到对应组件的 proto 文件，生成对应实现语言的 grpc 文件。
go 语言的 pb 文件已经生成并放在了 `spec/proto/pluggable/v1` 下，用户在使用时直接引用即可。
> 2. 组件除了需要实现 protobuf 文件中定义的接口外，还需要使用 socket 方式启动文件并将 sock 文件存放在 `/tmp/runtime/component-sockets` 默认路径下，
也可以通过环境变量 `LAYOTTO_COMPONENTS_SOCKETS_FOLDER` 修改 sock 存储路径位置。
> 3. 除此之外，用户还需要注册 reflection 服务到 grpc server 中，用于 layotto 服务发现时获取该 grpc 服务具体实现接口的 spec。 具体代码可以参考 `demo/pluggable/hello/main.go`

## step2. 启动 Layotto

```shell
cd cmd/layotto
go build -o layotto .
./layotto start -c ../../configs/config_hello_component.json
```

> 配置文件中填写组件的 type 为 `hello-grpc-demo`，由 socket 文件的前缀名决定。 配置项与注册普通 hello 组件一致。提供 metadata 项，便于用户设置自定义配置需求。

## step3. 组件校验

基于现有的组件测试代码，来测试用户实现的可插拔组件的正确性。

```shell
cd demo/hello/common
go run . -s helloworld
```

程序输出以下结果表示可插拔组件注册运行成功

```shell
runtime client initializing for: 127.0.0.1:34904
hello
```

## 了解 Layotto 可插拔组件的实现原理

如果您对实现原理感兴趣，或者想扩展一些功能，可以阅读[可插拔组件的设计文档](/docs/design/pluggable/design.md)
