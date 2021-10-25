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
import { Any } from "google-protobuf/google/protobuf/any_pb";
import {
  InvokeServiceRequest,
  CommonInvokeRequest,
  HTTPExtension,
  InvokeResponse,
} from '../proto/runtime_pb';
import { API, RequestMetadata } from './API';

export default class Invoker extends API {
  async invoke(id: string, methodName: string, httpVerb = HTTPExtension.Verb.GET, data: string | object = {}, meta?: RequestMetadata): Promise<object> {
    const message = new CommonInvokeRequest();
    message.setMethod(methodName);

    const httpExtension = new HTTPExtension();
    httpExtension.setVerb(httpVerb);
    message.setHttpExtension(httpExtension);

    const dataSerialized = new Any();
    if (typeof data === 'string') {
      message.setContentType('');
      dataSerialized.setValue(Buffer.from(data, 'utf8'));
    } else {
      message.setContentType('application/json');
      dataSerialized.setValue(Buffer.from(JSON.stringify(data), 'utf8'));
    }
    message.setData(dataSerialized);

    const req = new InvokeServiceRequest();
    req.setId(id);
    req.setMessage(message);

    return new Promise((resolve, reject) => {
      this.runtime.invokeService(req, this.createMetadata(meta), (err, res: InvokeResponse) => {
        if (err) return reject(err);
        const data = JSON.parse(res.getData());
        resolve(data);
      });
    });
  }
}
