# 使用Configuration API调用nacos配置中心

该示例展示了如何通过Layotto，对nacos配置中心进行增删改查以及watch的过程。

![](../../../img/configuration/nacos/layotto-nacos-configstore-component.png)

## 第一步：启动 Nacos 和 Layotto

### 方式一：通过Docker Compose方式启动

你可以使用 [docker-compose](https://docs.docker.com/compose/) 来快速启动 layotto 和 nacos。

```shell
cd docker/layotto-nacos
# Start nacos and layotto with docker-compose
docker-compose up -d
```

### 方式二：通过本地编译启动

您可以使用 [nacos 官网文档](https://nacos.io/zh-cn/docs/quick-start-docker.html) 提供的方式启动 nacos，然后本地编译、运行 Layotto。

当然你需要按照自己的 nacos 配置，修改 `configs/config_nacos.json` 文件。

> [!TIP|label: 不适合 Windows 用户]
> Layotto 在 Windows 下会编译失败。建议 Windows 用户使用 docker-compose 部署

## 第二步：启动客户端Demo，调用Layotto增删改查

```shell
 cd ${project_path}/demo/configuration/common
```

```shell @if.not.exist client
 go build -o client
```

```shell
 ./client -s "config_demo"
```

打印出如下信息则代表调用成功：

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

## 下一步

### 这个客户端Demo做了什么？

示例客户端程序中使用了Layotto提供的golang版本sdk，调用Layotto 的Configuration API对配置数据进行增删改查、订阅变更。

sdk位于`sdk`目录下，用户可以通过sdk调用Layotto提供的API。

除了使用sdk，您也可以用任何您喜欢的语言、通过grpc直接和Layotto交互。

其实sdk只是对grpc很薄的封装，用sdk约等于直接用grpc调。


### 细节以后再说，继续体验其他API

通过左侧的导航栏，继续体验别的API吧！