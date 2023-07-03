# 使用Configuration API调用配置中心

## 什么是Configuration API
应用程序在启动和运行的时候会去读取一些「配置信息」，比如：数据库连接参数、启动参数、接口的超时时间、应用程序的端口等。「配置」，基本上伴随着应用程序的整个生命周期。

应用演进到微服务架构后，会部署到很多台机器上，而配置也分散在集群各个机器上，难于管理。于是就出现了「配置中心」，集中管理配置，同时也解决一些新的问题，比如：版本管理（为了支持回滚），权限管理等。

常用的配置中心有很多，例如Spring Cloud Config，Apollo，Nacos，而且云厂商经常会提供自己的配置管理服务，例如AWS Parameter Store,Google RuntimeConfig

不幸的是，这些配置中心的API都不一样。当应用想要跨云部署，或者想要移植（比如从阿里云搬到腾讯云），应用需要重构代码。

Layotto Configuration API的设计目标是定义一套统一的配置中心API，应用只需要关心API、不需要关心具体用的哪个配置中心，让应用能够随意移植，让应用足够"云原生"。

## Configuration API和State API的区别是？
Q: 为啥要单独搞个Configuration API？和 State API主要区别是？感觉两者差不多

A: Configuration会有一些特殊的能力，比如sidecar做配置缓存，比如app订阅配置变更的消息，比如configuration有一些特殊的schema (tag,version,namespace之类的）

这就像配置中心和数据库的区别，都是存储，但是前者领域特定，有特殊功能

## 快速入门
- [使用Apollo配置中心](zh/start/configuration/start-apollo.md)
- [使用Etcd配置中心](zh/start/configuration/start.md)
- [使用nacos配置中心](zh/start/confguration/start-nacos.md)