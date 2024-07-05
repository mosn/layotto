# State API
## 什么是State API
State API是一套对Key/Value数据进行增删改查的API。您的应用程序可以使用同一套State API操作不同的数据库（或某种存储系统），对Key/Value模型的数据进行增删改查。

API支持批量CRUD操作，支持声明对并发安全和数据一致性的要求，由Layotto帮您处理复杂的并发安全和数据一致性问题。

## 何时使用State API
如果您的应用需要访问Key/Value存储、进行增删改查，那么使用State API是一个不错的选择，它有以下好处：

- 多（云）环境部署：同一套业务代码部署在不同环境

中立的API可以帮助您的应用和存储供应商、云厂商解耦，能够不改代码部署在不同的云上。

- 多语言复用中间件：同一个DB（和数据中间件）能支持不同语言的应用

如果您的公司内部有不同语言开发的应用（例如同时有java和python应用），那么传统做法是为每种语言开发一套数据中间件sdk（用于路由，容灾，流量管理等目的）。

使用State API可以帮助您免去维护多语言sdk的烦恼，不同语言的应用可以用同一套grpc API和Layotto交互。

## 如何使用State API
您可以通过grpc调用State API，接口定义在[runtime.proto](https://github.com/mosn/layotto/blob/main/spec/proto/runtime/v1/runtime.proto) 中。

使用前需要先对组件进行配置，详细的配置说明见[状态管理组件文档](zh/component_specs/state/common.md)

### 使用示例
Layotto client sdk封装了grpc调用的逻辑，使用sdk调用State API的示例可以参考[快速开始：使用State API](zh/start/state/start.md)


### Save state
用于保存一批状态数据

```protobuf
  // Saves an array of state objects
  rpc SaveState(SaveStateRequest) returns (google.protobuf.Empty) {}
```

#### 入参

```protobuf

// SaveStateRequest is the message to save multiple states into state store.
message SaveStateRequest {
  // Required. The name of state store.
  string store_name = 1;

  // Required. The array of the state key values.
  repeated StateItem states = 2;
}

// StateItem represents state key, value, and additional options to save state.
message StateItem {
  // Required. The state key
  string key = 1;

  // Required. The state data for key
  bytes value = 2;

  // (optional) The entity tag which represents the specific version of data.
  // The exact ETag format is defined by the corresponding data store. Layotto runtime only treats ETags as opaque strings. 
  Etag etag = 3;

  // (optional) additional key-value pairs to be passed to the state store.
  map<string, string> metadata = 4;

  // (optional) Options for concurrency and consistency to save the state.
  StateOptions options = 5;
}

// Etag represents a state item version
message Etag {
  // value sets the etag value
  string value = 1;
}


// StateOptions configures concurrency and consistency for state operations
message StateOptions {
  // Enum describing the supported concurrency for state.
  // The API server uses Optimized Concurrency Control (OCC) with ETags.
  // When an ETag is associated with an save or delete request, the store shall allow the update only if the attached ETag matches with the latest ETag in the database.
  // But when ETag is missing in the write requests, the state store shall handle the requests in the specified strategy(e.g. a last-write-wins fashion).
  enum StateConcurrency {
    CONCURRENCY_UNSPECIFIED = 0;
    // First write wins
    CONCURRENCY_FIRST_WRITE = 1;
    // Last write wins
    CONCURRENCY_LAST_WRITE = 2;
  }

  // Enum describing the supported consistency for state.
  enum StateConsistency {
    CONSISTENCY_UNSPECIFIED = 0;
    //  The API server assumes data stores are eventually consistent by default.A state store should:
    //
    // - For read requests, the state store can return data from any of the replicas
    // - For write request, the state store should asynchronously replicate updates to configured quorum after acknowledging the update request.
    CONSISTENCY_EVENTUAL = 1;

    // When a strong consistency hint is attached, a state store should:
    //
    // - For read requests, the state store should return the most up-to-date data consistently across replicas.
    // - For write/delete requests, the state store should synchronisely replicate updated data to configured quorum before completing the write request.
    CONSISTENCY_STRONG = 2;
  }

  StateConcurrency concurrency = 1;
  StateConsistency consistency = 2;
}
```

#### 返回

`google.protobuf.Empty`

### Get State

```protobuf
  // Gets the state for a specific key.
  rpc GetState(GetStateRequest) returns (GetStateResponse) {}
```

为避免文档和代码不一致，详细入参和返回值请参考[proto文件](https://github.com/mosn/layotto/blob/main/spec/proto/runtime/v1/runtime.proto)

### Get bulk state

```protobuf
  // Gets a bulk of state items for a list of keys
  rpc GetBulkState(GetBulkStateRequest) returns (GetBulkStateResponse) {}
```

为避免文档和代码不一致，详细入参和返回值请参考[proto文件](https://github.com/mosn/layotto/blob/main/spec/proto/runtime/v1/runtime.proto)

### Delete state

```protobuf
  // Deletes the state for a specific key.
  rpc DeleteState(DeleteStateRequest) returns (google.protobuf.Empty) {}
```

为避免文档和代码不一致，详细入参和返回值请参考[proto文件](https://github.com/mosn/layotto/blob/main/spec/proto/runtime/v1/runtime.proto)

### Delete bulk state

```protobuf
  // Deletes a bulk of state items for a list of keys
  rpc DeleteBulkState(DeleteBulkStateRequest) returns (google.protobuf.Empty) {}
```

为避免文档和代码不一致，详细入参和返回值请参考[proto文件](https://github.com/mosn/layotto/blob/main/spec/proto/runtime/v1/runtime.proto)

### State transactions

```protobuf
  // Executes transactions for a specified store
  rpc ExecuteStateTransaction(ExecuteStateTransactionRequest) returns (google.protobuf.Empty) {}
```

为避免文档和代码不一致，详细入参和返回值请参考[proto文件](https://github.com/mosn/layotto/blob/main/spec/proto/runtime/v1/runtime.proto)
