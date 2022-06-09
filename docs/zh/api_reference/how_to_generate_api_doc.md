# 如何基于proto文件生成接口文档


我们可以用[protoc-gen-doc](https://github.com/pseudomuto/protoc-gen-doc) 和docker来生成接口文档，相关命令如下：

(需要在layotto项目下运行命令)

<!-- tabs:start -->
#### **Make 命令生成**
```bash
make proto.doc
```
该命令会用 docker 启动 protoc-gen-doc，生成文档

#### **用 docker 启动 protoc-gen-doc**
`make proto.doc` 等价于执行以下命令:

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