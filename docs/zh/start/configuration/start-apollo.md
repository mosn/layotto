# 使用Configuration API调用apollo配置中心

## 什么是Configuration API
应用程序在启动和运行的时候会去读取一些「配置信息」，比如：数据库连接参数、启动参数、接口的超时时间、应用程序的端口等。「配置」，基本上伴随着应用程序的整个生命周期。

应用演进到微服务架构后，会部署到很多台机器上，而配置也分散在集群各个机器上，难于管理。于是就出现了「配置中心」，集中管理配置，同时也解决一些新的问题，比如：版本管理（为了支持回滚），权限管理等。

常用的配置中心有很多，例如Spring Cloud Config，Apollo，Nacos，而且云厂商经常会提供自己的配置管理服务，例如AWS Parameter Store,Google RuntimeConfig

不幸的是，这些配置中心的API都不一样。当应用想要跨云部署，或者想要移植（比如从腾讯云搬到阿里云），应用需要重构代码。

Layotto Configuration API的设计目标是定义一套统一的配置中心API，应用只需要关心API、不需要关心具体用的哪个配置中心，让应用能够随意移植，让应用足够"云原生"。

## 快速开始

该示例展示了如何通过Layotto，对apollo配置中心进行增删改查以及watch的过程。

该示例的架构如下图，启动的进程有：客户端程程序、Layotto、Apollo服务器

![img.png](../../../img/configuration/apollo/arch.png)

### 部署apollo配置中心并修改Layotto（可选）

您可以跳过这一步，使用本demo无需自己部署apollo服务器。本demo会使用[apollo官方](https://github.com/ctripcorp/apollo) 提供的演示环境http://106.54.227.205/

如果您自己部署了apollo，可以修改Layotto的[config文件](https://github.com/mosn/layotto/blob/main/configs/config_apollo.json) ，将apollo服务器地址改成您自己的。

### 运行Layotto server 端

将项目代码下载到本地后，切换代码目录、编译：

```bash
cd ${projectpath}/cmd/layotto
go build
#备注 如果发现构建失败无法下载,请进行如先设置
go env -w GOPROXY="https://goproxy.cn,direct"
```

完成后目录下会生成layotto文件，运行它：

```bash
./layotto start -c ../../configs/config_apollo.json
```

### 启动客户端Demo，调用Layotto增删改查

```bash
 cd ${projectpath}/demo/configuration/apollo
 go build -o apolloClientDemo
 ./apolloClientDemo
```

打印出如下信息则代表调用成功：

```bash
save key success
get configuration after save, &{Key:key1 Content:value1 Group:application Label:prod Tags:map[feature:print release:1.0.0] Metadata:map[]} 
get configuration after save, &{Key:haha Content:heihei Group:application Label:prod Tags:map[feature:haha release:1.0.0] Metadata:map[]} 
delete keys success
write start
receive subscribe resp store_name:"apollo" app_id:"apollo" items:<key:"heihei" content:"heihei1" group:"application" label:"prod" tags:<key:"feature" value:"haha" > tags:<key:"release" value:"16" > >
```

### 下一步

示例客户端Demo中使用了Layotto提供的golang版本sdk，sdk位于`sdk`目录下，用户可以通过对应的sdk直接调用Layotto提供的服务。

除了使用sdk，您也可以用任何您喜欢的语言、通过grpc直接和Layotto交互
