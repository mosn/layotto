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
import { KVString, RequestWithMeta } from './common';

export type GetConfigurationRequest = {
  storeName: string;
  appId: string;
  keys: string[];
  group?: string;
  label?: string;
  subscribeUpdate?: boolean;
  metadata?: KVString;
} & RequestWithMeta;

export type GetConfigurationItem = {
  key: string;
  content: string;
  group: string;
  label: string;
  tags: KVString;
  metadata: KVString;
};

export type SaveConfigurationItem = {
  key: string,
  content: string,
  group?: string,
  label?: string,
  tags?: KVString,
  metadata?: KVString,
};

export type SaveConfigurationRequest = {
  storeName: string;
  appId: string;
  items: SaveConfigurationItem[];
  metadata?: KVString;
} & RequestWithMeta;

export type DeleteConfigurationRequest = {
  storeName: string;
  appId: string;
  keys: string[];
  group?: string;
  label?: string;
  metadata?: KVString;
} & RequestWithMeta;

export type SubscribeConfigurationRequest = {
  storeName: string;
  appId: string;
  keys: string[];
  group?: string;
  label?: string;
  metadata?: KVString;
  onData(items: GetConfigurationItem[]): void;
  onClose(err?: Error): void;
} & RequestWithMeta;
