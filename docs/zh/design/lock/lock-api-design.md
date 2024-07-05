# 0. 太长，不看
添加TryLock和Unlock API. 

续租API有争议，第一版不加入续租API

# 1. 调研
| **系统** | **能否实现分布式锁** | **阻塞锁(基于watch)** | **可用性** | **写操作线性一致** | **sequencer([chubby论文里提出的feature](https://static.googleusercontent.com/media/research.google.com/zh-TW//archive/chubby-osdi06.pdf))** | **续租** |
| --- | --- | --- | --- | --- | --- | --- |
| 单机redis | yes | x | 单点失效时，锁服务不可用 | yes | yes(need poc) | yes |
| redis集群 | yes | x | yes | no. 故障转移可能导致丢锁 | yes(need poc) | yes |
| redis Redlock | yes |  |  |  |  |  |
| consul | yes |  |  |  |  |  |
| zookeeper | yes | yes | yes 有fo能力,[200 ms内完成选举](https://pdos.csail.mit.edu/6.824/papers/zookeeper.pdf) | yes | yes 使用zxid作为sequencer | yes |
| etcd | yes | yes | yes | yes | yes use revision | yes lease.KeepAlive |

可以看到能力还是有一定差异

# 2. High-level design
## 2.1. API
### 2.1.0. 设计原则
我们面临着很多诱惑，分布式锁其实有很多feature可以做（阻塞锁，可重入锁，读写锁，sequencer等等）

但是毕竟我们的目标是设计一套足够通用的API规范，那么在API制定上还是尽量保守些，start simple，先把简单、常用的功能抽象成API规范，等用户反馈后再考虑将更多功能抽象成API规范

### 2.1.1. TryLock/Unlock
最基础的加锁、解锁功能。TryLock非阻塞，如果没有抢到锁直接返回

proto:

```protobuf
// Distributed Lock API
rpc TryLock(TryLockRequest)returns (TryLockResponse) {}

rpc Unlock(UnlockRequest)returns (UnlockResponse) {}

message TryLockRequest {
  string store_name = 1;
  // resource_id is the lock key.
  string resource_id = 2;
  // lock_owner indicate the identifier of lock owner.
  // This field is required.You can generate a uuid as lock_owner.For example,in golang:
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
  // expire is the time before expire.The time unit is second.
  int32 expire = 4;
}

message TryLockResponse {

  bool success = 1;
}

message UnlockRequest {
  string store_name = 1;
  // resource_id is the lock key.
  string resource_id = 2;

  string lock_owner = 3;
}

message UnlockResponse {
  enum Status {
    SUCCESS = 0;
    LOCK_UNEXIST = 1;
    LOCK_BELONG_TO_OTHERS = 2;
    INTERNAL_ERROR = 3;
  }

  Status status = 1;
}

```


**Q: expire字段的时间单位？**

A: 秒。

**Q: 能否强制让用户把秒数配大点，别配太小？**

A: 没法在编译时或者启动时限制，算了

**Q: 如果多个客户端传相同的LockOwner会怎么样?**

case 1. 两个客户端app-id不一样，传的LockOwner相同，不会发生冲突

case 2. 两个客户端app-id一样，传的LockOwner相同，会发生冲突。可能会出现抢到别人的锁、释放别人的锁等异常情况

因此用户需要保证LockOwner的唯一性，例如给每个请求分配一个UUID,golang写法：

```go
import "github.com/google/uuid"
//...
req.LockOwner = uuid.New().String()
```

**Q: 为啥不加metadata**

A: 一开始尽量保守一些，等有人反馈有需求再加，或者实现组件过程中发现确实有需要再加

**Q: 以后要支持sequencer、可重入锁之类的feature咋加？**

A: 入参增加feature option，组件也要实现Support()接口

### 2.1.2. 续租
#### Solution A: add an API "LockKeepAlive"

```protobuf
rpc LockKeepAlive(stream LockKeepAliveRequest) returns (stream LockKeepAliveResponse){}
  
message LockKeepAliveRequest {
  // resource_id is the lock key.
  string resource_id = 1;

  string client_id = 2;
  // expire is the time to expire
  int64 expire = 3;
}

message LockKeepAliveResponse {
  enum Status {
    SUCCESS = 0;
    LOCK_UNEXIST = 1;
    LOCK_BELONG_TO_OTHERS = 2;
  }
  // resource_id is the lock key.
  string resource_id = 1;

  Status status = 2;
}
```

续租的入参、返回结果都是stream,这里参考etcd的实现，app和sidecar只需要维护一个连接，每次用锁需要续租的时候都复用该连接传递续租请求。

**Q: 为啥不把续租作为一个stream参数塞到tryLock里？**

A: 很多业务不用续租，让trylock尽量简单；

尽量单一职责，后续加上阻塞锁，可以公用续租API；

把带stream的部分拆开成一个接口，以便http实现简单。

**Q:续租逻辑太复杂，能否让用户不感知？**

A: sdk屏蔽掉这层逻辑，开个线程/协程/nodejs定时事件，自动续租


#### Solution B: 用户不感知续租逻辑，自动续租，App和sidecar维持统一心跳
缺点/难点：

1. 使用统一心跳的话，难以定制心跳间隔

解法是保证心跳间隔低，比如1秒1次

2. 如何保证可靠的故障检测？
   
例如以下java代码，unlock可能失败：

```java
try{

}finally{
  lock.unlock()   
}
```

如果是单机锁，unlock能保证成功（除非整个jvm故障），但是unlock走网络调用的话可能失败。调用失败后，怎么保证心跳断掉？

这就要求业务在开发时要往心跳检测里上报一些细粒度的状态。

我们可以定义http callback接口,由actuator轮询检测，约定callback返回的数据结构为：

```json
{
  "status": "UP",
  "details": {
    "lock": [
      {
        "resource_id": "res1",
        "client_id": "dasfdasfasdfa",
        "type": "unlock_fail"
      }
    ],
    "xxx": []
  }
}
```

应用要处理状态收集、上报、上报成功后清理、限制map容量（比如上报失败次数多了，map太大OOM怎么办），总归要求app实现一些复杂逻辑，也要放sdk里

3. 这个实现其实和续租一样的，都是开一个单独的连接做状态管理，使用过程中有需要就通过这个公用连接上报状态。
4. API spec依赖心跳逻辑。依赖心跳间隔、心跳返回的数据结构。相当于API spec依赖了Layotto的实现，除非我们能把心跳实现也标准化掉（包括间隔、返回的数据结构等）
5. 心跳检测失败一次，sidecar是否继续续租？如果sidecar停止续租了，下一次心跳检测又正常了，sidecar是否恢复续租？比较难界定，因为心跳检测会有个怀疑机制（比如心跳失败几次才算失败），等怀疑期满再停止续租太慢了，激进停止续租又会有反复问题
6. app心跳停了一次、后面又恢复的场景，sidecar如果继续续租，怎么知道app的锁还在不在（app可能只是短时间网络抖动、然后恢复了，也可能是干脆重启了没锁了）

#### Solution C:app自己重试unlock
如果unlock失败，app自己异步重试unlock

#### 结论
目前大家对续租方案有不同意见，一期先不实现续租功能。

个人倾向A方案,sdk屏蔽掉续租逻辑。虽然用grpc直接调需要处理续租逻辑，但是续租算是分布式锁常用方案，开发者的理解成本低。

抛出来看大家意见


## 2.2. 组件
### 2.2.1. 如何处理“组件不支持某个feature option”
模拟>运行报错>ignore

特殊情况才ignore，比如consistency传了是eventual consistency,但是存储系统本身是强一致的，可以ignore这个option

### 2.2.2. 如何处理“组件不支持某个API”
模拟>运行报错

以事务API举例



### 2.2.3. 选型
一期选择单机redis实现

### 2.2.4. "运行时报错"式设计会破坏可移植性，如何让用户更简单的评估能否移植？
A. 提供文档，文档上每种组件给的保证、feature不同，让用户自己评估

B. 在sidecar的配置里声明API支持的feature option,如果sidecar启动发现组件不match，启动时报错

C. sidecar打运行时日志，自动统计app用了哪些feature，待需要移植时进行日志分析

D. 先不管

E. 做一个静态分析工具，自动分析可移植性

结论：选A，因为简单

# 3. Future work

- 可重入

会有一些计数逻辑。需要考虑是所有锁默认支持可重入，还是传参里面加个feature option、标识用户需要可重入

- 阻塞锁
- sequencer

# 4. Reference
[How to do distributed locking](https://martin.kleppmann.com/2016/02/08/how-to-do-distributed-locking.html)

[The Chubby lock service for loosely-coupled distributed systems](https://static.googleusercontent.com/media/research.google.com/zh-TW//archive/chubby-osdi06.pdf)

[https://www.jianshu.com/p/6e72e3ee5623](https://www.jianshu.com/p/6e72e3ee5623)

[http://zhangtielei.com/posts/blog-redlock-reasoning.html](http://zhangtielei.com/posts/blog-redlock-reasoning.html)

[http://zhangtielei.com/posts/blog-redlock-reasoning-part2.html](http://zhangtielei.com/posts/blog-redlock-reasoning-part2.html)
