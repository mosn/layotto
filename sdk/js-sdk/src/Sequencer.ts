/*
 * Copyright 2021 Layotto Authors
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
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
