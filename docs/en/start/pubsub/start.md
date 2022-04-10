# Use Pub/Sub API to implement pub/sub pattern
## What is Pub/Sub API
Developers often use message queues (such as open source Rocket MQ, Kafka, such as AWS SNS/SQS provided by cloud vendors) to implement message publishing and subscription. The publish-subscribe model can help applications better decouple and cope with peak traffic.

Unfortunately, the APIs of these message queue products are different. When developers want to deploy their apps across clouds, or want their apps to be portable (for example, easily moving from Alibaba Cloud to Tencent Cloud), they have to refactor their code.

The design goal of Layotto Pub/Sub API is to define a unified message publish/subscribe API. The application only needs to care about the API, and does not need to care about the specific message queue product being used, so that the application can be transplanted at will, and the application is sufficiently "cloud native" .

## Quick start
This example shows how to call redis through Layotto to publish/subscribe messages.

The architecture of this example is shown in the figure below. The running processes are: redis, a Subscriber program that listens to events, Layotto, and a Publisher program that publishes events.

![img_1.png](../../../img/mq/start/img_1.png)

### Step 1. Deploy and Run Redis in Docker

1. Get the latest version of Redis image.
   
Here we pull the latest version of the official image:

```shell
docker pull redis:latest
```

2. Check local mirror

Use the following command to check whether Redis is installed:

```shell
docker images
```
![img.png](../../../img/mq/start/img.png)

3. Run the container

After the installation is complete, we can use the following command to run the Redis container:

```shell
docker run -itd --name redis-test -p 6380:6379 redis
```

Parameter Description:

`-p 6380:6379`: Map port 6379 of the container to port 6380 of the host. The outside can directly access the Redis service through the host ip:6380.

### Step 2. Start the Subscriber program and subscribe to events
Build:

```shell
 cd ${project_path}/demo/pubsub/redis/server/
 go build -o subscriber
 ```

Start subscriber:

```shell @background
 ./subscriber
```
If the following information is printed out, it means the startup is successful:

```bash
Start listening on port 9999 ...... 
```

Explanation:

The program will start a gRPC server and open two API:

- ListTopicSubscriptions

Calling this API will return the topics subscribed by the application. This program will return "topic1"

- OnTopicEvent

When a new event occurs, Layotto will call this API to notify the Subscriber of the new event.

After the program receives a new event, it will print the event to the command line.

### Step 3. Run Layotto

After downloading the project code to the local, switch the code directory and compile:

```shell
cd ${project_path}/cmd/layotto
```

```shell @if.not.exit layotto
go build
```

After completion, the layotto file will be generated in the directory, run it:

```shell @background
./layotto start -c ../../configs/config_redis.json
```

### Step 4. Run the Publisher program and call Layotto to publish events

```shell
 cd ${project_path}/demo/pubsub/redis/client/
 go build -o publisher
 ./publisher
```

If the following information is printed, the call is successful:

```bash
Published a new event.Topic: topic1 ,Data: value1 
```

### Step 5. Check the event message received by the subscriber

Go back to the subscriber's command line and you will see that a new message has been received:

```bash
Start listening on port 9999 ...... 
Received a new event.Topic: topic1 , Data:value1 
```

### Next Step
#### What did this client Demo do?
The demo client program uses the golang version SDK provided by Layotto, calls Layotto Pub/Sub API, and publishes events to redis. Later, Layotto received the new events in redis, and sent the new events back to the callback API opened by the Subscriber program to notify the Subscriber.

The sdk is located in the `sdk` directory, and users can call the API provided by Layotto through the sdk.

In addition to using sdk, you can also interact with Layotto directly through grpc in any language you like.

In fact, sdk is only a very thin package for grpc, using sdk is about equal to directly using grpc.

#### Details later, let's continue to experience other APIs
Explore other Quickstarts through the navigation bar on the left.


#### Understand the principle of Pub/Sub API implementation
If you are interested in the implementation principle, or want to extend some functions, you can read [Pub/Sub API design document](en/design/pubsub/pubsub-api-and-compability-with-dapr-component.md)
