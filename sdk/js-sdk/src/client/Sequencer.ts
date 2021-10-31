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
  GetNextIdRequest as GetNextIdRequestPB,
  GetNextIdResponse as GetNextIdResponsePB,
  SequencerOptions as SequencerOptionsPB,
} from '../../proto/runtime_pb';
import { API } from './API';
import { GetNextIdRequest } from '../types/Sequencer';

export default class Sequencer extends API {
  // Get next unique id with some auto-increment guarantee
  async getNextId(request: GetNextIdRequest): Promise<string> {
    const req = new GetNextIdRequestPB();
    req.setStoreName(request.storeName);
    req.setKey(request.key);
    if (request.options) {
      const sequencerOptions = new SequencerOptionsPB();
      sequencerOptions.setIncrement(request.options.increment);
      req.setOptions(sequencerOptions);
    }
    this.mergeMetadataToMap(req.getMetadataMap(), request.metadata);

    return new Promise((resolve, reject) => {
      this.runtime.getNextId(req, this.createMetadata(request), (err, res: GetNextIdResponsePB) => {
        if (err) return reject(err);
        resolve(res.getNextId());
      });
    });
  }
}
