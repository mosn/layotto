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
package spec.sdk.reactor.v1.domain.core.pubsub;

import java.util.Collections;
import java.util.HashMap;
import java.util.Map;

/**
 * A request to publish an event.
 */
public class PublishEventRequest {

    private final String pubsubName;

    private final String topic;

    private final Object data;

    private String contentType;

    private Map<String, String> metadata = new HashMap<>();

    /**
     * Constructor for PublishEventRequest.
     *
     * @param pubsubName name of the pubsub
     * @param topic      name of the topic in the pubsub
     * @param data       data to published
     */
    public PublishEventRequest(String pubsubName, String topic, Object data) {
        this.pubsubName = pubsubName;
        this.topic = topic;
        this.data = data;
    }

    public String getPubsubName() {
        return pubsubName;
    }

    public String getTopic() {
        return topic;
    }

    public Object getData() {
        return data;
    }

    public String getContentType() {
        return this.contentType;
    }

    public PublishEventRequest setContentType(String contentType) {
        this.contentType = contentType;
        return this;
    }

    public Map<String, String> getMetadata() {
        return metadata;
    }

    public PublishEventRequest setMetadata(Map<String, String> metadata) {
        this.metadata = metadata == null ? null : Collections.unmodifiableMap(metadata);
        return this;
    }
}
