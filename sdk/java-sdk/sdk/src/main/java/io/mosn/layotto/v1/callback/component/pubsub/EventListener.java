package io.mosn.layotto.v1.callback.component.pubsub;

import spec.sdk.runtime.v1.domain.pubsub.TopicEventRequest;

public interface EventListener {
    void onEvent(TopicEventRequest request) throws Exception;
}
