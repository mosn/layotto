# Tencent Cloud OSS

## metadata fields

Example：configs/config_file_qiniu_oss.json

| Field | Required | Description                  |
| --- |----------|------------------------------|
| endpoint | Y        | OSS bucket domin             |
| accessKeyID | Y        | accessKeyID                  |
| accessKeySecret | Y        | accessKeySecret              |
| bucket | Y        | bucket name                  |
| private | N        | whether to use private space |
| useHTTPS | N        | whether to use http domain   |
| useCdnDomains | N        | whether to use cdn           |

## Prepare

1.create bucket https://portal.qiniu.com/kodo/bucket

2.get keys https://portal.qiniu.com/user/key

After the above operation steps are completed, configure endpoint, AK and SK to `configs/config_file_qiniu_oss.json`
file

## Run layotto

````shell
cd cmd/layotto_multiple_api/
go build -o layotto
./layotto start -c ../../configs/config_file_qiniu_oss.json
````

## Run Demo

````shell
cd demo/file/qiniu
go build -o oss

./oss put dir/a.txt aaa 
./oss get dir/a.txt 
./oss list dir/
./oss del dir/a.txt
````