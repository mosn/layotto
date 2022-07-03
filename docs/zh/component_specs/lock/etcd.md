# Etcd

## 配置项说明
示例：configs/runtime_config.json

| 字段 | 必填 | 说明 |
| --- | --- | --- |
| endpoints | Y | etcd 的服务地址 ip+端口，多个地址使用英文分号（;）分隔 |
| dialTimeout | N | 建立连接超时，单位：秒，默认值：5 |
| username | N | etcd 认证用户名 |
| password | N | etcd 认证密码 |
| keyPrefixPath | N | 在 etcd 建立锁 key 的前缀，默认值：`/layotto/` |
| tlsCert | N | tls 证书路径 |
| tlsCertKey | N | tls 证书 key 路径 |
| tlsCa | N | tls ca 路径 |

## 怎么启动 etcd

etcd的启动方式可以参考etcd的[官方文档](https://etcd.io/docs/v3.5/quickstart/)

简单说明：

访问 https://github.com/etcd-io/etcd/releases 下载对应操作系统的 etcd（也可用 docker）

下载完成执行命令启动：

````shell
./etcd
````

默认监听地址为 `localhost:2379`

## 启动 layotto

````shell
cd ${project_path}/cmd/layotto
go build
````

>如果 build 报错，可以在项目根目录执行 `go mod vendor`

编译成功后执行:

````shell
./layotto start -c ../../configs/runtime_config.json
````

## 运行 Demo

````shell
 cd ${project_path}/demo/lock/common/
 go build -o client
 ./client -s "lock_demo"
````

