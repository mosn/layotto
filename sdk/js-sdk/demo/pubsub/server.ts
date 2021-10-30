import { Server } from 'layotto';

async function main() {
  const server = new Server();
  server.pubsub.subscribe('redis', 'topic1', async (data) => {
    console.log('topic1 event data: %j', data);
  });

  await server.start();
}

main();
