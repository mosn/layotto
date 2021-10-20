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

public class PublishEventRequest {
    // The name of the pubsub component
    private String              pubsubName;
    // The pubsub topic
    private String              topic;
    // The data which will be published to topic.
    private byte[]              data;

    // The content type for the data (optional).
    private String              contentType;

    // The metadata passing to pub components
    //
    // metadata property:
    // - key : the key of the message.
    private Map<String, String> metadata;

    /**
     * Getter method for property <tt>pubsubName</tt>.
     *
     * @return property value of pubsubName
     */
    public String getPubsubName() {
        return pubsubName;
    }

    /**
     * Setter method for property <tt>pubsubName</tt>.
     *
     * @param pubsubName value to be assigned to property pubsubName
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
     * Getter method for property <tt>data</tt>.
     *
     * @return property value of data
     */
    public byte[] getData() {
        return data;
    }

    /**
     * Setter method for property <tt>data</tt>.
     *
     * @param data value to be assigned to property data
     */
    public void setData(byte[] data) {
        this.data = data;
    }

    /**
     * Getter method for property <tt>contentType</tt>.
     *
     * @return property value of contentType
     */
    public String getContentType() {
        return contentType;
    }

    /**
     * Setter method for property <tt>contentType</tt>.
     *
     * @param contentType value to be assigned to property contentType
     */
    public void setContentType(String contentType) {
        this.contentType = contentType;
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
