# In-Memory

## Example

You can use `configs/config_standalone.json`

## Start Layotto
Download Layotto:

```shell
git clone https://github.com/mosn/layotto.git
```

Change directory:

```shell
cd cd layotto/cmd/layotto
```

Build:

```shell @if.not.exist layotto
go build
```

Run Layotto:

```shell @background
./layotto start -c ../../configs/config_standalone.json
```

## Run Demo

```shell
cd ${project_path}/demo/sequencer/in-memory/
 go build -o client
 ./client
```

And you will see:

```bash
runtime client initializing for: 127.0.0.1:34904
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