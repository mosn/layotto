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
import { Client, RumtimeTypes } from '../../../src';

describe('client/Sequencer.test.ts', () => {
  let client: Client;
  const storeName = 'etcd';

  beforeAll(() => {
    client = new Client();
  });

  it('should get next id success', async () => {
    let lastId = '0';
    let currentId = '0';
    const ids: string[] = [];
    for (let i = 0; i < 20; i++) {
      lastId = currentId;
      currentId = await client.sequencer.getNextId({
        storeName, 
        key: 'user_info',
      });
      assert(BigInt(currentId) > BigInt(lastId));
      ids.push(currentId);
    }
    assert.equal(ids.length, 20);
    // console.log('ids: %j', ids);
  });

  it('should get next id with options increment:WEAK success', async () => {
    let lastId = '0';
    let currentId = '0';
    const ids: string[] = [];
    for (let i = 0; i < 20; i++) {
      lastId = currentId;
      currentId = await client.sequencer.getNextId({
        storeName, 
        key: 'user_info',
        options: {
          increment: RumtimeTypes.SequencerOptions.AutoIncrement.WEAK,
        },
      });
      assert(BigInt(currentId) > BigInt(lastId));
      ids.push(currentId);
    }
    assert.equal(ids.length, 20);
    // console.log('ids: %j', ids);
  });
});
