# 基于S3协议实现对象存储的无感迁移

## 快速开始

Layotto提供了访问OSS的示例 [demo](https://github.com/mosn/layotto/blob/main/demo/oss/client.go) ,该示例基于S3协议实现了对象的一系列操作，当前
已支持部分接口。可以做到在不同的OSS实例之间进行无感迁移。

### step 1.  启动layotto

layotto提供了aws的配置文件`configs/config_oss.json`，配置文件内容如下所示:

```json
"grpc_config": {
  "oss": {
    "oss_demo": {
      "type": "aws.oss",
      "metadata": {
        "basic_config":{
          "region": "your-oss-resource-region",
          "endpoint": "your-oss-resource-endpoint",
          "accessKeyID": "your-oss-resource-accessKeyID",
          "accessKeySecret": "your-oss-resource-accessKeySecret"
        }
      }
    }
  }
}
```

配置中对应的字段，需要替换成自己的OSS账号的配置。type 支持多种类型，例如 `aliyun.oss`对应阿里云的OSS服务, `aws.oss` 对应亚马逊云的 S3 服务。
用户可以根据自己的实际场景进行配置。

配置好后，切换目录:

```shell
#备注 请将${project_path}替换成你的项目路径
cd ${project_path}/cmd/layotto
```

构建:

```shell @if.not.exist layotto
go build -o layotto
```

启动 Layotto:

```shell @background
./layotto start -c ../../configs/config_oss.json
```

### step 2. 启动测试demo

Layotto提供了访问文件的示例 [demo](https://github.com/mosn/layotto/blob/main/demo/oss/client.go)

```shell
cd ${project_path}/demo/oss/
go build client.go

# 上传名为test3.txt的文件到名为antsys-wenxuwan的bucket下，内容为"hello"
./client put antsys-wenxuwan test3.txt "hello"

# 获取antsys-wenxuwan bucket下名为test3.txt的文件
./client get antsys-wenxuwan test3.txt

# 删除antsys-wenxuwan bucket下名为test3.txt的文件
./client del antsys-wenxuwan test3.txt

# 返回antsys-wenxuwan bucket下的所有文件信息
./client list antsys-wenxuwan

```

#### 细节以后再说，继续体验其他API
通过左侧的导航栏，继续体验别的API吧！
