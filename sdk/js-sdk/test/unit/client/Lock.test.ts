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
import { Client, RumtimeTypes, utils } from '../../../src';

describe('Lock.test.ts', () => {
  let client: Client;
  const storeName = 'redis';

  beforeAll(() => {
    client = new Client();
  });

  it('should unLock not exists resource', async () => {
    const resourceId = 'js-sdk-lock-not-exists';
    const lockOwner = client.lock.uuid();
    const result = await client.lock.unLock({ storeName, resourceId, lockOwner });
    assert.equal(result.status, RumtimeTypes.UnlockResponse.Status.LOCK_UNEXIST);
  });

  it('should lock 2 seconds work', async () => {
    const resourceId = 'js-sdk-lock-2-seconds';
    const lockOwner = client.lock.uuid();
    const success1 = await client.lock.tryLock({ storeName, resourceId, lockOwner, expire: 2 });
    assert.equal(success1, true);
    const success2 = await client.lock.tryLock({ storeName, resourceId, lockOwner, expire: 2 });
    assert.equal(success2, false);
    await utils.sleep(2500);
    // wait for lock exipre after 2.5s
    const success3 = await client.lock.tryLock({ storeName, resourceId, lockOwner, expire: 2 });
    assert.equal(success3, true);

    // unlock by other owner will fail
    const lockOwner2 = client.lock.uuid();
    const status1 = await client.lock.unLock({ storeName, resourceId, lockOwner: lockOwner2 });
    assert.equal(status1.status, RumtimeTypes.UnlockResponse.Status.LOCK_BELONG_TO_OTHERS);

    // unlock success by owner
    const status2 = await client.lock.unLock({ storeName, resourceId, lockOwner });
    assert.equal(status2.status, RumtimeTypes.UnlockResponse.Status.SUCCESS);
  });
});
