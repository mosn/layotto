# Add TryLock and Unlock API

Add TryLock and Unlock API.

The renewal API is disputed and version 1 is not part of the renewal API

# Research

| **System**    | **Distributed Lock** | **Block lock (watch)** | **Availability**                                                                                                                                                                                                                                                    | **Write line consistent**                      | **sequencer ([chubby论文里提出的feature](https://static.googleusercontent.com/media/research.google.com/zh-TW//archive/chubby-osdi06.pdf)** | **Renewed**                          |
| ------------- | -------------------- | ----------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------ |
| Solo Redis    | yes                  | x                                         | Lock service is not available when single point expires                                                                                                                                                                                                             | yes                                            | yes(needed pic)                                                                                                                      | yes                                  |
| rediscluster  | yes                  | x                                         | yes                                                                                                                                                                                                                                                                 | no. Failures can cause locking | yes(needed pic)                                                                                                                      | yes                                  |
| redis Redlock | yes                  |                                           |                                                                                                                                                                                                                                                                     |                                                |                                                                                                                                                         |                                      |
| consul        | yes                  |                                           |                                                                                                                                                                                                                                                                     |                                                |                                                                                                                                                         |                                      |
| zookeeper     | yes                  | yes                                       | yes has fo, [elections within 200 ms] (https://pdos.csail.mit.edu/6.824/papers/zookeeper.pdf) | yes                                            | yes use zxid as sensor                                                                                                                                  | yes                                  |
| etcd          | yes                  | yes                                       | yes                                                                                                                                                                                                                                                                 | yes                                            | yes use revision                                                                                                                                        | yes lease. KeepAlive |

Can see some difference in capacity or

# High-level design

## API

### Design principles

We are faced with many temptations, distributive locks have many features that can be done (blocked, relocked, read, lock, sequencer, etc.)

But after all, our goal is to design a sufficient set of common API specifications, so as to be as conservative as possible on the API, to abstract simple and common features into API specifications before considering more features into API norms

### TryLock/Unlock

Base Lockout and unlock features.TryLock is not blocked, return directly if no locks are grabbed

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

**Q: the time unit for the expire?**

A: seconds.

**Q: Can you force users to make seconds large points without being too small?**

A: is not allowed on compilation or on startup.

**Q: What if multiple clients transmit the same LockOwner?**

Case 1. Two client app-id is different, with the same lockOwner passed and no conflict

Case 2. Like the two client app-ids, there is a conflict between LockOwner with the same passage.There may be exceptional circumstances in which someone can be robbed, frees someone else's lock.

The user therefore needs to ensure the uniqueness of LockOwner, such as assigning a UUID, golangprototype： per request

```go
import "github.com/google/uuid"
//...
req.LockOwner = uuid.New().String()
```

**Q: What does not meta**

A: Be as conservative as possible at the beginning, and some people get feedback about additional needs or find out in the implementation of the component that there is a real need for settlement

**Q: Do you want to support sequentier and relocking features like then?**

A: Add feature option to add components to support interface

### Renewal

#### Solvation A: add an API "LockKeepAlive"

```protobuf
rpc LockKeepAlive(stream LockKeepAliveRequest) returns (stream LockKeepAliveResponse) {}
  
message LockKeepAliveRequest L$
  // resource_id is the lock key.
  string resource_id = 1;

  string client_id = 2;
  // expire is the time to expire
  int64 expire = 3;
}

message LockKeepAliveResponse LO
  enum Status Fit
    SUCCESS = 0;
    LOCK_UNEXIST = 1;
    LOCK_BELONG_TO_OTHERS = 2;
  }
  // resource_id is the lock key.
  string resource_id = 1;

  Status status = 2;
}
```

The re-entry of the lease and the return result is stream, where the implementation of the reference etcd, the app and sidecar need to maintain only one connection and use the link to relay the renewal lease request every time the lock requires a renewal lease.

**Q: Why don't you use the renewal as a stream parameter to tryLock?**

A: Many businesses are not rented to make trylock as simple as possible;

To maximize single responsibilities, followed by blocklocks that can be used to renew the API;

Split parts with streams into an interface to make it easy for https.

**Q: The renewal logic is too complex to make users unknown?**

A: sdk block this logic, start thread/protocol/nodejs timing, automatic renewal

#### Solvolution B: Users are not aware of continued lease logic, auto-rent, app and sidecar maintain a unified heart jump

Disadvantages/Difficulty：

1. Use uniform heart jumps, it is difficult to customize interval

Solve is a guaranteed low jump intervals, for example 1 second time in 1 second

2. How to ensure reliable troubleshooting?

e.g. the java code below, unlock may fail：

```java
tryLO

}finallyLO
  lock.unlock()   
}
```

If a single machine lock, unlock guarantees success (unless the entire jvm failure), but unlocks can fail.How can you make sure to jump off when the call fails?

This requires business to jump to the state of some fine particle levels when developing the business.

We can define the HTML callback interface to be detected by actuator wheel and agree that the callback returned data structure is：

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

App handles state gathering, reporting and clearing after successful reporting, limiting map capacity (e.g. the number of failed reports and the size of the map's big OOOM). Altogether, it requires the app to perform some complicated logic, as well as to release the sdk

3. This implementation, like renewal of rent, is managed on a separate connection that needs to be reported through this public connection.
4. API speciality dependency jump logic.Dependency jumps interval and heart jump back data structure.The equivalent of the API speck relies on Layotto implementation, unless we can also standardize heart jump and drop (including intervals, returned data structures, etc.)
5. Heart jump failed once. Do sidecar continue rent?If the sidecar stopped renting anymore, will the next heart jump test be normal, will sidecar resume the lease?It is more difficult to define, because there is a skeptical mechanism (for example, a few failures of heart failure), when the skepticism ends too slowly and when the radical stop of the lease can be repeated
6. The app skipped once and then restored. sidecar knows how to know if the app continues to be locked (it may be just a short time to shake, then restored, or it may simply restart unlocked)

#### Solvation C:app retries unlock

If unlocks fail, apply your own asynchronous retry unlock

#### Conclusion

There are currently divergent views on the renewal of the rental scheme and the renewal of the rental function will not be achieved in one issue.

Personal inclination to Program A, sdk blocks the logic of renewal of rent.While a grpc direct transfer would require the processing of the renewal logic, the renewal of the lease would be a distributed localized locking programme and a low cost of understanding by the developers.

Throw out for everybody

## Components

### How to deal with "components do not support a feature option"

Simulate >Run Error>ignore

Exceptional circumstances, such as consistence, are eventual, but the storage system itself is strong and can be used as an option

### How to deal with "components do not support an API"

Simulate >Run Error

As transaction API examples

### Selection:

Select Single Machine Redisimplementation

### "Running missing" designs destroy portability and how can users be made simpler to assess whether they can transplant?

Providing a document, where each component of the document gives a different guarantee, feature, allowing users to assess themselves

Declares the feature option supported by the API in the sidecar configuration. If sidecar starts to discover that the component is not match, report an error on startup

C. Residecar running time logs, automatic statistical data on what features the app uses and log analysis when transplant is required

D. Profiling first

E. Providing a static analysis tool for automatic transplantation analysis

Conclusion：Option A because it is simple

# Future work

- Remindable

There will be some arithmetic logic.You need to consider whether all locks can be reentered by default, or add a feature to the upload.

- Block Lock
- sequencer

# Reference

[How to do distressed blocking](https://martin.kleppmann.com/2016/02/08/how-to-do-distributed-locking.html)

[The Chubby lock service for loosely-coupled distributed systems](https://static.googleusercontent.com/media/research.google.com/zh-TW//archive/chubby-osdi06.pdf)

[https://www.jianshu.com/p/6e72e3eee5623](https://www.jianshu.com/p/6e72eee5623)

[http://zhangtielei.com/posts/blog-redlock-reasoning.html](http://zhangtielei.com/posts/blog-redlock-reasoning.html)

[http://zhangtielei.com/posts/blog-redlock-reasoning-part2.html](http://zhangtielei.com/posts/blog-redlock-reasoning-part2.html)
