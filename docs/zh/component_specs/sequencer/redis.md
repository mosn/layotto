# Redis

## 配置项说明
示例：configs/config_sequencer_redis.json

| 字段 | 必填 | 说明 |
| --- | --- | --- |
| redisHost | Y | redis服务器地址,例如localhost:6379 |
| redisPassword | Y | redis密码 |
|maxRetries|N| 放弃前的最大重试次数，默认值为3|
|maxRetryBackoff|N|  每次重试之间的最大退避时间，默认值为2s |
|enableTLS |N| 客户端是否验证服务器的证书链和主机名，默认值为false|

## 如何避免生成重复id
redis组件在丢数据的情况下可能生成重复id，为了避免重复id需要使用单机redis，[需要特殊配置redis服务器，把两种落盘策略都打开、每次写操作都写磁盘](https://redis.io/topics/persistence) 避免丢数据。

## 怎么启动Redis
如果想启动redis的demo，需要先用Docker启动一个Redis
命令：
```shell
docker pull redis:latest
docker run -itd --name redis-test -p 6379:6379 redis
```

## 启动 layotto

````shell
cd ${projectpath}/cmd/layotto
go build
````
>如果 build 报错，可以在项目根目录执行 `go mod vendor`

编译成功后执行:
````shell
./layotto start -c ../../configs/config_sequencer_redis.json
````

## 运行 Demo

````shell
cd ${projectpath}/demo/sequencer/redis/
 go build -o client
 ./client
````