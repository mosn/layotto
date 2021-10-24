import { strict as assert } from 'assert';
import { Client } from '../../src';
import { UnlockResponseStatus } from '../../src/types/Lock';
import { sleep } from '.,/../../src/utils';

describe('Lock.test.ts', () => {
  let client: Client;
  const storeName = 'redis';

  beforeAll(() => {
    client = new Client();
  });

  it('should unLock not exists resource', async () => {
    const resourceId = 'js-sdk-lock-not-exists';
    const lockOwner = client.lock.uuid();
    const status = await client.lock.unLock(storeName, resourceId, lockOwner);
    assert.equal(status, UnlockResponseStatus.LOCK_UNEXIST);
  });

  it('should lock 2 seconds work', async () => {
    const resourceId = 'js-sdk-lock-2-seconds';
    const lockOwner = client.lock.uuid();
    const success1 = await client.lock.tryLock(storeName, resourceId, lockOwner, 2);
    assert.equal(success1, true);
    const success2 = await client.lock.tryLock(storeName, resourceId, lockOwner, 2);
    assert.equal(success2, false);
    await sleep(2500);
    // wait for lock exipre after 2.5s
    const success3 = await client.lock.tryLock(storeName, resourceId, lockOwner, 2);
    assert.equal(success3, true);

    // unlock by other owner will fail
    const lockOwner2 = client.lock.uuid();
    const status1 = await client.lock.unLock(storeName, resourceId, lockOwner2);
    assert.equal(status1, UnlockResponseStatus.LOCK_BELONG_TO_OTHERS);

    // unlock success by owner
    const status2 = await client.lock.unLock(storeName, resourceId, lockOwner);
    assert.equal(status2, UnlockResponseStatus.SUCCESS);
  });
});
