import {
  GetNextIdRequest,
  GetNextIdResponse,
  SequencerOptions as SequencerOptionsPB,
} from '../proto/runtime_pb';
import { API, RequestMetadata } from './API';
import { SequencerOptions } from './types/Sequencer';

export default class Sequencer extends API {
  // Get next unique id with some auto-increment guarantee
  async getNextId(storeName: string, key: string, options?: SequencerOptions, meta?: RequestMetadata): Promise<string> {
    const req = new GetNextIdRequest();
    req.setStoreName(storeName);
    req.setKey(key);
    if (options) {
      const sequencerOptions = new SequencerOptionsPB();
      sequencerOptions.setIncrement(options.increment);
      req.setOptions(sequencerOptions);
    }

    return new Promise((resolve, reject) => {
      this.runtime.getNextId(req, this.createMetadata(meta), (err, res: GetNextIdResponse) => {
        if (err) return reject(err);
        resolve(res.getNextId());
      });
    });
  }
}
