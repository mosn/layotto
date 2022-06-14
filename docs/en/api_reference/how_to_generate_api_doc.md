# How to generate `.pb.go` code and API reference
Note: the commands below should be executed under layotto directory

## How to compile the proto files into `.pb.go` code
<!-- tabs:start -->
### **Make cmmand(recommended)**
```bash
make proto.code
```
This command uses docker to run protoc and generate `.pb.go` code files.

### **Install protoc**
1. Install protoc version: [v3.17.3](https://github.com/protocolbuffers/protobuf/releases/tag/v3.17.3)

2. Install protoc-gen-go v1.28 and protoc-gen-go-grpc v1.2

3. Generate gRPC `.pb.go` code

```bash
cd spec/proto/runtime/v1
protoc -I. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=require_unimplemented_servers=false,paths=source_relative *.proto
```
<!-- tabs:end -->
## How to generate API reference doc according to the proto files
We can use [protoc-gen-doc](https://github.com/pseudomuto/protoc-gen-doc) and docker to generate api documents,the command is as follows:  

<!-- tabs:start -->
### **Make command(recommended)**
```bash
make proto.doc
```
This command uses docker to run protoc-gen-doc and generate docs.

### **Use docker to run protoc-gen-doc**
`make proto.doc` essentially run commands below:

```
docker run --rm \
-v  $(pwd)/docs/en/api_reference:/out \
-v  $(pwd)/spec/proto/runtime/v1:/protos \
pseudomuto/protoc-gen-doc  --doc_opt=/protos/template.tmpl,runtime_v1.md runtime.proto
```

and 

```shell
docker run --rm \
-v  $(pwd)/docs/en/api_reference:/out \
-v  $(pwd)/spec/proto/runtime/v1:/protos \
pseudomuto/protoc-gen-doc  --doc_opt=/protos/template.tmpl,appcallback_v1.md appcallback.proto
```

<!-- tabs:end -->