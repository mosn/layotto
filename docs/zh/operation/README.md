# 如何部署、升级 Layotto
## 1. 如何部署 Layotto
有两种部署 Layotto 的方式

- 直接部署
  
- 通过 Kubernetes 部署

### 直接部署
您可以直接运行 Layotto 二进制文件、启动 Layotto。参考[快速开始](zh/start)

### 通过 Kubernetes 部署
您可以准备自己的镜像、k8s 配置文件，通过 Kubernetes 部署 Layotto.

我们正在开发官方版 Layotto 镜像以及通过 Helm 部署到 Kubernetes 的方案，欢迎加入共建，详见 https://github.com/mosn/layotto/issues/392

## 2. MOSN 用户如何将已有 MOSN 切换成 Layotto 
把 sidecar 镜像里的 MOSN 换成 Layotto 即可。

解释：

Layotto 和 MOSN 跑在同一个进程里，可以理解成:

> Layotto == MOSN + 一个特殊的 grpcFilter 打包到一起

所以： 

> 将 MOSN 换成 Layotto == 将 MOSN 换成 "MOSN + 一个特殊的 grpcFilter"

没有本质区别，只需注意版本，Layotto 和 MOSN 的版本是一一对应的关系。

之前已发布的 Layotto v0.3.0 对应的 MOSN 版本为 v0.24.1

## 3. 如何升级 Layotto
有两种升级方案：

- 使用 k8s 原生方案升级 sidecar 容器
  
- [平滑升级，自动迁移长连接](https://mosn.io/docs/concept/smooth-upgrade/)

平滑升级的好处是能做到自动迁移长连接，详细介绍可以点击上述文档查看。

实现平滑升级的方案有：

- MOSN 对 SIGHUP 做了监听，发送 SIGHUP 信号给 MOSN 进程，通过 ForkExec 生成一个新的 MOSN 进程。
  
- 直接重新启动一个新 MOSN 进程。容器间升级需要 Operator 的支持。

- [OpenKruise 做原地热升级的分享](https://mosn.io/blog/posts/mosn-sidecarset-hotupgrade/)