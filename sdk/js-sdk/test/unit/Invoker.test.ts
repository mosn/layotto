import { strict as assert } from 'assert';
import { Client } from '../../src';
import { HTTPExtensionVerb } from '../../src/types/Invoker';

describe.skip('Invoker.test.ts', () => {
  let client: Client;

  beforeAll(() => {
    client = new Client();
  });

  describe('invoke()', () => {
    it('should invoke HelloService:1.0 with text success', async () => {
      const state = await client.invoker.invoke('HelloService:1.0', '/hello', undefined, 'hello runtime with js-sdk');
      assert.equal(state, null);
    });

    it('should invoke HelloService:1.0 with json success', async () => {
      const state = await client.invoker.invoke('HelloService:1.0', '/hello', HTTPExtensionVerb.POST, {
        'hello runtime': 'I am js-sdk client',
      });
      assert.equal(state, null);
    });
  });
});

