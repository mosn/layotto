# Zookeeper

## metadata fields
Example: configs/config_zookeeper.json

| Field | Required | Description |
| --- | --- | --- |
| zookeeperHosts | Y | zookeeper server address, such as localhost:6380 |
| zookeeperPassword | Y | zookeeper Password |
| sessionTimeout | N | Session timeout,Unit second, same as zookeeper's sessionTimeout|
|logInfo|N|`true` means zookeeper log messages with info level should be logged; `false` means only error messages should be logged|

## Warning 
The sequencer id component of zookeeper is implemented using the version provided by zk. The version cannot exceed int32, and when overflow happens, an error will be returned and error logs will be printed. Nothing else will be processed.

It is recommended that you monitor zookeeper carefully and prevent the overflow. 

## How to start Zookeeper
If you want to run the zookeeper demo, you need to start a Zookeeper server with Docker first.

command:

```shell
docker pull zookeeper
docker run --privileged=true -d --name zookeeper --publish 2181:2181  -d zookeeper:latest
```

## Run layotto

````shell
cd ${project_path}/cmd/layotto
go build
````

>If build reports an error, it can be executed in the root directory of the project `go mod vendor`

Execute after the compilation is successful:

````shell
./layotto start -c ../../configs/config_zookeeper.json
````

## Run Demo

````shell
cd ${project_path}/demo/sequencer/zookeeper/
 go build -o client
 ./client 
````
