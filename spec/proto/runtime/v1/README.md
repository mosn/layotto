[中文](README-zh.md)
## How to compile these proto files into golang code
```shell
cd ${your PROJECT path}

# compile
protoc -I=./spec/ \
  --go_out=./spec/ \
  --go_opt paths=source_relative \
  --go-grpc_out require_unimplemented_servers=false:./spec/\
  --go-grpc_opt paths=source_relative \
  --grpc-gateway_out=./spec/ \
  --grpc-gateway_opt paths=source_relative \
  ./spec/proto/runtime/v1/*.proto
```

my protoc version: 
```shell
$ protoc --version
libprotoc 3.17.3
```
