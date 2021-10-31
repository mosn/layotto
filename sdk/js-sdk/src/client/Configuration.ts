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
import {
  GetConfigurationRequest as GetConfigurationRequestPB,
  GetConfigurationResponse as GetConfigurationResponsePB,
  SaveConfigurationRequest as SaveConfigurationRequestPB,
  ConfigurationItem as ConfigurationItemPB,
  DeleteConfigurationRequest as DeleteConfigurationRequestPB,
  SubscribeConfigurationRequest as SubscribeConfigurationRequestPB,
  SubscribeConfigurationResponse as SubscribeConfigurationResponsePB,
} from '../../proto/runtime_pb';
import { API } from './API';
import {
  GetConfigurationRequest,
  GetConfigurationItem,
  SaveConfigurationRequest,
  DeleteConfigurationRequest,
  SubscribeConfigurationRequest,
} from '../types/Configuration';
import { convertArrayToKVString } from '../types/common';

const debug = debuglog('layotto:client:configuration');

export default class Configuration extends API {
  // GetConfiguration gets configuration from configuration store.
  async get(request: GetConfigurationRequest): Promise<GetConfigurationItem[]> {
    const req = new GetConfigurationRequestPB();
    req.setStoreName(request.storeName);
    req.setAppId(request.appId);
    req.setKeysList(request.keys);
    if (request.group) req.setGroup(request.group);
    if (request.label) req.setLabel(request.label);
    if (request.subscribeUpdate !== undefined) {
      req.setSubscribeUpdate(request.subscribeUpdate);
    }
    this.mergeMetadataToMap(req.getMetadataMap(), request.metadata);

    return new Promise((resolve, reject) => {
      this.runtime.getConfiguration(req, this.createMetadata(request), (err, res: GetConfigurationResponsePB) => {
        if (err) return reject(err);
        resolve(res.getItemsList().map(item => this.createGetConfigurationItem(item)));
      });
    });
  }

  // SaveConfiguration saves configuration into configuration store.
  async save(request: SaveConfigurationRequest): Promise<void> {
    const req = new SaveConfigurationRequestPB();
    req.setStoreName(request.storeName);
    req.setAppId(request.appId);
    req.setItemsList(request.items.map(item => {
      const configurationItem = new ConfigurationItemPB();
      configurationItem.setKey(item.key);
      configurationItem.setContent(item.content);
      if (item.group) configurationItem.setGroup(item.group);
      if (item.label) configurationItem.setLabel(item.label);
      return configurationItem;
    }));
    this.mergeMetadataToMap(req.getMetadataMap(), request.metadata);

    return new Promise((resolve, reject) => {
      this.runtime.saveConfiguration(req, this.createMetadata(request), (err) => {
        if (err) return reject(err);
        resolve();
      });
    });
  }

  // DeleteConfiguration deletes configuration from configuration store.
  async delete(request: DeleteConfigurationRequest): Promise<void> {
    const req = new DeleteConfigurationRequestPB();
    req.setStoreName(request.storeName);
    req.setAppId(request.appId);
    req.setKeysList(request.keys);
    if (request.group) req.setGroup(request.group);
    if (request.label) req.setLabel(request.label);
    this.mergeMetadataToMap(req.getMetadataMap(), request.metadata);

    return new Promise((resolve, reject) => {
      this.runtime.deleteConfiguration(req, this.createMetadata(request), (err) => {
        if (err) return reject(err);
        resolve();
      });
    });
  }

  // SubscribeConfiguration gets configuration from configuration store and subscribe the updates.
  subscribe(request: SubscribeConfigurationRequest) {
    const req = new SubscribeConfigurationRequestPB();
    req.setStoreName(request.storeName);
    req.setAppId(request.appId);
    req.setKeysList(request.keys);
    if (request.group) req.setGroup(request.group);
    if (request.label) req.setLabel(request.label);
    this.mergeMetadataToMap(req.getMetadataMap(), request.metadata);

    let lastError: Error;
    let isCloseOrEnd = false;
    const call = this.runtime.subscribeConfiguration(this.createMetadata(request));
    call.on('readable', () => {
      const res: SubscribeConfigurationResponsePB = call.read();
      debug('readable emit, has res: %s', !!res);
      if (!res) return;
      const items: GetConfigurationItem[] = res.getItemsList().map(item => this.createGetConfigurationItem(item));
      request.onData(items);
    });
    call.on('error', (err) => {
      debug('error emit, isCloseOrEnd: %s, %s', isCloseOrEnd, err);
      lastError = err;
    });
    call.on('close', () => {
      debug('close emit, isCloseOrEnd: %s', isCloseOrEnd);
      if (isCloseOrEnd) return;
      isCloseOrEnd = true;
      request.onClose(lastError);
    });
    call.on('end', () => {
      debug('end emit, isCloseOrEnd: %s', isCloseOrEnd);
      if (isCloseOrEnd) return;
      isCloseOrEnd = true;
      request.onClose(lastError);
    });
    call.write(req);
    return call;
  }

  private createGetConfigurationItem(item: ConfigurationItemPB): GetConfigurationItem {
    const obj = item.toObject();
    return {
      key: obj.key,
      content: obj.content,
      group: obj.group,
      label: obj.label,
      tags: convertArrayToKVString(obj.tagsMap),
      metadata: convertArrayToKVString(obj.metadataMap),
    };
  }
}
