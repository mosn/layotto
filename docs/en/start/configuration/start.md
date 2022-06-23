# Configuration API demo with Etcd

This example shows that when you are using etcd as a configuration center, how to add, delete, modify, and watch the etcd through Layotto. 

The architecture of this demo is shown in the figure below. The processes started are: client APP, Layotto, etcd

![](https://gw.alipayobjects.com/mdn/rms_5891a1/afts/img/A*dzGaSb78UCoAAAAAAAAAAAAAARQnAQ)

[Then config file](https://github.com/mosn/layotto/blob/main/configs/runtime_config.json) claims `etcd` in the `config_store` section, but users can change it to other configuration center they want (currently only support etcd and apollo).
## step 1. Deploy etcd and Layotto
<!-- tabs:start -->
### **With Docker Compose**
You can start etcd and Layotto with docker-compose

```bash
cd docker/layotto-etcd
# Start etcd and layotto with docker-compose
docker-compose up -d
```

### **Compile locally (not for Windows)**
You can run etcd with Docker, then compile and run Layotto locally.

> [!TIP|label: Not for Windows users]
> Layotto fails to compile under Windows. Windows users are recommended to deploy using docker-compose
### step 1.1 Start etcd
If you want to run this demo, you need to start a etcd server first.

You can download etcd from `https://github.com/etcd-io/etcd/releases` （You can also use docker.）

start it:

```shell @background
./etcd
```

Then you can access etcd with the address `localhost:2379`.

### step 1.2 Start Layotto
Build Layotto:

```shell
cd ${project_path}/cmd/layotto
```

```shell @if.not.exist layotto
go build
```

Run it:

```shell @background
./layotto start -c ../../configs/runtime_config.json
```

<!-- tabs:end -->

## step 2. Start client APP

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
save key success
get configuration after save, &{Key:key1 Content:value1 Group:application Label:prod Tags:map[feature:print release:1.0.0] Metadata:map[]} 
get configuration after save, &{Key:haha Content:heihei Group:application Label:prod Tags:map[feature:haha release:1.0.0] Metadata:map[]} 
delete keys success
write start
receive subscribe resp store_name:"config_demo" app_id:"apollo" items:<key:"heihei" content:"heihei1" group:"application" label:"prod" tags:<key:"feature" value:"haha" > tags:<key:"release" value:"16" > >
```

## step 3. Stop containers and release resources
<!-- tabs:start -->
### **Docker Compose**
If you started etcd and Layotto with docker-compose, you can shut them down as follows:

```bash
cd ${project_path}/docker/layotto-etcd
docker-compose stop
```

### **Destroy the etcd container**
If you started etcd with Docker, you can destroy the etcd container as follows:

```shell
docker rm -f etcd
```

<!-- tabs:end -->

## Next step
### What did this demo do?
The demo client uses the golang version SDK provided by Layotto, and invokes Layotto's Configuration API to add, delete, modify, and subscribe to configuration data.

The sdk is located in the `sdk` directory. Users can invoke the Layotto API using the sdk.

In addition to using sdk, you can also interact with Layotto directly through grpc in any language you like.

In fact, sdk is only a very thin package for grpc, using sdk is about equal to directly using grpc.


### Let's continue to experience other APIs
Explore other Quickstarts through the navigation bar on the left.
