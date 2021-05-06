<h2>LayOtto support configuration center</h2>

[查看中文版本](../../../zh/start/configuration/start.md)

## Quick start

This example shows how to add, delete, modify, and watch the etcd configuration center through LayOtto. Please install [Docker](https://www.docker.com/get-started) software on this machine in advance.
[config file](../../../../configs/runtime_config.json) defines etcd in config_stores, and users can change the configuration file to the configuration center they want (currently supports etcd and apollo).

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

### Start LayOtto

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

If the following information is printed out, it means the startup is complete and LayOtto is running now：

```bash
runtime client initializing for: 127.0.0.1:34904
receive hello response: greeting
get configuration after save, &{Key:hello1 Content:world1 Group:default Label:default Tags:map[] Metadata:map[]}
get configuration after save, &{Key:hello2 Content:world2 Group:default Label:default Tags:map[] Metadata:map[]}
receive watch event, &{Key:hello1 Content:world1 Group:default Label:default Tags:map[] Metadata:map[]}
receive watch event, &{Key:hello1 Content: Group:default Label:default Tags:map[] Metadata:map[]}
```

### Next step

LayOtto provides the golang version of the SDK, which is located in the runtime/sdk directory. Users can directly call the services provided by LayOtto through the corresponding SDK.

