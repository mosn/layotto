package spec.sdk.runtime.v1.domain.pubsub;

import java.util.Map;

public class TopicEventRequest {
    // id identifies the event. Producers MUST ensure that source + id
    // is unique for each distinct event. If a duplicate event is re-sent
    // (e.g. due to a network error) it MAY have the same id.
    private String id ;

    // source identifies the context in which an event happened.
    // Often this will include information such as the type of the
    // event source, the organization publishing the event or the process
    // that produced the event. The exact syntax and semantics behind
    // the data encoded in the URI is defined by the event producer.
    private String source ;

    // The type of event related to the originating occurrence.
    private String type ;

    // The version of the CloudEvents specification.
    private String specVersion ;

    // The content type of data value.
    private String contentType ;

    // The content of the event.
    private byte[] data ;

    // The pubsub topic which publisher sent to.
    private String topic ;

    // The name of the pubsub the publisher sent to.
    private String pubsubName;

    // add a map to pass some extra properties.
    private Map<String,String> metadata ;

    /**
     * Getter method for property <tt>id</tt>.
     *
     * @return property value of id
     */
    public String getId() {
        return id;
    }

    /**
     * Setter method for property <tt>id</tt>.
     *
     * @param id value to be assigned to property id
     */
    public void setId(String id) {
        this.id = id;
    }

    /**
     * Getter method for property <tt>source</tt>.
     *
     * @return property value of source
     */
    public String getSource() {
        return source;
    }

    /**
     * Setter method for property <tt>source</tt>.
     *
     * @param source value to be assigned to property source
     */
    public void setSource(String source) {
        this.source = source;
    }

    /**
     * Getter method for property <tt>type</tt>.
     *
     * @return property value of type
     */
    public String getType() {
        return type;
    }

    /**
     * Setter method for property <tt>type</tt>.
     *
     * @param type value to be assigned to property type
     */
    public void setType(String type) {
        this.type = type;
    }

    /**
     * Getter method for property <tt>specVersion</tt>.
     *
     * @return property value of specVersion
     */
    public String getSpecVersion() {
        return specVersion;
    }

    /**
     * Setter method for property <tt>specVersion</tt>.
     *
     * @param specVersion value to be assigned to property specVersion
     */
    public void setSpecVersion(String specVersion) {
        this.specVersion = specVersion;
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
