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
import { KV, RequestWithMeta } from './common';

export type GetConfigurationRequest = RequestWithMeta<{
  storeName: string;
  appId: string;
  keys: string[];
  group?: string;
  label?: string;
  subscribeUpdate?: boolean;
  metadata?: KV<string>;
}>;

export type GetConfigurationItem = {
  key: string;
  content: string;
  group: string;
  label: string;
  tags: KV<string>;
  metadata: KV<string>;
};

export type SaveConfigurationItem = {
  key: string,
  content: string,
  group?: string,
  label?: string,
  tags?: KV<string>,
  metadata?: KV<string>,
};

export type SaveConfigurationRequest = RequestWithMeta<{
  storeName: string;
  appId: string;
  items: SaveConfigurationItem[];
  metadata?: KV<string>;
}>;

export type DeleteConfigurationRequest = RequestWithMeta<{
  storeName: string;
  appId: string;
  keys: string[];
  group?: string;
  label?: string;
  metadata?: KV<string>;
}>;

export type SubscribeConfigurationRequest = RequestWithMeta<{
  storeName: string;
  appId: string;
  keys: string[];
  group?: string;
  label?: string;
  metadata?: KV<string>;
  onData(items: GetConfigurationItem[]): void;
  onClose(err?: Error): void;
}>;
