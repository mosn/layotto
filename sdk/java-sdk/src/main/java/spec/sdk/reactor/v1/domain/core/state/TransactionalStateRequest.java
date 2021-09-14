/*
 * Copyright (c) CloudRuntimes Contributors.
 * Licensed under the MIT License.
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
    private final Map<String, String> metadata;

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
