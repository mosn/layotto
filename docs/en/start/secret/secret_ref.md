# Reference secrets in component configuration

we can inject secrets to other components using secret store.

Use `secret_ref` in your config:

```json
{
  "sequencer": {
    "redis": {
      "type": "redis",
      "metadata": {
        "redisHost": "127.0.0.1:6380",
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

An example is [config_ref_example.json](https://github.com/mosn/layotto/blob/main/configs/config_ref_example.json)

## Quick start

This example shows how to inject redis password to sequencer component using redis store

### Step 0:  Run Redis with password

```shell
docker run --name redis -p 6380:6379 -d --restart=always redis:5.0.3 redis-server --appendonly yes --requirepass "redis123"
```

### Step 1:  Run Layotto

After downloading the project code to the local, switch the code directory and compile:

```shell
cd ${project_path}/cmd/layotto
```

build:

```shell @if.not.exist layotto
go build -o layotto
```

Once finished, the layotto file will be generated in the directory, run it:

```shell @background
./layotto start -c ../../configs/config_ref_example.json
```

### Step 2: Run the client program and call Layotto to get the sequence

```shell
 cd ${project_path}/demo/sequencer/common/
```

```shell @if.not.exist client
 go build -o client
```

```shell
 ./client -s "redis"
```

If the following information is printed, the demo is successful:

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