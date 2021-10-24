import * as crypto from 'crypto';
import { TryLockRequest, TryLockResponse, UnlockRequest, UnlockResponse } from '../proto/runtime_pb';
import { API, RequestMetadata } from './API';

export default class Lock extends API {
  // A non-blocking method trying to get a lock with ttl
  // expire is the time before expire. The time unit is second.
  async tryLock(storeName: string, resourceId: string, lockOwner: string, expire: number, meta?: RequestMetadata): Promise<boolean> {
    const req = new TryLockRequest();
    req.setStoreName(storeName);
    req.setResourceId(resourceId);
    req.setLockOwner(lockOwner);
    req.setExpire(expire);

    return new Promise((resolve, reject) => {
      this.runtime.tryLock(req, this.createMetadata(meta), (err, res: TryLockResponse) => {
        if (err) return reject(err);
        resolve(res.getSuccess());
      });
    });
  }

  async unLock(storeName: string, resourceId: string, lockOwner: string, meta?: RequestMetadata): Promise<UnlockResponse.Status> {
    const req = new UnlockRequest();
    req.setStoreName(storeName);
    req.setResourceId(resourceId);
    req.setLockOwner(lockOwner);

    return new Promise((resolve, reject) => {
      this.runtime.unlock(req, this.createMetadata(meta), (err, res: UnlockResponse) => {
        if (err) return reject(err);
        resolve(res.getStatus());
      });
    });
  }

  uuid() {
    return crypto.randomUUID();
  }
}
