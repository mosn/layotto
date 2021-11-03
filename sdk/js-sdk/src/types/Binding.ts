import { KV, RequestWithMeta } from './common';

export type InvokeBindingRequest = RequestWithMeta<{
  name: string,
  data: Uint8Array | string,
  metadata: KV<string>
  operation: string,
}>;

export type InvokeBindingResponse = {
  data: Uint8Array,
  metadata: KV<string>
};
