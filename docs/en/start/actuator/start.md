# Use Layotto Actuator for health check and metadata query

This example shows how to perform health check and metadata query through Layotto Actuator's Http API

## What is Layotto Actuator

In the production environment, the status of the application needs to be monitored, and Layotto has built-in a monitoring function, which is called Actuator. 

Using Layotto Actuator can help you monitor and manage Layotto and the applications behind Layotto, such as health checks, query runtime metadata, etc.

All these features can be accessed through the HTTP API.

## Quick start

### Run Layotto server

After downloading the project source code, change directory and compile:

```shell
cd ${project_path}/cmd/layotto
```

```shell @if.not.exist layotto
go build
```

After completion, the layotto file will be generated in the directory, run it:

```shell @background
./layotto start -c ../../configs/config_standalone.json
```

>Q: The demo report an error?
>
>A: With the default configuration, Layotto will connect to apollo's demo server, but the configuration in that demo server may be modified by others. So the error may be because some configuration has been modified.
>
> In this case, you can try other demos.

### Access the health check API

Visit /actuator/health/liveness

```shell
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

If the current state is unhealthy, the value of "status" will be "DOWN", and the Http status code returned will be 503.

### Query metadata

Visit /actuator/info

```shell
curl http://127.0.0.1:34999/actuator/info
```

return:

```json
{
  "app": {
    "name": "Layotto",
    "version": "0.1.0",
    "compiled": "2021-05-20T14:32:40.522057+08:00"
  }
}
```

[comment]: <> (### Simulate a configuration error scenario)

[comment]: <> (If a configuration error causes Layotto unavailable after startup, it can be discovered in time through the health check function.)

[comment]: <> (We can simulate a configuration error scenario by starting Layotto with an incorrect configuration file:)

[comment]: <> (```bash)

[comment]: <> (./layotto start -c ../../configs/wrong/config_apollo_health.json)

[comment]: <> (```)

[comment]: <> (There isn't an 'open_api_token' field in the configuration file,which is required to access apollo.)

[comment]: <> (Access the health check API &#40;note that the port configured here is 34888, which is different from the previous example&#41;:)

[comment]: <> (```bash)

[comment]: <> (curl http://127.0.0.1:34888/actuator/health/liveness)

[comment]: <> (```)

[comment]: <> (return:)

[comment]: <> (```json)

[comment]: <> ({)

[comment]: <> (  "components": {)

[comment]: <> (    "apollo": {)

[comment]: <> (      "status": "DOWN",)

[comment]: <> (      "details": {)

[comment]: <> (        "reason": "configuration illegal:no open_api_token")

[comment]: <> (      })

[comment]: <> (    },)

[comment]: <> (    "runtime_startup": {)

[comment]: <> (      "status": "DOWN",)

[comment]: <> (      "details": {)

[comment]: <> (        "reason": "configuration illegal:no open_api_token")

[comment]: <> (      })

[comment]: <> (    })

[comment]: <> (  },)

[comment]: <> (  "status": "DOWN")

[comment]: <> (})

[comment]: <> (```)

[comment]: <> ("status": "DOWN" in json means the current status is unhealthy. The Http status code returned this time is 503.)

## Next step

### Integrated into Kubernetes health check

Layotto provides two built-in health check API: /actuator/health/readiness and /actuator/health/liveness, corresponding to the two semantics of Readiness and Liveness in the Kubernetes health check feature.

Therefore, you can refer to [Kubernetes documentation](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/) to integrate these two API into the Kubernetes ecosystem.

### Add health check or metadata query capabilities to your components

If you are implementing your own Layotto component, you can add health check capabilities to it. You can refer to the implementation of the apollo component (the code is at components/configstores/apollo/indicator.go), implement the info.Indicator interface, and inject it into the Actuator.

### How it works

If you are interested in the implementation principle, or want to extend some functions in Actuator, you can read [Actuator's design document](en/design/actuator/actuator-design-doc.md)