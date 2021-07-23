# 基于redis使用分布式锁

## 快速开始

该示例展示了如何通过Layotto调用redis，进行分布式锁的抢锁、解锁操作。

该示例的架构如下图，启动的进程有：redis、Layotto、一个演示用的client程序（其中包含两个协程，并发抢锁）

![img.png](../../../img/lock/img.png)
### 部署redis

1. 取最新版的 Redis 镜像。
这里我们拉取官方的最新版本的镜像：

```shell
docker pull redis:latest
```

2. 查看本地镜像
   使用以下命令来查看是否已安装了 redis：

```shell
docker images
```
![img.png](../../../img/mq/start/img.png)

3. 运行容器

安装完成后，我们可以使用以下命令来运行 redis 容器：

```shell
docker run -itd --name redis-test -p 6380:6379 redis
```

参数说明：

-p 6380:6379：映射容器服务的 6379 端口到宿主机的 6380 端口。外部可以直接通过宿主机ip:6380 访问到 Redis 的服务。

### 运行Layotto

将项目代码下载到本地后，切换代码目录、编译：

```bash
cd ${projectpath}/cmd/layotto
go build
```

完成后目录下会生成layotto文件，运行它：

```bash
./layotto start -c ../../configs/config_lock_redis.json
```

### 运行客户端程序，调用Layotto抢锁/解锁

```bash
 cd ${projectpath}/demo/lock/redis/
 go build -o client
 ./client
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

### 下一步
#### 使用sdk或者grpc客户端
示例客户端程序中使用了Layotto提供的golang版本sdk，sdk位于`sdk`目录下，用户可以通过对应的sdk直接调用Layotto提供的服务。

除了使用sdk，您也可以用任何您喜欢的语言、通过grpc直接和Layotto交互

#### 了解分布式锁 API的实现原理

如果您对实现原理感兴趣，或者想扩展一些功能，可以阅读[分布式锁 API的设计文档](zh/design/lock/lock-api-design.md)