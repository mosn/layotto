package spec.sdk.runtime.v1.domain.state;

import java.util.Map;

/**
 * A request to delete a state by key.
 */
public class DeleteStateRequest {

    private final String stateStoreName;

    private final String key;

    private Map<String, String> metadata;

    private String etag;

    private StateOptions stateOptions;

    public DeleteStateRequest(String storeName, String key) {
        this.stateStoreName = storeName;
        this.key = key;
    }

    public String getStateStoreName() {
        return stateStoreName;
    }

    public String getKey() {
        return key;
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

    /**
     * Getter method for property <tt>etag</tt>.
     *
     * @return property value of etag
     */
    public String getEtag() {
        return etag;
    }

    /**
     * Setter method for property <tt>etag</tt>.
     *
     * @param etag value to be assigned to property etag
     */
    public void setEtag(String etag) {
        this.etag = etag;
    }

    /**
     * Getter method for property <tt>stateOptions</tt>.
     *
     * @return property value of stateOptions
     */
    public StateOptions getStateOptions() {
        return stateOptions;
    }

    /**
     * Setter method for property <tt>stateOptions</tt>.
     *
     * @param stateOptions value to be assigned to property stateOptions
     */
    public void setStateOptions(StateOptions stateOptions) {
        this.stateOptions = stateOptions;
    }
}
