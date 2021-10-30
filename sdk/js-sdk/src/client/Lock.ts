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
import * as crypto from 'crypto';
import {
  TryLockRequest as TryLockRequestPB,
  TryLockResponse as TryLockResponsePB,
  UnlockRequest as UnlockRequestPB,
  UnlockResponse as UnlockResponsePB,
} from '../../proto/runtime_pb';
import { API } from './API';
import { TryLockRequest, UnlockRequest } from '../types/Lock';

export default class Lock extends API {
  // A non-blocking method trying to get a lock with ttl
  // expire is the time before expire. The time unit is second.
  async tryLock(request: TryLockRequest): Promise<boolean> {
    const req = new TryLockRequestPB();
    req.setStoreName(request.storeName);
    req.setResourceId(request.resourceId);
    req.setLockOwner(request.lockOwner);
    req.setExpire(request.expire);

    return new Promise((resolve, reject) => {
      this.runtime.tryLock(req, this.createMetadata(request), (err, res: TryLockResponsePB) => {
        if (err) return reject(err);
        resolve(res.getSuccess());
      });
    });
  }

  async unLock(request: UnlockRequest): Promise<UnlockResponsePB.AsObject> {
    const req = new UnlockRequestPB();
    req.setStoreName(request.storeName);
    req.setResourceId(request.resourceId);
    req.setLockOwner(request.lockOwner);

    return new Promise((resolve, reject) => {
      this.runtime.unlock(req, this.createMetadata(request), (err, res: UnlockResponsePB) => {
        if (err) return reject(err);
        resolve(res.toObject());
      });
    });
  }

  uuid() {
    return crypto.randomUUID();
  }
}
