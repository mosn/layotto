# gRPC API 接口文档

Layotto 有多个 gRPC proto 文件, 对应的接口文档在：

[https://mosn.io/layotto/api/v1/runtime.html](https://mosn.io/layotto/api/v1/runtime.html)

这些 proto 里定义了 Layotto 的运行时 API, 包括：

  - Layotto 对 App 提供的 API
  - 需要由 App 来实现的 callback API。 Layotto 会回调 App、获取 pubsub 订阅消息   

除此之外，Layotto 还提供了一些扩展 API，包括:




s3: [spec/proto/extension/v1/s3](https://mosn.io/layotto/api/v1/s3.html) 

