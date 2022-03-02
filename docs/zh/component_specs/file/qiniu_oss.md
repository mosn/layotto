# 七牛云 对象存储

## 配置项说明

示例：configs/config_file_qiniu_oss.json

| 字段              | 必填  | 说明                                        |
|-----------------|-----|-------------------------------------------|
| endpoint        | Y   | 七牛云对象存储空间绑定的域名 不需要带协议 http:// 或者 https:// |
| accessKeyID     | Y   | 通行ID                                      |
| accessKeySecret | Y   | 通行密码                                      |
| bucket          | Y   | 存储空间名称                                    |
| private         | N   | 是否为私有空间                                   |
| useHTTPS        | N   | 是否使用 http                                 |
| useCdnDomains   | N   | 是否使用 cdn 加速                               |

## 启动准备

1.七牛云并创建存储空间 https://portal.qiniu.com/kodo/bucket

2.获取Keys https://portal.qiniu.com/user/key


以上操作步骤完成后，将 endpoint、AK、SK 配置到 `configs/config_file_qiniu_oss.json` 文件中

## 启动 layotto

````shell
cd cmd/layotto_multiple_api/
go build -o layotto
./layotto start -c ../../configs/config_file_qiniu_oss.json
````

## 运行 Demo

````shell
cd demo/file/qiniu
go build -o oss

./oss put dir/a.txt aaa #创建文件
./oss get dir/a.txt #获取文件
./oss list dir/ #列出目录下文件
./oss del dir/a.txt #删除文件
````