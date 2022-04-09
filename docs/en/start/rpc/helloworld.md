# Hello World

## Quick Start
![](https://user-images.githubusercontent.com/26001097/148895424-b286feb5-a122-4fe5-9012-0c235f16b9c7.png)

### step 1. compile and start layotto
build:
```shell
go build -o layotto cmd/layotto/main.go
```

run:
```shell background
./layotto -c demo/rpc/http/example.json
```

### step 2. start echoserver
```shell background
go run demo/rpc/http/echoserver/echoserver.go
```

### step 3. call runtime InvokerService api.
```shell
go run demo/rpc/http/echoclient/echoclient.go -d 'hello layotto'
```

![rpchello.png](../../../img/rpc/rpchello.png)

## Explanation

1. configure mosn to match http request header id equals HelloService:1.0, forward to localhost:8889
2. echoserver listen at localhost:8889
3. echoclient call the InvokeService grpc api.
