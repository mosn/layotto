# 使用Configuration API调用Etcd配置中心

## 什么是Configuration API
应用程序在启动和运行的时候会去读取一些「配置信息」，比如：数据库连接参数、启动参数、接口的超时时间、应用程序的端口等。「配置」，基本上伴随着应用程序的整个生命周期。

应用演进到微服务架构后，会部署到很多台机器上，而配置也分散在集群各个机器上，难于管理。于是就出现了「配置中心」，集中管理配置，同时也解决一些新的问题，比如：版本管理（为了支持回滚），权限管理等。

常用的配置中心有很多，例如Spring Cloud Config，Apollo，Nacos，而且云厂商经常会提供自己的配置管理服务，例如AWS Parameter Store,Google RuntimeConfig

不幸的是，这些配置中心的API都不一样。当应用想要跨云部署，或者想要移植（比如从阿里云搬到腾讯云），应用需要重构代码。

Layotto Configuration API的设计目标是定义一套统一的配置中心API，应用只需要关心API、不需要关心具体用的哪个配置中心，让应用能够随意移植，让应用足够"云原生"。

## 快速开始

该示例展示了如何通过Layotto，对etcd配置中心进行增删改查以及watch的过程。请提前在本机上安装[Docker](https://www.docker.com/get-started) 软件。
[config文件](https://github.com/mosn/layotto/blob/main/configs/runtime_config.json) 在config_stores中定义了etcd，用户可以更改配置文件为自己想要的配置中心（目前支持etcd和apollo）。


### 生成镜像

首先请确认把layotto项目放在如下目录：

```
$GOPATH/src/github/layotto/layotto
```

然后执行如下命令：

```bash
cd $GOPATH/src/github/layotto/layotto  
make image
```

运行结束后本地会生成两个镜像：

```bash

xxx@B-P59QMD6R-2102 img % docker images
REPOSITORY          TAG                 IMAGE ID            CREATED             SIZE
layotto/layotto     0.1.0-662eab0       0370527a51a1        10 minutes ago      431MB
```

### 运行Layotto

```bash
docker run -p 34904:34904 layotto/layotto:0.1.0-662eab0
```

Mac和Windows不支持--net=host, 如果是在linux上可以直接把 -p 34904:34904 替换成 --net=host。


### 启动本地client

```bash
cd layotto/demo/configuration/etcd
go build
./etcd
```

打印出如下信息则代表启动完成：

```bash
runtime client initializing for: 127.0.0.1:34904
receive hello response: greeting
get configuration after save, &{Key:hello1 Content:world1 Group:default Label:default Tags:map[] Metadata:map[]}
get configuration after save, &{Key:hello2 Content:world2 Group:default Label:default Tags:map[] Metadata:map[]}
receive watch event, &{Key:hello1 Content:world1 Group:default Label:default Tags:map[] Metadata:map[]}
receive watch event, &{Key:hello1 Content: Group:default Label:default Tags:map[] Metadata:map[]}
```

### 拓展

Layotto 提供了golang版本的sdk，位于runtime/sdk目录下，用户可以通过对应的sdk直接调用Layotto提供的服务。

