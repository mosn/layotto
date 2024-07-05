# Tencent Cloud OSS

## metadata fields

Example：configs/config_file_tencentcloud_oss.json

| Field | Required  | Description                     |
| --- |-----|---------------------------------|
| endpoint | Y   | OSS server address              |
| accessKeyID | Y   | accessKeyID                     |
| accessKeySecret | Y   | accessKeySecret                 |
| timeout | N   | request timeout in milliseconds |

## Prepare

1. Login https://cloud.tencent.com/

2. Create Bucket

visit https://console.cloud.tencent.com/cos/bucket to create bucket

![](../../../img/file/create_tencent_oss_bucket.png)

3.Create AK and SK

visit https://console.cloud.tencent.com/cam/capi to create AK and SK

After the above operation steps are completed, configure endpoint, AK and SK to `configs/config_file_tencentcloud_oss.JSON` file

## Run layotto

````shell
cd cmd/layotto_multiple_api/
go build -o layotto
./layotto start -c ../../configs/config_file_tencentcloud_oss.json
````

## Run Demo

````shell
cd demo/file/tencentcloud
go build -o oss

./oss put dir/a.txt aaa 
./oss get dir/a.txt 
./oss list dir/
./oss del dir/a.txt
````