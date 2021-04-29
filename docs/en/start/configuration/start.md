<h2>LayOtto support configuration center</h2>

## Quick start

This example shows how to add, delete, modify, and watch the etcd configuration center through LayOtto. Please install [Docker](https://www.docker.com/get-started) software on this machine in advance.
[config file](../../../../configs/runtime_config.json) defines etcd in config_stores, and users can change the configuration file to the configuration center they want (currently supports etcd and apollo).

###build docker image

```bash
  cd yourscoderectory
  
  make build-image
```

After make success, you can see two images with docker images command：

```bash

xxx@B-P59QMD6R-2102 img % docker images
REPOSITORY                                TAG                   IMAGE ID       CREATED        SIZE
runtime                                   0.1.0-94d61d8         8d0040e3e3b0   24 hours ago   439MB
mosnio/runtime                            0.1.0-94d61d8         8d0040e3e3b0   24 hours ago   439MB
```

###Start LayOtto

```bash
docker run -p 34904:34904 mosnio/runtime:0.1.0-94d61d8
```


Mac and Windows do not support --net=host, if it is on linux, you can directly replace -p 34904:34904 with --net=host.


###start client

```bash
 cd yourDir/runtime/demo/configuration/etcd
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

###Other

LayOtto provides the golang version of the SDK, which is located in the runtime/sdk directory. Users can directly call the services provided by LayOtto through the corresponding SDK.

