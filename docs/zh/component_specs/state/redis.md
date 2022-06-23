# Redis

## 配置项说明
示例：configs/config_redis.json

| 字段 | 必填 | 说明 |
| --- | --- | --- |
| redisHost | Y | redis服务器地址,例如localhost:6380 |
| redisPassword | Y | redis密码 |

## 怎么启动Redis
如果想启动redis的demo，需要先用Docker启动一个Redis
命令：

```shell
docker pull redis:latest
docker run -itd --name redis-test -p 6380:6379 redis
```