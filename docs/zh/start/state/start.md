# 使用State API进行状态管理
## 什么是State API
State API是一套对Key/Value数据进行增删改查的API。您的应用程序可以使用同一套State API操作不同的数据库（或某种存储系统），对Key/Value模型的数据进行增删改查。

API支持批量CRUD操作，支持声明对并发安全和数据一致性的要求，由Layotto帮您处理复杂的并发安全和数据一致性问题。

## 快速开始

该示例展示了如何通过Layotto调用 Redis，进行状态数据的增删改查。

该示例的架构如下图，启动的进程有：Redis、Layotto、客户端程程序

![img.png](../../../img/state/img.png)

### step 1. 启动 Redis 和 Layotto
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

#### step 1.1. 用 Docker 运行 Redis

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

#### step 1.2. 编译、运行 Layotto

将项目代码下载到本地后，切换代码目录：

```shell
# change directory to ${your project path}/cmd/layotto
cd cmd/layotto
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

### step 2. 运行客户端程序，调用Layotto进行增删改查
<!-- tabs:start -->
#### **Go**

构建、运行 go 语言 demo:

```shell
# open a new terminal tab
# change directory to ${your project path}/demo/state/redis/
 cd ${project_path}/demo/state/common/
 go build -o client
 ./client -s "state_demo"
```

打印出如下信息则代表调用成功：

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

#### **Java**

Download java sdk and examples:

```shell @if.not.exist java-sdk
git clone https://github.com/layotto/java-sdk
```

切换目录:

```shell
cd java-sdk
```

构建:

```shell @if.not.exist examples-state/target/examples-state-jar-with-dependencies.jar
# build example jar
mvn -f examples-state/pom.xml clean package
```

运行:

```
java -jar examples-state/target/examples-state-jar-with-dependencies.jar
```

打印出以下信息说明运行成功:

```bash
SaveState succeeded.key:key1 , value: v11
GetState succeeded. key:key1  value:v11
DeleteState succeeded. key:key1
GetState after delete. key:key1  value:
SaveBulkState succeeded. key:key1 , key2
GetBulkState succeeded. key:key2
```

<!-- tabs:end -->

### step 3. 销毁容器，释放资源
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
#### 这个客户端程序做了什么？
示例客户端程序中使用了Layotto提供的golang版本sdk，调用Layotto 的State API进行增删改查。

sdk位于`sdk`目录下，用户可以通过sdk调用Layotto提供的API。

除了使用sdk，您也可以用任何您喜欢的语言、通过grpc直接和Layotto交互。

其实sdk只是对grpc很薄的封装，用sdk约等于直接用grpc调。

#### 想要详细了解State API?
State API是干嘛的，解决什么问题，我应该在什么场景使用它？

如果您产生了这样的困惑，想要了解State API的更多细节，可以进一步阅读[State API使用文档](zh/building_blocks/state/reference) 

#### 细节以后再说，继续体验其他API
通过左侧的导航栏，继续体验别的API吧！