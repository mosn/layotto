
## How to compile these proto files
```shell
cd ${your PROJECT path}/spec/proto/runtime/v1
protoc --go_out=paths=source_relative:. --go-grpc_out=. --go-grpc_opt=require_unimplemented_servers=false,paths=source_relative  --go-grpc_opt=require_unimplemented_servers=false *.proto
```