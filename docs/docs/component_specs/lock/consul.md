# Consul

## 配置项说明

示例：configs/config_consul.json

| 字段 | 必填 | 说明 |
| --- | --- | --- |
| address | Y | consul服务器地址,例如localhost:8500 |
| scheme | Y | 客户端连接模式,HTTP/HTTPS |
| username | N | 指定用户名 |
| password | N | 指定密码 |

## 怎么启动Consul

如果想启动Consul的demo，需要先用Docker启动一个Consul 命令：

```shell
docker run --name consul -d -p 8500:8500 consul
```

