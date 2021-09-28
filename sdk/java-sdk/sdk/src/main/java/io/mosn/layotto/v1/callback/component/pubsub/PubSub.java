package io.mosn.layotto.v1.callback.component.pubsub;

import io.mosn.layotto.v1.callback.component.Component;
import spec.sdk.runtime.v1.domain.pubsub.TopicEventRequest;
import spec.sdk.runtime.v1.domain.pubsub.TopicEventResponse;
import spec.sdk.runtime.v1.domain.pubsub.TopicSubscription;

import java.util.Set;

/**
 * Pub/Sub client.
 */
public interface PubSub extends Component {

    Set<TopicSubscription> listTopicSubscriptions();

    TopicEventResponse onTopicEvent(TopicEventRequest request);
}
