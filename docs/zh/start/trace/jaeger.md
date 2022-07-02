# Jaeger trace 接入

## 配置

示例：configs/config_trace_jaeger.json

```json
{
  "tracing": {
    "enable": true,
    "driver": "jaeger",
    "config": {
      "service_name": "layotto"
    }
  }
}
```

| 字段           | 必填 | 说明                                                   |
|--------------|----|------------------------------------------------------|
| service_name | Y  | 服务名称                                                 |
| agent_host   | N  | agent组件端口                                            |
| strategy     | N  | 数据上报方式，默认使用 collector 方式. 可选的配置值有`collector`和`agent` |
|collector_endpoint | N  | collector的端口号，默认http://127.0.0.1:14268/api/traces    |

## 运行Jaeger

```shell
cd ${project_path}/diagnostics/jaeger

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

![img.png](https://gw.alipayobjects.com/mdn/rms_5891a1/afts/img/A*-f2LSLAR9YMAAAAAAAAAAAAAARQnAQ)


## 清理资源

```shell
cd ${project_path}/diagnostics/jaeger

docker-compose -f jaeger-docker-compose.yaml down
```