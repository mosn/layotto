# Zookeeper

## 配置项说明

示例：configs/config_lock_zookeeper.json

| 字段 | 必填 | 说明 |
| --- | --- | --- |
| zookeeperHosts | Y | zookeeper服务器地址,支持配置zk集群, 例如: 127.0.0.1:2181;127.0.0.2:2181 |
| zookeeperPassword | Y | zookeeper password|
| sessionTimeout | N | 会话的超时时间,单位秒,同zookeeper的sessionTimeout|
|logInfo|N|true会打印zookeeper操作的所有信息，false只会打印zookeeper的错误信息|

## 怎么启动Zookeeper

如果想启动zookeeper的demo，需要先用Docker启动一个Zookeeper 命令：

```shell
docker pull zookeeper
docker run --privileged=true -d --name zookeeper --publish 2181:2181  -d zookeeper:latest
```
