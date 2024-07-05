# Dubbo JSON RPC Example
*Note: This demo needs to run under go v1.17*
## Quick Start
### step 1. Edit config file，add `dubbo_json_rpc` filter

![jsonrpc.jpg](../../../img/rpc/jsonrpc.jpg)

### step 2. Compile and start layotto

```shell @if.not.exist layotto
go build -o layotto ./cmd/layotto
```

```shell @background
./layotto -c demo/rpc/dubbo_json_rpc/example.json
```

### step 3. Start dubbo server

use `dubbo-go-samples` repo's example server.

```shell @if.not.exist dubbo-go-samples
git clone https://github.com/apache/dubbo-go-samples.git
```

```shell
cd dubbo-go-samples
git reset --hard f0d1e1076397a4736de080ffb16cd0963c8c2f9d

# start zookeeper
cd rpc/jsonrpc/go-server
docker-compose -f docker/docker-compose.yml up -d

# prepare to build dubbo server
cd cmd
export DUBBO_GO_CONFIG_PATH="../conf/dubbogo.yml"
```

Build dubbo server:

```shell @if.not.exist server
go build -o server .
```

Start dubbo server:

```shell @background.sleep 3s
./server
```

### step 4. call runtime InvokerService api.

```shell @cd ${project_path}
go run demo/rpc/dubbo_json_rpc/dubbo_json_client/client.go -d '{"jsonrpc":"2.0","method":"GetUser","params":["A003"],"id":9527}'
```

![jsonrpc.jpg](../../../img/rpc/jsonrpcresult.jpg)

## Next Step

If you are interested in the implementation principle, or want to extend some functions, you can read [RPC design document](en/design/rpc/rpc-design-doc.md)
