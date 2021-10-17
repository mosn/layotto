import { strict as assert } from 'assert';
import { Client } from '../../src';

describe('State.test.ts', () => {
  let client: Client;
  const storeName = 'redis';

  beforeAll(async () => {
    client = new Client();
  });

  it('should save one item and get key success', async () => {
    const key = 'js-sdk-unit-' + Date.now();
    const value = `hello js-sdk, with 中文, 😄, at ${Date()}`;
    await client.state.save(storeName, [
      { key, value },
    ]);
    const resValue = await client.state.get(storeName, key);
    assert.equal(Buffer.from(resValue).toString(), value);

    await client.state.save(storeName, { key, value });
    const resValue2 = await client.state.get(storeName, key);
    assert.equal(Buffer.from(resValue2).toString(), value);
  });

  it('should save many items and get them success', async () => {
    const items: { key: string, value: string }[] = [];
    for (let i = 0; i < 20; i++) {
      const key = `key${i}:js-sdk-unit-${Date.now()}`;
      const value = `key${i}:hello js-sdk, with 中文, 😄, at ${Date()}`;
      items.push({ key, value });
    }
    await client.state.save(storeName, items);

    for (const { key, value } of items) {
      const resValue = await client.state.get(storeName, key);
      assert.equal(Buffer.from(resValue).toString(), value);
    }

    const keys = items.map(i => i.key);
    const states = await client.state.getBulk(storeName, keys);
    assert.equal(states.length, items.length);
    for (let i = 0; i < states.length; i++) {
      const state = states[i];
      const item = items.find(i => i.key === state.key);
      assert(item);
      assert.equal(state.key, item.key);
      assert.equal(Buffer.from(state.value).toString(), item.value);
    }
  });
});
