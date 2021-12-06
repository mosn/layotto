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
package io.mosn.layotto.v1.domain;

import com.fasterxml.jackson.annotation.JsonInclude;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.ObjectMapper;

import java.io.IOException;
import java.util.Arrays;
import java.util.Objects;

/**
 * A cloud event in Layotto.
 *
 * @param <T> The type of the payload.
 */
public final class CloudEvent<T> {

    /**
     * Mime type used for CloudEvent.
     */
    public static final String          CONTENT_TYPE  = "application/cloudevents+json";

    /**
     * Shared Json serializer/deserializer as per Jackson's documentation.
     */
    protected static final ObjectMapper OBJECT_MAPPER = new ObjectMapper()
                                                          .configure(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES,
                                                              false)
                                                          .setSerializationInclusion(JsonInclude.Include.NON_NULL);

    /**
     * Identifier of the message being processed.
     */
    private String                      id;

    /**
     * Event's source.
     */
    private String                      source;

    /**
     * Envelope type.
     */
    private String                      type;

    /**
     * Version of the specification.
     */
    private String                      specversion;

    /**
     * Type of the data's content.
     */
    private String                      datacontenttype;

    /**
     * Cloud event specs says data can be a JSON object or string.
     */
    private T                           data;

    /**
     * Cloud event specs says binary data should be in data_base64.
     */
    @JsonProperty("data_base64")
    private byte[]                      binaryData;

    /**
     * Instantiates a CloudEvent.
     */
    public CloudEvent() {
    }

    /**
     * Instantiates a CloudEvent.
     *
     * @param id              Identifier of the message being processed.
     * @param source          Source for this event.
     * @param type            Type of event.
     * @param specversion     Version of the event spec.
     * @param datacontenttype Type of the payload.
     * @param data            Payload.
     */
    public CloudEvent(
                      String id,
                      String source,
                      String type,
                      String specversion,
                      String datacontenttype,
                      T data) {
        this.id = id;
        this.source = source;
        this.type = type;
        this.specversion = specversion;
        this.datacontenttype = datacontenttype;
        this.data = data;
    }

    /**
     * Instantiates a CloudEvent.
     *
     * @param id          Identifier of the message being processed.
     * @param source      Source for this event.
     * @param type        Type of event.
     * @param specversion Version of the event spec.
     * @param binaryData  Payload.
     */
    public CloudEvent(
                      String id,
                      String source,
                      String type,
                      String specversion,
                      byte[] binaryData) {
        this.id = id;
        this.source = source;
        this.type = type;
        this.specversion = specversion;
        this.datacontenttype = "application/octet-stream";
        this.binaryData = binaryData == null ? null : Arrays.copyOf(binaryData, binaryData.length);
        ;
    }

    /**
     * Deserialize a message topic from Layotto.
     *
     * @param payload Payload sent from Layotto.
     * @return Message (can be null if input is null)
     * @throws IOException If cannot parse.
     */
    public static CloudEvent<?> deserialize(byte[] payload) throws IOException {
        if (payload == null) {
            return null;
        }

        return OBJECT_MAPPER.readValue(payload, CloudEvent.class);
    }

    /**
     * Gets the identifier of the message being processed.
     *
     * @return Identifier of the message being processed.
     */
    public String getId() {
        return id;
    }

    /**
     * Sets the identifier of the message being processed.
     *
     * @param id Identifier of the message being processed.
     */
    public void setId(String id) {
        this.id = id;
    }

    /**
     * Gets the event's source.
     *
     * @return Event's source.
     */
    public String getSource() {
        return source;
    }

    /**
     * Sets the event's source.
     *
     * @param source Event's source.
     */
    public void setSource(String source) {
        this.source = source;
    }

    /**
     * Gets the envelope type.
     *
     * @return Envelope type.
     */
    public String getType() {
        return type;
    }

    /**
     * Sets the envelope type.
     *
     * @param type Envelope type.
     */
    public void setType(String type) {
        this.type = type;
    }

    /**
     * Gets the version of the specification.
     *
     * @return Version of the specification.
     */
    public String getSpecversion() {
        return specversion;
    }

    /**
     * Sets the version of the specification.
     *
     * @param specversion Version of the specification.
     */
    public void setSpecversion(String specversion) {
        this.specversion = specversion;
    }

    /**
     * Gets the type of the data's content.
     *
     * @return Type of the data's content.
     */
    public String getDatacontenttype() {
        return datacontenttype;
    }

    /**
     * Sets the type of the data's content.
     *
     * @param datacontenttype Type of the data's content.
     */
    public void setDatacontenttype(String datacontenttype) {
        this.datacontenttype = datacontenttype;
    }

    /**
     * Gets the cloud event data.
     *
     * @return Cloud event's data. As per specs, data can be a JSON object or string.
     */
    public T getData() {
        return data;
    }

    /**
     * Sets the cloud event data. As per specs, data can be a JSON object or string.
     *
     * @param data Cloud event's data. As per specs, data can be a JSON object or string.
     */
    public void setData(T data) {
        this.data = data;
    }

    /**
     * Gets the cloud event's binary data.
     *
     * @return Cloud event's binary data.
     */
    public byte[] getBinaryData() {
        return this.binaryData == null ? null : Arrays.copyOf(this.binaryData, this.binaryData.length);
    }

    /**
     * Sets the cloud event's binary data.
     *
     * @param binaryData Cloud event's binary data.
     */
    public void setBinaryData(byte[] binaryData) {
        this.binaryData = binaryData == null ? null : Arrays.copyOf(binaryData, binaryData.length);
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (o == null || getClass() != o.getClass()) {
            return false;
        }
        CloudEvent<?> that = (CloudEvent<?>) o;
        return Objects.equals(id, that.id)
            && Objects.equals(source, that.source)
            && Objects.equals(type, that.type)
            && Objects.equals(specversion, that.specversion)
            && Objects.equals(datacontenttype, that.datacontenttype)
            && Objects.equals(data, that.data)
            && Arrays.equals(binaryData, that.binaryData);
    }

    /**
     * {@inheritDoc}
     */
    @Override
    public int hashCode() {
        return Objects.hash(id, source, type, specversion, datacontenttype, data, binaryData);
    }
}
