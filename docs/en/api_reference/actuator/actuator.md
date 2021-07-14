# Actuator Http API

Layotto Actuator API provides functions such as health check, view runtime metadata, and supports integration into open source infrastructure (for example, it can be integrated into k8s health check)

Similar to Spring Boot Actuator, Actuator API has more imagination in the future: Monitoring, Metrics, Auditing, and more.

## 1. Health Check
### /actuator/health/liveness
Used to check the health status and determine "whether restarting is needed"

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

### /actuator/health/readiness
Used to check the health status and determine "Do we need to temporarily cut off the traffic and make sure no user visit this machine"

Q: What is the difference with the above API?

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
## 2. View runtime metadata

### /actuator/info

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