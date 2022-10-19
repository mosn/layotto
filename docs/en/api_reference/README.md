# gRPC API reference

Layotto has multiple gRPC proto files, and the corresponding API reference are at:

[https://mosn.io/layotto/api/v1/runtime.html](https://mosn.io/layotto/api/v1/runtime.html)

These protos define Layotto's runtime API, including:

- API provided by Layotto for App
- The callback API that needs to be implemented by App. Layotto will call back the App and get the pubsub subscription message

In addition to this, Layotto also provides some extension APIs, including:










cryption: [spec/proto/extension/v1/cryption](https://mosn.io/layotto/api/v1/cryption.html) 

delay_queue: [spec/proto/extension/v1/delay_queue](https://mosn.io/layotto/api/v1/delay_queue.html) 

email: [spec/proto/extension/v1/email](https://mosn.io/layotto/api/v1/email.html) 

phone: [spec/proto/extension/v1/phone](https://mosn.io/layotto/api/v1/phone.html) 

s3: [spec/proto/extension/v1/s3](https://mosn.io/layotto/api/v1/s3.html) 

sms: [spec/proto/extension/v1/sms](https://mosn.io/layotto/api/v1/sms.html) 
