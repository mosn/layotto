# RPC API
## 什么是RPC API
Layotto的RPC API基于[Mosn](https://mosn.io/docs/overview/)的grpc handler设计，提供多协议、模块化、智能、安全的**服务调用**。

RPC API的接口与[Dapr](https://docs.dapr.io/zh-hans/developing-applications/building-blocks/service-invocation/service-invocation-overview/)一致，可以在代码[invoke.go](https://github.com/mosn/layotto/blob/3802c4591181fdbcfb7dd07cbbdbadeaaada650a/sdk/go-sdk/client/invoke.go)中看到接口的具体设计细节。

使用 Layotto RPC API 进行服务调用，您的应用程序可以使用标准 HTTP 或 [X-Protocol](https://cloudnative.to/blog/x-protocol-common-address-solution/) 协议可靠且安全地与其他应用程序通信.

![sidecar](https://mosn.io/en/docs/concept/sidecar-pattern/sidecar-pattern.jpg)

Layotto 使用 [sidecar](https://mosn.io/docs/concept/sidecar-pattern/) 架构。使用 Layotto ，您可以在任何 Layotto 实例上使用RPC API调用应用程序。 Sidecar 编程模型鼓励每个应用程序与其自己的 Layotto 实例进行对话。 Layotto 实例发现并相互通信。


## 何时使用 RPC API，有什么好处？
在许多具有多个服务需要相互通信的环境中，开发人员经常会问自己以下问题：

- 如何发现和调用不同服务的方法？
- 如何通过加密安全地调用其他服务并对方法应用访问控制？
- 如何处理重试和瞬时错误？
- 如何使用跟踪查看带有指标的调用图以诊断生产中的问题？

Layotto 通过提供服务调用 API 来解决这些挑战，该 API 充当反向代理与内置服务发现的组合，同时利用内置的分布式跟踪、度量、错误处理、加密等。

## 如何使用RPC API？
您可以通过 grpc 接口 **InvokeService** 进行 RPC 调用。 API 在 [api.go](https://github.com/mosn/layotto/blob/77e0a4b2af063ff9e365a933c4735655898de369/pkg/grpc/api.go) 中定义。

该组件需要在使用前进行配置。详细的配置说明见[配置参考](https://mosn.io/layotto/#/zh/configuration/overview)，目前Layotto使用的是集成MOSN的MOSN 4层过滤器，运行在MOSN上，所以配置文件Layotto使用的其实是一个MOSN配置文件。所以也可以参考[MOSN配置文件](https://mosn.io/docs/configuration/)的文档。

### 演示 1：Hello World -- 基本 Golang HTTP 服务器
此演示的快速入门文档：[Hello World](https://mosn.io/layotto/#/zh/start/rpc/helloworld)

[echoserver](https://github.com/mosn/layotto/blob/77e0a4b2af063ff9e365a933c4735655898de369/demo/rpc/http/echoserver/echoserver.go)在8889端口发布一个简单的接口，配置文件[example.json](https://github.com/mosn/layotto/blob/77e0a4b2af063ff9e365a933c4735655898de369/demo/rpc/http/example.json)利用mosn的路由能力转发http头中id字段等于**HelloService:1.0**的请求 到本地8889端口，然后在[echoclient](https://github.com/mosn/layotto/blob/b66b998f50901f8bd1cce035478579c1b47f986d/demo/rpc/http/echoclient/echoclient.go)使用接口 **InvokeService** 进行 RPC 调用。

```golang
resp, err := cli.InvokeService(
		ctx,
		&runtimev1pb.InvokeServiceRequest{
			Id: "HelloService:1.0",
			Message: &runtimev1pb.CommonInvokeRequest{
				Method:      "/hello",
				ContentType: "",
				Data:        &anypb.Any{Value: []byte(*data)}}},
	)
```

### Demo 2：Dubbo JSON RPC
本demo快速入门文档：[Dubbo JSON RPC Example](https://mosn.io/layotto/#/zh/start/rpc/dubbo_json_rpc)

服务端由dubbo示例程序[dubbo-go-samples](https://github.com/apache/dubbo-go-samples)充当，配置文件 [example.json](https://github.com/mosn/layotto/blob/77e0a4b2af063ff9e365a933c4735655898de369/demo/rpc/dubbo_json_rpc/example.json) 使用插件[dubbo_json_rpc](https://github.com/mosn/layotto/blob/8db7a2297bd05d1b0c4452cc980d8f6412a82f3a/components/rpc/callback/dubbo_json_rpc.go)，以产生请求头。随后请求端[client](https://github.com/mosn/layotto/blob/b66b998f50901f8bd1cce035478579c1b47f986d/demo/rpc/dubbo_json_rpc/dubbo_json_client/client.go) 使用接口 **InvokeService** 进行 RPC 调用。

```golang
resp, err := cli.InvokeService(
		ctx,
		&runtimev1pb.InvokeServiceRequest{
			Id: "com.ikurento.user.UserProvider",
			Message: &runtimev1pb.CommonInvokeRequest{
				Method:        "GetUser",
				ContentType:   "",
				Data:          &anypb.Any{Value: []byte(*data)},
				HttpExtension: &runtimev1pb.HTTPExtension{Verb: runtimev1pb.HTTPExtension_POST},
			},
		},
	)
```


## 实现原理
如果对实现原理感兴趣，或者想扩展一些功能，可以阅读[RPC设计文档](https://mosn.io/layotto/#/zh/design/rpc/rpc%E8%AE%BE%E8%AE%A1%E6%96%87%E6%A1%A3)。
