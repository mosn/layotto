import { Metadata } from '@grpc/grpc-js';
import { RuntimeClient } from '../proto/runtime_grpc_pb';

// gRPC meta data
export type RequestMetadata = {
  [key: string]: string;
};

export class API {
  readonly runtime: RuntimeClient;
  constructor(runtime: RuntimeClient) {
    this.runtime = runtime;
  }

  createMetadata(meta?: RequestMetadata) {
    const metadata = new Metadata();
    for (const key in meta) {
      metadata.add(key, meta[key]);
    }
    return metadata;
  }
}
