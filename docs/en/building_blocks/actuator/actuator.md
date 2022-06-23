# Actuator Http API

Layotto Actuator API provides functions such as health check, view runtime metadata.It can be used to view the health status and runtime metadata of Layotto and app, and supports integration into open source infrastructure (for example, it can be integrated into k8s health check)

Similar to Spring Boot Actuator, Actuator API has more imagination in the future: Monitoring, Metrics, Auditing, and more.

## 0. When to use Actuator Http API
Actuator API is generally used for operation and maintenance systems. For example, k8s calls Actuator API to monitor the health status of Layotto and App. If the status is not good, k8s will restart the Pod or temporarily cut off the traffic;

For another example, on the Dashboard for SRE, by calling Actuator API, you can clearly see the metadata of each Layotto instance and App (for example, what is the current effective configuration), which is convenient for troubleshooting.

## 1. Health Check API
### /actuator/health/liveness
Used to check the health status of Layotto and app. The health status can be used to determine "whether restarting is needed".

GET,no parameters.

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

Return field description:

HTTP status code 200 means success, other (status code above 400) means failure.

There are three types of status fields:

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

Note: By default, the API will only return the health status of Layotto. If you want the API to also return the health status of the App, you need to develop a plugin that calls back the App. You can refer to [Actuator's design document](en/design/actuator/actuator-design-doc.md), or contact us directly to provide you with a detailed explanation.

### /actuator/health/readiness
Used to check the health status of Layotto and app. The health status can be used to determine "Do we need to temporarily cut off the traffic and make sure no user visit this machine"

**Q: What is the difference with the above API?**

A: The liveness check is used to check some unrecoverable faults, "Do we need to restart it";
Readiness is used to check some temporary and recoverable states. For example, the application is warming up the cache. It needs to tell the infrastructure "Don't lead traffic to me now". After it finishes warming up, the infrastructure will reinvoke the API and get the result "I am ready to serve customers"

GET,no parameters.

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

Note: By default, the API will only return the health status of Layotto. If you want the API to also return the health status of the App, you need to develop a plugin that calls back the App. You can refer to [Actuator's design document](en/design/actuator/actuator-design-doc.md), or contact us directly to provide you with a detailed explanation.

## 2. Query runtime metadata API

### /actuator/info
Used to view the runtime metadata of Layotto and app. 

GET,no parameters.

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

**Q: Which runtime metadata will be returned?**

Currently only version number

We can add more information in the future:

- Callback app
- Runtime configuration parameters

Actuator adopts a plug-in architecture, you can also add your own plug-ins as needed, and let the API return the runtime metadata you care about.

Note: By default, the API will only return Layotto's runtime metadata. If you want the API to also return the App's runtime metadata, you need to develop a plugin that calls back the App. You can refer to [Actuator's design document](en/design/actuator/actuator-design-doc.md), or contact us directly to provide you with a detailed explanation.

## 3. Explanation for API path

Actuator API path adopts restful style. After different Endpoints are registered in Actuator, the path is:

```
/actuator/{endpoint_name}/{params}
```

For example:

```
/actuator/health/liveness
```

The 'health' element in the path above identifies the Endpoint name is health, and 'liveness' is the parameter passed to the health Endpoint.

Multiple parameters can be passed, such as /actuator/xxxendpoint/a/b/c/d, and the semantics of the parameters are determined by each Endpoint.


The paths registered by default are:

```
/actuator/health/liveness

/actuator/health/readiness

/actuator/info
```

## 4. API usage example
See [Quick start document](en/start/actuator/start.md)