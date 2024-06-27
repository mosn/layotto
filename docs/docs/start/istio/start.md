# Layotto集成 Istio 1.5.x 演示

## 一、背景介绍

Layotto支持的RPC能力底层是使用MOSN实现的，也就是说可以复用MOSN多年来在服务治理领域建设的能力。

MOSN作为Istio官方认可的数据面实现，这里就对Layotto如何跟Istio结合，通过XDS协议来动态调整服务路由进行说明。

## 二、准备工作

想要启动示例，必须安装以下组件：
1. [Docker Desktop](https://www.docker.com/products/docker-desktop)
   
    直接官网下载安装包安装即可。
   
2. [minikube](https://minikube.sigs.k8s.io/docs/start/) 
   
    按照官网操作即可。

3. [Istio-1.5.x](https://github.com/istio/istio/releases/tag/1.5.2)

    需要下载`1.5.x`版本的`istio`，解压后进行如下配置方便后续操作。

    ```
    export PATH=$PATH:${你的istio目录}/bin
    ```
   
## 三、启动示例

1. 启动Docker Desktop
2. 执行如下命令启动`minikube`
   
   ```
   minikube start
   ```
   
3. 执行如下命令启动demo中的client、server（所依赖的镜像已经全部上传docker hub）
   
   ```
   kubectl apply -f layotto-injected.yaml
   ```
   
   其中`layotto-injected.yaml`文件中的内容在[这里](https://github.com/mosn/layotto/blob/istio-1.5.x/demo/istio/layotto-injected.yaml) ，复制即可。
4. 执行命令`kubectl get pod`查看启动状态（首次启动需要下载依赖镜像，请耐心等待）
   
   ```
   NAME                         READY   STATUS    RESTARTS   AGE
   client-665c5cc4f-tfxrk       2/2     Running   0          49m
   server-v1-685966b499-8hnqp   2/2     Running   0          49m
   server-v2-6cfff5dbb5-4hlgb   2/2     Running   0          49m
   ```
   
   命令执行完后看到类似上述内容则表示启动成功，我们部署了一个client端以及一个server端，其中server端分为了v1,v2两个版本。
5. 由于原生的`Istio`如果想要从外部访问集群里面的服务需要配置`istio-ingressgateway`服务，这会增加大家使用演示的成本，因此这里我们使用代理命名进行访问，
   执行如下命令：
   
   ```
   kubectl port-forward svc/client 9080:9080
   ```
   
   然后直接访问如下链接即可，也可以直接在浏览器中访问。
   
   ```
   curl localhost:9080/grpc
   ```
   
   当看到如下响应时就表示示例启动成功。
   
   ```
   GET /hello 
   hello, i am layotto v1
   ```
   
## 四、使用Istio动态改变路由策略

### A、按version路由能力
1. 执行如下命令创建destination rules
   
   ```
   kubectl apply -f destination-rule-all.yaml
   ```
   
   其中`destination-rule-all.yaml`文件内容在[这里](https://github.com/mosn/layotto/blob/istio-1.5.x/demo/istio/layotto-destination-rule-all.yaml)

2. 执行如下命令指定只访问V1服务
   
   ```
   kubectl apply -f layotto-virtual-service-all-v1.yaml
   ```
   
   其中`layotto-virtual-service-all-v1.yaml`文件内容在[这里](https://github.com/mosn/layotto/blob/istio-1.5.x/demo/istio/layotto-virtual-service-all-v1.yaml)
3. 上述命令执行完以后，后续请求只会拿到v1的返回结果，如下：
   
   ```
   GET /hello 
   hello, i am layotto v1
   ```

### B、按header信息进行路由
1. 执行如下命令把路由规则修改为请求header中包含`name:layotto`时会访问v1服务，其他则访问v2服务
  
   ```
   kubectl apply -f layotto-header-route.yaml
   ```
   
2. 发送请求即可看到效果
   
   ```
   curl -H 'name: layotto' localhost:9080/grpc
   ```
   


## 五、注意事项

1. 由于示例中使用的是`istio 1.5.2`，属于一个比较老的版本，因此该演示不会合并到主干，而是以一个独立的分支`istio-1.5.x`存在。目前 main 分支代码已经集成了 `istio 1.10.x`。
2. 示例中使用的client、server源码可以参考[这里](https://github.com/mosn/layotto/tree/istio-1.5.x/demo/istio) 。
3. 为了上手简单，上述使用到的`layotto-injected.yaml`文件是已经通过istio完成注入的，整个注入过程如下：
   1. 执行如下命令指定`istio`使用`Layotto`作为数据面
   
   ```
   istioctl manifest apply  --set .values.global.proxy.image="mosnio/proxyv2:layotto"   --set meshConfig.defaultConfig.binaryPath="/usr/local/bin/mosn"
   ```
   
   2. 通过`kube-inject`的方式实现Sidecar注入
   
   ```
   istioctl kube-inject -f layotto.yaml > layotto-injected.yaml
   ```
   
   其中`layotto.yaml`文件内容在[这里](https://github.com/mosn/layotto/blob/istio-1.5.x/demo/istio/layotto.yaml)
   
   3. 把`layotto-injected.yaml`中所有的`/usr/local/bin/envoy`替换为`/usr/local/bin/mosn`
  
   ```
   sed -i "s/\/usr\/local\/bin\/envoy/\/usr\/local\/bin\/mosn/g" ./layotto-injected.yaml
   ```

