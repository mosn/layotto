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

#### step 1. 启动redis并写入测试数据

这里只是需要一个可以正常使用 Redis 即可，至于 Redis 安装在哪里没有特别限制，可以是虚拟机里，也可以是本机或者服务器。

这里以用 Docker 安装 Redis 为例，进行介绍。

启动 Redis 容器:

```shell
docker run -d --name redis-test -p 6379:6379 redis
```

调用 Redis 容器中的 redis-cli,执行`set book1 100` 

```shell
docker exec -i redis-test redis-cli set book1 100
```

这条命令将 key 为 `book1` 的缓存值设置为100。如果返回如下结果,代表 set 成功:

```bash
OK
```

我们可以执行 `get book1`,看看现在 `book1` 对应的值是多少:

```shell
docker exec -i redis-test redis-cli get book1
```

返回:

```bash
"100"
```

符合预期

#### step 2. 启动layotto

构建:

```shell @if.not.exist layotto_wasmer
go build -tags wasmer -o ./layotto_wasmer ./cmd/layotto/main.go
```

运行:

```shell @background
./layotto_wasmer start -c ./demo/faas/config.json
```

**注：需要把`./demo/faas/config.json`中的 redis 地址修改为实际地址，默认地址为：localhost:6379。**

#### step 3. 发送请求

```shell
curl -H 'id:id_1' 'localhost:2045?name=book1'
```

将返回如下结果:

```bash
There are 100 inventories for book1.
```

该http请求会访问 Layotto 中的wasm模块。该wasm模块会调用redis进行逻辑处理

#### step 4. 销毁容器，释放资源

```shell
docker rm -f redis-test
```

### 说明

该功能目前仍处于试验阶段，社区里对于WASM跟宿主的交互API也不够统一，因此如果您有该模块的需求欢迎发表在issue区，我们一起建设WASM！