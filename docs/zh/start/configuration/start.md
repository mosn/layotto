# 通过Layotto调用etcd配置中心

## 快速开始

该示例展示了如何通过Layotto，对etcd配置中心进行增删改查以及watch的过程。请提前在本机上安装[Docker](https://www.docker.com/get-started) 软件。
[config文件](https://github.com/mosn/layotto/blob/main/configs/runtime_config.json)在config_stores中定义了etcd，用户可以更改配置文件为自己想要的配置中心（目前支持etcd和apollo）。


### 生成镜像

首先请确认把layotto项目放在如下目录：

```
$GOPATH/src/github/layotto/layotto
```

然后执行如下命令：

```bash
cd $GOPATH/src/github/layotto/layotto  
make image
```

运行结束后本地会生成两个镜像：

```bash

xxx@B-P59QMD6R-2102 img % docker images
REPOSITORY          TAG                 IMAGE ID            CREATED             SIZE
layotto/layotto     0.1.0-662eab0       0370527a51a1        10 minutes ago      431MB
```

### 运行Layotto

```bash
docker run -p 34904:34904 layotto/layotto:0.1.0-662eab0
```

Mac和Windows不支持--net=host, 如果是在linux上可以直接把 -p 34904:34904 替换成 --net=host。


### 启动本地client

```bash
cd layotto/demo/configuration/etcd
go build
./etcd
```

打印出如下信息则代表启动完成：

```bash
runtime client initializing for: 127.0.0.1:34904
receive hello response: greeting
get configuration after save, &{Key:hello1 Content:world1 Group:default Label:default Tags:map[] Metadata:map[]}
get configuration after save, &{Key:hello2 Content:world2 Group:default Label:default Tags:map[] Metadata:map[]}
receive watch event, &{Key:hello1 Content:world1 Group:default Label:default Tags:map[] Metadata:map[]}
receive watch event, &{Key:hello1 Content: Group:default Label:default Tags:map[] Metadata:map[]}
```

### 拓展

Layotto 提供了golang版本的sdk，位于runtime/sdk目录下，用户可以通过对应的sdk直接调用Layotto提供的服务。

