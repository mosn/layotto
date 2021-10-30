import { RequestWithMeta, KVString } from './common';

export type PubSubCallback = (data: any) => Promise<any | void>;

export type PublishEventRequest = {
  pubsubName: string;
  topic: string;
  data: object;
  metadata?: KVString;
} & RequestWithMeta;
