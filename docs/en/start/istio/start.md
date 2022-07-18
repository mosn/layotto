# Demo of Istio 1.5.x integration

## 1. Background

The RPC capability supported by layotto is implemented by mosn, it means that the RPC feature can reuse the capabilities in mosn.

Mosn is officially recognized by istio as a data plane implementation. Here is how to combine layotto (build on mosn) with istio and dynamically adjust service routing through XDS protocol.

## 2. Preparation

before starting the demo，you must install some components as follows：
1. [Docker Desktop](https://www.docker.com/products/docker-desktop)

   Download from the official website and install it. 

2. [minikube](https://minikube.sigs.k8s.io/docs/start/)

   same as above, just follow the doc in official website.

3. [Istio-1.5.x](https://github.com/istio/istio/releases/tag/1.5.2)

   You need to download the `1.5.x` version of `istio`. After decompressing, configure it as follows to facilitate subsequent operations.

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
   kubectl apply -f layotto-injected.yaml
   ```
   
   The contents of `layotto-injected.yaml` is [here](https://github.com/mosn/layotto/blob/istio-1.5.x/demo/istio/layotto-injected.yaml) ，just copy it。
4. Run the command `kubectl get pod` to check the status (it needs to download the dependent images during the first startup,so please wait patiently)
   
   ```
   NAME                         READY   STATUS    RESTARTS   AGE
   client-665c5cc4f-tfxrk       2/2     Running   0          49m
   server-v1-685966b499-8hnqp   2/2     Running   0          49m
   server-v2-6cfff5dbb5-4hlgb   2/2     Running   0          49m
   ```
   
   When you see something similar to the above, it indicates that the startup is successful. We have deployed a client and a server. The server side is divided into V1 and V2 versions.
   
5. If you want to access the services in the `istio` cluster from the outside, you must configure the `istio ingress gateway` service, which will increase the cost of getting started. Therefore, the proxy method is used here to simplify this demo.

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

### A. route according to version
1. Run the following command to create destination rules
   
   ```
   kubectl apply -f destination-rule-all.yaml
   ```
   
   The contents of `destination-rule-all.yaml` is [here](https://github.com/mosn/layotto/blob/istio-1.5.x/demo/istio/layotto-destination-rule-all.yaml)

2. Run the following command to specify that only the V1 service is accessed
   
   ```
   kubectl apply -f layotto-virtual-service-all-v1.yaml
   ```
   
   The contents of `layotto-virtual-service-all-v1.yaml` is [here](https://github.com/mosn/layotto/blob/istio-1.5.x/demo/istio/layotto-virtual-service-all-v1.yaml)
3. After the above command is executed, subsequent requests will only get the return result of v1, as follows:
   
   ```
   GET /hello 
   hello, i am layotto v1
   ```
   
### B. route according to a specific header
1. Run the following command to modify the routing rules to access the v1 service when the request header contains `name:layotto`, and other access to the v2 service
   
   ```
   kubectl apply -f layotto-header-route.yaml
   ```
   
2. Send the request to see the result
   
   ```
   curl -H 'name: layotto' localhost:9080/grpc
   ```

## 5. Note

1. Since the example uses `istio 1.5.2`, which is an older version, the demo will not be merged into the main branch, but exists as a separate branch `istio-1.5.x`. Currently the main branch code has been integrated with `istio 1.10.x`.
   
2. For the source code of client and server used in the example, please refer to [here](https://github.com/mosn/layotto/tree/istio-1.5.x/demo/istio).
3. In order to get started simple, the `layotto-injected.yaml` file used above has been injected through istio already. This injection process is as follows:
   1. Run the following command to specify `istio` to use `Layotto` as the data plane
   
   ```
   istioctl manifest apply  --set .values.global.proxy.image="mosnio/proxyv2:layotto"   --set meshConfig.defaultConfig.binaryPath="/usr/local/bin/mosn"
   ```
   
   2. Sidecar injection is achieved through `kube-inject`
   
   ```
   istioctl kube-inject -f layotto.yaml > layotto-injected.yaml
   ```
   
   The contents of `layotto.yaml` is [here](https://github.com/mosn/layotto/blob/istio-1.5.x/demo/istio/layotto.yaml)
   
   3. Run the following command to replace all `/usr/local/bin/envoy` in `layotto-injected.yaml` with `/usr/local/bin/mosn`
   
   ```
   sed -i "s/\/usr\/local\/bin\/envoy/\/usr\/local\/bin\/mosn/g" ./layotto-injected.yaml
   ```

