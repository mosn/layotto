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
  SayHelloRequest as SayHelloRequestPB,
  SayHelloResponse as SayHelloResponsePB,
} from '../../proto/runtime_pb';
import { API } from './API';
import { SayHelloRequest } from '../types/Hello';

export default class Hello extends API {
  async sayHello(request?: SayHelloRequest): Promise<string> {
    const req = new SayHelloRequestPB();
    if (!request) request = {};
    req.setServiceName(request.serviceName ?? 'helloworld');
    if (request.name) req.setName(request.name);
    
    return new Promise((resolve, reject) => {
      this.runtime.sayHello(req, this.createMetadata(request as SayHelloRequest), (err, res: SayHelloResponsePB) => {
        if (err) return reject(err);
        resolve(res.getHello());
      });
    });
  }
}
