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