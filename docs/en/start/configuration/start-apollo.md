# Configuration API demo with apollo
This example shows how to add, delete, modify, and watch the [apollo configuration center](https://github.com/apolloconfig/apollo) through Layotto.

The architecture of this example is shown in the figure below. The processes started are: client APP, Layotto, Apollo server

![img.png](../../../img/configuration/apollo/arch.png)

## Step 1.Deploy Apollo (optional)

You can skip this step, you don't need to deploy the apollo server yourself to use this demo. This demo will use the demo environment provided by [apollo official](https://github.com/apolloconfig/apollo): http://81.68.181.139

If you have deployed apollo yourself, you can modify Layotto's config file (e.g. configs/config_apollo.json in the project) to change the apollo server address to your own.

## Step 2. Run Layotto server

Download the project code to the local:

```bash
git clone https://github.com/mosn/layotto.git
```

Switch the code directory and compile:

```shell
cd ${project_path}/cmd/layotto
```

```shell @if.not.exist layotto
go build
```

After success, a new layotto file will be generated in the directory. Let's run it:

```shell @background
./layotto start -c ../../configs/config_apollo.json
```

>Q: The demo report an error?
>
>A: With the default configuration, Layotto will connect to apollo's demo server, but the configuration in that demo server may be modified by others. So the error may be because some configuration has been modified.
>
> In this case, you can try other demos, such as [the etcd demo](en/start/configuration/start.md)

## Step 3. Run the client demo 

The client demo calls Layotto to add, delete, modify, and query configuration

```shell
 cd ${project_path}/demo/configuration/common
```

```shell @if.not.exist client
 go build -o client
```

```shell
 ./client -s "config_demo"
```

If the following information is printed, the call is successful：

```bash
save key success
get configuration after save, &{Key:key1 Content:value1 Group:application Label:prod Tags:map[feature:print release:1.0.0] Metadata:map[]} 
get configuration after save, &{Key:haha Content:heihei Group:application Label:prod Tags:map[feature:haha release:1.0.0] Metadata:map[]} 
delete keys success
write start
receive subscribe resp store_name:"config_demo" app_id:"apollo" items:<key:"heihei" content:"heihei1" group:"application" label:"prod" tags:<key:"feature" value:"haha" > tags:<key:"release" value:"16" > >
```

### Next step
#### What did this client Demo do?
The demo client program uses the golang version SDK provided by Layotto, and calls Layotto's Configuration API to add, delete, modify, and subscribe to configuration data.

The sdk is located in the `sdk` directory, and users can call the API provided by Layotto through the sdk.

In addition to using sdk, you can also interact with Layotto directly through grpc in any language you like.

In fact, sdk is only a very thin package for grpc, using sdk is about equal to directly using grpc.


#### Details later, let's continue to experience other APIs
Explore other Quickstarts through the navigation bar on the left.
