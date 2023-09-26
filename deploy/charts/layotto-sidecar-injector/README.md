# Introduction

This chart deploys the Layotto sidecar injector component on a Kubernetes cluster using the Helm package manager.

## Prerequisites

- Helm 3
- Kubernetes cluster

## Installing the Chart

Ensure Helm is initialized in your Kubernetes cluster.

For more details on initializing Helm, [read the Helm docs](https://helm.sh/docs/)

1. Add xiaoxiang10086.github.io as an helm repo

```
helm repo add layotto https://xiaoxiang10086.github.io/layotto-helm-charts/
helm repo update
```

2. Install the Layotto chart on your cluster in the Layotto-system namespace:

```
helm install layotto layotto/layotto-sidecar-injector --namespace layotto-system --wait
```


## Verify installation

Once the chart is installed, verify the Layotto sidecar injector component pods are running in the layotto-system namespace:

```
kubectl get pods --namespace layotto-system
```

## Uninstall the Chart

To uninstall/delete the layotto release:

```
helm uninstall layotto -n layotto-system
```

## Configuration
The following tables list the configurable parameters of the chart and their default values.


| **Parameter**                                      | **Description**                                              | **Default**           |
| -------------------------------------------------- | ------------------------------------------------------------ | --------------------- |
| `registry`                | Docker image registry                                        | `docker.io/layotto` |
| `image.name`              | Docker image name for Layotto runtime sidecar to inject into an application | `layotto`             |
| `sidecarImagePullPolicy`  | Layotto sidecar image pull policy                            | `IfNotPresent`        |
| `replicaCount`            | Number of replicas                                           | `1`                   |
| `injectorImage.name`      | Docker image name for sidecar injector component             | `injector`    |
| `injectorImagePullPolicy` | Layotto sidecar injector image pull policy                   | `IfNotPresent`        |
| `webhookFailurePolicy`    | Failure policy for the sidecar injector                      | `Ignore`              |
