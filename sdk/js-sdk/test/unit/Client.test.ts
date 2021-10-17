import { strict as assert } from 'assert';
import { Client } from '../../src';

describe('Client.test.ts', () => {
  let client: Client;
  beforeAll(async () => {
    client = new Client();
    const hello = await client.hello.sayHello();
    assert.equal(hello, 'greeting, ');
  });

  it('should create a Client with default port', () => {
    assert.equal(client.port, '34904');
    assert(client.runtime);
    assert(client.state);
  });
});
