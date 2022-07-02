# Skywalking trace 接入

## 配置

示例：configs/config_trace_skywalking.json

```json
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
```

| 字段 | 必填  | 说明                       |
| --- |-----|--------------------------|
| reporter | Y   | 上报方式 grpc                |
| backend_service | Y   | skywalking oap server 地址 |
| service_name | Y   | 服务名称                     |

## 运行 skywalking

```shell
docker-compose -f diagnostics/skywalking/skywalking-docker-compose.yaml up -d
```

## 运行 layotto
<!-- tabs:start -->
### **使用 Docker**
您可以用 docker 启动 Layotto

```bash
docker run -d \
  -v "$(pwd)/configs/config_trace_skywalking.json:/runtime/configs/config.json" \
  -p 34904:34904 --network=skywalking_default --name layotto \
  layotto/layotto start
```

### **本地编译（不适合 Windows)**
您可以本地编译、运行 Layotto。

> [!TIP|label: 不适合 Windows 用户]
> Layotto 在 Windows 下会编译失败。建议 Windows 用户使用 docker 部署

构建:

```shell
cd cmd/layotto_multiple_api/
```

```shell @if.not.exist layotto
# build it
go build -o layotto
```

运行:

```shell @background
./layotto start -c ../../configs/config_trace_skywalking.json
```
<!-- tabs:end -->

## 运行 Demo

```shell
cd ${project_path}/demo/flowcontrol
go run client.go
```

访问 http://127.0.0.1:8080

![](../../../img/trace/sky.png)

## 清理资源
如果您使用 Docker 启动 Layotto，记得删除容器：

```bash
docker rm -f layotto
```

记得关闭 skywalking:

```shell
cd ${project_path}/diagnostics/skywalking

docker-compose -f skywalking-docker-compose.yaml down
```
