# API 插件: 注册您自己的API
这是一个向您展示如何注册您自己的API的演示程序。

Layotto现有api-plugin的功能让你根据您的需要添加您自己的API

## step 1. 使用一个新的helloworld API运行Layotto
切换目录：

```shell
# 切换目录 
cd ${project_path}/cmd/layotto_multiple_api
```

编译Layotto：

```shell @if.not.exist layotto
# 编译命令
go build -o layotto
```

运行Layotto：

```shell @background
./layotto start -c ../../configs/config_standalone.json
```

Q: 这其中发生了什么？

检查[`main.go`](https://github.com/mosn/layotto/blob/d74ff0e8940e0eb9c73b1d3275a17d29be36bd5c/cmd/layotto_multiple_api/main.go#L203) 中的代码，然后你会发现Layotto在启动期间注册了一个新的API：

```go
		// 在这里注册您的grpc API
        runtime.WithGrpcAPI(
            // 默认的grpc API
            default_api.NewGrpcAPI,
            // 一个展示如何注册您自己的API的示例
            helloworld_api.NewHelloWorldAPI,
        ),
```

## step 2. 调用这个helloworld API

```shell
# 切换目录 
cd ${project_path}/cmd/layotto_multiple_api
# 运行客户端示例
go run client/main.go
```

这个结果将会是：

```bash
Greeting: Hello world
```

这个消息是您在步骤1中刚刚注册的helloworld API的响应结果。

## 下一步

您可以参考演示的代码来实现你自己的API。快来试试吧！

想要了解更多的详情，您可以参考[设计文档](zh/design/api_plugin/design.md)

为了简化 API 插件的开发，Layotto 社区提供了一套代码生成器，可以基于 proto 文件生成 API 插件相关代码，见 [文档](zh/start/api_plugin/generate.md)