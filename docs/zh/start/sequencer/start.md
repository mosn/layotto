# 使用Sequencer API生成分布式唯一、自增id
## 什么是Sequencer API
Sequencer API用于生成分布式唯一、自增id。

Sequencer API支持声明对自增的需求，包括趋势递增(WEAK)和严格递增(STRONG)

## 快速开始

该示例展示了如何通过Layotto调用Etcd，生成分布式唯一、自增id。

该示例的架构如下图，启动的进程有：Etcd、Layotto、客户端程程序

![img.png](../../../img/sequencer/etcd/img.png)

### 第一步：部署存储系统（Etcd）

etcd的启动方式可以参考etcd的[官方文档](https://etcd.io/docs/v3.5/quickstart/)

简单说明：

访问 https://github.com/etcd-io/etcd/releases 下载对应操作系统的 etcd（也可用 docker）

下载完成执行命令启动：

```shell @background
./etcd
```

默认监听地址为 `localhost:2379`

### 第二步：运行Layotto

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

### 第三步：运行客户端程序，调用Layotto生成唯一id

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

### 下一步
#### 这个客户端程序做了什么？
示例客户端程序中使用了Layotto提供的golang版本sdk，调用Layotto Sequencer API,生成分布式唯一、自增id。

sdk位于`sdk`目录下，用户可以通过sdk调用Layotto提供的API。

除了使用sdk，您也可以用任何您喜欢的语言、通过grpc直接和Layotto交互。

其实sdk只是对grpc很薄的封装，用sdk约等于直接用grpc调。


#### 想要详细了解Sequencer API?
Sequencer API是干嘛的，解决什么问题，我应该在什么场景使用它？

如果您产生了这样的困惑，想要了解Sequencer API的更多使用细节，可以进一步阅读[Sequencer API使用文档](zh/api_reference/sequencer/reference)

#### 细节以后再说，继续体验其他API
通过左侧的导航栏，继续体验别的API吧！


#### 了解Sequencer API的实现原理

如果您对实现原理感兴趣，或者想扩展一些功能，可以阅读[Sequencer API的设计文档](zh/design/sequencer/design.md)