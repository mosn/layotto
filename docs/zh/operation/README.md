# 如何部署、升级 Layotto
## 1. 如何部署 Layotto
有几种部署 Layotto 的方式

- 直接部署
- 使用 Docker 部署  
- 通过 Kubernetes 部署

### 直接部署
您可以直接运行 Layotto 二进制文件、启动 Layotto。参考[快速开始](zh/start/)

### 使用 Docker 部署
Layotto 提供了官方 Docker 镜像，包括：
- [layotto/layotto](https://hub.docker.com/repository/docker/layotto/layotto)
- [layotto/layotto.arm64](https://hub.docker.com/repository/docker/layotto/layotto.arm64)

镜像内不包含 `config.json` 配置文件，您可以将自己的配置文件挂载进镜像的`/runtime/configs/config.json`目录, 然后启动镜像。例如：

```shell
docker run -v "$(pwd)/configs/config.json:/runtime/configs/config.json" -d  -p 34904:34904 --name layotto layotto/layotto start
```

您也可以通过 docker-compose 同时启动 Layotto 和 其他系统（比如 Redis)，参考[快速开始](zh/start/state/start?id=step-1-%e5%90%af%e5%8a%a8-redis-%e5%92%8c-layotto)

### 在 Kubernetes 集群中部署
#### 方案1. 通过 Istio 部署
如果您是 Istio 用户，可以通过 Istio 部署 Layotto Sidecar。

为了集成 Istio，您需要将 Layotto Sidecar 构建成 proxyv2 镜像。有以下几种构建方法：

- 可以参考 [MOSN 的教程](https://mosn.io/docs/user-guide/start/istio/). 把教程中的 MOSN 镜像换成 Layotto 镜像即可。

- 为了让构建过程更方便些，Layotto 开发了 make 脚本。按需修改好 main 包(在 `cmd/layotto` 目录下)后，敲以下命令即可打包成 proxyv2镜像:

```shell
make image.proxyv2.build VERSION=latest
```

这个脚本的原理是：把 MOSN proxyv2 镜像中的二进制文件替换成 Layotto 的二进制文件。

- 社区构建好了[layotto/proxyv2](https://hub.docker.com/r/layotto/proxyv2) 镜像，用于[线上实验室](https://killercoda.com/mosn-tutorial/course/layotto) 演示 Istio+Layotto。这个镜像把 Layotto 所有组件都打包进去了，导致特别大。我们建议您还是按照自己的需求修改 main、自己构建镜像，这样构建出来的镜像会小很多。 

#### 方案2. 其他方式
您可以准备自己的镜像、K8s 配置文件，通过 Kubernetes 部署 Layotto.

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

- 使用 K8s 原生方案升级 sidecar 容器
  
- [平滑升级，自动迁移长连接](https://mosn.io/docs/concept/smooth-upgrade/)

平滑升级的好处是能做到自动迁移长连接，详细介绍可以点击上述文档查看。

实现平滑升级的方案有：

- MOSN 对 SIGHUP 做了监听，发送 SIGHUP 信号给 MOSN 进程，通过 ForkExec 生成一个新的 MOSN 进程。
  
- 直接重新启动一个新 MOSN 进程。容器间升级需要 Operator 的支持。

- [OpenKruise 做原地热升级的分享](https://mosn.io/blog/posts/mosn-sidecarset-hotupgrade/)