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
````shell
./etcd
````

默认监听地址为 `localhost:2379`

### 第二步：运行Layotto

将项目代码下载到本地后，切换代码目录、编译：

```bash
cd ${projectpath}/cmd/layotto
go build
```

完成后目录下会生成layotto文件，运行它：

```bash
./layotto start -c ../../configs/config_sequencer_etcd.json
```

### 第三步：运行客户端程序，调用Layotto进行增删改查

```bash
 cd ${projectpath}/demo/sequencer/common/
 go build -o client
 ./client -s "etcd"
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
#### 使用sdk或者grpc客户端
示例客户端程序中使用了Layotto提供的golang版本sdk，sdk位于`sdk`目录下，用户可以通过对应的sdk直接调用Layotto提供的服务。

除了使用sdk，您也可以用任何您喜欢的语言、通过grpc直接和Layotto交互
