# Jaeger trace 接入

## 运行Jaeger

```shell
cd ${project_path}/diagnostics/jaeger/jaeger-docker-compose.yaml

docker-compose -f jaeger-docker-compose.yaml up -d
```

## 运行layotto

可以按照如下方式启动一个layotto的server：

切换目录:

```shell
cd ${project_path}/cmd/layotto_multiple_api
```

构建:

```shell @if.not.exist layotto
go build -o layotto
```

运行:

```shell @background
./layotto start -c ../../configs/config_trace_jaeger.json 
```

## 运行 Demo

对应的调用端代码在[client.go](https://github.com/mosn/layotto/blob/main/demo/flowcontrol/client.go) 中，运行它会调用layotto的SayHello接口：

切换目录:

```shell
 cd ${project_path}/demo/flowcontrol/
``` 

构建:

```shell @if.not.exist client 
 go build -o client
```
运行:

```shell
./client
```

访问 http://localhost:16686

![img.png](../../../img/trace/jaeger.png)


## 清理资源

```shell
cd ${project_path}/diagnostics/jaeger/jaeger-docker-compose.yaml

docker-compose -f jaeger-docker-compose.yaml down
```