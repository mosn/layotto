   
# DelayQueue API demo

本示例演示如何调用 Layotto DelayQueue API.

## 什么是 DelayQueue API ?

DelayQueue is a special kind of message queue, which lets you postpone the delivery of new messages to consumers.
For example, you can invoke this API and tell the message queue "please send this message to the consumers after 5 minutes".

## step 1. 运行 Layotto
<!-- tabs:start -->
### **With Docker**
您可以用 Docker 启动 Layotto

```bash
docker run -v "$(pwd)/configs/config_standalone.json:/runtime/configs/config.json" -d  -p 34904:34904 --name layotto layotto/layotto start
```

### **本地编译（不适合 Windows)**
您可以本地编译、运行 Layotto。
> [!TIP|label: 不适合 Windows 用户]
> Layotto 在 Windows 下会编译失败。建议 Windows 用户使用 docker 部署

将项目代码下载到本地后，切换代码目录：

```shell
cd ${project_path}/cmd/layotto
```

构建:

```shell @if.not.exist layotto
go build
```

完成后目录下会生成 Layotto文件，运行它：

```shell @background
./layotto start -c ../../configs/config_standalone.json
```

<!-- tabs:end -->

## step 2. 运行客户端程序，调用 Layotto DelayQueue API
<!-- tabs:start -->
### **Go**

构建、运行 go 语言 demo:

```shell
 cd ${project_path}/demo/delay_queue/common/
 go build -o client
 ./client -s "demo"
```

打印出如下信息则代表调用成功：

```bash
TODO
```

### **Java**

下载 java sdk 和示例代码:

```shell @if.not.exist java-sdk
git clone https://github.com/layotto/java-sdk
```

```shell
cd java-sdk
```

构建 examples:

```shell @if.not.exist examples-delay_queue/target/examples-delay_queue-jar-with-dependencies.jar
# build example jar
mvn -f examples-delay_queue/pom.xml clean package
```

运行:

```shell
java -jar examples-delay_queue/target/examples-delay_queue-jar-with-dependencies.jar
```

打印出以下信息说明运行成功:

```bash
TODO
```

<!-- tabs:end -->

## step 3. 销毁容器,释放资源
<!-- tabs:start -->
### **销毁 Docker container**
如果您是用 Docker 启动的 Layotto，可以按以下方式销毁容器：

```bash
docker rm -f layotto
```

<!-- tabs:end -->

## 下一步
### 这个客户端程序做了什么？
示例客户端程序中使用了 Layotto 提供的多语言 sdk，调用Layotto DelayQueue API。

go sdk位于`sdk`目录下，java sdk 在 https://github.com/layotto/java-sdk

除了使用sdk调用Layotto提供的API，您也可以用任何您喜欢的语言、通过grpc直接和Layotto交互。

### 细节以后再说，继续体验其他API
通过左侧的导航栏，继续体验别的API吧！

### Reference

[API Reference](https://mosn.io/layotto/api/v1/delay_queue.html)

<!--design_doc_url-->

 <!-- end services -->

