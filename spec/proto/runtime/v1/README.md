
## How to compile these proto files into golang code
```shell
cd ${your PROJECT path}/spec/proto/runtime/v1
protoc -I.  --go_out=plugins=grpc,paths=source_relative:. *.proto
```

my protoc version: 
```shell
$ protoc --version
libprotoc 3.11.2
```
