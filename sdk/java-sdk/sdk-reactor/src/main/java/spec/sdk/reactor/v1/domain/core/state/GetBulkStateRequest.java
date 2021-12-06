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

import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import java.util.Map;

/**
 * A request to get bulk state by keys.
 */
public class GetBulkStateRequest {

    private final String        storeName;

    private final List<String>  keys;

    private Map<String, String> metadata;

    private int                 parallelism = 1;

    public GetBulkStateRequest(String storeName, List<String> keys) {
        this.storeName = storeName;
        this.keys = keys == null ? null : Collections.unmodifiableList(keys);
    }

    public GetBulkStateRequest(String storeName, String... keys) {
        this.storeName = storeName;
        this.keys = keys == null ? null : Collections.unmodifiableList(Arrays.asList(keys));
    }

    public String getStoreName() {
        return storeName;
    }

    public List<String> getKeys() {
        return keys;
    }

    public int getParallelism() {
        return parallelism;
    }

    public GetBulkStateRequest setParallelism(int parallelism) {
        this.parallelism = parallelism;
        return this;
    }

    public Map<String, String> getMetadata() {
        return metadata;
    }

    public GetBulkStateRequest setMetadata(Map<String, String> metadata) {
        this.metadata = metadata == null ? null : Collections.unmodifiableMap(metadata);
        return this;
    }
}
