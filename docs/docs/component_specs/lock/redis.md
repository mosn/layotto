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

# Redis集群

## 配置项说明
示例：configs/config_lock_redis_cluter.json

| 字段 | 必填 | 说明 |
| --- | --- | --- |
| redisHosts | Y | redis节点地址,多个节点以','隔开，例如localhost:6381,localhost:6382 |
| redisPassword | Y | redis密码,多个节点共用 |
| concurrency | N | redis集群锁操作中协程池并发度,用于控制并发加锁协程数,默认为cpu核数 |

## 关于集群锁
使用红锁算法实现
配置3个redis节点只能容忍1个节点异常，配置5个节点能容忍两个节点异常，配置1个节点时退化成单点锁，建议使用5个节点

## 如何启动多个redis节点
如果想启动redis集群锁的demo，需要先用Docker启动5个Redis
命令：

```shell
docker pull redis:latest
docker run -itd --name redis1 -p 6381:6379 redis
docker run -itd --name redis2 -p 6382:6379 redis
docker run -itd --name redis3 -p 6383:6379 redis
docker run -itd --name redis4 -p 6384:6379 redis
docker run -itd --name redis5 -p 6385:6379 redis
```
