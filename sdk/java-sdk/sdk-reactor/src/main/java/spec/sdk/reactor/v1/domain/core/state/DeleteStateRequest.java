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
package spec.sdk.reactor.v1.domain.core.state;

import java.util.Collections;
import java.util.Map;

/**
 * A request to delete a state by key.
 */
public class DeleteStateRequest {

    private final String        stateStoreName;

    private final String        key;

    private Map<String, String> metadata;

    private String              etag;

    private StateOptions        stateOptions;

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

    public StateOptions getStateOptions() {
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
