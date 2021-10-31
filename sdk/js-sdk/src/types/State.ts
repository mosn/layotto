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
  StateOptions as StateOptionsPB,
} from '../../proto/runtime_pb';
import { KVString, RequestWithMeta } from './common';

export type StateOptions = {
  concurrency: StateOptionsPB.StateConcurrency;
  consistency: StateOptionsPB.StateConsistency;
};

export type StateItem = {
  key: string;
  value: Uint8Array | string;
  etag?: string;
  options?: StateOptions;
}

export type DeleteStateItem = {
  key: string;
  etag?: string;
  options?: StateOptions;
};

export type ResponseStateItem = {
  key: string;
  value: Uint8Array;
  etag: string;
  // metadata;
}

export enum StateOperationType {
  Upsert = 'upsert',
  Delete = 'delete',
}

export type StateOperation = {
  operationType: StateOperationType;
  request: StateItem | DeleteStateItem;
};

export type SaveStateRequest = {
  storeName: string;
  states: StateItem[] | StateItem;
} & RequestWithMeta;

export type GetStateRequest = {
  storeName: string;
  key: string;
  metadata?: KVString;
} & RequestWithMeta;

export type GetBulkStateRequest = {
  storeName: string;
  keys: string[];
  parallelism?: number;
  metadata?: KVString;
} & RequestWithMeta;

export type DeleteStateRequest = {
  storeName: string;
  key: string;
  etag?: string;
  options?: StateOptions
  metadata?: KVString;
} & RequestWithMeta;

export type DeleteBulkStateRequest = {
  storeName: string;
  states: DeleteStateItem[];
} & RequestWithMeta;

export type ExecuteStateTransactionRequest = {
  storeName: string;
  operations: StateOperation[];
  metadata?: KVString;
} & RequestWithMeta;
