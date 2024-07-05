# 基于MinIO OSS实现文件的增删改查

## 快速开始

Layotto提供了访问文件的示例 [demo](https://github.com/mosn/layotto/blob/main/demo/file/client.go) ,该示例实现了文件的增删改查操作。

### step 1. 启动 MinIO 和 Layotto
<!-- tabs:start -->
#### **使用 Docker Compose**
您可以使用 docker-compose 启动 MinIO 和 Layotto

```bash
cd docker/layotto-minio
# Start MinIO and layotto with docker-compose
docker-compose up -d
```

#### **本地编译（不适合 Windows)**
您可以使用 Docker 运行 MinIO，然后本地编译、运行 Layotto。

> [!TIP|label: 不适合 Windows 用户]
> Layotto 在 Windows 下会编译失败。建议 Windows 用户使用 docker-compose 部署
#### step 1.1. 启动 MinIO 服务
您可以使用 Docker 启动本地MinIO服务, 参考[官方文档](https://min.io/docs/minio/container/index.html)

```shell
docker run -d -p 9000:9000 -p 9090:9090 --name minio \
-e "MINIO_ROOT_USER=layotto" \
-e "MINIO_ROOT_PASSWORD=layotto_secret" \
--restart=always \
minio/minio server /data --console-address ':9090'
```

#### step 1.2. 启动layotto

layotto提供了minio的配置文件[oss配置](https://github.com/mosn/layotto/blob/main/configs/config_file.json) ，如下所示

```json
                      "file": {
                        "minio": {
                          "metadata":[
                            {
                              "endpoint": "play.min.io",
                              "accessKeyID": "Q3AM3UQ867SPQQA43P2F",
                              "accessKeySecret": "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG",
                              "SSL":true,
                              "region":"us-east-1"
                            }
                          ]
                        }
                      }
```

默认配置会连接`play.min.io`, 如果您自己部署了 Minio, 可以按需修改其中的配置。

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
./layotto start --config ../../configs/config_file.json
```
<!-- tabs:end -->

### step 2. 启动测试demo

Layotto提供了访问文件的示例 [demo](https://github.com/mosn/layotto/blob/main/demo/file/client.go)

```shell
cd ${project_path}/demo/file
go build client.go

# 创建名为test的bucket
./client bucket test
# 上传文件到test bucket，前缀为hello，内容为"hello layotto"
./client put test/hello/layotto.txt "hello layotto"
# 获取 layotto.txt的内容
./client get test/hello/layotto.txt
# 获取test bucket下的前缀为hello的所有文件列表
./client list test/hello
# 获取layotto.txt文件的元数据
./client stat test/hello/layotto.txt
# 删除layotto.txt文件
./client del test/hello/layotto.txt
```

### step 3. 销毁容器，释放资源
<!-- tabs:start -->
#### **关闭 Docker Compose**
如果您是用 docker-compose 启动的 MinIO 和 Layotto，可以按以下方式关闭：

```bash
cd ${project_path}/docker/layotto-minio
docker-compose stop
```

#### **销毁 MinIO Docker 容器**
如果您是用 Docker 启动的 MinIO，可以按以下方式销毁 MinIO 容器：

```shell
docker rm -f minio
```
<!-- tabs:end -->

#### 细节以后再说，继续体验其他API
通过左侧的导航栏，继续体验别的API吧！

#### 了解File API的实现原理

如果您对实现原理感兴趣，或者想扩展一些功能，可以阅读[File API的设计文档](zh/design/file/file-design.md)
