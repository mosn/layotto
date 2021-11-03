import { Readable } from 'stream';
import { RequireExactlyOne } from 'type-fest';
import { KV, RequestWithMeta } from './common';

export type GetFileRequest = RequestWithMeta<{
  storeName: string;
  name: string;
  metadata?: KV<string>;
}>;

export type PutFileRequest = RequestWithMeta<RequireExactlyOne<{
  storeName: string;
  name: string;
  stream?: Readable,
  data?: Uint8Array,
  metadata?: KV<string>;
}, 'stream' | 'data'>>;

export type ListFileResponse = {
  names: string[];
};
