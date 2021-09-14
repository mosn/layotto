/*
 * Copyright (c) CloudRuntimes Contributors.
 * Licensed under the MIT License.
 */

package spec.sdk.reactor.v1.domain.core.state;

import java.util.Map;

/**
 * This class reprent what a State is.
 *
 * @param <T> The type of the value of the sate
 */
public class State<T> {

    /**
     * The value of the state.
     */
    private final T value;

    /**
     * The key of the state.
     */
    private final String key;

    /**
     * The ETag to be used
     * Keep in mind that for some state stores (like redis) only numbers are supported.
     */
    private final String etag;

    /**
     * The metadata which will be passed to state store component.
     */
    private final Map<String, String> metadata;

    /**
     * The error in case the key could not be retrieved.
     */
    private final String error;

    /**
     * The options used for saving the state.
     */
    private final StateOptions options;

    /**
     * Create an immutable state reference to be retrieved or deleted.
     * This Constructor CAN be used anytime you need to retrieve or delete a state.
     *
     * @param key - The key of the state
     */
    public State(String key) {
        this.key = key;
        this.value = null;
        this.etag = null;
        this.metadata = null;
        this.options = null;
        this.error = null;
    }

    /**
     * Create an immutable state reference to be retrieved or deleted.
     * This Constructor CAN be used anytime you need to retrieve or delete a state.
     *
     * @param key     - The key of the state
     * @param etag    - The etag of the state - Keep in mind that for some state stores (like redis) only numbers
     *                are supported.
     * @param options - REQUIRED when saving a state.
     */
    public State(String key, String etag, StateOptions options) {
        this.value = null;
        this.key = key;
        this.etag = etag;
        this.metadata = null;
        this.options = options;
        this.error = null;
    }

    /**
     * Create an immutable state.
     * This Constructor CAN be used anytime you want the state to be saved.
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
        this.error = null;
    }

    /**
     * Create an immutable state.
     * This Constructor CAN be used anytime you want the state to be saved.
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
        this.error = null;
    }

    /**
     * Create an immutable state.
     * This Constructor CAN be used anytime you want the state to be saved.
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
        this.error = null;
    }

    /**
     * Create an immutable state.
     * This Constructor MUST be used anytime the key could not be retrieved and contains an error.
     *
     * @param key   - The key of the state.
     * @param error - Error when fetching the state.
     */
    public State(String key, String error) {
        this.value = null;
        this.key = key;
        this.etag = null;
        this.metadata = null;
        this.options = null;
        this.error = error;
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
     * Retrieve the error for this state.
     *
     * @return The error for this state.
     */

    public String getError() {
        return error;
    }

    /**
     * Retrieve the Options used for saving the state.
     *
     * @return The options to save the state
     */
    public StateOptions getOptions() {
        return options;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }

        if (!(o instanceof State)) {
            return false;
        }

        State<?> that = (State<?>) o;

        if (getValue() != null ? !getValue().equals(that.getValue()) : that.getValue() != null) {
            return false;
        }

        if (getKey() != null ? !getKey().equals(that.getKey()) : that.getKey() != null) {
            return false;
        }

        if (getEtag() != null ? !getEtag().equals(that.getEtag()) : that.getEtag() != null) {
            return false;
        }

        if (getError() != null ? !getError().equals(that.getError()) : that.getError() != null) {
            return false;
        }

        if (getMetadata() != null ? !getMetadata().equals(that.getMetadata()) : that.getMetadata() != null) {
            return false;
        }

        if (getOptions() != null ? !getOptions().equals(that.getOptions()) : that.getOptions() != null) {
            return false;
        }

        return true;
    }

    @Override
    public int hashCode() {
        int result = getValue() != null ? getValue().hashCode() : 0;
        result = 31 * result + (getKey() != null ? getKey().hashCode() : 0);
        result = 31 * result + (getEtag() != null ? getEtag().hashCode() : 0);
        result = 31 * result + (getMetadata() != null ? getMetadata().hashCode() : 0);
        result = 31 * result + (getError() != null ? getError().hashCode() : 0);
        result = 31 * result + (getOptions() != null ? options.hashCode() : 0);
        return result;
    }

    @Override
    public String toString() {
        return "StateKeyValue{"
                + "key='" + key + "'"
                + ", value=" + value
                + ", etag='" + etag + "'"
                + ", metadata={'" + (metadata != null ? metadata.toString() : null) + "'}"
                + ", error='" + error + "'"
                + ", options={'" + (options != null ? options.toString() : null) + "'}"
                + "}";
    }
}
