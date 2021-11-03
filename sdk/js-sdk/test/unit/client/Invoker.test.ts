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
import { strict as assert } from 'assert';
import { Client } from '../../../src';
import { RumtimeTypes } from '../../../src';

describe.skip('client/Invoker.test.ts', () => {
  let client: Client;

  beforeAll(() => {
    client = new Client();
  });

  describe('invoke()', () => {
    it('should invoke HelloService:1.0 with text success', async () => {
      const state = await client.invoker.invoke({
        id: 'HelloService:1.0', 
        method: '/hello',
        data: 'hello runtime with js-sdk',
      });
      assert.equal(state, null);
    });

    it('should invoke HelloService:1.0 with json success', async () => {
      const state = await client.invoker.invoke({
        id: 'HelloService:1.0', 
        method: '/hello', 
        httpVerb: RumtimeTypes.HTTPExtension.Verb.POST, 
        data: {
          'hello runtime': 'I am js-sdk client',
        },
      });
      assert.equal(state, null);
    });

    it.skip('should invoke rpc success', async () => {
      const res = await client.invoker.invoke({
        id: 'com.alipay.rpc.common.service.facade.SampleService:1.0', 
        method: 'echoStr',
        // contentType: 'json',
        data: {
          signatures: ['java.lang.String'],
          arguments: ['abc'],
        },
        requestMeta: {
          content_type: 'json',
          upstream_content_type: 'hessian',
        },
      });
      console.log(res);
    });
  });
});

