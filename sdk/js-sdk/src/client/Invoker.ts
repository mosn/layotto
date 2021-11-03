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
import { Any } from 'google-protobuf/google/protobuf/any_pb';
import {
  InvokeServiceRequest as InvokeServiceRequestPB,
  CommonInvokeRequest as CommonInvokeRequestPB,
  HTTPExtension,
  InvokeResponse as InvokeResponsePB,
} from '../../proto/runtime_pb';
import { API } from './API';
import { InvokeServiceRequest, InvokeResponse } from '../types/Invoker';

export default class Invoker extends API {
  async invoke(request: InvokeServiceRequest): Promise<InvokeResponse> {
    const message = new CommonInvokeRequestPB();
    message.setMethod(request.method);

    const httpVerb = request.httpVerb ?? HTTPExtension.Verb.GET;
    const httpExtension = new HTTPExtension();
    httpExtension.setVerb(httpVerb);
    message.setHttpExtension(httpExtension);

    if (request.data) {
      const dataSerialized = new Any();
      if (typeof request.data === 'string') {
        message.setContentType(request.contentType ?? 'text/plain; charset=UTF-8');
        dataSerialized.setValue(Buffer.from(request.data, 'utf8'));
      } else {
        message.setContentType(request.contentType ?? 'application/json');
        dataSerialized.setValue(Buffer.from(JSON.stringify(request.data), 'utf8'));
      }
      message.setData(dataSerialized);
    }

    console.log(message);

    const req = new InvokeServiceRequestPB();
    req.setId(request.id);
    req.setMessage(message);

    return new Promise((resolve, reject) => {
      this.runtime.invokeService(req, this.createMetadata(request), (err, res: InvokeResponsePB) => {
        if (err) return reject(err);
        const contentType = res.getContentType().split(';', 1)[0].toLowerCase();
        const rawData = res.getData();
        let content;
        if (contentType === 'application/json') {
          if (rawData) {
            content = JSON.parse(Buffer.from(rawData.getValue_asU8()).toString());
          } else {
            content = {}; 
          }
        } else if (contentType === 'text/plain') {
          if (rawData) {
            content = Buffer.from(rawData.getValue_asU8()).toString();
          } else {
            content = '';
          }
        } else {
          content = rawData ? rawData.getValue_asU8() : [];
        }
        const response: InvokeResponse = { contentType, content };
        resolve(response);
      });
    });
  }
}
