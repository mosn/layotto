import { 
  SaveStateRequest,
  StateItem as StateItemPB,
  Etag as ETagTB,
  StateOptions as StateOptionsPB,
  GetStateRequest,
  GetBulkStateRequest,
} from '../proto/runtime_pb';
import API from './API';

export type StateConcurrency = StateOptionsPB.StateConcurrency;
export type StateConsistency = StateOptionsPB.StateConsistency;

type StateOptions = {
  concurrency: StateConcurrency;
  consistency: StateConsistency;
};

type Metadata = {
  [key: string]: string;
};

type StateItem = {
  key: string;
  value: Uint8Array | string;
  etag?: string;
  metadata?: Metadata;
  options?: StateOptions;
}

type BulkResponseStateItem = {
  key: string;
  value: Uint8Array;
  etag: string;
}

export default class State extends API {
  // Saves an array of state objects
  async save(storeName: string, states: StateItem[] | StateItem): Promise<void> {
    const stateList: StateItemPB[] = [];
    if (!Array.isArray(states)) {
      states = [states];
    }
    for (const item of states) {
      const stateItem = new StateItemPB();
      stateItem.setKey(item.key);
      if (typeof item.value === 'string') {
        stateItem.setValue(Buffer.from(item.value, 'utf8'));
      } else {
        stateItem.setValue(item.value);
      }
      if (item.etag !== undefined) {
        const etag = new ETagTB();
        etag.setValue(item.etag);
        stateItem.setEtag(etag);
      }
      if (item.options !== undefined) {
        const options = new StateOptionsPB();
        options.setConcurrency(item.options.concurrency);
      }
      stateList.push(stateItem);
    }

    const req = new SaveStateRequest();
    req.setStoreName(storeName);
    req.setStatesList(stateList);

    return new Promise((resolve, reject) => {
      this.runtime.saveState(req, (err) => {
        if (err) return reject(err);
        resolve();
      });
    });
  }

  // Gets the state for a specific key
  async get(storeName: string, key: string): Promise<Uint8Array> {
    const req = new GetStateRequest();
    req.setStoreName(storeName);
    req.setKey(key);

    return new Promise((resolve, reject) => {
      this.runtime.getState(req, (err, res) => {
        if (err) return reject(err);
        resolve(res.getData_asU8());
      });
    });
  }

  // Gets a bulk of state items for a list of keys
  async getBulk(storeName: string, keys: string[], parallelism = 10): Promise<BulkResponseStateItem[]> {
    const req = new GetBulkStateRequest();
    req.setStoreName(storeName);
    req.setKeysList(keys);
    req.setParallelism(parallelism);

    return new Promise((resolve, reject) => {
      this.runtime.getBulkState(req, (err, res) => {
        if (err) return reject(err);
        const states: BulkResponseStateItem[] = [];
        const itemsList = res.getItemsList();
        for (const item of itemsList) {
          states.push({
            key: item.getKey(),
            value: item.getData_asU8(),
            etag: item.getEtag(),
          });
        }
        resolve(states);
      });
    });
  }

  // Deletes the state for a specific key
  async delete() {

  }

  // Deletes a bulk of state items for a list of keys
  async deleteBulk() {}

  // async transaction() {}
}
