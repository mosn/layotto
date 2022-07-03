# Distributed Lock API
## What is distributed lock API
The distributed lock API is based on a certain storage system (such as Etcd, Zookeeper) to provide developers with a simple and easy-to-use distributed lock API. Developers can use the API to obtain locks and protect shared resources from concurrency problems.

## How to use distributed lock API
You can call the distributed lock API through grpc. The API is defined in [runtime.proto](https://github.com/mosn/layotto/blob/main/spec/proto/runtime/v1/runtime.proto).

The component needs to be configured before use. For detailed configuration instructions, see [Distributed Lock Component Document](en/component_specs/lock/common.md)

### Example
Layotto client sdk encapsulates the logic of grpc calling. For an example of using sdk to call distributed lock API, please refer to [Quick Start: Using Distributed Lock API](en/start/lock/start.md)


### TryLock

```protobuf
// A non-blocking method trying to get a lock with ttl.
rpc TryLock(TryLockRequest)returns (TryLockResponse) {}

message TryLockRequest {
  // Required. The lock store name,e.g. `redis`.
  string store_name = 1;
  // Required. resource_id is the lock key. e.g. `order_id_111`
  // It stands for "which resource I want to protect"
  string resource_id = 2;
  
  // Required. lock_owner indicate the identifier of lock owner.
  // You can generate a uuid as lock_owner.For example,in golang:
  //
  // req.LockOwner = uuid.New().String()
  //
  // This field is per request,not per process,so it is different for each request,
  // which aims to prevent multi-thread in the same process trying the same lock concurrently.
  //
  // The reason why we don't make it automatically generated is:
  // 1. If it is automatically generated,there must be a 'my_lock_owner_id' field in the response.
  // This name is so weird that we think it is inappropriate to put it into the api spec
  // 2. If we change the field 'my_lock_owner_id' in the response to 'lock_owner',which means the current lock owner of this lock,
  // we find that in some lock services users can't get the current lock owner.Actually users don't need it at all.
  // 3. When reentrant lock is needed,the existing lock_owner is required to identify client and check "whether this client can reenter this lock".
  // So this field in the request shouldn't be removed.
  string lock_owner = 3;
  
  // Required. expire is the time before expire.The time unit is second.
  int32 expire = 4;
}


message TryLockResponse {
  bool success = 1;
}
```

**Q: What is the time unit of the expire field?**

A: Seconds.

**Q: What would happen if different applications pass the same lock_owner?**

case 1. If two apps with different app-id pass the same lock_owner,they won't conflict because lock_owner is grouped by 'app-id ',while 'app-id' is configurated in sidecar's static config(configurated in config.json)

case 2.If two apps with same app-id pass the same lock_owner,they will conflict and the second app will obtained the same lock already used by the first app.Then the correctness property will be broken.

So user has to care about the uniqueness property of lock_owner.You can generate a uuid as lock_owner.For example,in golang:

```go
req.LockOwner = uuid.New().String()
```

### Unlock

```protobuf
  rpc Unlock(UnlockRequest)returns (UnlockResponse) {}
```

To avoid inconsistencies between the documentation and the code, please refer to [proto file](https://github.com/mosn/layotto/blob/main/spec/proto/runtime/v1/runtime.proto) for detailed input parameters and return values

## Why is the distributed lock API designed like this
If you are interested in the implementation principle and design logic, you can refer to [Distributed Lock API Design Document](en/design/lock/lock-api-design)
