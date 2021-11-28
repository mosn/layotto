## 将业务逻辑通过WASM下沉进sidecar
### 功能介绍
service mesh 和 multi-runtime 的 sidecar 是全公司通用的基础设施，但实践中,业务系统也会有自己的sdk，也会有推动用户升级难、版本碎片的问题.

比如某中台系统以jar包形式开发了sdk，供上层业务系统使用。他们的feature不算全公司通用，因此没法说服中间件团队、开发到公司统一的sidecar里。

![img_1.png](../../../img/wasm/img_1.png)

而如果变成这样：

![img.png](../../../img/wasm/img.png)

如果开发者不再开发sdk(jar包），改成开发.wasm文件、支持独立升级部署，就没有推动业务方升级的痛苦了,想要升级的时候在运维平台上操作发布即可，不需要app和sidecar重启

Layotto支持加载编译好的WASM文件，并通过`proxy_abi_version_0_2_0`版本的API与目标WASM进行交互。

### 快速开始

1. 启动redis并写入测试数据

这里只是需要一个可以正常使用 Redis 即可，至于 Redis 安装在哪里没有特别限制，可以是虚拟机里，也可以是本机或者服务器，这里以安装在 mac 为例进行介绍。

```
> brew install redis
> redis-server /usr/local/etc/redis.conf
```

```
> redis-cli
127.0.0.1:6379> set book1 100
OK
```

2. 启动layotto

```
go build -tags wasmer -o ./layotto ./cmd/layotto/main.go
./layotto start -c ./demo/wasm/config.json
```
**注：需要把`./demo/faas/config.json`中的 redis 地址修改为实际地址，默认地址为：localhost:6379。**

3. 发送请求

```
curl -H 'id:id_1' 'localhost:2045?name=book1'
There are 100 inventories for book1.
```

该http请求会访问Layotto中的wasm模块。该wasm模块会调用redis进行逻辑处理

### 说明

该功能目前仍处于试验阶段，社区里对于WASM跟宿主的交互API也不够统一，因此如果您有该模块的需求欢迎发表在issue区，我们一起建设WASM！