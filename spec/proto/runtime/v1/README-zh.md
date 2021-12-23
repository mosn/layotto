## 如何把这些proto文件编译成go代码
在项目根目录下，编译成go代码:
```shell
protoc -I=./spec/ \
  --go_out=./spec/ \
  --go_opt paths=source_relative \
  --go-grpc_out require_unimplemented_servers=false:./spec/\
  --go-grpc_opt paths=source_relative \
  --grpc-gateway_out=./spec/ \
  --grpc-gateway_opt paths=source_relative \
  ./spec/proto/runtime/v1/*.proto
```

注：我的 protoc 版本是:
```shell
$ protoc --version
libprotoc 3.17.3
```

### 注意：踩过的坑

1. 相对路径会影响编译出来的东西。比如相对路径用`/spec/proto/`,这么写：
```shell
protoc -I=./spec/proto/ \
  --go_out=./spec/proto/ \
  --go_opt paths=source_relative \
  --go-grpc_out require_unimplemented_servers=false:./spec/proto/\
  --go-grpc_opt paths=source_relative \
  --grpc-gateway_out=./spec/proto/ \
  --grpc-gateway_opt paths=source_relative \
  ./spec/proto/runtime/v1/*.proto
```

这么写编译出来的东西会启动报错、遇到protobuf的一个bug：同名Service、路径不一样的话，注册grpc服务器会报错，这是个bug,但是新版本的protobuf还没修，见
https://stackoverflow.com/questions/67693170/proto-file-is-already-registered-with-different-packages 

解决方案是修改编译时的相对路径，改成`/spec`，这样Service的名字就不一样了，绕过了这个Bug