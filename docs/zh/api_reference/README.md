# gRPC API 接口文档

Layotto 有多个 gRPC proto 文件, 对应的接口文档在：

[https://mosn.io/layotto/api/v1/runtime.html](https://mosn.io/layotto/api/v1/runtime.html)

这些 proto 里定义了 Layotto 的运行时 API, 包括：

  - Layotto 对 App 提供的 API
  - 需要由 App 来实现的 callback API。 Layotto 会回调 App、获取 pubsub 订阅消息   

除此之外，Layotto 还提供了一些扩展 API，包括:










cryption: [spec/proto/extension/v1/cryption](https://mosn.io/layotto/api/v1/cryption.html) 

delay_queue: [spec/proto/extension/v1/delay_queue](https://mosn.io/layotto/api/v1/delay_queue.html) 

email: [spec/proto/extension/v1/email](https://mosn.io/layotto/api/v1/email.html) 

phone: [spec/proto/extension/v1/phone](https://mosn.io/layotto/api/v1/phone.html) 

s3: [spec/proto/extension/v1/s3](https://mosn.io/layotto/api/v1/s3.html) 

sms: [spec/proto/extension/v1/sms](https://mosn.io/layotto/api/v1/sms.html) 
