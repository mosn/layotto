# Redis

## metadata fields
Example: configs/config_sequencer_redis.json

| Field | Required | Description |
| --- | --- | --- |
| redisHost | Y | redis server address, such as localhost:6380 |
| redisPassword | Y | redis Password |
|maxRetries|N| maximum number of retries before giving upy,default value is 3|
|maxRetryBackoff|N|  maximum backoff between each retry,default value is 2s |
|enableTLS |N|  controls whether a client verifies the server's certificate chain and host name,default value is false|

## How to avoid generating duplicate id
Redis components may generate duplicate IDs in the case of data loss. 

In order to avoid data loss and duplicate IDs, you need to use stand-alone redis and [use both persistence methods to get a degree of data safety comparable to what PostgreSQL can provide you.](https://redis.io/topics/persistence)

## How to start Redis
If you want to run the redis demo, you need to start a Redis server with Docker first.

command:
```shell
docker pull redis:latest
docker run -itd --name redis-test -p 6379:6379 redis
```

## Run layotto

````shell
cd ${projectpath}/cmd/layotto
go build
````
>If build reports an error, it can be executed in the root directory of the project `go mod vendor`

Execute after the compilation is successful:
````shell
./layotto start -c ../../configs/config_sequencer_redis.json
````

## Run Demo

````shell
cd ${projectpath}/demo/sequencer/redis/
 go build -o client
 ./client
````