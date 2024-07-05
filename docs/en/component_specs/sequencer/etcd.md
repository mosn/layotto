# Etcd

## metadata fields
Example: configs/runtime_config.json

| Field | Required | Description |
| --- | --- | --- |
| endpoints | Y | etcd server address, multiple address use `;` separate |
| dialTimeout | N | dialTimeout is the timeout for failing to establish a connection in seconds. default: 5 |
| username | N | etcd auth username |
| password | N | etcd auth password |
| keyPrefixPath | N | sequencer key prefix in etcd, default: `/layotto_sequencer/` |
| tlsCert | N | tls certificate path |
| tlsCertKey | N | tls certificate key path |
| tlsCa | N | tls ca path |

## How to start etcd
If you want to run the etcd demo, you need to start a etcd server.

Steps：

download etcd from `https://github.com/etcd-io/etcd/releases` （You can also use docker.）

start：

````shell
./etcd
````

default listen address `localhost:2379`

## Run layotto

````shell
cd ${project_path}/cmd/layotto
go build
````

>If build reports an error, it can be executed in the root directory of the project `go mod vendor`

Execute after the compilation is successful:

````shell
./layotto start -c ../../configs/runtime_config.json
````

## Run Demo

````shell
cd ${project_path}/demo/sequencer/etcd/
 go build -o client
 ./client -s "sequencer_demo"
````

