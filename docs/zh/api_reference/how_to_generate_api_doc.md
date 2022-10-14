# 如何基于proto文件生成代码、文档

假设你写了一个新的proto文件`spec/proto/extension/v1/email/email.proto`，并且想在Layotto中实现这个API。

```protobuf
// EmailService 用来发送电子邮件
service EmailService {

  // 发送带有模板的电子邮件
  rpc SendEmailWithTemplate(SendEmailWithTemplateRequest) returns (SendEmailWithTemplateResponse) {}

  // 不使用模板，发送带有原始信息的电子邮件
  rpc SendEmail(SendEmailRequest) returns (SendEmailResponse) {}

}

// 不同的message类型......
```

这将会是一项非常繁琐的工作，你必须为此编写大量的代码和文档。

幸运的是，Layotto通过工具可以自动生成代码/文档/CI配置，这将为你节省大量的时间！

## 步骤一：请确保你的proto文件满足以下要求
- 文件路径应该为`spec/proto/extension/v1/{api short name}/{api short name}.proto`
- 每个proto文件中只能有一个`service`。下面是一个错误的示例:

```protobuf
//  EmailService 用来发送电子邮件
service EmailService {
  // ...
}

// Wrong: 在.proto文件中应该只有一个service。
service EmailService2 {
  // ...
}

// 不同的message类型......
```

- 如果你不想为proto生成quickstart docs，添加注释`/* @exclude skip quickstart_generator */` 。
- 如果你不想为proto生成sdk和sidecar代码，添加注释`/* @exclude skip code_generator */` 。

你可以把 `spec/proto/extension/v1/s3/oss.proto` 作为一个例子 :

```protobuf
/* @exclude skip quickstart_generator */
/* @exclude skip code_generator */
// ObjectStorageService是对blob存储或所谓的 "对象存储 "的抽象，例如阿里云OSS，AWS S3。
// 调用ObjectStorageService API对二进制文件进行一些CRUD操作，例如，对文件进行查询，删除操作等。
service ObjectStorageService{
  //......
}
```

这些特殊的注释被称为 "Master's commands"，还有许多其他的命令，你可以查看[文档](https://github.com/layotto/protoc-gen-p6#masters-commands)了解更多细节。

## 步骤二：检查环境

要运行生成器，你需要满足如下条件：
- Go语言版本 >=1.16
- 启用Docker

## 步骤三：开始生成

```shell
make proto
```

执行命令后，你将会得到：

- 生成的代码
  - `.pb.go`
  - `_grpc.pb.go`
  - layotto go-sdk
  - layotto sidecar 代码（实现了新的 API ）
- 生成的文档
  - API 参考文档
  - 自动更新 API 文档列表
  - quickstart 文档（包括中文和英文）
  - 自动更新侧边栏（该工具将把生成的快速入门文档添加到https://mosn.io/layotto 的侧边栏中）
- 自动更新 CI（该工具将把生成的快速入门文档添加到CI脚本`etc/script/test-quickstart.sh`中）
## 步骤四：编写其余的代码
现在，你需要完成如下代码的编写工作：

- Layotto component
- go examples

![image](https://user-images.githubusercontent.com/26001097/188782762-bc1404a8-b891-45d3-a1ac-f86cafdbc0ab.png)

- java examples

![image](https://user-images.githubusercontent.com/26001097/188782989-9aec893f-9d12-4ee6-9a64-940b0ba1ba1b.png)

## 实现原理
我们有一个叫做[protoc-gen-p6](https://github.com/layotto/protoc-gen-p6)的protoc插件，用于为Layotto生成代码。 

## 如果只想生成pb/documentataion怎么办？
上面的步骤生成了所有的文件，但如果只想生成`.pb.go`代码怎么办？如果只想生成文档呢？

### 如何把 proto 文件编译成`.pb.go`代码
<!-- tabs:start -->
#### **Make 命令生成(推荐)**
本地启动 docker 后，在项目根目录下运行：

```bash
make proto-code
```

该命令会用 docker 启动 protoc，生成`.pb.go`代码。

这种方式更方便，开发者不需要修改本地 protoc 版本，省去了很多烦恼。

#### **手动安装工具**
1. 安装 protoc version: [v3.17.3](https://github.com/protocolbuffers/protobuf/releases/tag/v3.17.3)

2. 安装 protoc-gen-go v1.28 和 protoc-gen-go-grpc v1.2
 
3. 生成gRPC `.pb.go`

```bash
cd spec/proto/runtime/v1
protoc -I. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=require_unimplemented_servers=false,paths=source_relative *.proto
```
<!-- tabs:end -->
### 如何基于proto文件生成接口文档

我们可以用[protoc-gen-doc](https://github.com/pseudomuto/protoc-gen-doc) 和docker来生成接口文档，相关命令如下：

<!-- tabs:start -->
#### **Make 命令生成(推荐)**
本地启动 docker 后，在项目根目录下运行：

```bash
make proto-doc
```

该命令会用 docker 启动 protoc-gen-doc，生成文档

#### **用 docker 启动 protoc-gen-doc**
`make proto-doc` 调用了脚本 `etc/script/generate-doc.sh`,这个脚本的作用是使用docker运行protoc-gen-doc.

你可以在 `etc/script/generate-doc.sh` 查看更多细节。

<!-- tabs:end -->
