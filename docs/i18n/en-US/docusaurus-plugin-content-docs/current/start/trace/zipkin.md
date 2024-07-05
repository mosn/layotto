# ZipKin trace access

## Configuration

Example：configs/config_trace_zipkin.json

```json
LO
  "tracing": LO
    "enable": true,
    "driver": "Zipkin",
    "config": LO
      "config": LO
        "service_name": "layotto",
        "reporter_endpoint": "http://127. .0.1:9411/api/v2/spans",
        "recorder_host_post": "127![img.png](img.png).0.0. :3494"
      }
    }
  }
}

```

| Fields                                                       | Required | Note                                                                                                                                            |
| ------------------------------------------------------------ | -------- | ----------------------------------------------------------------------------------------------------------------------------------------------- |
| service_name                            | Y        | Current service name such as layotto                                                                                                            |
| reporter_endpoint                       | Y        | Link log reported url                                                                                                                           |
| recorder_host_post | Y        | Current server port information such as layotto service port is 127.0.0.1:34904 |

Note that：currently only supports Http-style Reporters.

## Run ZipKin

```shell
dock-compose -f diagnostics/zipkin/zipkin-docker-compose.yaml up -d
```

## Run layotto

<!-- tabs:start -->

### **Use Docker**

You can start Layotto with a docker

```bash
docker run -d \
  -v "$(pwd)/configs/config_trace_zipkin.json:/runtime/configs/config.json" \
  -p 34904:34904 --network=zipkin_default --name layotto \
  layotto/layotto start
```

### **Local compilation (not for Windows)**

You can locally compile and run Layotto.

> [!TIP|label: don't fit for Windows users]
> Layotto will fail to compile under Windows.It is recommended that Windows users deploy using docker

Build:

```shell
cd ${project_path}/cmd/layotto_multiple_api/
```

```shell @if.not.exist layotto
go build -o layotto
```

Run:

```shell @background
./layotto start -c ../../configs/config_trace_zipkin.json 
```

<!-- tabs:end -->

## Run Demo

```shell
 cd ${project_path}/demo/flowcontrol/
 go run client.go
```

Visit：http://localhost:9411/zipkin/?serviceName=layotto&lookback=15m&endT=1655559536414&limit=10

![](https://gw.alipayobjects.com/mdn/rms_5891a1/afts/img/A*WodlQKsN5UcAAAAAAAAAAAAAARQnAQ)

## Clean up resources

If you start Layotto using Docker, delete container：

```bash
docker rm -f layotto
```

Remember to close zipkin:

```shell
cd ${project_path}/diagnostics/zipkin

docker-compose -f zipkin-docker-compose.yaml down
```
