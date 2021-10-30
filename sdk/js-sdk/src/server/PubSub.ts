import { debuglog } from 'node:util';
import { PubSubCallback } from '../types/PubSub';
import GRPCServerImpl from './GRPCServerImpl';

const debug = debuglog('layotto:server:pubsub');

export default class PubSub {
  readonly server: GRPCServerImpl;

  constructor(server: GRPCServerImpl) {
    this.server = server;
  }
  
  async subscribe(pubsubName: string, topic: string, cb: PubSubCallback): Promise<void> {
    debug('Registering onTopicEvent Handler: PubSub = %s, Topic = %s', pubsubName, topic);
    this.server.registerPubSubSubscriptionHandler(pubsubName, topic, cb);
  }
}
