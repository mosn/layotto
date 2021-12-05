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
import java.util.List;
import java.util.Map;

public class TransactionalStateRequest<T> {

    /**
     * Transactional operations list.
     */
    private final List<TransactionalStateOperation<T>> operations;

    /**
     * Metadata used for transactional operations.
     */
    private final Map<String, String>                  metadata;

    /**
     * Constructor to create immutable transactional state request object.
     *
     * @param operations List of operations to be performed.
     * @param metadata   Metadata used for transactional operations.
     */
    public TransactionalStateRequest(List<TransactionalStateOperation<T>> operations, Map<String, String> metadata) {
        this.operations = operations;
        this.metadata = metadata;
    }

    public List<TransactionalStateOperation<T>> getOperations() {
        return Collections.unmodifiableList(operations);
    }

    public Map<String, String> getMetadata() {
        return metadata;
    }
}
