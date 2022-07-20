# 使用 Secret Ref 注入secret

我们可以用secret store将secrets注入到其他组件。

用 `secret_ref` 来配置:

```json
{
  "sequencer": {
    "redis": {
      "type": "redis",
      "metadata": {
        "redisHost": "127.0.0.1:6380",
        "redisPassword": ""
      },
      "secret_ref": [
        {
          "store_name": "local.file",
          "key": "db-user-pass:password",
          "sub_key": "db-user-pass:password",
          "inject_as": "redisPassword"
        }
      ]
    }
  }
}
```

一个示例是 [config_ref_example.json](https://github.com/mosn/layotto/blob/main/configs/config_ref_example.json)

## 快速开始

该示例展示了如何注入redis password到sequencer组件

### Step 0:  运行redis并初试密码

```shell
docker run --name redis -p 6380:6379 -d --restart=always redis:5.0.3 redis-server --appendonly yes --requirepass "redis123"
```

### Step 1:  运行 Layotto

将项目代码下载到本地后，切换代码目录、编译：

```shell
cd ${project_path}/cmd/layotto
```

build:

```shell @if.not.exist layotto
go build -o layotto
```

完成后目录下会生成layotto文件，运行它：

```shell @background
./layotto start -c ../../configs/config_ref_example.json
```

### 第二步：运行客户端程序，调用 Layotto 获取sequence

```shell
 cd ${project_path}/demo/sequencer/common/
```

```shell @if.not.exist client
 go build -o client
```

```shell
 ./client -s "redis"
```

打印出如下信息则代表调用成功：

```bash
Try to get next id.Key:key666 
Next id:next_id:1 
Next id:next_id:2 
Next id:next_id:3 
Next id:next_id:4 
Next id:next_id:5 
Next id:next_id:6 
Next id:next_id:7 
Next id:next_id:8 
Next id:next_id:9 
Next id:next_id:10 
Demo success!

```