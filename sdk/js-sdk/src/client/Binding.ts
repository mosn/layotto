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
  InvokeBindingRequest as InvokeBindingRequestPB,
  InvokeBindingResponse as InvokeBindingResponsePB,
} from '../../proto/runtime_pb';
import { API } from './API';
import {
  InvokeBindingRequest,
  InvokeBindingResponse,
} from '../types/Binding';
import { convertMapToKVString } from '../utils';

export default class Binding extends API {
  async invoke(request: InvokeBindingRequest): Promise<InvokeBindingResponse> {
    const req = new InvokeBindingRequestPB();
    req.setName(request.name);
    req.setOperation(request.operation);
    if (typeof request.data === 'string') {
      req.setData(Buffer.from(request.data, 'utf8'));
    } else {
      req.setData(request.data);
    }
    this.mergeMetadataToMap(req.getMetadataMap(), request.metadata);
    
    return new Promise((resolve, reject) => {
      this.runtime.invokeBinding(req, this.createMetadata(request), (err, res: InvokeBindingResponsePB) => {
        if (err) return reject(err);
        resolve({
          data: res.getData_asU8(),
          metadata: convertMapToKVString(res.getMetadataMap()),
        });
      });
    });
  }
}
