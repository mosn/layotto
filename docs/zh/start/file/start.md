# 基于阿里云OSS实现文件的增删改查

## 快速开始

Layotto提供了访问文件的示例 [demo](../../../../demo/file/client.go),该示例实现了文件的增删改查操作。

### 第一步：启动layotto

layotto提供了file的配置文件[oss配置](../../../../configs/config_file.json)，如下图所示

![img.png](../../../img/file/img.png)

上述配置信息需要开通[阿里云OSS](https://www.aliyun.com/product/oss) 服务。

### 第二步：启动测试demo

Layotto提供了访问文件的示例 [demo](../../../../demo/file/client.go)

```go

go build client.go //编译生成client可执行文件

./client put fileName //上传文件
./client get fileName //下载文件
./client del fileName //删除文件
./client list fileName //查看文件

```

#### 细节以后再说，继续体验其他API
通过左侧的导航栏，继续体验别的API吧！

#### 了解分布式锁 API的实现原理

如果您对实现原理感兴趣，或者想扩展一些功能，可以阅读[File API的设计文档](../../design/file/file-design.md)