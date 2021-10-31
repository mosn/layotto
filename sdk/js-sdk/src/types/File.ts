import { Readable } from 'stream';
import { KVString, RequestWithMeta, RequireOnlyOne } from './common';

export type GetFileRequest = {
  storeName: string;
  name: string;
  metadata?: KVString;
} & RequestWithMeta;

export type PutFileRequest = RequireOnlyOne<{
  storeName: string;
  name: string;
  stream?: Readable,
  data?: Uint8Array,
  metadata?: KVString;
}, 'stream' | 'data'> & RequestWithMeta;

export type ListFileResponse = {
  names: string[];
};
