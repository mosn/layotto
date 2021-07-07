# 通过Layotto调用redis，进行状态管理

## 快速开始

该示例展示了如何通过Layotto调用redis，进行状态数据的增删改查。

该示例的架构如下图，启动的进程有：redis、Layotto、客户端程程序

![img.png](https://raw.githubusercontent.com/seeflood/layotto/main/docs/img/state/img.png)

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
![img.png](https://raw.githubusercontent.com/seeflood/layotto/main/docs/img/mq/start/img.png)

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
./layotto start -c ../../configs/config_state_redis.json
```

### 运行客户端程序，调用Layotto进行增删改查

```bash
 cd ${projectpath}/demo/state/redis/
 go build -o client
 ./client
```

打印出如下信息则代表调用成功：

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

### 下一步
#### 使用sdk或者grpc客户端
示例客户端程序中使用了Layotto提供的golang版本sdk，sdk位于`sdk`目录下，用户可以通过对应的sdk直接调用Layotto提供的服务。

除了使用sdk，您也可以用任何您喜欢的语言、通过grpc直接和Layotto交互