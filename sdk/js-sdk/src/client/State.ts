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
  SaveStateRequest as SaveStateRequestPB,
  StateItem as StateItemPB,
  Etag as EtagPB,
  StateOptions as StateOptionsPB,
  GetStateRequest as GetStateRequestPB,
  GetStateResponse as GetStateResponsePB,
  GetBulkStateRequest as GetBulkStateRequestPB,
  DeleteStateRequest as DeleteStateRequestPB,
  DeleteBulkStateRequest as DeleteBulkStateRequestPB,
  ExecuteStateTransactionRequest as ExecuteStateTransactionRequestPB,
  TransactionalStateOperation as TransactionalStateOperationPB,
} from '../../proto/runtime_pb';
import { API } from './API';
import {
  DeleteBulkStateRequest,
  DeleteStateItem,
  DeleteStateRequest,
  ExecuteStateTransactionRequest,
  GetBulkStateRequest,
  GetStateRequest,
  ResponseStateItem,
  SaveStateRequest,
  StateItem,
} from '../types/State';
import { isEmptyPBMessage, convertMapToKVString } from '../utils';

export default class State extends API {
  // Saves an array of state objects
  async save(request: SaveStateRequest): Promise<void> {
    let states = request.states;
    if (!Array.isArray(states)) {
      states = [states];
    }
    const stateList = this.createStateItemPBList(states);
    const req = new SaveStateRequestPB();
    req.setStoreName(request.storeName);
    req.setStatesList(stateList);

    return new Promise((resolve, reject) => {
      this.runtime.saveState(req, this.createMetadata(request), (err) => {
        if (err) return reject(err);
        resolve();
      });
    });
  }

  // Gets the state for a specific key
  async get(request: GetStateRequest): Promise<ResponseStateItem | null> {
    const req = new GetStateRequestPB();
    req.setStoreName(request.storeName);
    req.setKey(request.key);
    this.mergeMetadataToMap(req.getMetadataMap(), request.metadata);

    return new Promise((resolve, reject) => {
      this.runtime.getState(req, this.createMetadata(request), (err, res: GetStateResponsePB) => {
        if (err) return reject(err);
        if (isEmptyPBMessage(res)) {
          return resolve(null);
        }
        resolve({
          key: request.key,
          value: res.getData_asU8(),
          etag: res.getEtag(),
          metadata: convertMapToKVString(res.getMetadataMap()),
        });
      });
    });
  }

  // Gets a bulk of state items for a list of keys
  async getBulk(request: GetBulkStateRequest): Promise<ResponseStateItem[]> {
    const req = new GetBulkStateRequestPB();
    req.setStoreName(request.storeName);
    req.setKeysList(request.keys);
    if (request.parallelism) req.setParallelism(request.parallelism);
    this.mergeMetadataToMap(req.getMetadataMap(), request.metadata);

    return new Promise((resolve, reject) => {
      this.runtime.getBulkState(req, this.createMetadata(request), (err, res) => {
        if (err) return reject(err);
        const states: ResponseStateItem[] = [];
        const itemsList = res.getItemsList();
        for (const item of itemsList) {
          // pb.message.array[0] is key, pb.message.array[1] is value
          if (isEmptyPBMessage(item, 1)) {
            continue;
          }
          states.push({
            key: item.getKey(),
            value: item.getData_asU8(),
            etag: item.getEtag(),
            metadata: convertMapToKVString(item.getMetadataMap()),
          });
        }
        resolve(states);
      });
    });
  }

  // Deletes the state for a specific key
  async delete(request: DeleteStateRequest): Promise<void> {
    const req = new DeleteStateRequestPB();
    req.setStoreName(request.storeName);
    req.setKey(request.key);
    if (request.etag) {
      const etagInstance = new EtagPB();
      etagInstance.setValue(request.etag);
      req.setEtag(etagInstance);
    }
    if (request.options) {
      const optionsInstance = new StateOptionsPB();
      optionsInstance.setConcurrency(request.options.concurrency);
      optionsInstance.setConsistency(request.options.consistency);
      req.setOptions(optionsInstance);
    }
    this.mergeMetadataToMap(req.getMetadataMap(), request.metadata);

    return new Promise((resolve, reject) => {
      this.runtime.deleteState(req, this.createMetadata(request), (err) => {
        if (err) return reject(err);
        resolve();
      });
    });
  }

  // Deletes a bulk of state items for a list of keys
  async deleteBulk(request: DeleteBulkStateRequest): Promise<void> {
    const req = new DeleteBulkStateRequestPB();
    req.setStoreName(request.storeName);
    const stateList = this.createStateItemPBList(request.states);
    req.setStatesList(stateList);

    return new Promise((resolve, reject) => {
      this.runtime.deleteBulkState(req, this.createMetadata(request), (err) => {
        if (err) return reject(err);
        resolve();
      });
    });
  }

  // Executes transactions for a specified store
  async executeTransaction(request: ExecuteStateTransactionRequest): Promise<void> {
    const req = new ExecuteStateTransactionRequestPB();
    req.setStorename(request.storeName);
    const operationsList: TransactionalStateOperationPB[] = [];
    for (const operation of request.operations) {
      const ops = new TransactionalStateOperationPB();
      ops.setOperationtype(operation.operationType);
      const stateItem = this.createStateItemPB(operation.request);
      ops.setRequest(stateItem);
      operationsList.push(ops);
    }
    req.setOperationsList(operationsList);
    this.mergeMetadataToMap(req.getMetadataMap(), request.metadata);

    return new Promise((resolve, reject) => {
      this.runtime.executeStateTransaction(req, this.createMetadata(request), (err, _res) => {
        if (err) return reject(err);
        resolve();
      });
    });
  }

  private createStateItemPB(item: StateItem | DeleteStateItem): StateItemPB {
    const stateItem = new StateItemPB();
    stateItem.setKey(item.key);
    if ('value' in item) {
      if (typeof item.value === 'string') {
        stateItem.setValue(Buffer.from(item.value, 'utf8'));
      } else {
        stateItem.setValue(item.value);
      }
    }
    if (item.etag !== undefined) {
      const etag = new EtagPB();
      etag.setValue(item.etag);
      stateItem.setEtag(etag);
    }
    if (item.options !== undefined) {
      const options = new StateOptionsPB();
      options.setConcurrency(item.options.concurrency);
      options.setConsistency(item.options.consistency);
      stateItem.setOptions(options);
    }
    return stateItem;
  }

  private createStateItemPBList(items: StateItem[] | DeleteStateItem[]): StateItemPB[] {
    const list: StateItemPB[] = [];
    for (const item of items) {
      list.push(this.createStateItemPB(item));
    }
    return list;
  }
}
