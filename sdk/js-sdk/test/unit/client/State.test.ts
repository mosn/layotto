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
import { Client } from '../../../src';
import { StateOperation, StateOperationType } from '../../../src/types/State';

describe('client/State.test.ts', () => {
  let client: Client;
  const storeName = 'redis';

  beforeAll(async () => {
    client = new Client();
  });

  describe('get(), getBulk()', () => {
    it('should throw error when storeName not exists', async () => {
      await assert.rejects(
        async () => {
          await client.state.get({ storeName: 'notexists-store', key: 'foo' });
        },
        (err: any) => {
          // console.error(err);
          assert.equal(err.code, 3);
          assert.equal(err.details, 'state store notexists-store is not found');
          return true;
        }
      );
    });

    it('should get not exists key', async () => {
      const key = 'js-sdk-unit-notexists-' + Date.now();
      const state = await client.state.get({ storeName, key });
      assert.equal(state, null);
    });

    it('should get empty value', async () => {
      const key = 'js-sdk-unit-empty-' + Date.now();
      await client.state.save({
        storeName, 
        states: {
          key,
          value: '',
        },
      });
      const state = await client.state.get({ storeName, key });
      assert(state);
      assert.equal(state.value.length, 0);
      assert.equal(state.key, key);
      assert.equal(state.etag, '1');
    });

    it('should save one item and get key success', async () => {
      const key = 'js-sdk-unit-' + Date.now();
      const value = `hello js-sdk, with 中文, 😄, at ${Date()}`;
      await client.state.save({
        storeName, 
        states: { key, value },
        requestMeta: { traceid: `mock-traceid-unittest-${Date.now()}` },
      });
      const state = await client.state.get({ storeName, key });
      assert(state);
      assert.equal(Buffer.from(state.value).toString(), value);
  
      await client.state.save({
        storeName, 
        states: { key, value },
      });
      const state2 = await client.state.get({ storeName, key });
      assert(state2);
      assert.equal(Buffer.from(state2.value).toString(), value);
    });
  
    it('should save bulk items and get them success', async () => {
      const items: { key: string, value: string }[] = [];
      for (let i = 0; i < 20; i++) {
        const key = `key${i}:js-sdk-unit-${Date.now()}`;
        const value = `key${i}:hello js-sdk, with 中文, 😄, at ${Date()}`;
        items.push({ key, value });
      }
      await client.state.save({
        storeName, 
        states: items,
      });
  
      for (const { key, value } of items) {
        const state = await client.state.get({ storeName, key });
        assert(state);
        assert.equal(Buffer.from(state.value).toString(), value);
      }
  
      const keys = items.map(i => i.key);
      const states = await client.state.getBulk({ storeName, keys });
      assert.equal(states.length, items.length);
      for (let i = 0; i < states.length; i++) {
        const state = states[i];
        const item = items.find(i => i.key === state.key);
        assert(item);
        assert.equal(state.key, item.key);
        assert.equal(Buffer.from(state.value).toString(), item.value);
      }

      const states2 = await client.state.getBulk({
        storeName, keys, requestMeta: { foo: 'bar' },
      });
      assert.equal(states2.length, items.length);
    });
  });

  describe('delete(), deleteBulk()', () => {
    it('should delete one key success', async () => {
      const key = 'js-sdk-unit-delete-' + Date.now();
      // delete not exists
      await client.state.delete({ storeName, key });

      const value = `hello js-sdk, with 中文, 😄, at ${Date()}`;
      await client.state.save({ storeName, states: { key, value } });

      // delete exists
      await client.state.delete({ storeName, key });

      const state = await client.state.get({ storeName, key });
      assert.equal(state, null);
    });

    it('should delete bulk keys success', async () => {
      const items: { key: string, value: string }[] = [];
      const keys: { key: string }[] = [];
      for (let i = 0; i < 20; i++) {
        const key = `key${i}:js-sdk-unit-delete-bulk-${Date.now()}`;
        const value = `key${i}:hello js-sdk, with 中文, 😄, at ${Date()}`;
        items.push({ key, value });
        keys.push({ key });
      }
      // delete not exists
      await client.state.deleteBulk({ storeName, states: keys });
      let states = await client.state.getBulk({
        storeName, 
        keys: keys.map(i => i.key),
      });
      assert.equal(states.length, 0);

      await client.state.save({ storeName, states: items });
      states = await client.state.getBulk({
        storeName, 
        keys: keys.map(i => i.key),
      });
      assert.equal(states.length, 20);

      // delete 5 keys
      await client.state.deleteBulk({ storeName, states: keys.slice(0, 5) });

      states = await client.state.getBulk({
        storeName, 
        keys: keys.map(i => i.key),
      });
      assert.equal(states.length, 15);

      // delete all exists
      await client.state.deleteBulk({ storeName, states: keys });

      states = await client.state.getBulk({
        storeName, 
        keys: keys.map(i => i.key),
      });
      assert.equal(states.length, 0);
    });
  });

  describe('executeTransaction()', () => {
    it('should upsert then delete success', async () => {
      const operations: StateOperation[] = [];
      const keys: { key: string }[] = [];
      for (let i = 0; i < 20; i++) {
        const key = `key${i}:js-sdk-unit-upsert-bulk-${Date.now()}`;
        const value = `key${i}:hello js-sdk, with 中文, 😄, at ${Date()}`;
        operations.push({
          operationType: StateOperationType.Upsert,
          request: { key, value },
        });
        keys.push({ key });
      }
      // not work now: mosn error
      // 2021-10-24 14:06:52,590 [ERROR] [mosn.proxy.panic] [grpc] [unary] grpc unary handle panic: interface conversion: interface {} is *state.SetRequest, not state.SetRequest, method: /spec.proto.runtime.v1.Runtime/ExecuteStateTransaction, stack:goroutine 8213 [running]:
      await client.state.executeTransaction({ storeName, operations });
      // const states = await client.state.getBulk(storeName, keys.map(i => i.key));
      // assert.equal(states.length, 20);
    });
  });
});
