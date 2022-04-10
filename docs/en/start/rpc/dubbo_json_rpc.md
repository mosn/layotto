# Dubbo JSON RPC Example

## Quick Start
### step 1. Edit config file，add `dubbo_json_rpc` filter

![jsonrpc.jpg](../../../img/rpc/jsonrpc.jpg)

### step 2. Compile and start layotto
```shell
go build -o layotto cmd/layotto/main.go
```

```shell @background
./layotto -c demo/rpc/dubbo_json_rpc/example.json
```

### step 3. Start dubbo server

use `dubbo-go-samples` repo's example server.

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

```shell @background
sh ./bin/load.sh start
```

### step 4. call runtime InvokerService api.
```shell
go run demo/rpc/dubbo_json_rpc/dubbo_json_client/client.go -d '{"jsonrpc":"2.0","method":"GetUser","params":["A003"],"id":9527}'
```

![jsonrpc.jpg](../../../img/rpc/jsonrpcresult.jpg)

## Next Step

If you are interested in the implementation principle, or want to extend some functions, you can read [RPC design document](en/design/rpc/rpc-design-doc.md)
