/*
 * Copyright (c) CloudRuntimes Contributors.
 * Licensed under the MIT License.
 */

package spec.sdk.runtime.v1.domain.core.binding;

import java.util.Collections;
import java.util.Map;

/**
 * A request to invoke binding.
 */
public class InvokeBindingRequest {

    private final String name;

    private final String operation;

    private Object data;

    private Map<String, String> metadata;

    public InvokeBindingRequest(String bindingName, String operation) {
        this.name = bindingName;
        this.operation = operation;
    }

    public String getName() {
        return name;
    }

    public String getOperation() {
        return operation;
    }

    public Object getData() {
        return data;
    }

    public InvokeBindingRequest setData(Object data) {
        this.data = data;
        return this;
    }

    public Map<String, String> getMetadata() {
        return metadata;
    }

    public InvokeBindingRequest setMetadata(Map<String, String> metadata) {
        this.metadata = metadata == null ? null : Collections.unmodifiableMap(metadata);
        return this;
    }
}
