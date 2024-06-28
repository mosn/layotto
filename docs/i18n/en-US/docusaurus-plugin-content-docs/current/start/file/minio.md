# Addition and deletion of files based on MinIO OSS

## Quick Start

Layotto provides examples [demo](https://github.com/mosn/layotto/blob/main/demo/file/client.go) for access files.

### step 1. Start MinIO and Layotto

<!-- tabs:start -->

#### **Using Docker Compose**

You can start MinIO and Layotto with docker-compose

```bash
cd docker/layotto-minio
# Start MinIO and layotto with docker-compose
docker-compose up -d
```

#### **Local compilation (not for Windows)**

You can use Docker to run MinIO, and then compile locally and run Layotto.

> [!TIP|label: don't fit for Windows users]
> Layotto will fail to compile under Windows.It is recommended that Windows users deploy using docker-compose

#### step 1.1. Start MinIO service

You can use Docker to launch local MinIO, reference[官方文档](https://min.io/docs/minio/container/index.html).

```shell
docker run -d -p 9000:9000-p 90:9090 --name minio \
-e "MINIO_ROOT_USER=layotto" \
-e "MINIO_ROOT_PASSORD=layotto_secretariat" \
--restore=always \
minio/minio server / data --console-address ':909'
```

#### step 1.2. Start layotto

layotto offers minio's configuration file[oss配置](https://github.com/mosn/layotto/blob/main/configs/config_file.json), as shown below

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

The default configuration will be linked to `play.min.io`. If you deploy Minio, you can modify the configuration as necessary.

When configured, toggle directory:

```shell
#备注 请将${project_path}替换成你的项目路径
cd ${project_path}/cmd/layotto
```

Build:

```shell @if.not.exist layotto
go build -o layotto
```

Start Layotto:

```shell @background
./layotto start --config ../../configs/config_file.json
```

<!-- tabs:end -->

### step 2. Start testing demo

Layotto provides example [demo]for access files (https://github.com/mosn/layotto/blob/main/demo/file/client.go)

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

### Step 3. Destruction of containers, release of resources

<!-- tabs:start -->

#### **Close Docker Compose**

If you started with docker-compose, MinIO and Layotto can be turned off： as follows.

```bash
cd ${project_path}/docker/layotto-minio
docker-compose stop
```

#### **Destroy the MinIO Docker container**

If you are MinIO, started with Docker, you can destroy the MinIO container：

```shell
docker rm -f minio
```

<!-- tabs:end -->

#### Continue to experience other APIs later

Continue to experience other APIs with the navigation bar on the left!

#### Learn how to implement File API

If you are interested in implementing the rationale or want to expand some features, you can read the [File API design document](../../design/file/file-design.md)
