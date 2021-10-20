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
 *
 */
package spec.sdk.runtime.v1.domain.pubsub;

import java.util.Map;

public class TopicSubscription {
    // Required. The name of the pubsub containing the topic below to subscribe to.
    private String pubsubName;

    // Required. The name of topic which will be subscribed
    private String topic;

    // The optional properties used for this topic's subscription e.g. session id
    private Map<String, String> metadata;

    /**
     * Getter method for property <tt>pubsub_name</tt>.
     *
     * @return property value of pubsub_name
     */
    public String getPubsubName() {
        return pubsubName;
    }

    /**
     * Setter method for property <tt>pubsub_name</tt>.
     *
     * @param pubsubName value to be assigned to property pubsub_name
     */
    public void setPubsubName(String pubsubName) {
        this.pubsubName = pubsubName;
    }

    /**
     * Getter method for property <tt>topic</tt>.
     *
     * @return property value of topic
     */
    public String getTopic() {
        return topic;
    }

    /**
     * Setter method for property <tt>topic</tt>.
     *
     * @param topic value to be assigned to property topic
     */
    public void setTopic(String topic) {
        this.topic = topic;
    }

    /**
     * Getter method for property <tt>metadata</tt>.
     *
     * @return property value of metadata
     */
    public Map<String, String> getMetadata() {
        return metadata;
    }

    /**
     * Setter method for property <tt>metadata</tt>.
     *
     * @param metadata value to be assigned to property metadata
     */
    public void setMetadata(Map<String, String> metadata) {
        this.metadata = metadata;
    }
}
