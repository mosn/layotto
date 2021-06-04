
## How to compile these proto files
```shell
cd ${your PROJECT path}/spec/proto/runtime/v1
protoc -I.  --go_out=plugins=grpc,paths=source_relative:. *.proto
```