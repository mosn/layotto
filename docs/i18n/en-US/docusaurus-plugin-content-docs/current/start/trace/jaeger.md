# Jaeger trace

## Configuration

Example：configs/config_trace_jaeger.json

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

| Fields                                  | Required | Note                                                                                                                                                  |
| --------------------------------------- | -------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- |
| service_name       | Y        | Service Name                                                                                                                                          |
| agent_host         | N        | agent component port                                                                                                                                  |
| Strategy                                | N        | Data reporting, default using collector method. Optional configuration values include `collector` and `agent`                         |
| collector_endpoint | N        | port number for collector, default http:///127.0.0.1:14268/api/traces |

## Run Jaeger

```shell
cd ${project_path}/diagnostics/jaeger

docker-compose -f jaeger-docker-compose.yaml up -d
```

## Run layotto

A layoto's server： can be started as follows.

Switch directory:

```shell
cd ${project_path}/cmd/layotto_multiple_api
```

Build:

```shell @if.not.exist layotto
go build -o layotto
```

Run:

```shell @background
./layotto start -c ../../configs/config_trace_jaeger.json 
```

## Run Demo

The corresponding call end code is in[client.go](https://github.com/mosn/layotto/blob/main/demo/flowcontrol/client.go), which runs the Sayhello interface with layotto：

Switch directory:

```shell
 cd ${project_path}/demo/flowcontrol/
```

Build:

```shell @if.not.exist client 
 go build -o customer
```

Run:

```shell
./client
```

Visit http://localhost:16686

![img.png](https://gw.alipayobjects.com/mdn/rms_5891a1/afts/img/A*-f2LSLAR9YMAAAAAAAAAAAAAARQnAQ)

## Clean up resources

```shell
cd ${project_path}/diagnostics/jaeger

docker-compose -f jaeger-docker-compose.yaml down
```
