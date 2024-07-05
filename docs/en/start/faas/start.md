## FaaS QuickStart

### 1. Features

Layotto supports loading and running functions in the form of wasm, and supports calling each other between functions and accessing infrastructure, such as Redis.

Detailed design documents can refer to：[FaaS design](en/design/faas/faas-poc-design.md)

### 2. Dependent software

The following software needs to be installed to run this demo:

1. [Docker Desktop](https://www.docker.com/products/docker-desktop)

   Download the installation package from the official website and install it.

2. [minikube](https://minikube.sigs.k8s.io/docs/start/)

   Follow the instructions on the official website.

3. [virtualbox](https://www.oracle.com/virtualization/technologies/vm/virtualbox.html)
   
   Download the installation package from the official website and install it. You can also use [homebrew](https://brew.sh/) to install it on mac. If the startup fails after installation, please refer to [The host-only adapter we just created is not visible](https://github.com/kubernetes/minikube/issues/3614).


### 3. Setup

#### A、Install & run Redis

The example only needs a Redis server that can be used normally. As for where it is installed, there is no special restriction. It can be a virtual machine, a local machine or a server. Here, the installation on mac is used as an example to introduce.

```
> brew install redis
> redis-server /usr/local/etc/redis.conf
```

**Note: If you want external services to connect to redis, you need to modify the protected-mode in redis.conf to no,At the same time, add bind * -::* to let it monitor all interfaces.**

#### B、Start minikube in virtualbox + containerd mode

```
> minikube start --driver=virtualbox --container-runtime=containerd
```

#### C、Compile & install Layotto

```
> git clone https://github.com/mosn/layotto.git
> cd layotto
> make wasm-build
> minikube cp ./_output/linux/amd64/layotto /home/docker/layotto
> minikube cp ./demo/faas/config.json /home/docker/config.json
> minikube ssh
> sudo chmod +x layotto
> sudo mv layotto /usr/bin/
```

**Note1: You need to modify the redis address as needed, the default address is: localhost:6379**

**Note2: Need to modify the path of the wasm file in `./demo/faas/config.json` to `/home/docker/function_1.wasm` and `/home/docker/function_2.wasm`**

**Note3: We can also load WASM file dynamically. For details, see [WASM Dynamic Load](https://mosn.io/layotto/#/en/start/wasm/start?id=dynamic-load)**

#### D、Compile & install containerd-shim-layotto-v2

```
> git clone https://github.com/layotto/containerd-wasm.git
> cd containerd-wasm
> sh build.sh
> minikube cp containerd-shim-layotto-v2 /home/docker/containerd-shim-layotto-v2
> minikube ssh
> sudo chmod +x containerd-shim-layotto-v2
> sudo mv containerd-shim-layotto-v2 /usr/bin/
```

#### E、Modify & restart containerd

Add laytto runtime configuration.

```
> minikube ssh
> sudo vi /etc/containerd/config.toml
[plugins.cri.containerd.runtimes.layotto]
  runtime_type = "io.containerd.layotto.v2"
```

Restart containerd for the latest configuration to take effect

```
sudo systemctl restart containerd
```

#### F、Install wasmer (If the vm engine uses wasmer, you need to execute the following command)

```
> curl -L -O https://github.com/wasmerio/wasmer/releases/download/2.0.0/wasmer-linux-amd64.tar.gz
> tar zxvf wasmer-linux-amd64.tar.gz
> sudo cp lib/libwasmer.so /usr/lib/libwasmer.so
```

### 4. Quickstart

#### A、Start Layotto

```
> minikube ssh
> layotto start -c /home/docker/config.json
```

#### B、Create Layotto runtime

```
> kubectl apply -f ./demo/faas/layotto-runtimeclass.yaml
runtimeclass.node.k8s.io/layotto created
```

#### C、Create Function
This operation will automatically inject function_1.wasm and function_2.wasm into the Virtualbox virtual machine.

```
> kubectl apply -f ./demo/faas/function-1.yaml
pod/function-1 created

> kubectl apply -f ./demo/faas/function-2.yaml
pod/function-2 created
```

#### D、Write inventory to Redis

```
> redis-cli
127.0.0.1:6379> set book1 100
OK
```

#### E、Send request

```
> minikube ip
192.168.99.117

> curl -H 'id:id_1' '192.168.99.117:2045?name=book1'
There are 100 inventories for book1.
```

### 5. Process introduction

![img.png](../../../img/faas/faas-request-process.jpg)

1. send http request to func1
2. func1 calls func2 through Runtime ABI
3. func2 calls redis through Runtime ABI
4. Return results

### Common problem description

1. Virtualbox failed to start, "The host-only adapter we just created is not visible":

    refer  [The host-only adapter we just created is not visible](https://github.com/kubernetes/minikube/issues/3614)

2. When Layotto is started, the redis connection fails, and "occurs an error: redis store: error connecting to redis at" is printed:

   Check the redis configuration to see if it is caused by a redis configuration error.

### 6. Note

The FaaS model is currently in the POC stage, and the features are not complete. It will be improved in the following aspects in the future:
1. Limit the maximum resources that can be used when the function is running, such as cpu, heap, stack, etc.
2. The functions loaded and run by different Layotto can call each other.
3. Fully integrate into the k8s ecology, such as reporting the use of resources to k8s, so that k8s can perform better scheduling.
4. Add more Runtime ABI.

If you are interested in FaaS or have any questions or ideas, please leave a message in the issue area, we can build FaaS together!
