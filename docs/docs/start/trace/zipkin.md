# ZipKin trace 接入

## 配置

示例：configs/config_trace_zipkin.json

```json
{
  "tracing": {
    "enable": true,
    "driver": "Zipkin",
    "config": {
      "config": {
        "service_name": "layotto",
        "reporter_endpoint": "http://127.0.0.1:9411/api/v2/spans",
        "recorder_host_post": "127.0.0.1:34904"
      }
    }
  }
}

```

| 字段   | 必填  | 说明                       |
|------|-----|--------------------------|
| service_name | Y   | 当前服务名称，例如layotto         |
| reporter_endpoint | Y   | 链路日志上报url                |
| recorder_host_post     | Y   | 当前服务端口信息，例如layotto服务的端口为127.0.0.1:34904 |

注意：目前只支持Http方式的Reporter。

## 运行ZipKin

```shell
docker-compose -f diagnostics/zipkin/zipkin-docker-compose.yaml up -d
```

## 运行layotto

<!-- tabs:start -->

### **使用 Docker**

您可以用 docker 启动 Layotto

```bash
docker run -d \
  -v "$(pwd)/configs/config_trace_zipkin.json:/runtime/configs/config.json" \
  -p 34904:34904 --network=zipkin_default --name layotto \
  layotto/layotto start
```

### **本地编译（不适合 Windows)**
您可以本地编译、运行 Layotto。

> [!TIP|label: 不适合 Windows 用户]
> Layotto 在 Windows 下会编译失败。建议 Windows 用户使用 docker 部署


构建:

```shell
cd ${project_path}/cmd/layotto_multiple_api/
```

```shell @if.not.exist layotto
go build -o layotto
```

运行:

```shell @background
./layotto start -c ../../configs/config_trace_zipkin.json 
```
<!-- tabs:end -->

## 运行 Demo

```shell
 cd ${project_path}/demo/flowcontrol/
 go run client.go
``` 

访问：http://localhost:9411/zipkin/?serviceName=layotto&lookback=15m&endTs=1655559536414&limit=10

![](https://gw.alipayobjects.com/mdn/rms_5891a1/afts/img/A*WodlQKsN5UcAAAAAAAAAAAAAARQnAQ)

## 清理资源

如果您使用 Docker 启动 Layotto，记得删除容器：

```bash
docker rm -f layotto
```

记得关闭 zipkin:

```shell
cd ${project_path}/diagnostics/zipkin

docker-compose -f zipkin-docker-compose.yaml down
```