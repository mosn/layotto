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
import { ServerCredentials, Server as GRPCServer } from '@grpc/grpc-js';
import { AppCallbackService } from '../../proto/appcallback_grpc_pb';
import GRPCServerImpl from './GRPCServerImpl';
import PubSub from './PubSub';
import { sleep } from '../utils';

const debug = debuglog('layotto:server:main');

export default class Server {
  readonly port: string;
  readonly pubsub: PubSub;
  private readonly _serverImpl: GRPCServerImpl;
  private readonly _server: GRPCServer;

  constructor(port: string = process.env.appcallback_GRPC_PORT ?? '9999') {
    this.port = port;
    this._serverImpl = new GRPCServerImpl();
    this.pubsub = new PubSub(this._serverImpl);

    this._server = new GRPCServer();
    // @ts-ignore
    this._server.addService(AppCallbackService, this._serverImpl);
    debug('AppCallbackService start and listen on port:%s', this.port);
  }

  async start(): Promise<void> {
    await this._bind();
    this._server.start();
    await sleep(250);
  }

  async close(): Promise<void> {
    return new Promise((resolve, reject) => {
      this._server.tryShutdown((err) => {
        if (err) return reject(err);
        debug('Closed Server');
        resolve();
      });
    });
  }

  private async _bind(): Promise<void> {
    debug('Starting to listen on 127.0.0.1:%s', this.port);
    return new Promise((resolve, reject) => {
      const serverCredentials = ServerCredentials.createInsecure();
      this._server.bindAsync(`127.0.0.1:${this.port}`, serverCredentials, (err, port) => {
        if (err) return reject(err);
        debug('Listening on 127.0.0.1:%s', port);
        resolve();
      });
    })
  }
}
