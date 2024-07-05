# In-Memory

## 配置项说明

直接使用配置：configs/config_standalone.json


## 启动 layotto

````shell
cd ${project_path}/cmd/layotto
go build
````

编译成功后执行:

````shell
./layotto start -c ../../configs/config_standalone.json
````

## 运行 Demo

````shell
cd ${project_path}/demo/lock/in-memory/
 go build -o client
 ./client
````