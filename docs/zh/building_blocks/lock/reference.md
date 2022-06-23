# Distributed Lock API
## 什么是分布式锁 API
分布式锁 API基于某种存储系统（比如Etcd,Zookeeper)为开发者提供简单、易用的分布式锁API接口，开发者可以用该API获取锁、保护共享资源免受并发问题的烦扰。

## 如何使用分布式锁 API
您可以通过grpc调用分布式锁 API，接口定义在[runtime.proto](https://github.com/mosn/layotto/blob/main/spec/proto/runtime/v1/runtime.proto) 中。

使用前需要先对组件进行配置，详细的配置说明见[分布式锁组件文档](zh/component_specs/lock/common.md)

### 使用示例
Layotto client sdk封装了grpc调用的逻辑，使用sdk调用分布式锁 API的示例可以参考[快速开始：使用分布式锁API](zh/start/lock/start.md)

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

**Q: expire字段的时间单位？**

A: 秒。

**Q: 如果多个客户端传相同的LockOwner会怎么样?**

case 1. 两个客户端app-id不一样，传的LockOwner相同，不会发生冲突

case 2. 两个客户端app-id一样，传的LockOwner相同，会发生冲突。可能会出现抢到别人的锁、释放别人的锁等异常情况

因此用户需要保证LockOwner的唯一性，例如给每个请求分配一个UUID,golang写法：

```go
import "github.com/google/uuid"
//...
req.LockOwner = uuid.New().String()
```

### Unlock

```protobuf
  rpc Unlock(UnlockRequest)returns (UnlockResponse) {}
```

为避免文档和代码不一致，详细入参和返回值请参考[proto文件](https://github.com/mosn/layotto/blob/main/spec/proto/runtime/v1/runtime.proto)

## 为什么分布式锁 API被设计成这样
如果您对实现原理、设计逻辑感兴趣，可以查阅[分布式锁API设计文档](zh/design/lock/lock-api-design)