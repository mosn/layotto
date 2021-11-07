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
import { tmpdir } from 'os';
import { existsSync, createWriteStream, createReadStream } from 'fs';
import { mkdtemp, rm } from 'fs/promises';
import { pipeline } from 'stream/promises';
import { join } from 'path';
import { Client, utils } from '../../../src';

describe.skip('client/File.test.ts', () => {
  let client: Client;
  let tmpfileDir: string;
  const storeName = 'aliOSS';
  const bucket = 'layotto-js-sdk-local-test';

  beforeAll(async () => {
    client = new Client();
    tmpfileDir = await mkdtemp(join(tmpdir(), 'layotto-js-sdk-file-unittest-tmp-'));
  });

  afterAll(async () => {
    // rm -rf
    await rm(tmpfileDir, { recursive: true, force: true });
  });

  describe('get()', () => {
    it('should get a file stream success', async () => {
      const filepath = join(tmpfileDir, 'getFileStream.jpg');
      console.log(filepath);
      assert(!existsSync(filepath));

      const stream = await client.file.get({
        storeName: 'aliOSS',
        name: 'foo/tab3.jpg',
      });
      assert(stream);
      await pipeline(
        stream,
        createWriteStream(filepath),
      );
      assert(existsSync(filepath));
    }, 30000);

    it('should get a file stream from minio', async () => {
      const filepath = join(tmpfileDir, 'minio.jpg');
      console.log(filepath);
      assert(!existsSync(filepath));

      const stream = await client.file.get({
        storeName: 'minioOSS',
        name: '1.jpg',
        metadata: {
          bucket,
          endpoint: 'http://127.0.0.1:9000',
        },
      });
      assert(stream);
      await pipeline(
        stream,
        createWriteStream(filepath),
      );
      assert(existsSync(filepath));
    });

    it('should throw error when storeName not exists', async () => {
      await assert.rejects(
        async () => {
          const filepath = join(tmpfileDir, 'storeName-not-exists.jpg');
          const stream = await client.file.get({
            storeName: 'notexists',
            name: 'foo.jpg',
          });
          await pipeline(
            stream,
            createWriteStream(filepath),
          );
        },
        (err: any) => {
          assert.equal(err.code, 3);
          assert.equal(err.details, 'not supported store type: notexists');
          return true;
        }
      )
    });
  });

  describe('put()', () => {
    it('should put a file stream to server success', async () => {
      const stream = createReadStream(__filename);
      await client.file.put({
        storeName: 'aliOSS',
        name: 'layotto_js_sdk_unittest.test.ts',
        stream,
      });
    }, 10000);

    it('should put bytes data to server success', async () => {
      await client.file.put({
        storeName: 'aliOSS',
        name: 'layotto_js_sdk_unittest.txt',
        data: Buffer.from('哈哈😄，text still work', 'utf8'),
      });
    }, 10000);

    it.skip('should work with minio', async () => {
      await client.file.put({
        storeName: 'minio',
        name: 'layotto_js_sdk_unittest.txt',
        data: Buffer.from('哈哈😄，text still work', 'utf8'),
      });
    });
  });

  describe('list()', () => {
    it('should list dir success', async () => {
      const { names } = await client.file.list({
        storeName,
        name: 'layotto-js-sdk-local-test',
        metadata: {
          Prefix: 'foo/',
        },
      });
      console.log(names);
      assert(names.length > 0);
    });

    it.skip('should work with minio', async () => {
      const { names } = await client.file.list({
        storeName: 'minio',
        name: 'layotto-js-sdk-local-test',
      });
      console.log(names);
      assert(names.length > 0);
    });
  });

  describe('delete()', () => {
    it('should put bytes data to server and delete it success', async () => {
      const name = `layotto_js_sdk_unittest_need_to_delete_${Date.now()}.txt`;
      await client.file.put({
        storeName,
        name,
        data: Buffer.from('哈哈😄，text still work, delete it now!', 'utf8'),
      });

      await utils.sleep(500);
      const filepath = join(tmpfileDir, name);
      const stream = await client.file.get({
        storeName,
        name,
      });
      assert(stream);
      await pipeline(
        stream,
        createWriteStream(filepath),
      );

      // delete it
      await client.file.delete({
        storeName,
        name,
      });

      // throw error
      await assert.rejects(
        async () => {
          const filepath = join(tmpfileDir, name);
          const stream = await client.file.get({
            storeName,
            name,
          });
          await pipeline(
            stream,
            createWriteStream(filepath),
          );
        },
        (err: any) => {
          // console.error(err);
          assert.equal(err.code, 13);
          assert.match(err.message, /StatusCode=404, ErrorCode=NoSuchKey/);
          return true;
        }
      );
    }, 10000);
  });
});

