# Sequencer API
## 什么是Sequencer API
Sequencer API用于生成分布式唯一、自增id。

Sequencer API支持声明对自增的需求，包括趋势递增(WEAK)和严格递增(STRONG)

## 何时使用Sequencer API
### 需要生成全局唯一id时
Q: 什么时候需要生成全局唯一id?

A: db不帮你自动生成的时候。比如：
- db做了分库分表，没帮你自动生成id，你又需要一个全局唯一的业务id(例如"订单id")
- 请求没走db，比如请求到了后端要生成一个traceId

### 希望生成的全局唯一id是递增的
具体来说有很多种：
- 不需要递增。这种情况UUID能解决，虽然缺点是比较长。**如果是这种情况建议自行用UUID解决，不需要调用本API**
- “趋势递增”。不追求一定递增，"大部分情况是递增的"就行。这种情况建议使用本API

Q: 什么场景需要趋势递增？

1. 对b+树类的db(例如MYSQL)来说,趋势递增的主键能更好的利用缓存（cache friendly）。

2. 拿来排序查最新数据。比如需求是查最新的100条消息，开发者不想新增个时间戳字段、建索引，如果id本身是递增的，那么查最新的100条消息时直接按id排序即可：

```
select * from message order by message-id limit 100
```

这在使用nosql的时候很常见，因为nosql在时间戳字段上加索引很难

- 全局单调递增

希望生成的id一定递增，没有任何倒退的情况。这种情况建议使用本API

## 如何使用Sequencer API
您可以通过grpc调用Sequencer API，接口定义在[runtime.proto](https://github.com/mosn/layotto/blob/main/spec/proto/runtime/v1/runtime.proto) 中。

Layotto client sdk封装了grpc调用的逻辑，使用sdk调用Sequencer API的示例可以参考[快速开始：使用Sequencer API](zh/start/sequencer/start.md)

使用前需要先对组件进行配置，详细的配置说明见[Sequencer组件文档](zh/component_specs/sequencer/common.md)

### Get next unique id

```protobuf
// Sequencer API
// Get next unique id with some auto-increment guarantee
rpc GetNextId(GetNextIdRequest)returns (GetNextIdResponse) {}
  
message GetNextIdRequest {
  // Required. Name of sequencer storage
  string store_name = 1;
  // Required. key is the identifier of a sequencer namespace,e.g. "order_table".
  string key = 2;
  // (optional) SequencerOptions configures requirements for auto-increment guarantee
  SequencerOptions options = 3;
  // (optional) The metadata which will be sent to the component.
  map<string, string> metadata = 4;
}

// SequencerOptions configures requirements for auto-increment guarantee
message SequencerOptions {
  // requirements for auto-increment guarantee
  enum AutoIncrement {
    // (default) WEAK means a "best effort" incrementing service.But there is no strict guarantee of global monotonically increasing.
    //The next id is "probably" greater than current id.
    WEAK = 0;
    // STRONG means a strict guarantee of global monotonically increasing.
    //The next id "must" be greater than current id.
    STRONG = 1;
  }

  AutoIncrement increment = 1;
}
  
message GetNextIdResponse{
  // The next unique id
  int64 next_id = 1;
}
```

为避免文档和代码不一致，详细入参和返回值请参考[proto文件](https://github.com/mosn/layotto/blob/main/spec/proto/runtime/v1/runtime.proto)
