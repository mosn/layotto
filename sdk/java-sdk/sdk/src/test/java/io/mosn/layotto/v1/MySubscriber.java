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
package io.mosn.layotto.v1;

import io.mosn.layotto.v1.callback.component.pubsub.Subscriber;
import spec.sdk.runtime.v1.domain.pubsub.TopicEventRequest;
import spec.sdk.runtime.v1.domain.pubsub.TopicEventResponse;
import spec.sdk.runtime.v1.domain.pubsub.TopicEventResponseStatus;
import spec.sdk.runtime.v1.domain.pubsub.TopicSubscription;

import java.util.HashSet;
import java.util.Set;

public class MySubscriber implements Subscriber {
    private final String pubsubName;
    private final String topic;

    public MySubscriber(String pubsubName, String topic) {
        this.pubsubName = pubsubName;
        this.topic = topic;
    }

    @Override
    public Set<TopicSubscription> listTopicSubscriptions() {
        TopicSubscription topicSub = new TopicSubscription();
        topicSub.setPubsubName(pubsubName);
        topicSub.setTopic(topic);
        Set<TopicSubscription> set = new HashSet<>();
        set.add(topicSub);
        return set;
    }

    @Override
    public TopicEventResponse onTopicEvent(TopicEventRequest request) {
        TopicEventResponse resp = new TopicEventResponse();
        resp.setStatus(TopicEventResponseStatus.SUCCESS);
        if (!topic.equals(request.getTopic())) {
            resp.setStatus(TopicEventResponseStatus.DROP);
        }
        if (!pubsubName.equals(request.getPubsubName())) {
            resp.setStatus(TopicEventResponseStatus.DROP);
        }
        return resp;
    }

    @Override
    public String getComponentName() {
        return "redis";
    }
}