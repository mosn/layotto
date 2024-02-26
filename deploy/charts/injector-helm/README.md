# Introduction

This chart deploys the Layotto injector component on a Kubernetes cluster using the Helm package manager.

## Prerequisites

- Helm 3
- Kubernetes cluster

## Source Code
- https://github.com/mosn/layotto

## Install the Chart

Ensure Helm is initialized in your Kubernetes cluster.

For more details on initializing Helm, [read the Helm docs](https://helm.sh/docs/)

You can choose install the helm chart from DockerHub.
```
helm install injector oci://docker.io/layotto/injector-helm --version v0.5.0 -n layotto-system --create-namespace --wait
```
You can also install from the source code:
```
make helm-install VERSION="v0.5.0"
```

## Verify installation

Once the chart is installed, verify the Layotto sidecar injector pod is running in the layotto-system namespace:

```
kubectl get pods --namespace layotto-system
```

## Uninstall the Chart

To uninstall/delete the layotto release:

```
helm uninstall injector -n layotto-system
```

## Configuration
The following tables list the configurable parameters of the chart and their default values.


| **Parameter**                                      | **Description**                                              | **Default**         |
| -------------------------------------------------- | ------------------------------------------------------------ |---------------------|
| `registry`                | Docker image registry                                        | `docker.io/layotto` |
| `image.name`              | Docker image name for Layotto runtime sidecar to inject into an application | `layotto`           |
| `sidecarImagePullPolicy`  | Layotto sidecar image pull policy                            | `IfNotPresent`      |
| `replicaCount`            | Number of replicas                                           | `1`                 |
| `injectorImage.name`      | Docker image name for sidecar injector component             | `layotto_injector`  |
| `injectorImagePullPolicy` | Layotto sidecar injector image pull policy                   | `IfNotPresent`      |
| `webhookFailurePolicy`    | Failure policy for the sidecar injector                      | `Ignore`            |
