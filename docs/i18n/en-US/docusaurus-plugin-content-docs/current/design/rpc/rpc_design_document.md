# RPC Design Document

## API Design

[layotto rpc API](https://github.com/mosn/layotto/blob/f70cdc6196963ad762cf809daf0579403c341def1/spec/proto/runtime/v1/runtime.proto) is consistent with Dapr.

### Core Abstract

In order to decouple with pb definition, a layer of RFC core abstraction has been added.

- Invoker： provides full RPC capability, currently only Mosn
- callback：before/after filter, can execute custom logic before and after request execution (e.g. add request head, such as protocol conversion)
- channel：send requests, receive responses, and interact with different transmission protocols

Since Mosn already has full RPC capacity support, layotto provides only a very light RPC framework

![img.png](/img/rpc/rpc-layer.png)

### Mosn Integration

The layotto RPC, based on Mosn grpc handler, works on 7 floors, while Mosn's proxy, as well as various filters, work on 4 floors and cannot be interacted by simple functional calls.

For **Full Replicate**Mosn, layotto uses new ideas for Mosn integration.

1. Channel will re-encode the request from L7 to L4
2. Create a virtual connection (net.Pipe), layotto holds one end local, mosn holding remote at the other end
3. layotto write to local, mosn will receive data
4. mosn reads from remote, executes filter and redisseminates it to remote
5. layotto read from remote, get responses

#### xprotocol

Mosn supports the popular RPC protocol via xprotocol.
designed a corresponding extension in Layotto. Simply complete the interface between RPC requests and xprotocol frame, it will be easy to support the xprotocol protocol.

#### Configure Parameters

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
