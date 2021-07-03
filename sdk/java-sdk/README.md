
## How to generate a Java PROTO file

### 1. 下载对应[proto](https://github.com/protocolbuffers/protobuf/releases/tag/v3.6.1) 文件到本地

my protoc version: 
```shell
$ protoc --version
libprotoc 3.17.3
```

### 2. 修改对应`proto`文件生成类名包名等信息

appcallback.proto:
```protobuf
option java_outer_classname = "LayottoAppCallbackProtos";
option java_package = "io.mosn.layotto.v1";
```
runtime.proto
```protobuf
option java_outer_classname = "LayottoProtos";
option java_package = "io.mosn.layotto.v1";
```

### 3. 编译其对应`JAVA`文件
```shell
cd ${your PROJECT path}/spec/proto/runtime/v1
protoc -I=. --java_out=./  orderInfo.proto
```
