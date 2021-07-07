# Dubbo JSON RPC Example

## Quick Start
1. Edit config file，add `dubbo_json_rpc` filter

![jsonrpc.jpg](../../../img/rpc/jsonrpc.jpg)

2. Compile and start layotto
```sh
go build -o layotto cmd/layotto/main.go
./layotto -c demo/rpc/dubbo_json_rpc/example.json
```

3. Start dubbo server

use `dubbo-go-samples` repo's example server.

```sh
git clone git@github.com:apache/dubbo-go-samples.git
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
sh ./bin/load.sh start
```

4. call runtime InvokerService api.
```sh
go run demo/rpc/dubbo_json_rpc/dubbo_json_client/client.go -d '{"jsonrpc":"2.0","method":"GetUser","params":["A003"],"id":9527}'
```

![jsonrpc.jpg](../../../img/rpc/jsonrpcresult.jpg)

## Next Step

If you are interested in the implementation principle, or want to extend some functions, you can read [RPC design document](en/design/rpc/rpc-design-doc.md)
