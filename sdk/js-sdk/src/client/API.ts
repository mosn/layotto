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
import { Metadata } from '@grpc/grpc-js';
import { RuntimeClient } from '../../proto/runtime_grpc_pb';
import { KV, RequestWithMeta, Map } from '../types/common';

export class API {
  readonly runtime: RuntimeClient;
  constructor(runtime: RuntimeClient) {
    this.runtime = runtime;
  }

  createMetadata(request: RequestWithMeta<{}>): Metadata {
    const metadata = new Metadata();
    if (!request.requestMeta) return metadata;
    for (const key of Object.keys(request.requestMeta)) {
      metadata.add(key, request.requestMeta[key]);
    }
    return metadata;
  }

  mergeMetadataToMap(map: Map<string>, ...metadatas: (KV<string> | undefined)[]) {
    for (const metadata of metadatas) {
      if (!metadata) continue;
      for (const key of Object.keys(metadata)) {
        map.set(key, metadata[key]);
      }
    }
  }
}
