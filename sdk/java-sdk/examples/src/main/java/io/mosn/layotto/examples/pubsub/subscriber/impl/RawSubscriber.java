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
package io.mosn.layotto.examples.pubsub.subscriber.impl;

import io.mosn.layotto.v1.callback.component.pubsub.Subscriber;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import spec.sdk.runtime.v1.domain.pubsub.TopicEventRequest;
import spec.sdk.runtime.v1.domain.pubsub.TopicEventResponse;
import spec.sdk.runtime.v1.domain.pubsub.TopicEventResponseStatus;
import spec.sdk.runtime.v1.domain.pubsub.TopicSubscription;

import java.util.HashSet;
import java.util.Map;
import java.util.Set;
import java.util.concurrent.ConcurrentHashMap;

/**
 * Raw pubsub client.
 */
public class RawSubscriber implements Subscriber {

    private static final Logger LOG = LoggerFactory.getLogger(RawSubscriber.class);

    private final Map<String, EventListener> listeners = new ConcurrentHashMap<>();
    private final String                     componentName;

    public RawSubscriber(String pubsubName) {
        componentName = pubsubName;
    }

    @Override
    public String getComponentName() {
        return componentName;
    }

    public void subscribe(String topic, EventListener listener) {
        if (listeners.putIfAbsent(topic, listener) != null) {
            throw new IllegalArgumentException("Listener for topic " + topic + " already exists!");
        }
    }

    @Override
    public Set<TopicSubscription> listTopicSubscriptions() {
        final HashSet<TopicSubscription> subscriptions = new HashSet<>();
        for (String topic : listeners.keySet()) {
            final TopicSubscription subscription = new TopicSubscription();
            subscription.setTopic(topic);
            subscription.setPubsubName(componentName);
            subscriptions.add(subscription);
        }
        return subscriptions;
    }

    @Override
    public TopicEventResponse onTopicEvent(TopicEventRequest request) {
        final String topic = request.getTopic();
        final EventListener eventListener = listeners.get(topic);
        if (eventListener == null) {
            LOG.error("Cannot find listener for topic:[{}]", topic);
            TopicEventResponse resp = new TopicEventResponse();
            resp.setStatus(TopicEventResponseStatus.DROP);
        }
        try {
            eventListener.onEvent(request);
            final TopicEventResponse response = new TopicEventResponse();
            response.setStatus(TopicEventResponseStatus.SUCCESS);
            return response;
        } catch (Exception e) {
            final TopicEventResponse response = new TopicEventResponse();
            response.setStatus(TopicEventResponseStatus.RETRY);
            return response;
        }
    }
}
