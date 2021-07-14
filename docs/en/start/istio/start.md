# Layotto integrate with Istio

## 1. Background

The RPC capability supported by layotto is implemented by mosn, it means that the RPC feature can reuse the capabilities in mosn.

Mosn is officially recognized by istio as a data plane implementation. Here is how to combine layotto (build on mosn) with istio and dynamically adjust service routing through XDS protocol.

## 2. Preparation

before start the demo，you must install some components as follows：
1. [Docker Desktop](https://www.docker.com/products/docker-desktop)

   Download from the official website and install it. 

2. [minikube](https://minikube.sigs.k8s.io/docs/start/)

   same as above, just follow the doc in official website.

3. [Istio-1.5.x](https://github.com/istio/istio/releases/tag/1.5.2)

   Currently, mosn only supports `istio 1.5.X` (the support for `istio 1.10.X` is already in CR), so you need to download the corresponding version of `istio`. After decompressing, configure it as follows to facilitate subsequent operations.
   ```
   export PATH=$PATH:${your istio directory}/bin
   ```

## 3. Start the demo

1. Run Docker Desktop
2. Run the following command to start `minikube`
   ```
   minikube start
   ```
3. Run the following command to start the services in the demo (all the dependent images have been uploaded to the docker hub)
   ```
   kubectl apply -f layotto.yaml
   ```
   The contents of `layotto.yaml` is [here](https://github.com/mosn/layotto/blob/istio-1.5.x/demo/istio/layotto-injected.yaml) ，just copy it。
4. Run the command `kubectl get pod` to check the status (the first startup needs to download the dependent image, please wait patiently)
   ```
   NAME                         READY   STATUS    RESTARTS   AGE
   client-665c5cc4f-tfxrk       2/2     Running   0          49m
   server-v1-685966b499-8hnqp   2/2     Running   0          49m
   server-v2-6cfff5dbb5-4hlgb   2/2     Running   0          49m
   ```
   When you see something similar to the above, it indicates that the startup is successful. We have deployed a client and a server. The server side is divided into V1 and V2 versions.
   
5. If you want to access the services in the `istio` cluster from the outside, you must configure the `istio ingress gateway` service, which will increase the cost of getting started. Therefore, the proxy method is used here to simplify.
   Run the following command
   ```
   kubectl port-forward svc/client 9080:9080
   ```
   Then you can directly access the following links, or you can directly access them in the browser.
   ```
   curl localhost:9080/grpc
   ```
   When you see the following response, the example starts successfully.
   ```
   GET /hello 
   hello, i am layotto v1
   ```
## 4. Using istio to dynamically change routing policy

#### A、按version路由能力
1. Run the following command to create destination rules
   ```
   kubectl apply -f destination-rule-all.yaml
   ```
   The contents of `destination-rule-all.yaml` is [here](https://github.com/mosn/layotto/blob/istio-1.5.x/demo/istio/layotto-destination-rule-all.yaml)

2. 执行如下命令指定只访问V1服务
   ```
   kubectl apply -f layotto-virtual-service-all-v1.yaml
   ```
   The contents of `layotto-virtual-service-all-v1.yaml` is [here](https://github.com/mosn/layotto/blob/istio-1.5.x/demo/istio/layotto-virtual-service-all-v1.yaml)


## 5. Note

1. Since `istio 1.5.2` is used in this demo, which belongs to an older version, the demo will not be merged into the `main` branch. Instead, it exists as an independent branch `istio-1.5.X`. After 'mosn' integrates with `istio 1.10.X`, it will be merged.
2. For the source code of client and server used in the example, please refer to [here](https://github.com/mosn/layotto/tree/istio-1.5.x/demo/istio).

