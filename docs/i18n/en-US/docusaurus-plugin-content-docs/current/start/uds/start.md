# Layotto Support ODS Communications

## Quick Start

Layotto provides the ability to communicate on ODS and has a better performance than the TCP.

### step 1. Start layotto

layotto provides the configuration file `configs/config_uds.json` to support UDS communications. The configuration file reads as follows:

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

There are two main differences compared to the TCP configuration, the type of network changed from tcp to unix, addresses from IP to unix socket file.

When configured, toggle directory:

```shell
#备注 请将${project_path}替换成你的项目路径
cd ${project_path}/cmd/layotto
```

Build:

```shell @if.not.exist layotto
go build -o layotto
```

Start Layotto:

```shell @background
./layotto start -c ../../configs/config_uds.json
```

### step 2. Start testing demo

<!-- tabs:start -->

#### **Go**

Build and run go language demo:

Layotto provides examples [demo]to call gRPC interfaces via ODS (https://github.com/mosn/layotto/blob/main/demo/uds/client.go)

```shell
cd ${project_path}/demo/uds/
go build client.go

# 通过UDS访问layotto的hellos组件
./client 
```

#### **Java**

Build, run java language demo:

Layotto java-sdk has supported calling gRPC via ODS

```shell @if.not.exist java-sdk
git clone https://github.com/layotto/java-sdk
```

Switch directory:

```shell
cd java-sdk
```

Build:

```shell @if.not.exist examples-uds/target/examples-uds-jar-with-dependencies.jar
# build example jar
mvn -f examples-uds/pom.xml clean package
```

Run:

```
java -jar examples-uds/target/examples-uds-in-with-dependencies.jar
```

The following information was printed and run successfully:

```bash
greeting, helloowold
```

<!-- tabs:end -->
