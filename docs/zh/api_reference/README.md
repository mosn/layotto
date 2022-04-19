# gRPC API 接口文档

Layotto 有两个 gRPC proto 文件, 对应的接口文档在：

- [spec/proto/runtime/v1/runtime.proto](https://github.com/mosn/layotto/blob/main/docs/en/api_reference/runtime_v1.md)

该 proto 定义的 gRPC API, 就是 Layotto 对 App 提供的 API。
  
- [spec/proto/runtime/v1/appcallback.proto](https://github.com/mosn/layotto/blob/main/docs/en/api_reference/appcallback_v1.md)

该接口需要由 App 来实现，用来处理 pubsub 订阅消息   
