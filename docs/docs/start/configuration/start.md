# 使用Configuration API调用Etcd配置中心
本示例展示了使用 etcd 作为配置中心时，如何通过Layotto，对 etcd 配置中心进行增删改查以及 watch。

本示例架构如下图，启动的进程有：客户端程程序、Layotto、etcd 。

![](https://gw.alipayobjects.com/mdn/rms_5891a1/afts/img/A*dzGaSb78UCoAAAAAAAAAAAAAARQnAQ)
## step 1. 启动 etcd 和 Layotto
<!-- tabs:start -->
### **使用 Docker Compose**
您可以使用 docker-compose 启动 etcd 和 Layotto

```bash
cd docker/layotto-etcd
# Start etcd and layotto with docker-compose
docker-compose up -d
```

### **本地编译（不适合 Windows)**
您可以使用 Docker 运行 etcd，然后本地编译、运行 Layotto。

> [!TIP|label: 不适合 Windows 用户]
> Layotto 在 Windows 下会编译失败。建议 Windows 用户使用 docker-compose 部署
### step 1.1 启动 etcd

etcd的启动方式可以参考etcd的[官方文档](https://etcd.io/docs/v3.5/quickstart/)

简单说明：

访问 https://github.com/etcd-io/etcd/releases 下载对应操作系统的 etcd（也可以用 docker，但是下载官方编译好的 etcd 更简单）

例如，如果是 macOS amd64 用户，可以点击下载：

![](https://gw.alipayobjects.com/mdn/rms_5891a1/afts/img/A*sc_HQaMXg4YAAAAAAAAAAAAAARQnAQ)

下载完成执行命令启动：

```shell @background
./etcd
```

默认监听地址为 `localhost:2379`

### step 1.2 启动 layotto

```shell
cd ${project_path}/cmd/layotto
```

构建:

```shell @if.not.exist layotto
go build -o layotto
```

编译成功后执行:

```shell @background
./layotto start -c ../../configs/runtime_config.json
```

> 解释：[runtime_config.json](https://github.com/mosn/layotto/blob/main/configs/runtime_config.json) 是 Layotto 的配置文件，它在 `config_store` 中声明了使用 etcd 作为配置中心。用户可以更改配置文件，改成使用自己想要用的其他配置中心（目前支持 etcd 和 apollo）。
<!-- tabs:end -->

## step 2. 启动测试demo

```shell
 cd ${project_path}/demo/configuration/common
```

```shell @if.not.exist client
 go build -o client
```

```shell
 ./client -s "config_demo"
```

打印出如下信息则代表启动完成：

```bash
save key success
get configuration after save, &{Key:key1 Content:value1 Group:application Label:prod Tags:map[feature:print release:1.0.0] Metadata:map[]} 
get configuration after save, &{Key:haha Content:heihei Group:application Label:prod Tags:map[feature:haha release:1.0.0] Metadata:map[]} 
delete keys success
write start
receive subscribe resp store_name:"config_demo" app_id:"apollo" items:<key:"heihei" content:"heihei1" group:"application" label:"prod" tags:<key:"feature" value:"haha" > tags:<key:"release" value:"16" > >
```

## step 3.销毁容器,释放资源
<!-- tabs:start -->
### **关闭 Docker Compose**
如果您是用 docker-compose 启动的 etcd 和 Layotto，可以按以下方式关闭：

```bash
cd ${project_path}/docker/layotto-etcd
docker-compose stop
```

### **销毁 etcd Docker 容器**
如果您是用 Docker 启动的 etcd，可以按以下方式销毁 etcd 容器：

```shell
docker rm -f etcd
```
<!-- tabs:end -->

## 下一步
### 这个客户端Demo做了什么？
示例客户端程序中使用了Layotto提供的golang版本sdk，调用Layotto 的Configuration API对配置数据进行增删改查、订阅变更。

sdk位于`sdk`目录下，用户可以通过sdk调用Layotto提供的API。

除了使用sdk，您也可以用任何您喜欢的语言、通过grpc直接和Layotto交互。

其实sdk只是对grpc很薄的封装，用sdk约等于直接用grpc调。


### 细节以后再说，继续体验其他API
通过左侧的导航栏，继续体验别的API吧！
