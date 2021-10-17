import { ChannelCredentials } from '@grpc/grpc-js';
import { RuntimeClient } from '../proto/runtime_grpc_pb';
import State from './State';
import Hello from './Hello';

export default class Client {
  readonly port: string;
  readonly runtime: RuntimeClient;
  readonly hello: Hello;
  // readonly pubsub: IClientPubSub;
  readonly state: State;
  // readonly binding: IClientBinding;
  // readonly invoker: IClientInvoker;
  // readonly secret: IClientSecret;
  // readonly actor: IClientActor;

  constructor(port: string = process.env.runtime_GRPC_PORT || '34904') {
    this.port = port;
    const clientCredentials = ChannelCredentials.createInsecure();

    console.log(`[Layotto-JS] Start connection to port:${this.port}`);
    this.runtime = new RuntimeClient(`127.0.0.1:${this.port}`, clientCredentials);
    this.hello = new Hello(this.runtime);
    this.state = new State(this.runtime);
    // this.pubsub = new GRPCClientPubSub(client);
    // this.binding = new GRPCClientBinding(client);
    // this.invoker = new GRPCClientInvoker(client);
    // this.secret = new GRPCClientSecret(client);
    // this.actor = new GRPCClientActor(client);
  }
}
