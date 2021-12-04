package spec.sdk.runtime.v1.domain.state;

import java.util.Map;

/**
 * This class reprent what a State is.
 *
 * @param <T> The type of the value of the state
 */
public class State<T> {

    /**
     * The key of the state.
     */
    private final String key;

    /**
     * The value of the state.
     */
    private final T value;

    /**
     * The ETag to be used Keep in mind that for some state stores (like redis) only numbers are supported.
     */
    private final String etag;

    /**
     * The metadata which will be passed to state store component.
     */
    private final Map<String, String> metadata;

    /**
     * The options used for saving the state.
     */
    private final StateOptions options;

    /**
     * Create an immutable state reference to be retrieved or deleted. This Constructor CAN be used anytime you need to retrieve or delete a
     * state.
     *
     * @param key - The key of the state
     */
    public State(String key) {
        this.key = key;
        this.value = null;
        this.etag = null;
        this.metadata = null;
        this.options = null;
    }

    /**
     * Create an immutable state reference to be retrieved or deleted. This Constructor CAN be used anytime you need to retrieve or delete a
     * state.
     *
     * @param key     - The key of the state
     * @param etag    - The etag of the state - Keep in mind that for some state stores (like redis) only numbers are supported.
     * @param options - REQUIRED when saving a state.
     */
    public State(String key, String etag, StateOptions options) {
        this.value = null;
        this.key = key;
        this.etag = etag;
        this.metadata = null;
        this.options = options;
    }

    /**
     * Create an immutable state. This Constructor CAN be used anytime you want the state to be saved.
     *
     * @param key     - The key of the state.
     * @param value   - The value of the state.
     * @param etag    - The etag of the state - for some state stores (like redis) only numbers are supported.
     * @param options - REQUIRED when saving a state.
     */
    public State(String key, T value, String etag, StateOptions options) {
        this.value = value;
        this.key = key;
        this.etag = etag;
        this.metadata = null;
        this.options = options;
    }

    /**
     * Create an immutable state. This Constructor CAN be used anytime you want the state to be saved.
     *
     * @param key      - The key of the state.
     * @param value    - The value of the state.
     * @param etag     - The etag of the state - for some state stores (like redis) only numbers are supported.
     * @param metadata - The metadata of the state.
     * @param options  - REQUIRED when saving a state.
     */
    public State(String key, T value, String etag, Map<String, String> metadata, StateOptions options) {
        this.value = value;
        this.key = key;
        this.etag = etag;
        this.metadata = metadata;
        this.options = options;
    }

    /**
     * Create an immutable state. This Constructor CAN be used anytime you want the state to be saved.
     *
     * @param key   - The key of the state.
     * @param value - The value of the state.
     * @param etag  - The etag of the state - some state stores (like redis) only numbers are supported.
     */
    public State(String key, T value, String etag) {
        this.value = value;
        this.key = key;
        this.etag = etag;
        this.metadata = null;
        this.options = null;
    }

    /**
     * Retrieves the Value of the state.
     *
     * @return The value of the state
     */
    public T getValue() {
        return value;
    }

    /**
     * Retrieves the Key of the state.
     *
     * @return The key of the state
     */
    public String getKey() {
        return key;
    }

    /**
     * Retrieve the ETag of this state.
     *
     * @return The etag of the state
     */
    public String getEtag() {
        return etag;
    }

    /**
     * Retrieve the metadata of this state.
     *
     * @return the metadata of this state
     */
    public Map<String, String> getMetadata() {
        return metadata;
    }

    /**
     * Retrieve the Options used for saving the state.
     *
     * @return The options to save the state
     */
    public StateOptions getOptions() {
        return options;
    }
}
