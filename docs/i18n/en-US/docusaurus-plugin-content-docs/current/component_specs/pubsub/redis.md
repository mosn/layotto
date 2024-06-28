# Redis

## metadata fields
Example: configs/config_redis.json

| Field | Required | Description |
| --- | --- | --- |
| redisHost | Y | redis server address, such as localhost:6380 |
| redisPassword | Y | redis Password |

## How to start Redis
If you want to run the redis demo, you need to start a Redis server with Docker first.

command:

```shell
docker pull redis:latest
docker run -itd --name redis-test -p 6380:6379 redis
```