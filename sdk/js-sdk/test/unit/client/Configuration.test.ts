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
import { Client, utils } from '../../../src';

describe('client/Configuration.test.ts', () => {
  let client: Client;
  const storeName = 'etcd';
  const appId = 'js-sdk-unittest';
  const key1 = `hello1-${Date.now()}`;
  const key2 = `hello2-${Date.now()}`;

  beforeAll(() => {
    client = new Client();
  });

  it('should save/get/delete configuration work', async () => {
    let configs = await client.configuration.get({
      storeName,
      appId,
      keys: [key1, key2],
    });
    assert.equal(configs.length, 0);
    // save success

    await client.configuration.save({
      storeName,
      appId,
      items: [
        {
          key: key1,
          content: 'world1 哈哈，😄',
        },
        {
          key: key2,
          content: 'world2 哈哈，😄',
        },
      ],
    });

    configs = await client.configuration.get({
      storeName,
      appId,
      keys: [key1, key2],
    });
    assert.equal(configs.length, 2);
    assert.equal(configs[0].key, key1);
    assert.equal(configs[0].content, 'world1 哈哈，😄');
    assert.equal(configs[1].key, key2);
    assert.equal(configs[1].content, 'world2 哈哈，😄');

    // delete them
    await client.configuration.delete({
      storeName,
      appId,
      keys: [key1, key2],
    });
    configs = await client.configuration.get({
      storeName,
      appId,
      keys: [key1, key2],
    });
    assert.equal(configs.length, 0);
  });

  it('should subscribe work', async () => {
    await client.configuration.save({
      storeName,
      appId,
      items: [
        {
          key: key1,
          content: 'world1 哈哈，😄',
        },
        {
          key: key2,
          content: 'world2 哈哈，😄',
        },
      ],
    });
    const configs = await client.configuration.get({
      storeName,
      appId,
      keys: [key1, key2],
    });
    assert.equal(configs.length, 2);

    let lastConfig = {};
    const call = client.configuration.subscribe({
      storeName,
      appId,
      keys: [key1, key2],
      onData(items) {
        // console.log('get items', items);
        for (const item of items) {
          lastConfig[item.key] = item.content;
        }
      },
      onClose(err) {
        assert(!err);
        // console.log('close, error: %s', err);
      },
    });
    // console.log('call send', call.destroyed);
    await utils.sleep(500);
    await client.configuration.save({
      storeName,
      appId,
      items: [
        {
          key: key1,
          content: 'world1111 哈哈，😄',
        },
        {
          key: key2,
          content: 'world2222 哈哈，😄',
        },
      ],
    });
    await utils.sleep(500);
    assert.equal(lastConfig[key1], 'world1111 哈哈，😄');
    assert.equal(lastConfig[key2], 'world2222 哈哈，😄');

    await client.configuration.save({
      storeName,
      appId,
      items: [
        {
          key: key2,
          content: 'world2222-update2 哈哈，😄',
        },
      ],
    });

    await utils.sleep(500);
    assert.equal(lastConfig[key1], 'world1111 哈哈，😄');
    assert.equal(lastConfig[key2], 'world2222-update2 哈哈，😄');

    call.end();
    call.destroy();
    await utils.sleep(500);
  });
});
