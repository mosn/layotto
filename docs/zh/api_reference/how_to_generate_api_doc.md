# How to generate api documents

## Step 1. generate english document
We can use [protoc-gen-doc](https://github.com/pseudomuto/protoc-gen-doc) and docker to generate api documents,the command is as follows:  
(Run in layotto directory)
```
docker run --rm \
-v  $(pwd)/docs/en/api_reference:/out \
-v  $(pwd)/spec/proto/runtime/v1:/protos \
pseudomuto/protoc-gen-doc  --doc_opt=markdown,api_reference_v1.md
```
## Step 2. copy it to the directory for chinese docs
(Run in layotto directory)
```shell
cp docs/en/api_reference/api_reference_v1.md docs/zh/api_reference/api_reference_v1.md
```