# Actuator Http API

Layotto Actuator API提供健康检查、查看运行时元数据等功能，支持集成进开源基础设施（例如可以集成进k8s健康检查）

类似于Spring Boot Actuator，Actuator API未来有更多的想象空间：Monitoring, Metrics, Auditing, and more.

## 1. 健康检查
### /actuator/health/liveness

用于检查健康状态，判断"是否需要重启"

GET

不需要传参
```json
// http://localhost:8080/actuator/health/liveness
// HTTP/1.1 200 OK

{
  "status": "UP",
  "components": {
    "livenessProbe": {
      "status": "UP",
      "details":{
				 
      }
    }
  }
}
```
返回字段说明：
HTTP状态码200代表成功，其他(400以上的状态码)代表失败

status字段有三种：
```go
var (
	// INIT means it is starting
	INIT = Status("INIT")
	// UP means it is healthy
	UP   = Status("UP")
	// DOWN means it is unhealthy
	DOWN = Status("DOWN")
)
```

### /actuator/health/readiness

用于检查健康状态，"是否需要暂时把流量切走、别访问这台机器"

Q: 和上面的接口的区别是?

A: liveness检查用于检查一些不可恢复的故障，"是否需要重启"；
而readiness用于检查一些临时性、可恢复的状态，比如应用正在预热缓存，需要告诉基础设施"先别把流量引到我这里来"，等过会预热好了，基础设施再调readiness检查的接口，会得到结果"我准备好了，可以接客了"

GET,不需要传参
```json
// http://localhost:8080/actuator/health/readiness
// HTTP/1.1 503 SERVICE UNAVAILABLE

{
  "status": "DOWN",
  "components": {
    "readinessProbe": {
      "status": "DOWN"
    }
  }
}
```
## 2. 查询运行时元数据

### /actuator/info

GET
```json
// http://localhost:8080/actuator/health/liveness
// HTTP/1.1 200 OK

{
    "app" : {
        "version" : "1.0.0",
        "name" : "Layotto"
    }
}
```

**Q: 会返回哪些运行时元数据？**

目前返回版本号

后续可以加上：

- 回调app
- 运行时配置参数

Actuator采用插件化架构，您也可以按需添加自己的插件，让API返回您关注的运行时元数据

## 3. 路径解释

Actuator API的路径采用restful风格，不同的Endpoint注册进Actuator后，路径是

```
/actuator/{endpoint_name}/{params}  
```

比如

```
/actuator/health/liveness
```

其中health标识Endpoint的名称是health，liveness是传给该Endpoint的参数。

参数支持传多个，形如 /a/b/c/d，具体传几个、参数的语义由每个Endpoint自己定


默认注册的路径有：

```
/actuator/health/liveness
/actuator/health/readiness
/actuator/info
```

## 4. API使用示例
您可以查看[Quick start文档](zh/start/actuator/start.md)
