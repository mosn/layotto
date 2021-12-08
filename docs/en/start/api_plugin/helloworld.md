# API plugin: register your own API
This is a demo to show you how to register your own API.

Layotto has the api-plugin feature to let you add your own API based on your need.

## step 0. change directory 
```shell
cd ${projectpath}/cmd/layotto_multiple_api
```

## step 1. start Layotto with a new helloworld API
Build and run Layotto :

```shell
go build -o layotto
# run it
./layotto start -c ../../configs/config_in_memory.json
```

Q: What happened?

Check the code in `main.go` and you will find a new API was registered during startup:

```go
		// register your grpc API here
		runtime.WithGrpcAPI(
			// Layotto API
			l8grpc.NewLayottoAPI,
			// a demo to show how to register your own API
			helloworld_api.NewHelloWorldAPI,
		),
```

## step 2. invoke the helloworld API
```shell
go run client/main.go
```
The result will be:

```shell
Greeting: Hello world
```

This message is the response of the helloworld API you just registered in step 1.

## Next
You can refer to the demo code to implement your own API.

Have a try !