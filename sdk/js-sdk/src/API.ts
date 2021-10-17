import { RuntimeClient } from '../proto/runtime_grpc_pb';

export default class API {
  readonly runtime: RuntimeClient;
  constructor(runtime: RuntimeClient) {
    this.runtime = runtime;
  }
}
