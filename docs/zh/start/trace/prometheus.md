# Prometheus metrics 接入

## 运行prometheus

```shell
cd diagnostics/prometheus

docker-compose -f prometheus-docker-compose.yaml up -d
```

## 运行layotto

可以按照如下方式启动一个layotto的server：

```
./layotto start -c ../../configs/runtime_config.json
```

## 运行 Demo

对应的调用端代码在[client.go](https://github.com/mosn/layotto/blob/main/demo/flowcontrol/client.go) 中，运行它会调用layotto的SayHello接口：
```
 cd ${projectpath}/demo/flowcontrol/
 go build -o client
 ./client
```
访问 http://127.0.0.1:9090

![](https://gw.alipayobjects.com/mdn/rms_5891a1/afts/img/A*mEVNSZMvtvEAAAAAAAAAAAAAARQnAQ)


## 清理资源

````shell
cd diagnostics/prometheus

docker-compose -f prometheus-docker-compose.yaml down
````