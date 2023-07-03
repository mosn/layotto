# Configuration API demo with Nacos

TThis example shows how to add, delete, modify, and watch the [nacos configuration center](https://nacos.io/zh-cn/index.html) through Layotto.

The architecture of this example is shown in the figure below. The processes started are: client APP, Layotto, Nacos server

![](../../../img/configuration/nacos/layotto-nacos-configstore-component.png)

[Then config file](https://github.com/mosn/layotto/blob/main/configs/config_nacos.json) claims `nacos` in the `config_store` section, but users can change it to other configuration center they want (currently only support nacos and nacos).
## Deploy  Nacos and Layotto
<!-- tabs:start -->
### Method 1: Start through Docker Compose
You can start nacos and Layotto with [docker-compose](https://docs.docker.com/compose/)

```bash
cd docker/layotto-nacos
# Start nacos and layotto with docker-compose
docker-compose up -d
```

### Method 2: Start through local compilation

You can start Nacos using the methods provided in the [Nacos official documentation](https://nacos.io/zh-cn/docs/quick-start-docker.html) and then compile and run Layotto locally.

> [!TIP|label: Not for Windows users]
> Layotto fails to compile under Windows. Windows users are recommended to deploy using docker-compose


## Start client APP

```shell
 cd ${project_path}/demo/configuration/common
```

```shell @if.not.exist client
 go build -o client
```

```shell
 ./client -s "config_demo"
```

If the following information is printed out, it means the client app has done all the CRUD operations successfully：

```bash
runtime client initializing for: 127.0.0.1:34904
save key success
get configuration after save, &{Key:key1 Content:value1 Group:application Label: Tags:map[] Metadata:map[]} 
get configuration after save, &{Key:haha Content:heihei Group:application Label: Tags:map[] Metadata:map[]} 
delete keys success
write start
receive subscribe resp store_name:"config_demo"  app_id:"testApplication_yang"  items:{key:"heihei"  content:"heihei1"  group:"application"}
write start
receive subscribe resp store_name:"config_demo"  app_id:"testApplication_yang"  items:{key:"heihei"  content:"heihei2"  group:"application"}
write start
receive subscribe resp store_name:"config_demo"  app_id:"testApplication_yang"  items:{key:"heihei"  content:"heihei3"  group:"application"}
write start
receive subscribe resp store_name:"config_demo"  app_id:"testApplication_yang"  items:{key:"heihei"  content:"heihei4"  group:"application"}
```


## Next step
### What did this demo do?
The demo client uses the golang version SDK provided by Layotto, and invokes Layotto's Configuration API to add, delete, modify, and subscribe to configuration data.

The sdk is located in the `sdk` directory. Users can invoke the Layotto API using the sdk.

In addition to using sdk, you can also interact with Layotto directly through grpc in any language you like.

In fact, sdk is only a very thin package for grpc, using sdk is about equal to directly using grpc.


### Let's continue to experience other APIs
Explore other Quickstarts through the navigation bar on the left.
