<h2>Layotto support configuration center</h2>
## What is Configuration API
When the application is started and running, it will read some "configuration information", such as: database connection parameters, startup parameters, RPC timeout, application port, etc. "Configuration" basically accompanies the entire life cycle of the application.

After the application evolves to the microservice architecture, it will be deployed on many machines, and the configuration will be scattered on each machine in the cluster, which is difficult to manage. So there is a "configuration center", which centrally manages the configuration, and also solves some new problems, such as: version management (in order to support rollback), authority management, etc.

There are many commonly used configuration centers, such as Spring Cloud Config, Apollo, Nacos, and cloud vendors often provide their own configuration management services, such as AWS Parameter Store, Google RuntimeConfig

Unfortunately, the APIs of these configuration centers are different. When an application wants to be deployed across clouds, or if it wants to be transplanted (for example, moving from Tencent Cloud to Alibaba Cloud), the application needs to refactor the code.

The design goal of Layotto Configuration API is to define a unified configuration center API. Applications only need to care about the API, not which configuration center is used, so that the application can be transplanted at will, and the application is sufficiently "cloud native".

## Quick start
This example shows how to add, delete, modify, and watch the etcd configuration center through Layotto. 

Please install [Docker](https://www.docker.com/get-started) software on your machine in advance.

[Config file](https://github.com/mosn/layotto/blob/main/configs/runtime_config.json) defines using etcd in config_stores section, and users can change the configuration file to the configuration center they want (currently supports etcd and apollo).

### Build docker image

At first, please make sure your layotto PATH is same as below:

```
$GOPATH/src/github/layotto/layotto
```

then execute `CMD` below:

```bash
cd $GOPATH/src/github/layotto/layotto  
make image
```

After make success, you can see two images with docker images command：

```bash

xxx@B-P59QMD6R-2102 img % docker images
REPOSITORY          TAG                 IMAGE ID            CREATED             SIZE
layotto/layotto     0.1.0-662eab0       0370527a51a1        10 minutes ago      431MB
```

### Start Layotto

```bash
docker run -p 34904:34904 layotto/layotto:0.1.0-662eab0
```


Mac and Windows do not support --net=host, if it is on linux, you can directly replace -p 34904:34904 with --net=host.


### Start client

```bash
cd layotto/demo/configuration/etcd
go build
./etcd
```

If the following information is printed out, it means the startup is complete and Layotto is running now：

```bash
runtime client initializing for: 127.0.0.1:34904
receive hello response: greeting
get configuration after save, &{Key:hello1 Content:world1 Group:default Label:default Tags:map[] Metadata:map[]}
get configuration after save, &{Key:hello2 Content:world2 Group:default Label:default Tags:map[] Metadata:map[]}
receive watch event, &{Key:hello1 Content:world1 Group:default Label:default Tags:map[] Metadata:map[]}
receive watch event, &{Key:hello1 Content: Group:default Label:default Tags:map[] Metadata:map[]}
```

### Next step

Layotto provides the golang version of the SDK, which is located in the runtime/sdk directory. Users can directly call the services provided by Layotto through the corresponding SDK.

