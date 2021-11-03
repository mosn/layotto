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
import { debuglog } from 'node:util';
import { ChannelCredentials } from '@grpc/grpc-js';
import { RuntimeClient } from '../../proto/runtime_grpc_pb';
import State from './State';
import Hello from './Hello';
import Invoker from './Invoker';
import Lock from './Lock';
import Sequencer from './Sequencer';
import Configuration from './Configuration';
import PubSub from './PubSub';
import File from './File';
import Binding from './Binding';

const debug = debuglog('layotto:client:main');

export default class Client {
  readonly host: string;
  readonly port: string;
  private _runtime: RuntimeClient;
  private _hello: Hello;
  private _state: State;
  private _invoker: Invoker;
  private _lock: Lock;
  private _sequencer: Sequencer;
  private _configuration: Configuration;
  private _pubsub: PubSub;
  private _file: File;
  private _binding: Binding;

  constructor(port: string = process.env.runtime_GRPC_PORT ?? '34904',
              host: string = process.env.runtime_GRPC_HOST ?? '127.0.0.1') {
    this.host = host;
    this.port = port;
    const clientCredentials = ChannelCredentials.createInsecure();
    this._runtime = new RuntimeClient(`${this.host}:${this.port}`, clientCredentials);
    debug('Start connection to %s:%s', this.host, this.port);    
  }

  get hello() {
    if (!this._hello) this._hello = new Hello(this._runtime);
    return this._hello;
  }

  get state() {
    if (!this._state) this._state = new State(this._runtime);
    return this._state;
  }

  get invoker() {
    if (!this._invoker) this._invoker = new Invoker(this._runtime);
    return this._invoker;
  }

  get lock() {
    if (!this._lock) this._lock = new Lock(this._runtime);
    return this._lock;
  }

  get sequencer() {
    if (!this._sequencer) this._sequencer = new Sequencer(this._runtime);
    return this._sequencer;
  }

  get configuration() {
    if (!this._configuration) this._configuration = new Configuration(this._runtime);
    return this._configuration;
  }

  get pubsub() {
    if (!this._pubsub) this._pubsub = new PubSub(this._runtime);
    return this._pubsub;
  }

  get file() {
    if (!this._file) this._file = new File(this._runtime);
    return this._file;
  }

  get binding() {
    if (!this._binding) this._binding = new Binding(this._runtime);
    return this._binding;
  }
}
