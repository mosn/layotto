# Configuration API demo with Etcd

This example shows that when you are using etcd as a configuration center, how to add, delete, modify, and watch the etcd through Layotto. 

The architecture of this demo is shown in the figure below. The processes started are: client APP, Layotto, etcd

![](https://gw.alipayobjects.com/mdn/rms_5891a1/afts/img/A*dzGaSb78UCoAAAAAAAAAAAAAARQnAQ)

[Then config file](https://github.com/mosn/layotto/blob/main/configs/runtime_config.json) claims `etcd` in the `config_store` section, and users can change the file to the configuration center they want (currently supports etcd and apollo).

## How to start etcd
If you want to run the etcd demo, you need to start a etcd server.

Steps：

download etcd from `https://github.com/etcd-io/etcd/releases` （You can also use docker.）

start：
````shell
./etcd
````

Then you can access etcd with the address `localhost:2379`.

## Run layotto

````shell
cd ${your project path}/cmd/layotto
go build
````

Execute after the compilation is successful:
````shell
./layotto start -c ../../configs/runtime_config.json
````

### Start client

```bash
cd ${your project path}/demo/configuration/etcd
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

## Next step
### What did this client Demo do?
The demo client uses the golang version SDK provided by Layotto, and invokes Layotto's Configuration API to add, delete, modify, and subscribe to configuration data.

The sdk is located in the `sdk` directory. Users can invoke the Layotto API using the sdk.

In addition to using sdk, you can also interact with Layotto directly through grpc in any language you like.

In fact, sdk is only a very thin package for grpc, using sdk is about equal to directly using grpc.


### Let's continue to experience other APIs
Explore other Quickstarts through the navigation bar on the left.
