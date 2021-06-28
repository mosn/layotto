# 0. tl;dl
Add TryLock and Unlock API. The Lock Renewal API is controversial and will not be added into the first version

# 1. Evaluation of products on the market
| **System** | **try lock** | **Blocking lock(based on watch)** | **Availability** | **Write operations are linearizable** | **sequencer([chubby's feature](https://static.googleusercontent.com/media/research.google.com/zh-TW//archive/chubby-osdi06.pdf))** | **Lock renewal** |
| --- | --- | --- | --- | --- | --- | --- |
| Stand-alone redis | yes | x | unavailable when single failure | yes | yes(need poc) | yes |
| redis cluster | yes | x | yes | no. Locks will be unsafe when fail-over happens | yes(need poc) | yes |
| redis Redlock | yes | | | | | |
| nacos | × | | | | | |
| consul | yes | | | | | |
| eureka | × | | | | | |
| zookeeper | yes | yes | yes. [the election completes within 200 ms](https://pdos.csail.mit.edu/6.824/papers/zookeeper.pdf) | yes | yes use zxid as sequencer | yes |
| etcd | yes | yes | yes | yes | yes use revision | yes lease.KeepAlive |

There are some differences in feature supporting.

# 2. High-level design
## 2.1. API
### 2.1.0. Design principles
We are faced with many temptations. In fact, there are many lock related features that can be supported (blocking locks, reentrant locks, read-write locks, sequencer, etc.)

But after all, our goal is to design a general API specification, so we should be as conservative as possible in API definition.Start simple, abstract the simplest and most commonly used functions into API specifications, and wait for user feedback before considering adding more abstraction into API specification.

### 2.1.1. TryLock/Unlock API
The most basic locking and unlocking API. 

TryLock is non-blocking, it return directly if the lock is not obtained.


proto:
```protobuf
// Distributed Lock API
rpc TryLock(TryLockRequest)returns (TryLockResponse) {}

rpc Unlock(UnlockRequest)returns (UnlockResponse) {}

message TryLockRequest {
  // resource_id is the lock key.
  string resource_id = 1;
  // client_id will be automatically generated if not set
  optional string client_id = 2;
  // expire is the time before expire
  int64 expire = 3;
}

message TryLockResponse {

  bool success = 1;

  string client_id = 2;
}

message UnlockRequest {
  // resource_id is the lock key.
  string resource_id = 1;

  string client_id = 2;
}

message UnlockResponse {
  enum Status {
    SUCCESS = 0;
    LOCK_UNEXIST = 1;
    LOCK_BELONG_TO_OTHERS = 2;
  }

  Status status = 1;
}

```
**Q: What is the time unit of the expire field?**

A: Seconds.

**Q: Can we force the user to set the number of seconds to be large enough(instead of too small)?**

A: There is no way to limit it at compile time or startup, forget it

**Q: Why not add metadata field**

A: Try to be conservative at the beginning, wait until someone feedbacks that there is a need, or find that there is a need to be added in the process of implementing the component

**Q: How to add features such as sequencer and reentrant locks in the future?**

A: Add feature options in the API parameters,and the component must also implement the Support() function

### 2.1.2. Lock Renewal
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
The input parameters and return results of this API are all streams. App and sidecar only need to maintain one connection. Each time the lock needs to be renewed, the connection is reused to transfer the renewal request.

**Q: Why not put the lock renewal as a stream parameter into tryLock?**

A: Many businesses do not need to renew leases, so we want trylock to be as simple as possible;

Single responsibility principle.When we want to add a blocking lock,the renewal API can be reused;

**Q: The renewal logic is too complicated, can we make it transparent to users?**

A: sdk shields this layer of logic, starts a thread/coroutine/nodejs timing event, and automatically renews the lease


#### Solution B: Users do not perceive the renewal logic, and automatically renew the lease. App and sidecar maintain a heartbeat for failure detect
Disadvantages/difficulties:

1. If you reuse a public heartbeat, it is difficult to customize the heartbeat interval

The solution is to ensure that the heartbeat interval is low enough, such as 1 time per second

2. How to ensure reliable failure detection?

For example, the following java code, unlock method may fail:
```java
try{

}finally{
  lock.unlock()
}
```
If it is a lock in JVM, unlock can guarantee success (unless the entire JVM fails), but unlock may fail if it is called via the network. How to ensure that the heartbeat is interrupted after the call fails?

This requires the app to report some fine-grained status to the heartbeat detection.

We can define the http callback SPI, which is polled and detected by the Layotto actuator, and the data structure returned by the callback is as follows:
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
The application has to handle status collection, reporting, cleaning up after the report is successful, and limiting the map capacity (for example, what if the map is too large when report fails too much times?), which requires the app to implement some complex logic, and it must be put in the SDK.

3. This implementation is actually the same as lease renewal. It opens a separate connection for status management, and user reports the status through this public connection when necessary.
   
4. API spec relies on heartbeat logic. It relies on the heartbeat interval and the data structure returned by the heartbeat. It is equivalent to that the API spec relies on the implementation of Layotto, unless we can also standardize the heartbeat API (including interval, returned data structure, etc.)

#### in conclusion
At present, everyone has different opinions on the solution of lease renewal, and the lease renewal function will not be added in the first version.

Personally I prefer the solution A.Let the SDK shields the renewal logic. Although users have to directly deal with the lease renewal logic when using grpc, lease renewal is a common solution for distributed locks, and it is not hard for developers to understand.

I put it here to see everyone's opinions

# 3. Future work

- Reentrant Lock

There will be some counting logic.We need to consider whether all locks support reentrancy by default, or add a feature option in the parameter to identify that the user needs it to be reentrant

- Blocking lock

- sequencer

# 4. Reference

[How to do distributed locking](https://martin.kleppmann.com/2016/02/08/how-to-do-distributed-locking.html)

[The Chubby lock service for loosely-coupled distributed systems](https://static.googleusercontent.com/media/research.google.com/zh-TW//archive/chubby-osdi06.pdf)