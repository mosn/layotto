# Skywalking trace 接入

## 配置

示例：configs/config_trace_skywalking.json

````json
{
  "tracing": {
    "enable": true,
    "driver": "SkyWalking",
    "config": {
      "reporter": "gRPC",
      "backend_service": "127.0.0.1:11800",
      "service_name": "layotto"
    }
  }
}
````

| 字段 | 必填  | 说明                       |
| --- |-----|--------------------------|
| reporter | Y   | 上报方式 grpc                |
| backend_service | Y   | skywalking oap server 地址 |
| service_name | Y   | 服务名称                     |

## 运行 skywalking

````shell
cd diagnostics/skywalking

docker-compose -f skywalking-docker-compose.yaml up -d
````

## 运行 layotto

````shell
cd cmd/layotto_multiple_api/
go build -o layotto
./layotto start -c ../../configs/config_trace_skywalking.json
````

## 运行 Demo

````shell
cd demo/flowControl
go run client.go
````

访问 http://127.0.0.1:8080

![](../../../img/trace/sky.png)

## 清理资源

````shell
cd diagnostics/skywalking

docker-compose -f skywalking-docker-compose.yaml down
````