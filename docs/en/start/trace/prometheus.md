# Prometheus metrics 

## Run prometheus

```shell
cd diagnostics/prometheus

docker-compose -f prometheus-docker-compose.yaml up -d
```

## Run layotto

A layotto server can be started as follows.

```
./layotto start -c ../../configs/runtime_config.json
```

## Run Demo

The corresponding call-side code is in [client.go](https://github.com/mosn/layotto/blob/main/demo/flowcontrol/client.go), and running it calls layotto's SayHello interface.

```
 cd ${projectpath}/demo/flowcontrol/
 go build -o client
 ./client
```
Access http://127.0.0.1:9090

![](../../../img/trace/prometheus.png)


## Clearance resources

````shell
cd diagnostics/prometheus

docker-compose -f prometheus-docker-compose.yaml down
````