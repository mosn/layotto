# Configuration API demo with apollo
## What is Configuration API

When the application is started and running, it will read some "configuration information", such as: database connection parameters, startup parameters, RPC timeout, application port, etc. "Configuration" basically accompanies the entire life cycle of the application.

After the application evolves to the microservice architecture, it will be deployed on many machines, and the configuration will be scattered on each machine in the cluster, which is difficult to manage. So there is a "configuration center", which centrally manages the configuration, and also solves some new problems, such as: version management (in order to support rollback), authority management, etc.

There are many commonly used configuration centers, such as Spring Cloud Config, Apollo, Nacos, and cloud vendors often provide their own configuration management services, such as AWS Parameter Store, Google RuntimeConfig

Unfortunately, the APIs of these configuration centers are different. When developers want to deploy their apps across clouds, or want their apps to be portable (for example, easily moving from Alibaba Cloud to Tencent Cloud), they have to refactor their code.

The design goal of Layotto Configuration API is to define a unified configuration center API. Applications only need to care about the API, not which configuration center is used, so that the application can be transplanted at will, and the application is sufficiently "cloud native".

## Quick start

This example shows how to add, delete, modify, and watch the [apollo configuration center](https://github.com/ctripcorp/apollo) through Layotto.

The architecture of this example is shown in the figure below. The processes started are: client APP, Layotto, Apollo server

![img.png](../../../img/configuration/apollo/arch.png)

### Deploy apollo and modify Layotto's config (optional)

You can skip this step, you don't need to deploy the apollo server yourself to use this demo. This demo will use the demo environment provided by [apollo official](https://github.com/ctripcorp/apollo): http://106.54.227.205/

If you have deployed apollo yourself, you can modify Layotto's config file (e.g. configs/config_apollo.json in the project) to change the apollo server address to your own.

### Run Layotto server

Download the project code to the local:
```bash
git clone https://github.com/mosn/layotto.git
```

Switch the code directory and compile:

```bash
cd ${projectpath}/cmd/layotto
go build
```

After success, a new layotto file will be generated in the directory. Let's run it:

```bash
./layotto start -c ../../configs/config_apollo.json
```

### Run the client demo 

The client demo calls Layotto to add, delete, modify, and query configuration

```bash
 cd ${projectpath}/demo/configuration/apollo
 go build -o apolloClientDemo
 ./apolloClientDemo
```

If the following information is printed, the call is successful：

```bash
save key success
get configuration after save, &{Key:key1 Content:value1 Group:application Label:prod Tags:map[feature:print release:1.0.0] Metadata:map[]} 
get configuration after save, &{Key:haha Content:heihei Group:application Label:prod Tags:map[feature:haha release:1.0.0] Metadata:map[]} 
delete keys success
write start
receive subscribe resp store_name:"apollo" app_id:"apollo" items:<key:"heihei" content:"heihei1" group:"application" label:"prod" tags:<key:"feature" value:"haha" > tags:<key:"release" value:"16" > >
```

### Next step

The client demo uses the golang version SDK provided by Layotto. The SDK is located in the `sdk` directory. Users can directly call the APIs provided by Layotto through the corresponding SDK.

Besides the SDK,you can also call Layotto server directly using grpc,which makes it easy for different language to interact with Layotto.

