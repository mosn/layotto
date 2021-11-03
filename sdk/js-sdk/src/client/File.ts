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
import { debuglog } from 'node:util';
import { Transform, Readable } from 'stream';
import { pipeline as pipelinePromise } from 'stream/promises';
import { 
  GetFileRequest as GetFileRequestPB,
  GetFileResponse as GetFileResponsePB,
  PutFileRequest as PutFileRequestPB,
  ListFileRequest as ListFileRequestPB,
  DelFileRequest as DelFileRequestPB,
} from '../../proto/runtime_pb';
import { API } from './API';
import { GetFileRequest, ListFileResponse, PutFileRequest } from '../types/File';

const debug = debuglog('layotto:client:file');

export default class File extends API {
  // Get a file stream
  async get(request: GetFileRequest): Promise<Readable> {
    const req = new GetFileRequestPB();
    req.setStoreName(request.storeName);
    req.setName(request.name);
    this.mergeMetadataToMap(req.getMetadataMap(), request.metadata);

    // convert GetFileResponsePB to Uint8Array
    const converter = new Transform({
      objectMode: true,
      transform(res: GetFileResponsePB, _, done) {
        const data = res.getData_asU8();
        done(null, data);
      }
    });
    const callStream = this.runtime.getFile(req, this.createMetadata(request));
    // Make sure callStream error handle by converter
    callStream.on('error', (err) => converter.emit('error', err));
    return callStream.pipe(converter);
  }

  async put(request: PutFileRequest): Promise<void> {
    const stream = request.stream ?? Readable.from(request.data);

    const ac = new AbortController();
    const signal = ac.signal;
    const writeStream = this.runtime.putFile(this.createMetadata(request), (err) => {
      if (err) {
        debug('putFile %j got server error: %s', request, err);
        // abort request and throw error
        // FIXME: should tell caller the real err reason, not only AbortError
        ac.abort();
      }
    });
    await pipelinePromise(
      stream,
      new Transform({
        objectMode: true,
        transform: (chunk: Uint8Array, _, done) => {
          const req = new PutFileRequestPB();
          req.setStoreName(request.storeName);
          req.setName(request.name);
          req.setData(chunk);
          this.mergeMetadataToMap(req.getMetadataMap(), request.metadata);
          done(null, req);
        }
      }),
      // @ts-ignore
      writeStream,
      { signal },
    );
  }

  async list(request: GetFileRequest): Promise<ListFileResponse> {
    const req = new GetFileRequestPB();
    req.setStoreName(request.storeName);
    req.setName(request.name);
    this.mergeMetadataToMap(req.getMetadataMap(), { storageType: 'Standard' }, request.metadata);
    const listReq = new ListFileRequestPB();
    listReq.setRequest(req);

    return new Promise((resolve, reject) => {
      this.runtime.listFile(listReq, this.createMetadata(request), (err, res) => {
        if (err) return reject(err);
        debug('listFile: %j, res: %j', request, res);
        resolve({
          names: res.getFileNameList(),
        });
      });
    });
  }

  async delete(request: GetFileRequest): Promise<void> {
    const req = new GetFileRequestPB();
    req.setStoreName(request.storeName);
    req.setName(request.name);
    const delReq = new DelFileRequestPB();
    delReq.setRequest(req);

    return new Promise((resolve, reject) => {
      this.runtime.delFile(delReq, this.createMetadata(request), (err) => {
        if (err) return reject(err);
        resolve();
      });
    });
  }
}
