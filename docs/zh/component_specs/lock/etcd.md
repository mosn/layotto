# Etcd

## 配置项说明
示例：configs/config_lock_etcd.json

| 字段 | 必填 | 说明 |
| --- | --- | --- |
| endpoints | Y | etcd 的服务地址 ip+端口，多个地址使用英文分号（;）分隔 |
| dialTimeout | N | 建立连接超时，单位：秒，默认值：5 |
| username | N | etcd 认证用户名 |
| password | N | etcd 认证密码 |
| keyPrefix | N | 在 etcd 建立锁 key 的前缀，默认值：`/layotto/` |

## 怎么启动 etcd
要先用Docker启动一个etcd
命令：

访问 https://github.com/etcd-io/etcd/releases 下载对应操作系统的 etcd（也可用 docker）

下载完成执行命令启动：
````shell
./etcd
````

默认监听地址为 `localhost:2379`

## 启动 layotto

````shell
cd ${projectpath}/cmd/layotto
go build
````
>如果 build 报错，可以在项目根目录执行 `go mod vendor`

编译成功后执行:
````shell
./layotto start -c ../../configs/config_lock_etcd.json
````

## 运行 Demo

````shell
cd ${projectpath}/demo/lock/etcd/
 go build -o client
 ./client
````

