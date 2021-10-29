import { RequestWithMeta } from './common';

export type SayHelloRequest = {
  serviceName?: string;
  name?: string;
} & RequestWithMeta;
