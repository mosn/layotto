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
  SayHelloRequest
} from '../proto/runtime_pb';
import { API, RequestMetadata } from './API';

export default class Hello extends API {
  async sayHello(serviceName = 'helloworld', name = '', meta?: RequestMetadata): Promise<string> {
    const req = new SayHelloRequest();
    req.setServiceName(serviceName);
    if (name) {
      req.setName(name);
    }
    
    return new Promise((resolve, reject) => {
      this.runtime.sayHello(req, this.createMetadata(meta), (err, res) => {
        if (err) return reject(err);
        resolve(res.getHello());
      });
    });
  }
}
