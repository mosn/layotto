# How to generate `.pb.go` code and corresponding documentation
Note: the commands below should be executed under layotto directory

```shell
make proto
```

Then you get:
- Generated code
    - `.pb.go` code
    - `_grpc.pb.go` code
    - layotto go-sdk code
    - layotto sidecar code
- Generated documentation
    - API reference docs
    - updated API reference list
    - quickstart document (both chinese and english)
    - updated sidebar (The tool will add the generated quickstart doc into the sidebar of https://mosn.io/layotto )
- Updated CI (The tool will add the generated quickstart doc into the CI script `etc/script/test-quickstart.sh`)

That's all :)

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
`make proto.doc` invokes the script `etc/script/generate-doc.sh`, which uses docker to run protoc-gen-doc.

You can check `etc/script/generate-doc.sh` for more details.

<!-- tabs:end -->