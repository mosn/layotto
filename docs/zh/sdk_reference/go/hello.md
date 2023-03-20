# Hello

用于与 Layotto 进行 grpc 通信，测试建连，作为入门的 demo 接口。类似 redis 中的 `ping-pong` 模式。

## SayHello

```go
SayHello(ctx context.Context, in *SayHelloRequest) (*SayHelloResp, error)

type SayHelloRequest struct {
	ServiceName string
}

type SayHelloResp struct {
	Hello string
}
```

```go
res, err := cli.SayHello(context.Background(), &client.SayHelloRequest{
    ServiceName: "helloworld",
})

// 返回的 res.Hello 为 "greeting"
```

> ServiceName 需要与配置文件中 `hellos` 组件下设置的 instance 名称一致。 具体可以参考 [start](zh/sdk_reference/go/start.md)

