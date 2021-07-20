# Use State API to manage state
## What is State API
Your application can use the same State API to operate different databases (or a certain storage system) to add, delete, modify and query the data of the Key/Value model.

API supports batch CRUD operations and supports the declaration of requirements for concurrency safety and data consistency. Layotto will deal with complex concurrency safety and data consistency issues for you.
## Quick start
This example shows how to call redis through Layotto to add, delete, modify and query status data.

The architecture of this example is shown in the figure below, and the started processes are: redis, Layotto, client program

![img.png](../../../img/state/img.png)
### Deploy redis using Docker

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

### Run Layotto

After downloading the project code to the local, enter the code directory and compile:

```bash
cd ${projectpath}/cmd/layotto
go build
```

The layotto file will be generated in the directory, run it:

```bash
./layotto start -c ../../configs/config_state_redis.json
```

### Run the client program, call Layotto to add, delete, modify and query

```bash
 cd ${projectpath}/demo/state/redis/
 go build -o client
 ./client
```

If the following information is printed, the call is successful:

```bash
SaveState succeeded.key:key1 , value: hello world 
GetState succeeded.[key:key1 etag:1]: hello world
SaveBulkState succeeded.[key:key1 etag:2]: hello world
SaveBulkState succeeded.[key:key2 etag:2]: hello world
GetBulkState succeeded.key:key1,value:hello world
GetBulkState succeeded.key:key3,value:
GetBulkState succeeded.key:key2,value:hello world
GetBulkState succeeded.key:key5,value:
GetBulkState succeeded.key:key4,value:
DeleteState succeeded.key:key1
DeleteState succeeded.key:key2
```
### Next step

The client demo uses the golang version SDK provided by Layotto. The SDK is located in the `sdk` directory. Users can directly call the APIs provided by Layotto through the corresponding SDK.

Besides the SDK,you can also call Layotto server directly using grpc,which makes it easy for different language to interact with Layotto.