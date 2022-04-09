# Distributed Lock API demo with redis

This example shows how to call redis through Layotto to trylock/unlock.

The architecture of this example is shown in the figure below, and the started processes are: redis, Layotto, a client program with two goroutines trying the same lock concurrently.

![img.png](../../../img/lock/img.png)
### Step 1. Deploy redis using Docker

1. Get the latest version of Redis docker image

Here we pull the latest version of the official image:

```shell
docker pull redis:latest
```

2. View the local mirror 
   
Use the following command to check if redis is installed:
   
```shell
docker images
```

![img.png](../../../img/mq/start/img.png)

3. Run the container

After the installation is complete, we can use the following command to run the redis container:

```shell
docker run -itd --name redis-test -p 6380:6379 redis
```

Parameter Description:

`-p 6380:6379`: Map port 6379 of the container to port 6380 of the host. The outside can directly access the Redis service through the host ip:6380.

### Step 2. Run Layotto

After downloading the project code to the local, enter the code directory and compile:

```shell
cd ${projectpath}/cmd/layotto
go build
```

The layotto file will be generated in the directory, run it:

```bash
./layotto start -c ../../configs/config_lock_redis.json
```

<!--
```shell
nohup ./layotto start -c ../../configs/config_lock_redis.json &
```
-->

### Step 3. Run the client program, call Layotto to add, delete, modify and query

```shell
 cd ${projectpath}/demo/lock/redis/
 go build -o client
 ./client
```

If the following information is printed, the call is successful:

```bash
client1 prepare to tryLock...
client1 got lock!ResourceId is resource_a
client2 prepare to tryLock...
client2 failed to get lock.ResourceId is resource_a
client1 prepare to unlock...
client1 succeeded in unlocking
client2 prepare to tryLock...
client2 got lock.ResourceId is resource_a
client2 succeeded in unlocking
Demo success!
```

### Next Step
#### What did this client Demo do?
The demo client program uses the golang version SDK provided by Layotto, calls the Layotto distributed lock API, and starts multiple goroutines to do locking and unlocking operations.

The sdk is located in the `sdk` directory, and users can call the API provided by Layotto through the sdk.

In addition to using sdk, you can also interact with Layotto directly through grpc in any language you like.

In fact, sdk is only a very thin package for grpc, using sdk is about equal to directly using grpc.

#### Details later, let's continue to experience other APIs
Explore other Quickstarts through the navigation bar on the left.


#### Understand the design principle of Distributed Lock API

If you are interested in the design principle, or want to extend some functions, you can read [Distributed Lock API design document](en/design/lock/lock-api-design.md)