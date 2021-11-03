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
import { setTimeout } from 'timers/promises';
import { Map as MapPB } from 'google-protobuf';
import { KV } from './types/common';

export async function sleep(ms: number) {
  await setTimeout(ms);
}

// jspb.Message
// https://github.com/protocolbuffers/protobuf/blob/master/js/message.js#L233
//
// use array to detect not exists or empty
// console.log(res)
// not exists
// {
//   wrappers_: null,
//   messageId_: undefined,
//   arrayIndexOffset_: -1,
//   array: [],
//   pivot_: 1.7976931348623157e+308,
//   convertedPrimitiveFields_: {}
// }
// empty data
// {
//   wrappers_: null,
//   messageId_: undefined,
//   arrayIndexOffset_: -1,
//   array: [ <1 empty item>, '1' ],
//   pivot_: 1.7976931348623157e+308,
//   convertedPrimitiveFields_: {}
// }
export function isEmptyPBMessage(item, emptyLength = 0) {
  if (item.array.length === emptyLength) return true;
  return false;
}

export function convertMapToKVString(map: MapPB<string, string>) {
  const kv: KV<string> = {};
  for (const [k, v] of map.entries()) {
    kv[k] = v;
  }
  return kv;
};
