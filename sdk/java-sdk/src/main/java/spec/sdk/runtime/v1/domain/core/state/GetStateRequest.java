/*
 * Copyright (c) CloudRuntimes Contributors.
 * Licensed under the MIT License.
 */

package spec.sdk.runtime.v1.domain.core.state;

import java.util.Collections;
import java.util.Map;

/**
 * A request to get a state by key.
 */
public class GetStateRequest {

    private final String storeName;

    private final String key;

    private Map<String, String> metadata;

    private spec.sdk.runtime.v1.domain.core.state.StateOptions stateOptions;

    public GetStateRequest(String storeName, String key) {
        this.storeName = storeName;
        this.key = key;
    }

    public String getStoreName() {
        return storeName;
    }

    public String getKey() {
        return key;
    }

    public spec.sdk.runtime.v1.domain.core.state.StateOptions getStateOptions() {
        return stateOptions;
    }

    public GetStateRequest setStateOptions(StateOptions stateOptions) {
        this.stateOptions = stateOptions;
        return this;
    }

    public Map<String, String> getMetadata() {
        return metadata;
    }

    public GetStateRequest setMetadata(Map<String, String> metadata) {
        this.metadata = metadata == null ? null : Collections.unmodifiableMap(metadata);
        return this;
    }
}
