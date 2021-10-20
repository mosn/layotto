package io.mosn.layotto.v1.callback.component.pubsub;

import java.util.Collection;
import java.util.HashSet;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

public class PubSubClientRegistryImpl implements PubSubRegistry {

    public final Map<String, PubSub> pubSubClients = new ConcurrentHashMap<>();

    @Override
    public void registerPubSubCallback(String pubsubName, PubSub callback) {
        if (pubSubClients.putIfAbsent(pubsubName, callback) != null) {
            throw new IllegalArgumentException("Pub/sub callback with name " + pubsubName + " already exists!");
        }
    }

    @Override
    public PubSub getCallbackByPubSubName(String pubSubName) {
        final PubSub pubSub = pubSubClients.get(pubSubName);
        if (pubSub != null) {
            return pubSub;
        }

        throw new IllegalArgumentException("Cannot find pubsub callback by name " + pubSubName);
    }

    @Override
    public Collection<PubSub> getAllPubSubCallbacks() {
        return new HashSet<>(pubSubClients.values());
    }
}
