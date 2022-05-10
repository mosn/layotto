# Skywalking trace 

## Configuration

Example: configs/config_trace_skywalking.json

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

| Field            | Required fields | Description  |
|------------------|-----|--------------------------|
| reporter         | Y   | Reporting method grpc               |
| backend_service  | Y   | skywalking oap server address |
| service_name     | Y   | Service Name                     |

## Run skywalking

````shell
cd ${project_path}/diagnostics/skywalking

docker-compose -f skywalking-docker-compose.yaml up -d
````

## Run layotto

````shell
cd ${project_path}/cmd/layotto_multiple_api/
go build -o layotto
./layotto start -c ../../configs/config_trace_skywalking.json
````

## Run Demo

````shell
cd ${project_path}/demo/flowcontrol
go run client.go
````

Access http://127.0.0.1:8080

![](../../../img/trace/sky.png)

## Clearance resources

````shell
cd ${project_path}/diagnostics/skywalking

docker-compose -f skywalking-docker-compose.yaml down
````