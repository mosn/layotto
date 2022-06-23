# Sequencer API
## What is Sequencer API
The Sequencer API is used to generate distributed unique, self-incrementing IDs.

The Sequencer API supports the declaration of demand for self-increment, including trend increase (WEAK) and strictly global increment (STRONG)
## When to use Sequencer API
### When you need to generate a globally unique id
Q: When do I need to generate a globally unique id?

A: When db does not automatically generate it for you. for example:

- After you do mysql sharding yourself, mysql will never automatically generate a globally unique id for you again, but you do need a globally unique business id (for example, "order id")
- The request does not go to the db, for example, a traceId is generated when the request reaches the backend servers.

### And you want the generated id increases automatically
Specifically, there are many types of requirement:

- No increment is required. UUID fits this situation, although the disadvantage is that it is relatively long. **If this is the case, it is recommended to use UUID to solve it by yourself, no need to call this API**
  
- "The trend is increasing". It means you do not need a strictly global increase, just want "in most cases it is increasing". In this case, it is recommended to use this API

Q: What scenarios will I need an increasing trend?

1. For b+ tree type db (such as MYSQL), the primary key with increasing trend can make better use of cache (cache friendly).

2. When you want to use the id to sort and query the latest data. For example, the requirement is to check the latest 100 messages, and the developer does not want to add a timestamp field and build an index on it. If the id itself is incremented, then the latest 100 messages can be sorted by id directly:

```
select * from message order by message-id limit 100
```

This is very common when using nosql, because it is difficult for nosql to index on another timestamp field

- Global monotonically increasing

When you want the generated id incremental without any regression, it is recommended to use this API

## How to use Sequencer API
You can call the Sequencer API through grpc. The API is defined in [runtime.proto](https://github.com/mosn/layotto/blob/main/spec/proto/runtime/v1/runtime.proto).

Layotto client sdk encapsulates the logic of grpc calling. For an example of using sdk to call Sequencer API, please refer to [Quick Start: Use Sequencer API](en/start/sequencer/start.md)

The components need to be configured before use. For detailed configuration options, see [Sequencer component document](en/component_specs/sequencer/common.md)
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

To avoid inconsistencies between the documentation and the code, please refer to [proto file](https://github.com/mosn/layotto/blob/main/spec/proto/runtime/v1/runtime.proto) for detailed input parameters and return values
