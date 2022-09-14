# Use Pub/Sub API to implement pub/sub pattern
## What is Pub/Sub API
Developers often use message queues (such as open source Rocket MQ, Kafka, such as AWS SNS/SQS provided by cloud vendors) to implement message publishing and subscription. The publish-subscribe model can help applications better decouple and cope with peak traffic.

Unfortunately, the APIs of these message queue products are different. When developers want to deploy their apps across clouds, or want their apps to be portable (for example, easily moving from Alibaba Cloud to Tencent Cloud), they have to refactor their code.

The design goal of Layotto Pub/Sub API is to define a unified message publish/subscribe API. The application only needs to care about the API, and does not need to care about the specific message queue product being used, so that the application can be transplanted at will, and the application is sufficiently "cloud native" .

## Quick start
This example shows how to call redis through Layotto to publish/subscribe messages.

The architecture of this example is shown in the figure below. The running processes are: redis, a Subscriber program that listens to events, Layotto, and a Publisher program that publishes events.

![img_1.png](../../../img/mq/start/img_1.png)

### Step 1. Start the Subscriber
<!-- tabs:start -->
#### **Go**
Build the golang subscriber

```shell
 cd demo/pubsub/server/
 go build -o subscriber
 ```

Start subscriber:

```shell @background
 ./subscriber -s pub_subs_demo
```

#### **Java**

Download the java sdk and examples:

```bash
git clone https://github.com/layotto/java-sdk
```

```bash
cd java-sdk
```

Build and run it:

```bash
# build example jar
mvn -f examples-pubsub-subscriber/pom.xml clean package
# run the example
java -jar examples-pubsub-subscriber/target/examples-pubsub-subscriber-jar-with-dependencies.jar
```

<!-- tabs:end -->

If the following information is printed out, it means the startup is successful:

```bash
Start listening on port 9999 ...... 
```

> [!TIP|label: What did this subscriber do]
> The Subscriber program started a gRPC server and exported two gRPC API:
>
> - ListTopicSubscriptions
>
> Calling this API will return the topics subscribed by the application. This program will return "topic1" and "hello"
>
> - OnTopicEvent
>
> When a new event occurs, Layotto will call this API to notify the Subscriber of the new event.
>
> After the program receives a new event, it will print the event to the command line.

### Step 2. Deploy Redis and Layotto
<!-- tabs:start -->
#### **with Docker Compose**
You can start Redis and Layotto with docker-compose

```bash
cd docker/layotto-redis
# Start redis and layotto with docker-compose
docker-compose up -d
```

#### **Compile locally (not for Windows)**
You can run Redis with Docker, then compile and run Layotto locally.

> [!TIP|label: Not for Windows users]
> Layotto fails to compile under Windows. Windows users are recommended to deploy using docker-compose

#### step 2.1. Run Redis with Docker

We can use the following command to run the Redis container:

```shell
docker run -itd --name redis-test -p 6380:6379 redis
```

Parameter Description:

`-p 6380:6379`: Map port 6379 of the container to port 6380 of the host. The outside can directly access the Redis service through the host ip:6380.

#### Step 2.2. Run Layotto

After downloading the project code to the local, switch the code directory and compile:

```shell
cd ${project_path}/cmd/layotto
```

```shell @if.not.exist layotto
go build
```

After completion, the layotto file will be generated in the directory, run it:

```shell @background
./layotto start -c ../../configs/config_redis.json
```

<!-- tabs:end -->

### Step 3. Run the Publisher program and call Layotto to publish events
<!-- tabs:start -->
#### **Go**
Build the golang publisher:

```shell
 cd ${project_path}/demo/pubsub/client/
 go build -o publisher
 ./publisher -s pub_subs_demo
```

#### **Java**

Download the java sdk and examples:

```shell @if.not.exist java-sdk
git clone https://github.com/layotto/java-sdk
```

```shell
cd java-sdk
```

Build:

```shell @if.not.exist examples-pubsub-publisher/target/examples-pubsub-publisher-jar-with-dependencies.jar
# build example jar
mvn -f examples-pubsub-publisher/pom.xml clean package
```

Run it:

```shell
# run the example
java -jar examples-pubsub-publisher/target/examples-pubsub-publisher-jar-with-dependencies.jar
```

<!-- tabs:end -->

If the following information is printed, the call is successful:

```bash
Published a new event.Topic: hello ,Data: world
Published a new event.Topic: topic1 ,Data: value1
```

### Step 4. Check the event message received by the subscriber

Go back to the subscriber's command line and you will see that a new message has been received:

```bash
Start listening on port 9999 ......
Received a new event.Topic: topic1 , Data: value1
Received a new event.Topic: hello , Data: world
```

### step 5. Stop containers and release resources
<!-- tabs:start -->
#### **Docker Compose**
If you started Redis and Layotto with docker-compose, you can shut them down as follows:

```bash
cd ${project_path}/docker/layotto-redis
docker-compose stop
```

#### **Destroy the Redis container**
If you started Redis with Docker, you can destroy the Redis container as follows:

```shell
docker rm -f redis-test
```

<!-- tabs:end -->

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
