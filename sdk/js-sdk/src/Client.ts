import { ChannelCredentials } from '@grpc/grpc-js';
import { RuntimeClient } from '../proto/runtime_grpc_pb';
import State from './State';
import Hello from './Hello';
import Invoker from './Invoker';
import Lock from './Lock';

export default class Client {
  readonly port: string;
  readonly runtime: RuntimeClient;
  readonly hello: Hello;
  // readonly pubsub: IClientPubSub;
  readonly state: State;
  readonly invoker: Invoker;
  readonly lock: Lock;
  // readonly secret: IClientSecret;
  // readonly actor: IClientActor;

  constructor(port: string = process.env.runtime_GRPC_PORT || '34904') {
    this.port = port;
    const clientCredentials = ChannelCredentials.createInsecure();

    console.log(`[Layotto-JS] Start connection to port:${this.port}`);
    this.runtime = new RuntimeClient(`127.0.0.1:${this.port}`, clientCredentials);
    this.hello = new Hello(this.runtime);
    this.state = new State(this.runtime);
    this.invoker = new Invoker(this.runtime);
    this.lock = new Lock(this.runtime);
  }
}
