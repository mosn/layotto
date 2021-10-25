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
import { ChannelCredentials } from '@grpc/grpc-js';
import { RuntimeClient } from '../proto/runtime_grpc_pb';
import State from './State';
import Hello from './Hello';
import Invoker from './Invoker';
import Lock from './Lock';
import Sequencer from './Sequencer';

export default class Client {
  readonly port: string;
  readonly runtime: RuntimeClient;
  readonly hello: Hello;
  // readonly pubsub: IClientPubSub;
  readonly state: State;
  readonly invoker: Invoker;
  readonly lock: Lock;
  readonly sequencer: Sequencer;
  // readonly secret: IClientSecret;
  // readonly actor: IClientActor;

  constructor(port: string = process.env.runtime_GRPC_PORT || '34904') {
    // debug with env runtime_GRPC_HOST
    const host = process.env.runtime_GRPC_HOST || '127.0.0.1';
    this.port = port;
    const clientCredentials = ChannelCredentials.createInsecure();

    console.log(`[Layotto-JS] Start connection to ${host}:${this.port}`);
    this.runtime = new RuntimeClient(`${host}:${this.port}`, clientCredentials);
    this.hello = new Hello(this.runtime);
    this.state = new State(this.runtime);
    this.invoker = new Invoker(this.runtime);
    this.lock = new Lock(this.runtime);
    this.sequencer = new Sequencer(this.runtime);
  }
}
