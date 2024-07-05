RPC 设计文档

### API 设计
[layotto rpc API](https://github.com/mosn/layotto/blob/f70cdc619693ad762cf809daf0579403c341def1/spec/proto/runtime/v1/runtime.proto#L19https://github.com/mosn/layotto/blob/f70cdc619693ad762cf809daf0579403c341def1/spec/proto/runtime/v1/runtime.proto#L19) 与Dapr保持一致.

### 核心抽象
为了与pb定义解耦，添加了一层RPC核心抽象.

- invoker： 提供完整的 RPC能力， 目前只对接了Mosn
- callback：before/after filter, 可以在请求执行前后执行自定义的逻辑（例如添加请求头，例如协议转换)
- channel：发送请求，接收响应，负责与不同传输协议交互

由于Mosn已经有了完整的RPC能力支持，layotto只提供了非常轻量的RPC框架

![img.png](../../../img/rpc/rpc-layer.png)

### Mosn集成

layotto的RPC是基于Mosn grpc handler的，工作在7层，而Mosn的代理能力，以及各种filter都是工作在4层的, 无法通过简单的函数调用来交互.

为了**完整复用**Mosn的全套能力，layotto使用了新的思路与Mosn集成.

1. channel会将请求重新编码，从L7重新回到L4
2. 创建一对虚拟连接(net.Pipe)，layotto持有一端local，mosn持有另一端remote
3. layotto向local写入，mosn会收到数据
4. mosn从remote读取，执行filter并进行代理转发，将响应写到remote
5. layotto从remote读取，获得响应


#### xprotocol
Mosn通过xprotocol支持了流行的RPC协议.
在Layotto里设计了对应的扩展机制，只需要完成RPC请求响应与xprotocol frame的互相转换，就可以方便的支持xprotocl协议.

#### 配置参数

```bigquery
{
  "mosn": {
    "config": {
      "before_invoke": [{
        "name": "xxx" // rpc调用前的filter
      }],
      "after_invoke": [{
        "name": "xxx" // rpc调用后的filter
      }],
      "channel": {
        "size": 1, // 与mosn通信使用的通道数量，可以简单理解成连接数
        "protocol": "http", // 与mosn通信使用的协议
        "listener": "egress_runtime_http" // mosn对应的listener端口
      }
    }
  }
}
```