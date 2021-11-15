## WASM on Layotto

### What is WASM on Layotto?

Layotto supports load the compiled WASM file, and interacts with it through the API of the `proxy_abi_version_0_2_0` version.

### Quick start

1. start redis server and write test data

The example only needs a Redis server that can be used normally. As for where it is installed, there is no special restriction. It can be a virtual machine, a local machine or a server. Here, the installation on mac is used as an example to introduce.

```
> brew install redis
> redis-server /usr/local/etc/redis.conf
```

```
> redis-cli
127.0.0.1:6379> set book1 100
OK
```


2. start Layotto server

```
go build -tags wasmer -o ./layotto ./cmd/layotto/main.go
./layotto start -c ./demo/faas/config.json
```

**Note: You need to modify the redis address as needed, the default address is: localhost:6379**

3. send request

```
curl -H 'id:id_1' 'localhost:2045?name=book1'
There are 100 inventories for book1.
```

### Note

This feature is still in the experimental stage, and the implementation of the WASM interactive API in the community is not uniform enough, so if you have any needs for this module, please post it in the issue area, we will build WASM together!