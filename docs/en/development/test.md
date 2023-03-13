Kitex 代码生成依赖于 thriftgo 和 protoc，需要先安装相应的编译器：[thriftgo](https://github.com/cloudwego/thriftgo) 或 [protoc](https://github.com/protocolbuffers/protobuf/releases)。

安装完上述工具后，通过 go 命令安装命令行工具本身

```
go install github.com/cloudwego/kitex/tool/cmd/kitex@latest
```

你也可以自己下载 Kitex 源码后，进入 `tool/cmd/kitex` 目录执行 `go install` 进行安装

完成后，可以通过执行  `kitex -version`  查看工具版本，或者  `kitex -help` 查看使用帮助。

## 生成代码

生成代码分两部分，一部分是结构体的编解码序列化代码，由底层编译器 thriftgo 或 protoc 生成；另一部分由 kitex 工具在前者产物上叠加，生成用于创建和发起 RPC 调用的桩代码。用户只需要执行 Kitex 代码生成工具，底层会自动完成所有代码的生成。