# Hello World

## Quick Start
![](https://user-images.githubusercontent.com/26001097/148895424-b286feb5-a122-4fe5-9012-0c235f16b9c7.png)

1. compile and start layotto
```sh
go build -o layotto cmd/layotto/main.go
./layotto -c demo/rpc/http/example.json
```

2. start echoserver
```sh
go run demo/rpc/http/echoserver/echoserver.go
```

3. call runtime InvokerService api.
```sh
go run demo/rpc/http/echoclient/echoclient.go -d 'hello layotto'
```

![rpchello.png](../../../img/rpc/rpchello.png)

## Explanation

1. configure mosn to match http request header id equals HelloService:1.0, forward to localhost:8889
2. echoserver listen at localhost:8889
3. echoclient call the InvokeService grpc api.
