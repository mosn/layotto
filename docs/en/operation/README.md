# How to deploy and upgrade Layotto

## 1. Deploy Layotto

There are some ways to deploy Layotto that you can find below.

- Deploy using released binaries
- Deploy using Docker
- Deploy on Kubernetes

### Deploy Layotto using released binaries

You can start Layotto directly via executing the binary file. Refer to the [Quick start](en/start) guide.

### Deploy using Docker

You can run Layotto using the official Docker images.Currently include:

- [layotto/layotto](https://hub.docker.com/repository/docker/layotto/layotto)
- [layotto/layotto.arm64](https://hub.docker.com/repository/docker/layotto/layotto.arm64)

It does not contain a `config.json` configuration file in the image, you can mount your own configuration file into the `/runtime/configs/` directory of the image. For example.

```shell
docker run -v "$(pwd)/configs/config.json:/runtime/configs/config.json" -d  -p 34904:34904 --name layotto layotto/layotto start
```

Of course, you can also run Layotto and other systems (such as Redis) at the same time via docker-compose. Refer to the [Quick start](en/start/state/start?id=step-1-deploy-redis-and-layotto)

### Deploy on Kubernetes

#### Option 1. Deploy using Istio

If you are using Istio, you can deploy the Layotto sidecar using Istio.

In order to integrate with Istio, you need to build the Layotto sidecar as a proxyv2 image. There are several ways to build:

- You can refer to [MOSN guide](https://mosn.io/docs/user-guide/start/istio/). Just replace the MOSN image in the tutorial with a Layotto image.

- To make the build process easier, Layotto developed the make script. After modifying the main package as needed (in the `cmd/layotto` directory), run the following command to package it into a proxyv2 image:

```shell
make image.proxyv2.build VERSION=latest
```

The principle of this script is to replace the binary file in the MOSN proxyv2 image with the binary file of Layotto.

- The community has built the [layotto/proxyv2](https://hub.docker.com/r/layotto/proxyv2) image for the [Online Lab](https://killercoda.com/mosn-tutorial/course /layotto) to demonstrate Istio+Layotto. This image packs all Layotto components into it, resulting in a very large size. We recommend that you modify the main package according to your own needs and build the image yourself, so that the built image will be much smaller.

#### Option 2. Other ways

You can prepare your own image and K8s configuration file, then deploy Layotto via Kubernetes.

We are working on the official Layotto image and the solution for deploying to Kubernetes using Helm, so feel free to join us to build it. More details in <https://github.com/mosn/layotto/issues/392>

## 2.Toggle existing MOSN to Layotto for MOSN users

Existing MOSN can be upgraded to Layotto by replacing the MOSN sidecar image with Layotto image.

Explanation:

Layotto and MOSN are running in the same process, which can be understood as:

> Layotto == MOSN + a special grpcFilter packaged together

So.

> replace MOSN with Layotto == replace MOSN with "MOSN + a special grpcFilter"

There is no essential difference between Layotto and MOSN, just pay attention to the versions of them, which must correspond to each other.

The previously released Layotto v0.3.0 corresponds to MOSN version v0.24.1

## 3. How to upgrade Layotto

There are two options to upgrade.

- Upgrade sidecar container using K8s native solution
  
- [Hot upgrade: upgrade the sidecar without affecting the business](https://mosn.io/en/docs/concept/smooth-upgrade/)

The advantage of hot upgrade is that it can automatically migrate persistent connections, which can be seen in detail by clicking the above document.

The solutions to achieve hot upgrade are as follows:

- Register a SIGHUP event listener with MOSN, and send a SIGHUP signal to the MOSN process to call ForkExec to generate a new MOSN process.
  
- Directly start a new MOSN process.

- [Hot upgrade using OpenKruise](https://mosn.io/blog/posts/mosn-sidecarset-hotupgrade/)
