# 使用Sequencer API生成分布式唯一、自增id
## 什么是Sequencer API
Sequencer API用于生成分布式唯一、自增id。

Sequencer API支持声明对自增的需求，包括趋势递增(WEAK)和严格递增(STRONG)

## 快速开始

该示例展示了如何通过Layotto调用Etcd，生成分布式唯一、自增id。

该示例的架构如下图，启动的进程有：Etcd、Layotto、客户端程程序

![img.png](../../../img/sequencer/etcd/img.png)

### step 1. 启动 etcd 和 Layotto
<!-- tabs:start -->
#### **使用 Docker Compose**
您可以使用 docker-compose 启动 etcd 和 Layotto

```bash
cd docker/layotto-etcd
# Start etcd and layotto with docker-compose
docker-compose up -d
```

#### **本地编译（不适合 Windows)**
您可以使用 Docker 运行 etcd，然后本地编译、运行 Layotto。
> [!TIP|label: 不适合 Windows 用户]
> Layotto 在 Windows 下会编译失败。建议 Windows 用户使用 docker-compose 部署
### step 1.1：部署存储系统（Etcd）

etcd的启动方式可以参考etcd的[官方文档](https://etcd.io/docs/v3.5/quickstart/)

简单说明：

访问 https://github.com/etcd-io/etcd/releases 下载对应操作系统的 etcd（也可用 docker）

下载完成执行命令启动：

```shell @background
./etcd
```

默认监听地址为 `localhost:2379`

### step 1.2：运行Layotto

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
./layotto start -c ../../configs/runtime_config.json
```

<!-- tabs:end -->

### step 2. 运行客户端程序，调用Layotto生成唯一id
<!-- tabs:start -->
#### **Go**

构建、运行 go 语言 demo:

```shell
 cd ${project_path}/demo/sequencer/common/
 go build -o client
 ./client -s "sequencer_demo"
```

打印出如下信息则代表调用成功：

```bash
Try to get next id.Key:key666 
Next id:next_id:1  
Next id:next_id:2  
Next id:next_id:3  
Next id:next_id:4  
Next id:next_id:5  
Next id:next_id:6  
Next id:next_id:7  
Next id:next_id:8  
Next id:next_id:9  
Next id:next_id:10  
Demo success!
```

#### **Java**

下载 java sdk 和示例代码:

```shell @if.not.exist java-sdk
git clone https://github.com/layotto/java-sdk
```

切换目录:

```shell
cd java-sdk
```

构建:

```shell @if.not.exist examples-sequencer/target/examples-sequencer-jar-with-dependencies.jar
# build example jar
mvn -f examples-sequencer/pom.xml clean package
```

运行:

```shell
java -jar examples-sequencer/target/examples-sequencer-jar-with-dependencies.jar
```

打印出以下信息说明运行成功:

```bash
Try to get next id.Key: examples
Next id: 1
Try to get next id.Key: examples
Next id: 2
Try to get next id.Key: examples
Next id: 3
Try to get next id.Key: examples
Next id: 4
Try to get next id.Key: examples
Next id: 5
Try to get next id.Key: examples
Next id: 6
Try to get next id.Key: examples
Next id: 7
Try to get next id.Key: examples
Next id: 8
Try to get next id.Key: examples
Next id: 9
Try to get next id.Key: examples
Next id: 10
```

<!-- tabs:end -->

### step 3.销毁容器,释放资源
<!-- tabs:start -->
#### **关闭 Docker Compose**
如果您是用 docker-compose 启动的 etcd 和 Layotto，可以按以下方式关闭：

```bash
cd ${project_path}/docker/layotto-etcd
docker-compose stop
```

#### **销毁 etcd Docker 容器**
如果您是用 Docker 启动的 etcd，可以按以下方式销毁 etcd 容器：

```shell
docker rm -f etcd
```

<!-- tabs:end -->

### 下一步
#### 这个客户端程序做了什么？
示例客户端程序中使用了Layotto提供的多语言 sdk，调用Layotto Sequencer API,生成分布式唯一、自增id。

go sdk位于`sdk`目录下，java sdk 在 https://github.com/layotto/java-sdk

除了使用sdk调用Layotto提供的API，您也可以用任何您喜欢的语言、通过grpc直接和Layotto交互。

其实sdk只是对grpc很薄的封装，用sdk约等于直接用grpc调。


#### 想要详细了解Sequencer API?
Sequencer API是干嘛的，解决什么问题，我应该在什么场景使用它？

如果您产生了这样的困惑，想要了解Sequencer API的更多使用细节，可以进一步阅读[Sequencer API使用文档](zh/api_reference/sequencer/reference)

#### 细节以后再说，继续体验其他API
通过左侧的导航栏，继续体验别的API吧！


#### 了解Sequencer API的实现原理

如果您对实现原理感兴趣，或者想扩展一些功能，可以阅读[Sequencer API的设计文档](zh/design/sequencer/design.md)