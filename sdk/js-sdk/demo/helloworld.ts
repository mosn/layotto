import { strict as assert } from 'assert';
import { Client } from 'layotto';

const client = new Client();
assert(client);

async function main() {
  const hello = await client.hello.sayHello('helloworld', 'js-sdk');
  console.log('%s', hello);
}

main();
