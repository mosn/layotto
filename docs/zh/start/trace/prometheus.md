# Prometheus metrics 接入

## 运行prometheus

window用户需要将prometheus.yml中的layotto改成'docker.for.windows.localhost:34903'

![](https://gw.alipayobjects.com/mdn/rms_5891a1/afts/img/A*mMAeSa8VQ-UAAAAAAAAAAAAAARQnAQ)

```shell
cd ${project_path}/demo/prometheus

docker-compose -f prometheus-docker-compose.yaml up -d
```

## 运行layotto

可以按照如下方式启动一个layotto的server：

切换目录:

```shell
cd ${project_path}/cmd/layotto
```

构建:

```shell @if.not.exist layotto
go build -o layotto
```

运行:

```shell @background
./layotto start -c ../../configs/config_standalone.json
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

访问 http://127.0.0.1:9090/graph?g0.expr=grpc_request_total

![](https://gw.alipayobjects.com/mdn/rms_5891a1/afts/img/A*mEVNSZMvtvEAAAAAAAAAAAAAARQnAQ)


## 清理资源

```shell
cd ${project_path}/demo/prometheus

docker-compose -f prometheus-docker-compose.yaml down
```