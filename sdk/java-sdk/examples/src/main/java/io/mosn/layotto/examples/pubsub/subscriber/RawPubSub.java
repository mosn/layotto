package io.mosn.layotto.examples.pubsub.subscriber;

import io.mosn.layotto.v1.callback.component.pubsub.PubSub;
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
public class RawPubSub implements PubSub {
    private static final Logger LOG = LoggerFactory.getLogger(RawPubSub.class);

    private final Map<String, EventListener> listeners = new ConcurrentHashMap<>();
    private final String                     componentName;

    public RawPubSub(String pubsubName) {
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
