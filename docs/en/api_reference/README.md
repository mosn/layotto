# gRPC API reference

Layotto has multiple gRPC proto files, and the corresponding API reference are at:

[https://mosn.io/layotto/api/v1/runtime.html](https://mosn.io/layotto/api/v1/runtime.html)

These protos define Layotto's runtime API, including:

- API provided by Layotto for App
- The callback API that needs to be implemented by App. Layotto will call back the App and get the pubsub subscription message

In addition to this, Layotto also provides some extension APIs, including:

s3: [spec/proto/extension/v1/s3](https://mosn.io/layotto/api/v1/s3.html) 
