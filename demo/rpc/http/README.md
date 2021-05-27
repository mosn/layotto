### Hello World Example

#### Build And Start Runtime
```sh
go build -o layotto cmd/layotto/main.go
./layotto -c demo/rpc/http/example.json
```

#### Start Backend Http Server
```sh
go run demo/rpc/http/echoserver/echoserver.go
```

#### Grpc Invoke
```sh
go run demo/rpc/http/echoclient/echoclient.go -d 'hello layotto'
```

![img.png](img.png)