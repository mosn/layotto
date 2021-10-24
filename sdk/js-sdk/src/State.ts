import { 
  SaveStateRequest,
  StateItem as StateItemPB,
  Etag,
  StateOptions as StateOptionsPB,
  GetStateRequest,
  GetBulkStateRequest,
  DeleteStateRequest,
  DeleteBulkStateRequest,
  ExecuteStateTransactionRequest,
  TransactionalStateOperation,
} from '../proto/runtime_pb';
import { API, RequestMetadata } from './API';
import {
  DeleteStateItem,
  ResponseStateItem,
  StateItem,
  StateOperation,
  StateOptions,
} from './types/State';

export default class State extends API {
  // Saves an array of state objects
  async save(storeName: string, states: StateItem[] | StateItem, meta?: RequestMetadata): Promise<void> {
    if (!Array.isArray(states)) {
      states = [states];
    }
    const stateList = this.createStateItemPBList(states);
    const req = new SaveStateRequest();
    req.setStoreName(storeName);
    req.setStatesList(stateList);

    return new Promise((resolve, reject) => {
      this.runtime.saveState(req, this.createMetadata(meta), (err) => {
        if (err) return reject(err);
        resolve();
      });
    });
  }

  // Gets the state for a specific key
  async get(storeName: string, key: string, meta?: RequestMetadata): Promise<ResponseStateItem | null> {
    const req = new GetStateRequest();
    req.setStoreName(storeName);
    req.setKey(key);

    return new Promise((resolve, reject) => {
      this.runtime.getState(req, this.createMetadata(meta), (err, res) => {
        if (err) return reject(err);
        if (this.isEmpty(res)) {
          return resolve(null);
        }
        resolve({
          key,
          value: res.getData_asU8(),
          etag: res.getEtag(),
        });
      });
    });
  }

  // Gets a bulk of state items for a list of keys
  async getBulk(storeName: string, keys: string[], parallelism = 10, meta?: RequestMetadata): Promise<ResponseStateItem[]> {
    const req = new GetBulkStateRequest();
    req.setStoreName(storeName);
    req.setKeysList(keys);
    req.setParallelism(parallelism);

    return new Promise((resolve, reject) => {
      this.runtime.getBulkState(req, this.createMetadata(meta), (err, res) => {
        if (err) return reject(err);
        const states: ResponseStateItem[] = [];
        const itemsList = res.getItemsList();
        for (const item of itemsList) {
          if (this.isEmpty(item)) {
            continue;
          }
          states.push({
            key: item.getKey(),
            value: item.getData_asU8(),
            etag: item.getEtag(),
            // metadata: item.getMetadataMap(),
          });
        }
        resolve(states);
      });
    });
  }

  // Deletes the state for a specific key
  async delete(storeName: string, key: string, etag = '', options?: StateOptions, meta?: RequestMetadata): Promise<void> {
    const req = new DeleteStateRequest();
    req.setStoreName(storeName);
    req.setKey(key);
    if (etag) {
      const etagInstance = new Etag();
      etagInstance.setValue(etag);
      req.setEtag(etagInstance);
    }
    if (options) {
      const optionsInstance = new StateOptionsPB();
      optionsInstance.setConcurrency(options.concurrency);
      optionsInstance.setConsistency(options.consistency);
      req.setOptions(optionsInstance);
    }

    return new Promise((resolve, reject) => {
      this.runtime.deleteState(req, this.createMetadata(meta), (err) => {
        if (err) return reject(err);
        resolve();
      });
    });
  }

  // Deletes a bulk of state items for a list of keys
  async deleteBulk(storeName: string, states: DeleteStateItem[], meta?: RequestMetadata): Promise<void> {
    const req = new DeleteBulkStateRequest();
    req.setStoreName(storeName);
    const stateList = this.createStateItemPBList(states);
    req.setStatesList(stateList);

    return new Promise((resolve, reject) => {
      this.runtime.deleteBulkState(req, this.createMetadata(meta), (err) => {
        if (err) return reject(err);
        resolve();
      });
    });
  }

  // Executes transactions for a specified store
  async executeTransaction(storeName: string, operations: StateOperation[], meta?: RequestMetadata): Promise<void> {
    const req = new ExecuteStateTransactionRequest();
    req.setStorename(storeName);
    const operationsList: TransactionalStateOperation[] = [];
    for (const operation of operations) {
      const ops = new TransactionalStateOperation();
      ops.setOperationtype(operation.operationType);
      const stateItem = this.createStateItemPB(operation.request);
      ops.setRequest(stateItem);
      operationsList.push(ops);
    }
    req.setOperationsList(operationsList);

    return new Promise((resolve, reject) => {
      this.runtime.executeStateTransaction(req, this.createMetadata(meta), (err, _res) => {
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
      const etag = new Etag();
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

  private isEmpty(obj: { getEtag(): string }) {
    return obj.getEtag() === '';
  }
}
