package spec.sdk.runtime.v1.domain.core;

import spec.sdk.runtime.v1.domain.core.pubsub.PublishEventRequest;
import reactor.core.publisher.Mono;

import java.util.Map;

/**
 * Publish and Subscribe Runtimes standard API defined.
 */
public interface PubSubRuntimes {

    /**
     * Publish an event.
     *
     * @param pubsubName the pubsub name we will publish the event to
     * @param topicName  the topicName where the event will be published.
     * @param data       the event's data to be published, use byte[] for skipping serialization.
     * @return a Mono plan of type Void.
     */
    Mono<Void> publishEvent(String pubsubName, String topicName, Object data);

    /**
     * Publish an event.
     *
     * @param pubsubName the pubsub name we will publish the event to
     * @param topicName  the topicName where the event will be published.
     * @param data       the event's data to be published, use byte[] for skipping serialization.
     * @param metadata   The metadata for the published event.
     * @return a Mono plan of type Void.
     */
    Mono<Void> publishEvent(String pubsubName, String topicName, Object data, Map<String, String> metadata);

    /**
     * Publish an event.
     *
     * @param request the request for the publish event.
     * @return a Mono plan of a CloudRuntimes's void response.
     */
    Mono<Void> publishEvent(PublishEventRequest request);
}
