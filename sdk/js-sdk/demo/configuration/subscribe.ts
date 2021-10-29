/*
 * Copyright 2021 Layotto Authors
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
import { strict as assert } from 'assert';
import { Client, utils } from 'layotto';

const client = new Client();
assert(client);

async function main() {
  const hello = await client.hello.sayHello({
    name: 'js-sdk',
  });
  console.log('%s', hello);

  await client.configuration.save({
    storeName: 'etcd',
    appId: 'js-sdk-demo',
    items: [
      {
        key: 'hello',
        content: 'world',
      },
    ],
  });

  const call = client.configuration.subscribe({
    storeName: 'etcd',
    appId: 'js-sdk-demo',
    keys: ['hello'],
    onData(items) {
      console.log('get items: %j', items);
    },
    onClose(err) {
      console.error('close with error: %s', err);
    }
  });
  await utils.sleep(100);

  await client.configuration.save({
    storeName: 'etcd',
    appId: 'js-sdk-demo',
    items: [
      {
        key: 'hello',
        content: 'world first',
      },
    ],
  });

  await utils.sleep(500);
  await client.configuration.save({
    storeName: 'etcd',
    appId: 'js-sdk-demo',
    items: [
      {
        key: 'hello',
        content: 'world second after 500ms',
      },
    ],
  });

  await utils.sleep(500);
  await client.configuration.save({
    storeName: 'etcd',
    appId: 'js-sdk-demo',
    items: [
      {
        key: 'hello',
        content: 'world third after 1000ms',
      },
    ],
  });

  await utils.sleep(500);
  call.end();
  call.destroy();
}

main();
