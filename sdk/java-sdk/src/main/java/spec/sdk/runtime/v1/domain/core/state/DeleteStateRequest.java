/*
 * Copyright (c) CloudRuntimes Contributors.
 * Licensed under the MIT License.
 */

package spec.sdk.runtime.v1.domain.core.state;

import java.util.Collections;
import java.util.Map;

/**
 * A request to delete a state by key.
 */
public class DeleteStateRequest {

    private final String stateStoreName;

    private final String key;

    private Map<String, String> metadata;

    private String etag;

    private spec.sdk.runtime.v1.domain.core.state.StateOptions stateOptions;

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

    public String getEtag() {
        return etag;
    }

    public DeleteStateRequest setEtag(String etag) {
        this.etag = etag;
        return this;
    }

    public spec.sdk.runtime.v1.domain.core.state.StateOptions getStateOptions() {
        return stateOptions;
    }

    public DeleteStateRequest setStateOptions(StateOptions stateOptions) {
        this.stateOptions = stateOptions;
        return this;
    }

    public Map<String, String> getMetadata() {
        return metadata;
    }

    public DeleteStateRequest setMetadata(Map<String, String> metadata) {
        this.metadata = metadata == null ? null : Collections.unmodifiableMap(metadata);
        return this;
    }
}
