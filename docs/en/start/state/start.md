# Use State API to manage state
## What is State API
Your application can use the same State API to operate different databases (or a certain storage system) to add, delete, modify and query the data of the Key/Value model.

API supports batch CRUD operations and supports the declaration of requirements for concurrency safety and data consistency. Layotto will deal with complex concurrency safety and data consistency issues for you.
## Quick start
This example shows how to call redis through Layotto to add, delete, modify and query status data.

The architecture of this example is shown in the figure below, and the started processes are: redis, Layotto, client program

![img.png](../../../img/state/img.png)
### step 1. Deploy redis using Docker

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

### step 2. Run Layotto

After downloading the project code to the local, change the code directory:

```shell
# change directory to ${your project path}/cmd/layotto
cd cmd/layotto
```

and then build layotto:

```shell @if.not.exist layotto
go build -o layotto
```

The layotto file will be generated in the directory, run it:

```shell @background
./layotto start -c ../../configs/config_redis.json
```

### step 3. Run the client program, call Layotto to add, delete, modify and query

```shell
# open a new terminal tab
# change directory to ${your project path}/demo/state/common/
 cd ${project_path}/demo/state/common/
 go build -o client
 ./client -s "state_demo"
```

If the following information is printed, the demo succeeded:

```bash
SaveState succeeded.key:key1 , value: hello world 
GetState succeeded.[key:key1 etag:3]: hello world
SaveBulkState succeeded.[key:key1 etag:2]: hello world
SaveBulkState succeeded.[key:key2 etag:2]: hello world
GetBulkState succeeded.key:key1 ,value:hello world ,etag:4 ,metadata:map[] 
GetBulkState succeeded.key:key4 ,value: ,etag: ,metadata:map[] 
GetBulkState succeeded.key:key2 ,value:hello world ,etag:2 ,metadata:map[] 
GetBulkState succeeded.key:key3 ,value: ,etag: ,metadata:map[] 
GetBulkState succeeded.key:key5 ,value: ,etag: ,metadata:map[] 
DeleteState succeeded.key:key1
DeleteState succeeded.key:key2
```

### step 4. Stop redis and release resources

```shell
docker rm -f redis-test
```

### Next step
#### What did this client Demo do?
The demo client program uses the golang version SDK provided by Layotto, and calls Layotto's State API to add, delete, modify, and read status data.

The sdk is located in the `sdk` directory, and users can call the API provided by Layotto through the sdk.

In addition to using sdk, you can also interact with Layotto directly through grpc in any language you like.

In fact, sdk is only a very thin package for grpc, using sdk is about equal to directly using grpc.

#### Want to learn more about State API?
What does the State API do, what problems it solves, and in what scenarios should I use it?

If you have such confusion and want to know more details about State API, you can read [State API Usage Document](zh/api_reference/state/reference)

#### Details later, let's continue to experience other APIs
Explore other Quickstarts through the navigation bar on the left.
