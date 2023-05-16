# Layotto 支持UNIX通信

## 快速开始

Layotto提供了基于UNIX通信的能力，相对于TCP方式来说，UNIX方式具有更好的性能。

### step 1.  启动layotto

layotto提供了支持UNIX通信的配置文件`configs/config_uds.json`，配置文件内容如下所示:

```json
{
  "servers": [
    {
      "default_log_path": "stdout",
      "default_log_level": "DEBUG",
      "routers": [
        {
          "router_config_name": "actuator_dont_need_router"
        }
      ],
      "listeners": [
        {
          "name": "grpc",
          "address": "/tmp/client-proxy.sock",
          "bind_port": true,
          "network": "unix",
          "filter_chains": [
            {
              "filters": [
                {
                  "type": "grpc",
                  "config": {
                    "server_name": "runtime",
                    "grpc_config": {
                      "hellos": {
                        "helloworld": {
                          "type": "helloworld",
                          "hello": "greeting"
                        }
                      }
                    }
                  }
                }
              ]
            }
          ]
        }
      ]
    }
  ]
}

```
与TCP配置相比主要有两个不同，network的类型从tcp变为unix，address从ip地址变为unix套接字文件。

配置好后，切换目录:

```shell
#备注 请将${project_path}替换成你的项目路径
cd ${project_path}/cmd/layotto
```

构建:

```shell @if.not.exist layotto
go build -o layotto
```

启动 Layotto:

```shell @background
./layotto start -c ../../configs/config_uds.json
```

### step 2. 启动测试demo

Layotto提供了访问文件的示例 [demo](https://github.com/mosn/layotto/blob/main/demo/uds/client.go)

```shell
cd ${project_path}/demo/oss/
go build client.go

# 通过uds访问layotto的hellos组件
./client 
```

