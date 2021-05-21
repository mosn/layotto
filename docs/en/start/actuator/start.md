# Use LayOtto Actuator for health check and metadata query

This example shows how to perform health check and metadata query through LayOtto Actuator's Http API

## What is LayOtto Actuator

In the production environment, the status of the application needs to be monitored, and LayOtto has built-in a monitoring function, which is called Actuator. 

Using LayOtto Actuator can help you monitor and manage LayOtto and the applications behind LayOtto, such as health checks, query runtime metadata, etc.

All these features can be accessed through the HTTP API.

## Quick start

### Run LayOtto server

After downloading the project source code, change directory and compile:

```bash
cd ${projectpath}/cmd/layotto
go build
```

After completion, the layotto file will be generated in the directory, run it:

```bash
./layotto start -c ../../configs/config_apollo_health.json
```

### Access the health check API

Visit /actuator/health/liveness

```bash
curl http://127.0.0.1:34999/actuator/health/liveness
```

return:

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

In the above json,"status": "UP" means the status is healthy. The Http status code returned is 200.

### Query metadata

Visit /actuator/info

```shell
curl http://127.0.0.1:34999/actuator/info
```

return:

```json
{
  "app": {
    "name": "LayOtto",
    "version": "0.1.0",
    "compiled": "2021-05-20T14:32:40.522057+08:00"
  }
}
```

### Simulate a configuration error scenario

If a configuration error causes LayOtto unavailable after startup, it can be discovered in time through the health check function.

We can simulate a configuration error scenario by starting LayOtto with an incorrect configuration file:

```shell
./layotto start -c ../../configs/wrong/config_apollo_health.json
```

There isn't an 'open_api_token' field in the configuration file,which is required to access apollo.

Access the health check API (note that the port configured here is 34888, which is different from the previous example):

```shell
curl http://127.0.0.1:34888/actuator/health/liveness
```

return:

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

"status": "DOWN" in json means the current status is unhealthy. The Http status code returned this time is 503.

## Next step

### Integrated into Kubernetes health check

LayOtto provides two built-in health check API: /actuator/health/readiness and /actuator/health/liveness, corresponding to the two semantics of Readiness and Liveness in the Kubernetes health check feature.

Therefore, you can refer to [Kubernetes documentation](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/) to integrate these two API into the Kubernetes ecosystem.

### Add health check or metadata query capabilities to your components

If you are implementing your own LayOtto component, you can add health check capabilities to it. You can refer to the implementation of the apollo component (the code is at pkg/services/configstores/apollo/indicator.go), implement the info.Indicator interface, and inject it into the Actuator.

### How it works

If you are interested in the implementation principle, or want to extend some functions in Actuator, you can read [Actuator's design document](../../design/actuator-design-doc.md)