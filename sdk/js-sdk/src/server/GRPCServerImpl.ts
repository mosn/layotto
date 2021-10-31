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
 *
 */
import { debuglog } from 'node:util';
import * as grpc from '@grpc/grpc-js';
import { Empty } from 'google-protobuf/google/protobuf/empty_pb';
import { IAppCallbackServer } from '../../proto/appcallback_grpc_pb';
import {
  ListTopicSubscriptionsResponse,
  TopicSubscription,
  TopicEventRequest,
  TopicEventResponse,
} from '../../proto/appcallback_pb';
import { PubSubCallback } from '../types/PubSub';

const debug = debuglog('layotto:server:grpc');

// @ts-ignore
export default class GRPCServerImpl implements IAppCallbackServer {
  private readonly _handlersTopics: { [key: string]: PubSubCallback };
  constructor() {
    this._handlersTopics = {};
  }

  private createPubSubHandlerKey(pubsubName: string, topic: string): string {
    return `${pubsubName}|${topic}`.toLowerCase();
  }

  registerPubSubSubscriptionHandler(pubsubName: string, topic: string, callback: PubSubCallback): void {
    const handlerKey = this.createPubSubHandlerKey(pubsubName, topic);
    if (this._handlersTopics[handlerKey]) {
      throw new Error(`Topic: "${handlerKey}" handler was exists`);
    }
    this._handlersTopics[handlerKey] = callback;
    debug('PubSub Event from topic: "%s" is registered', handlerKey);
  }

  async onTopicEvent(call: grpc.ServerUnaryCall<TopicEventRequest, TopicEventResponse>,
      callback: grpc.sendUnaryData<TopicEventResponse>): Promise<void> {
    const req = call.request;
    const res = new TopicEventResponse();
    const handlerKey = this.createPubSubHandlerKey(req.getPubsubName(), req.getTopic());
    
    const handler = this._handlersTopics[handlerKey];
    if (!handler) {
        debug('PubSub Event from topic: "%s" was not handled, drop now', handlerKey);
        // FIXME: should retry?
        res.setStatus(TopicEventResponse.TopicEventResponseStatus.DROP);
        return callback(null, res);
    }

    // https://mosn.io/layotto/#/zh/design/pubsub/pubsub-api-and-compability-with-dapr-component
    // PublishRequest.Data 和 NewMessage.Data 里面放符合 CloudEvent 1.0 规范的 json 数据（能反序列化放进 map[string]interface{}）
    const rawData = Buffer.from(req.getData_asU8()).toString();
    debug('PubSub Event from topic: "%s" raw data: %j, typeof %s', handlerKey, rawData, typeof rawData);
    let data: string | object;
    try {
      data = JSON.parse(rawData);
    } catch {
      data = rawData;
    }
    
    try {
      await handler(data);
      res.setStatus(TopicEventResponse.TopicEventResponseStatus.SUCCESS);
    } catch (e) {
      // FIXME: should retry?
      debug('PubSub Event from topic: "%s" handler throw error: %s, drop now', handlerKey, e);
      res.setStatus(TopicEventResponse.TopicEventResponseStatus.DROP);
    }

    callback(null, res);
  }

  async listTopicSubscriptions(_call: grpc.ServerUnaryCall<Empty, ListTopicSubscriptionsResponse>,
      callback: grpc.sendUnaryData<ListTopicSubscriptionsResponse>): Promise<void> {
    const res = new ListTopicSubscriptionsResponse();
    const subscriptionsList = Object.keys(this._handlersTopics).map(key => {
      const splits = key.split('|');
      const sub = new TopicSubscription();
      sub.setPubsubName(splits[0]);
      sub.setTopic(splits[1]);
      return sub;
    });
    debug('listTopicSubscriptions call: %j', subscriptionsList);
    res.setSubscriptionsList(subscriptionsList);
    callback(null, res);
  }
}
