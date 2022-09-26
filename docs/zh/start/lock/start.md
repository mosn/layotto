# 快速开始: 基于 Redis 使用分布式锁

该示例展示了如何通过Layotto调用 Redis，进行分布式锁的抢锁、解锁操作。

该示例的架构如下图，启动的进程有：Redis、Layotto、一个演示用的client程序（其中包含两个协程，并发抢锁）

![img.png](../../../img/lock/img.png)
## step 1. 部署 Redis 和 Layotto

<!-- tabs:start -->
### **使用 Docker Compose**
您可以用 docker-compose 启动 Redis 和 Layotto

```bash
cd docker/layotto-redis
# Start redis and layotto with docker-compose
docker-compose up -d
```

### **本地编译（不适合 Windows)**
您可以使用 Docker 运行 Redis，然后本地编译、运行 Layotto。

> [!TIP|label: 不适合 Windows 用户]
> Layotto 在 Windows 下会编译失败。建议 Windows 用户使用 docker-compose 部署

### step 1.1. 用 Docker 运行 Redis

1. 取最新版的 Redis 镜像。
这里我们拉取官方的最新版本的镜像：

```shell
docker pull redis:latest
```

2. 查看本地镜像
使用以下命令来查看是否已安装了 Redis：

```shell
docker images
```

![img.png](../../../img/mq/start/img.png)

3. 运行容器

安装完成后，我们可以使用以下命令来运行 Redis 容器：

```shell
docker run -itd --name redis-test -p 6380:6379 redis
```

参数说明：

-p 6380:6379：映射容器服务的 6379 端口到宿主机的 6380 端口。外部可以直接通过宿主机ip:6380 访问到 Redis 的服务。

### step 1.2. 运行 Layotto

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

## step 2. 运行客户端程序，调用Layotto抢锁/解锁
<!-- tabs:start -->
### **Go**

```shell
 cd ${project_path}/demo/lock/common/
 go build -o client
 ./client -s "lock_demo"
```

打印出如下信息则代表调用成功：

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

### **Java**

下载 java sdk 和示例代码:

```shell @if.not.exist java-sdk
git clone https://github.com/layotto/java-sdk
```

切换目录:

```shell
cd java-sdk
```

构建:

```shell @if.not.exist examples-lock/target/examples-lock-jar-with-dependencies.jar
# build example jar
mvn -f examples-lock/pom.xml clean package
```

运行:

```shell
java -jar examples-lock/target/examples-lock-jar-with-dependencies.jar
```

打印出以下信息说明运行成功:

```bash
TryLockResponse{success=true}
TryLockResponse{success=true}
TryLockResponse{success=true}
UnlockResponse{status=SUCCESS}
TryLockResponse{success=true}
UnlockResponse{status=LOCK_UNEXIST}
```

<!-- tabs:end -->

## 下一步
### 这个客户端程序做了什么？
示例客户端程序中使用了Layotto提供的golang版本sdk，调用Layotto 分布式锁API,启动多个协程进行抢锁、解锁操作。

sdk位于`sdk`目录下，用户可以通过sdk调用Layotto提供的API。

除了使用sdk，您也可以用任何您喜欢的语言、通过grpc直接和Layotto交互。

其实sdk只是对grpc很薄的封装，用sdk约等于直接用grpc调。


### 细节以后再说，继续体验其他API
通过左侧的导航栏，继续体验别的API吧！

### 了解分布式锁 API的实现原理

如果您对实现原理感兴趣，或者想扩展一些功能，可以阅读[分布式锁 API的设计文档](zh/design/lock/lock-api-design.md)