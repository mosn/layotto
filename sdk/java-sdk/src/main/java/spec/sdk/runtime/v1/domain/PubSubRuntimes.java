package spec.sdk.runtime.v1.domain;

import java.util.Map;

public interface PubSubRuntimes {

    void publishEvent(String pubsubName, String topicName, byte[] data);

    /**
     * Publish an event.
     *
     * @param pubsubName the pubsub name we will publish the event to
     * @param topicName  the topicName where the event will be published.
     * @param data       the event's data to be published, use byte[] for skipping serialization.
     */
    void publishEvent(String pubsubName, String topicName, Object data);

    /**
     * Publish an event.
     *
     * @param pubsubName the pubsub name we will publish the event to
     * @param topicName  the topicName where the event will be published.
     * @param data       the event's data to be published, use byte[] for skipping serialization.
     * @param metadata   The metadata for the published event.
     */
    void publishEvent(String pubsubName, String topicName, Object data, Map<String, String> metadata);

    void publishEvent(String pubsubName, String topicName, byte[] data, Map<String, String> metadata);
}
