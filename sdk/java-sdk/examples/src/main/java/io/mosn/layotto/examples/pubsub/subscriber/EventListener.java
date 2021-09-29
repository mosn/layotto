package io.mosn.layotto.examples.pubsub.subscriber;

import spec.sdk.runtime.v1.domain.pubsub.TopicEventRequest;

public interface EventListener {
    void onEvent(TopicEventRequest request) throws Exception;
}
