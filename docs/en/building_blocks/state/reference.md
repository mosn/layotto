# State API
## What is State API
State API is a set of APIs for adding, deleting, modifying and querying Key/Value data. Your application can use the same State API to operate different databases (or a certain storage system) with Key/Value model.

API supports batch CRUD operations and supports the declaration of requirements for concurrency security and data consistency. Layotto will help you deal with complex concurrency control and data consistency issues.

## When to use State API and what are the benefits?
If your application needs to do some CRUD operations on Key/Value storage, then using the State API is a good choice. It has the following benefits:

- Multi (cloud) environment deployment: the same application code can be deployed in different environments

A neutral API can help your application decouple from storage vendors and cloud vendors, and be able to deploy on different clouds without changing the code.

- Multi-language reuse middleware: the same DB (and data middleware) can support applications in different languages

If your company has applications developed in different languages (for example, both java and python applications), then the traditional approach is to develop a set of data middleware SDKs for each language(used for routing,traffic control or some other custom purposes).

Using State API can help you avoid the trouble of maintaining multilingual SDKs. Applications in different languages can interact with Layotto using the same set of grpc API.

## How to use State API
You can call the State API through grpc. The API is defined in [runtime.proto](https://github.com/mosn/layotto/blob/main/spec/proto/runtime/v1/runtime.proto).

The component needs to be configured before use. For detailed configuration items, see [State Component Document](en/component_specs/state/common.md)

### Example
Layotto client sdk encapsulates the logic of grpc call. For examples of using sdk to call State API, please refer to [Quick Start: Use State API](en/start/state/start.md)


### Save state
Used to save a batch of status data

```protobuf
  // Saves an array of state objects
  rpc SaveState(SaveStateRequest) returns (google.protobuf.Empty) {}
```

#### parameters

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

#### return

`google.protobuf.Empty`

### Get State

```protobuf
  // Gets the state for a specific key.
  rpc GetState(GetStateRequest) returns (GetStateResponse) {}
```

To avoid inconsistencies between this document and the code, please refer to [the newest proto file](https://github.com/mosn/layotto/blob/main/spec/proto/runtime/v1/runtime.proto) for detailed input parameters and return values.

### Get bulk state

```protobuf
  // Gets a bulk of state items for a list of keys
  rpc GetBulkState(GetBulkStateRequest) returns (GetBulkStateResponse) {}
```

To avoid inconsistencies between this document and the code, please refer to [the newest proto file](https://github.com/mosn/layotto/blob/main/spec/proto/runtime/v1/runtime.proto) for detailed input parameters and return values.

### Delete state

```protobuf
  // Deletes the state for a specific key.
  rpc DeleteState(DeleteStateRequest) returns (google.protobuf.Empty) {}
```

To avoid inconsistencies between this document and the code, please refer to [the newest proto file](https://github.com/mosn/layotto/blob/main/spec/proto/runtime/v1/runtime.proto) for detailed input parameters and return values.

### Delete bulk state

```protobuf
  // Deletes a bulk of state items for a list of keys
  rpc DeleteBulkState(DeleteBulkStateRequest) returns (google.protobuf.Empty) {}
```

To avoid inconsistencies between this document and the code, please refer to [the newest proto file](https://github.com/mosn/layotto/blob/main/spec/proto/runtime/v1/runtime.proto) for detailed input parameters and return values.

### State transactions

```protobuf
  // Executes transactions for a specified store
  rpc ExecuteStateTransaction(ExecuteStateTransactionRequest) returns (google.protobuf.Empty) {}
```

To avoid inconsistencies between this document and the code, please refer to [the newest proto file](https://github.com/mosn/layotto/blob/main/spec/proto/runtime/v1/runtime.proto) for detailed input parameters and return values.
