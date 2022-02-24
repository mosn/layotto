# 腾讯云 对象存储

## 配置项说明

示例：configs/config_file_tencentcloud_oss.json

| 字段 | 必填  | 说明                                        |
| --- |-----|-------------------------------------------|
| endpoint | Y   | oss bucket url 不需要带协议 http:// 或者 https:// |
| accessKeyID | Y   | 通行ID                                      |
| accessKeySecret | Y   | 通行密码                                      |
| timeout | N   | 和腾讯云OSS服务交互超时时间，单位毫秒，默认值100秒              |

## 启动准备

1. 登录腾讯云 https://cloud.tencent.com/

2. 创建存储桶

访问 https://console.cloud.tencent.com/cos/bucket 创建存储桶

![](../../../img/file/create_tencent_oss_bucket.png)

3.创建 AK 和 SK

访问 https://console.cloud.tencent.com/cam/capi 进行创建

以上操作步骤完成后，将 endpoint、AK、SK 配置到 `configs/config_file_tencentcloud_oss.json` 文件中

## 启动 layotto

````shell
cd cmd/layotto_multiple_api/
go build -o layotto
./layotto start -c ../../configs/config_file_tencentcloud_oss.json
````

## 运行 Demo

````shell
cd demo/file/tencentcloud
go build -o oss

./oss put dir/a.txt aaa #创建文件
./oss get dir/a.txt #获取文件
./oss list dir/ #列出目录下文件
./oss del dir/a.txt #删除文件
````