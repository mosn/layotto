## WASM概述

### 功能介绍

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

### 说明

该功能目前仍处于试验阶段，社区里对于WASM跟宿主的交互API也不够统一，因此如果您有该模块的需求欢迎发表在issue区，我们一起建设WASM！