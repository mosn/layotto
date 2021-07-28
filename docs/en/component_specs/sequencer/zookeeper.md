# Zookeeper

## metadata fields
Example: configs/config_sequencer_zookeeper.json

| Field | Required | Description |
| --- | --- | --- |
| zookeeperHosts | Y | zookeeper server address, such as localhost:6380 |
| zookeeperPassword | Y | zookeeper Password |
| sessionTimeout | N | Session timeout,Unit second, same as zookeeper's sessionTimeout|
|logInfo|N|true if zookeeper information messages are logged; false if only zookeeper errors are logged|

## How to start Zookeeper
If you want to run the zookeeper demo, you need to start a Zookeeper server with Docker first.

command:
```shell
docker pull zookeeper
docker run --privileged=true -d --name zookeeper --publish 2181:2181  -d zookeeper:latest
```

## Run layotto

````shell
cd ${projectpath}/cmd/layotto
go build
````
>If build reports an error, it can be executed in the root directory of the project `go mod vendor`

Execute after the compilation is successful:
````shell
./layotto start -c ../../configs/config_sequencer_zookeeper.json
````

## Run Demo

````shell
cd ${projectpath}/demo/sequencer/zookeeper/
 go build -o client
 ./client 
````
