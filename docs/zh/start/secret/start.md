# 使用Secret API获取secret
## 什么是Secret API
Secret API用于从file、env、k8s等获取secret

Secret API支持获取单个和所有secret
## 快速开始

该示例展示了如何通过Layotto Secret API 获取 file、env、k8s中的secret


### 第一步：运行Layotto

将项目代码下载到本地后，切换代码目录、编译：

```shell
cd ${project_path}/cmd/layotto
```
构建:
```shell @if.not.exist layotto
go build -o layotto
```
完成后目录下会生成layotto文件，运行它：

```shell @background
./layotto start -c ../../configs/config_standalone.json
```

### 第二步：运行客户端程序，调用 Layotto 获取 secret

```shell
 cd ${project_path}/demo/secret/common/
```

```shell @if.not.exist client
 go build -o client
```
```shell
 ./client -s "secret_demo"
```

打印出如下信息则代表调用成功：

```bash
data:{key:"db-user-pass:password" value:"S!S*d$zDsb="}
data:{key:"db-user-pass:password" value:{secrets:{key:"db-user-pass:password" value:"S!S*d$zDsb="}}} data:{key:"db-user-pass:username" value:{secrets:{key:"db-user-pass:username" value:"devuser"}}}
```


## 想要详细了解Secret API?
Layotto复用了Dapr的Secret API，了解更多：https://docs.dapr.io/operations/components/setup-secret-store/
