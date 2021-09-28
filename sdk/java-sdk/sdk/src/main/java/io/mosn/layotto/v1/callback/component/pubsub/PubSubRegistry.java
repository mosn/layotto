package io.mosn.layotto.v1.callback.component.pubsub;

import java.util.Collection;

public interface PubSubRegistry {

    void registerPubSubCallback(String pubsubName, PubSub callback);

    PubSub getCallbackByPubSubName(String pubSubName);

    Collection<PubSub> getAllPubSubCallbacks();
}
