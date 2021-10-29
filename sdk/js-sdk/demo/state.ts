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
import { strict as assert } from 'assert'
import { Client } from 'layotto';

const client = new Client();
assert(client);

async function main() {
  const storeName = 'redis';
  const key = 'foo-js-sdk';
  const value = `bar, from js-sdk, ${Date()}`;

  await client.state.save({
    storeName, 
    states: [
      { key, value },
    ], 
    requestMeta: { traceid: 'mock-tracerid-123123' },
  });
  console.log('saveState success, key: %j, value: %j', key, value);

  // await client.state.save(storeName, { key, value });
  // console.log('saveState success, key: %j, value: %j', key, value);

  const state = await client.state.get({ storeName, key });
  assert(state);
  console.log('getState success, key: %j, value: %j, toString: %j',
    key, state.value, Buffer.from(state.value).toString('utf8'));
}

main();
