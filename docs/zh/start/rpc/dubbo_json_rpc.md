# Dubbo JSON RPC Example

## 快速开始
### step 1. 修改配置文件，加入`dubbo_json_rpc`插件

![jsonrpc.jpg](../../../img/rpc/jsonrpc.jpg)

### step 2. 编译运行layotto
```shell
go build -o layotto cmd/layotto/main.go
```

```shell background
./layotto -c demo/rpc/dubbo_json_rpc/example.json
```

### step 3. 启动dubbo服务端

这里使用了`dubbo-go-samples`提供的示例服务
```shell
git clone https://github.com/apache/dubbo-go-samples.git
cd dubbo-go-samples

# start zookeeper
cd attachment/go-server
make -f ../../build/Makefile docker-up 
cd -

# build dubbo server
cd general/jsonrpc/go-server
sh assembly/mac/dev.sh

# start dubbo server
cd target/darwin/{generate_folder}/
```

```shell background
sh ./bin/load.sh start
```

### step 4. 通过GPRC接口发起调用
```shell
go run demo/rpc/dubbo_json_rpc/dubbo_json_client/client.go -d '{"jsonrpc":"2.0","method":"GetUser","params":["A003"],"id":9527}'
```

![jsonrpc.jpg](../../../img/rpc/jsonrpcresult.jpg)

### 下一步

如果您对实现原理感兴趣，或者想扩展一些功能，可以阅读[RPC的设计文档](zh/design/rpc/rpc设计文档.md)
