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
 *
 */
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
