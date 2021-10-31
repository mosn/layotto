# Layotto Node.js SDK

The Layotto Node.js SDK to build your application.

## Usage

### State

`demo/state.ts`

```ts
import { Client } from 'layotto';

const storeName = 'redis';
const key = 'foo-js-sdk';
const value = `bar, from js-sdk, ${Date()}`;

await client.state.save({
  storeName, 
  states: { key, value },
});
console.log('saveState success, key: %j, value: %j', key, value);

const resValue = await client.state.get({ storeName, key });
console.log('getState success, key: %j, value: %j, toString: %j',
  key, resValue, Buffer.from(resValue).toString('utf8'));
```

## Development

### Install dependencies

```bash
npm install
```

### Generate gRPC files

```bash
npm run build:grpc
```

### Run Tests
### step 1. Set up the environment

- Running redis under Docker

```bash
docker pull redis:latest
docker run -itd --name redis-test -p 6380:6379 redis
```

- Running etcd under Docker

```bash
docker pull quay.io/coreos/etcd
docker run -itd -p 2379:2379 --name etcd quay.io/coreos/etcd /usr/local/bin/etcd -advertise-client-urls http://0.0.0.0:2379 -listen-client-urls http://0.0.0.0:2379
```

- Start a echoserver for testing the rpc api
```shell
go run demo/rpc/http/echoserver/echoserver.go
```

If you want to know more about this,check https://mosn.io/layotto/#/zh/start/rpc/helloworld

- Start Layotto, see [How to run layotto](https://mosn.io/layotto/#/zh/start/state/start?id=%e7%ac%ac%e4%ba%8c%e6%ad%a5%ef%bc%9a%e8%bf%90%e8%a1%8clayotto)

```bash
cd ${projectpath}/cmd/layotto
go build

./layotto start -c ../../configs/config_integration_redis_etcd.json
```

### step 2: Run the tests
- Then, run the test script by npm

```bash
npm run test:unit
```

Enable trace debug log for grpc-js:

```bash
GRPC_TRACE=compression GRPC_VERBOSITY=debug GRPC_TRACE=all npm run test test/unit/client/Invoker.test.ts
```

## Reference

- [Core concepts, architecture and lifecycle](https://grpc.io/docs/what-is-grpc/core-concepts/)
- [Get Start with gRPC on Node.js](https://grpc.io/docs/languages/node/quickstart/)
- [Node.js gRPC Library](https://grpc.github.io/grpc/node/)
- [Understanding Streams in Node.js](https://nodesource.com/blog/understanding-streams-in-nodejs/)
- [Go Proxy in China](https://learnku.com/go/wikis/38122)
- [How to build a bi-directional streaming gRPC service with Node.js and Java](https://medium.com/@Mark.io/bi-directional-streaming-grpc-with-node-js-and-java-7cbe0f1e0693)
