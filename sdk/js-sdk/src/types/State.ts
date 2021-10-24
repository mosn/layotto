import { 
  StateOptions as StateOptionsPB,
} from '../../proto/runtime_pb';

export type StateConcurrency = StateOptionsPB.StateConcurrency;
export type StateConsistency = StateOptionsPB.StateConsistency;

export type StateOptions = {
  concurrency: StateConcurrency;
  consistency: StateConsistency;
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

