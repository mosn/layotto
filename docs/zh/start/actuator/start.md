# 使用Layotto Actuator进行健康检查和元数据查询

该示例展示了如何通过Layotto Actuator的Http API进行健康检查和元数据查询

## 什么是Layotto Actuator

在生产环境中，需要对应用程序的状态进行监控，而Layotto已经内置了一个监控功能，它叫Actuator。 使用Layotto Actuator可以帮助你监控和管理Layotto和Layotto服务的应用，比如健康检查、查询运行时元数据等。
所有的这些特性可以通过HTTP接口来访问。

## 快速开始

### 运行Layotto server 端

将项目代码下载到本地后，切换代码目录、编译：

```bash
cd ${projectpath}/cmd/layotto
go build
```

完成后目录下会生成layotto文件，运行它：

```bash
./layotto start -c ../../configs/config_apollo_health.json
```

### 访问健康检查接口

访问 /actuator/health/liveness

```bash
curl http://127.0.0.1:34999/actuator/health/liveness
```

返回：

```json
{
  "components": {
    "apollo": {
      "status": "UP"
    },
    "runtime_startup": {
      "status": "UP",
      "details": {
        "reason": ""
      }
    }
  },
  "status": "UP"
}
```

其中"status": "UP"代表状态健康。此时返回的Http状态码是200。

### 查询元数据

访问 /actuator/info

```shell
curl http://127.0.0.1:34999/actuator/info
```

返回：

```json
{
  "app": {
    "name": "Layotto",
    "version": "0.1.0",
    "compiled": "2021-05-20T14:32:40.522057+08:00"
  }
}
```

### 模拟配置错误的场景

如果Layotto配置错误导致启动后不能正常提供服务，通过健康检查功能可以及时发现。

我们可以模拟一下配置错误的场景，使用一个错误的配置文件启动Layotto:

```shell
./layotto start -c ../../configs/wrong/config_apollo_health.json
```

该配置文件中忘记配置了访问apollo需要的open_api_token。

访问健康检查接口（注意这里配置的端口是34888，和上一个例子中不一样）：

```shell
curl http://127.0.0.1:34888/actuator/health/liveness
```

返回：

```json
{
  "components": {
    "apollo": {
      "status": "DOWN",
      "details": {
        "reason": "configuration illegal:no open_api_token"
      }
    },
    "runtime_startup": {
      "status": "DOWN",
      "details": {
        "reason": "configuration illegal:no open_api_token"
      }
    }
  },
  "status": "DOWN"
}
```

json中"status": "DOWN"代表当前状态不健康。此时返回的Http状态码是503。


## 下一步

### 集成进Kubernetes健康检查

Layotto内置提供了/actuator/health/readiness和/actuator/health/liveness 两个健康检查接口，对应Kubernetes健康检查功能中Readiness和Liveness两个语义。

因此，您可以参照[Kubernetes的文档](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/) ，将这两个接口集成进Kubernetes健康检查。

### 为您的组件添加健康检查或元数据查询能力

如果您实现了自己的Layotto组件，可以为其添加健康检查能力。可以参考apollo组件的实现（文件在github.com/layotto/L8-components/configstores/apollo/indicator.go），实现info.Indicator接口，并将其注入进Actuator即可。

### 了解Actuator实现原理

如果您对实现原理感兴趣，或者想在Actuator扩展一些功能，可以阅读[Actuator的设计文档](../../design/actuator/actuator-design-doc.md)