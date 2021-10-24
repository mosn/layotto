import { strict as assert } from 'assert';
import { Client } from '../../src';
// import { SequencerOptionsAutoIncrement } from '../../src/types/Sequencer';

describe('Sequencer.test.ts', () => {
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
      currentId = await client.sequencer.getNextId(storeName, 'user_info');
      assert(BigInt(currentId) > BigInt(lastId));
      ids.push(currentId);
    }
    assert.equal(ids.length, 20);
    console.log('ids: %j', ids);
  });
});
