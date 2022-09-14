# 使用Pub/Sub API进行消息发布/订阅
## 什么是Pub/Sub API
开发者经常使用消息队列等中间件产品（比如开源的Rocket MQ,Kafka,比如云厂商提供的AWS SNS/SQS）来实现消息的发布、订阅。发布订阅模式可以帮助应用更好的解耦、应对流量洪峰。

不幸的是，这些消息队列产品的API都不一样。当应用想要跨云部署，或者想要移植（比如从阿里云搬到腾讯云），应用需要重构代码。

Layotto Pub/Sub API的设计目标是定义一套统一的消息发布/订阅API，应用只需要关心API、不需要关心具体用的哪个消息队列产品，让应用能够随意移植，让应用足够"云原生"。

## 快速开始

该示例展示了如何通过Layotto调用redis，进行消息发布/订阅。

该示例的架构如下图，启动的进程有：redis、一个监听事件的Subscriber程序、Layotto、一个发布事件的Publisher程序

![img_1.png](../../../img/mq/start/img_1.png)

### step 1. 启动 Subscriber 程序,订阅事件
<!-- tabs:start -->
#### **Go**
编译 golang 写的 subscriber:

```shell
 cd demo/pubsub/server/
 go build -o subscriber
```

运行:

```shell @background
 ./subscriber -s pub_subs_demo
```

#### **Java**

下载 java sdk 和 examples:

```bash
git clone https://github.com/layotto/java-sdk
```

切换目录:

```bash
cd java-sdk
```

构建、运行:

```bash
# build example jar
mvn -f examples-pubsub-subscriber/pom.xml clean package
# run the example
java -jar examples-pubsub-subscriber/target/examples-pubsub-subscriber-jar-with-dependencies.jar
```

<!-- tabs:end -->

打印出以下信息说明运行成功:

```bash
Start listening on port 9999 ......
```

> [!TIP|label: Subscriber 程序做了什么？]
> 该程序会启动一个gRPC服务器，开放两个接口：
> - ListTopicSubscriptions
>
> 调用该接口会返回应用订阅的Topic。本程序会返回"topic1"和 "hello"
>
> - OnTopicEvent
>
> 当有新的事件发生时，Layotto会调用该接口，将新事件通知给Subscriber。
>
> 本程序接收到新事件后，会将事件打印到命令行。

### step 2. 部署 Redis 和 Layotto
<!-- tabs:start -->
#### **使用 Docker Compose**
您可以用 docker-compose 启动 Redis 和 Layotto

```bash
cd docker/layotto-redis
# Start redis and layotto with docker-compose
docker-compose up -d
```

#### **本地编译（不适合 Windows)**
您可以使用 Docker 运行 Redis，然后本地编译、运行 Layotto。

> [!TIP|label: 不适合 Windows 用户]
> Layotto 在 Windows 下会编译失败。建议 Windows 用户使用 docker-compose 部署

#### step 2.1. 用 Docker 运行 Redis
我们可以使用以下命令来运行 Redis 容器：

```shell
docker run -itd --name redis-test -p 6380:6379 redis
```

参数说明：

-p 6380:6379：映射容器服务的 6379 端口到宿主机的 6380 端口。外部可以直接通过宿主机ip:6380 访问到 Redis 的服务。

#### step 2.2. 编译、运行 Layotto

将项目代码下载到本地后，切换代码目录：

```shell
cd ${project_path}/cmd/layotto
```

构建:

```shell @if.not.exist layotto
go build -o layotto
```

完成后目录下会生成layotto文件，运行它：

```shell @background
./layotto start -c ../../configs/config_redis.json
```

<!-- tabs:end -->

### step 3. 运行Publisher程序，调用Layotto发布事件
<!-- tabs:start -->
#### **Go**
编译 golang 写的 publisher:

```shell
 cd ${project_path}/demo/pubsub/client/
 go build -o publisher
 ./publisher -s pub_subs_demo
```

#### **Java**

下载 java sdk 和 examples:

```shell @if.not.exist java-sdk
git clone https://github.com/layotto/java-sdk
```

切换目录:

```shell
cd java-sdk
```

构建:

```shell @if.not.exist examples-pubsub-publisher/target/examples-pubsub-publisher-jar-with-dependencies.jar
# build example jar
mvn -f examples-pubsub-publisher/pom.xml clean package
```

运行:

```shell
# run the example
java -jar examples-pubsub-publisher/target/examples-pubsub-publisher-jar-with-dependencies.jar
```


<!-- tabs:end -->

打印出如下信息则代表调用成功：

```bash
Published a new event.Topic: hello ,Data: world
Published a new event.Topic: topic1 ,Data: value1
```

### step 4. 检查Subscriber收到的事件消息

回到subscriber的命令行，会看到接收到了新消息：

```bash
Start listening on port 9999 ......
Received a new event.Topic: topic1 , Data: value1
Received a new event.Topic: hello , Data: world
```

### step 5. 销毁容器，释放资源
<!-- tabs:start -->
#### **关闭 Docker Compose**
如果您是用 docker-compose 启动的 Redis 和 Layotto，可以按以下方式关闭：

```bash
cd ${project_path}/docker/layotto-redis
docker-compose stop
```

#### **销毁 Redis Docker 容器**
如果您是用 Docker 启动的 Redis，可以按以下方式销毁 Redis 容器：

```shell
docker rm -f redis-test
```

<!-- tabs:end -->


### 下一步
#### 这个Publisher程序做了什么？
示例Publisher程序中使用了Layotto提供的golang版本sdk，调用Layotto Pub/Sub API,发布事件到redis。随后Layotto监听到redis有新事件，将新事件回调Subscriber程序开放的接口，通知Subscriber。

sdk位于`sdk`目录下，用户可以通过sdk调用Layotto提供的API。

除了使用sdk，您也可以用任何您喜欢的语言、通过grpc直接和Layotto交互。

其实sdk只是对grpc很薄的封装，用sdk约等于直接用grpc调。


#### 细节以后再说，继续体验其他API
通过左侧的导航栏，继续体验别的API吧！

#### 了解Pub/Sub API实现原理

如果您对实现原理感兴趣，或者想扩展一些功能，可以阅读[Pub/Sub API的设计文档](zh/design/pubsub/pubsub-api-and-compability-with-dapr-component.md)