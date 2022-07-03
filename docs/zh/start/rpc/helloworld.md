# Hello World

## 快速开始
![](https://user-images.githubusercontent.com/26001097/148895424-b286feb5-a122-4fe5-9012-0c235f16b9c7.png)

### step 1. 编译运行layotto
下载 layotto 源码后，切换目录:

```shell
cd ${project_path}/cmd/layotto
```

构建:

```shell @if.not.exist layotto
go build -o layotto
```

运行：

```shell @background
./layotto -c ../../demo/rpc/http/example.json
```

### step 2. 启动echoserver服务端

```shell @background
go run ${project_path}/demo/rpc/http/echoserver/echoserver.go
```

### step 3. 通过GPRC接口发起调用

```shell
go run ${project_path}/demo/rpc/http/echoclient/echoclient.go -d 'hello layotto'
```

![rpchello.png](../../../img/rpc/rpchello.png)

#### 解释

1. example.json配置文件中, 利用mosn的路由能力，将http header中id字段等于HelloService:1.0的请求，转发到本地8889端口
2. echoserver会listen本地的8889端口
3. echoclient中会发起GRPC请求到layotto，
