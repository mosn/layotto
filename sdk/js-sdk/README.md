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

await client.state.save(storeName, { key, value });
console.log('saveState success, key: %j, value: %j', key, value);

const resValue = await client.state.get(storeName, key);
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

run the unit tests in your local env:

- Start Layotto first, see [How to run layotto](https://mosn.io/layotto/#/zh/start/state/start?id=%e7%ac%ac%e4%ba%8c%e6%ad%a5%ef%bc%9a%e8%bf%90%e8%a1%8clayotto)

```bash
cd ${projectpath}/cmd/layotto
go build

./layotto start -c ../../configs/config_state_redis.json
```

- Then, run unit tests script by npm

```bash
npm run test:unit
```
