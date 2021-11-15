## 如何使用java sdk
### 1. import sdk
对于 Maven 项目，将以下配置添加进 `pom.xml` 文件:
```xml
<project>
  ...
  <dependencies>
    ...
    <dependency>
        <groupId>io.mosn.layotto</groupId>
        <artifactId>runtime-sdk-parent</artifactId>
        <version>1.0.0-SNAPSHOT</version>
    </dependency>
    ...
  </dependencies>
  ...
</project>
```

### 2. 运行 examples 示例
可以本地部署redis和Layotto，然后运行java应用示例，通过java sdk调Layotto，Layotto转发给redis

#### 第一步：部署redis

1. 取最新版的 Redis 镜像。
   这里我们拉取官方的最新版本的镜像：

```shell
docker pull redis:latest
```

2. 查看本地镜像
   使用以下命令来查看是否已安装了 redis：

```shell
docker images
```

3. 运行容器

安装完成后，我们可以使用以下命令来运行 redis 容器：

```shell
docker run -itd --name redis-test -p 6380:6379 redis
```

参数说明：

-p 6380:6379：映射容器服务的 6379 端口到宿主机的 6380 端口。外部可以直接通过宿主机ip:6380 访问到 Redis 的服务。

#### 第二步：构建并运行Layotto

clone仓库到本地:

```sh
git clone https://github.com/mosn/layotto.git
```

构建并运行Layotto:

```bash
# make sure you replace this` ${projectpath}` with your own project path.
cd ${projectpath}/cmd/layotto
go build
./layotto start -c ../../configs/config_redis.json
```

构建java-sdk [Maven](https://maven.apache.org/install.html) (Apache Maven version 3.x) 项目:

```sh
# make sure you replace this` ${projectpath}` with your own project path.
cd ${projectpath}/sdk/java-sdk
mvn clean install
```

#### 第三步：运行java sdk示例
通过以下Examples示例来了解如何使用SDK:
* [Hello world](./examples/src/main/java/io/mosn/layotto/examples/helloworld)
* [State management](./examples/src/main/java/io/mosn/layotto/examples/state)
* [Pubsub API](./examples/src/main/java/io/mosn/layotto/examples/pubsub)


## 如何将proto文件编译成java代码

### 1. 下载编译工具 [protoc](https://github.com/protocolbuffers/protobuf/releases)
my protoc version:
```shell
$ protoc --version
libprotoc 3.11.2
```

### 2. 修改对应`proto`文件生成类名包名等信息
(需先修改文件内部service名)
`spec/proto/runtime/v1/appcallback.proto` : 
```protobuf
option java_outer_classname = "AppCallbackProto";
option java_package = "spec.proto.runtime.v1";
```
`spec/proto/runtime/v1/runtime.proto` :
```protobuf
option java_outer_classname = "RuntimeProto";
option java_package = "spec.proto.runtime.v1";
```

### 3. 编译其对应`JAVA`文件
```shell
cd ${your PROJECT path}/spec/proto/runtime/v1
protoc -I=. --java_out=./  runtime.proto
```

PS: 建议用maven插件`protoc-gen-grpc-java`生成protobuf和grpc的java代码
