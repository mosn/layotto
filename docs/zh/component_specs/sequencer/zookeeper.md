# Zookeeper

## 配置项说明

示例：configs/config_zookeeper.json

| 字段 | 必填 | 说明 |
| --- | --- | --- |
| zookeeperHosts | Y | zookeeper服务器地址,支持配置zk集群, 例如: 127.0.0.1:2181;127.0.0.2:2181 |
| zookeeperPassword | Y | zookeeper password|
| sessionTimeout | N | 会话的超时时间,单位秒,同zookeeper的sessionTimeout|
|logInfo|N|true会打印zookeeper操作的所有信息，false只会打印zookeeper的错误信息|

## 警告
zookeeper的自增id组件使用zk的version实现, version不能超过int32(虽然我们的sequencer API设计成返回int64),超过会溢出。当GetNextId方法发生溢出时，会产生error且打印错误日志，除此之外不会做任何处理。

建议您监控zookeeper中的version，避免溢出发生

## 怎么启动Zookeeper

如果想启动zookeeper的demo，需要先用Docker启动一个Zookeeper 命令：

```shell
docker pull zookeeper
docker run --privileged=true -d --name zookeeper --publish 2181:2181  -d zookeeper:latest
```

## 启动 layotto

````shell
cd ${project_path}/cmd/layotto
go build
````

> 如果 build 报错，可以在项目根目录执行 `go mod vendor`

编译成功后执行:

````shell
./layotto start -c ../../configs/config_zookeeper.json
````

## 运行 Demo

````shell
cd ${project_path}/demo/losequencerck/zookeeper/
 go build -o client
 ./client
````
