/*
 * Copyright (c) CloudRuntimes Contributors.
 * Licensed under the MIT License.
 */

package spec.sdk.runtime.v1.domain.core.state;


import java.util.Objects;

public class TransactionalStateOperation<T> {

    /**
     * The type of operation to be executed.
     */
    private final OperationType operation;

    /**
     * State values to be operated on.
     */
    private final spec.sdk.runtime.v1.domain.core.state.State<T> request;

    /**
     * Construct an immutable transactional state operation object.
     *
     * @param operationType The type of operation done.
     * @param state         The required state.
     */
    public TransactionalStateOperation(OperationType operationType, spec.sdk.runtime.v1.domain.core.state.State<T> state) {
        this.operation = operationType;
        this.request = state;
    }

    public OperationType getOperation() {
        return operation;
    }

    public State<T> getRequest() {
        return request;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) {
            return true;
        }
        if (o == null || getClass() != o.getClass()) {
            return false;
        }
        TransactionalStateOperation<?> that = (TransactionalStateOperation<?>) o;
        return operation.equals(that.operation)
                && request.equals(that.request);
    }

    @Override
    public int hashCode() {
        return Objects.hash(operation, request);
    }

    @Override
    public String toString() {
        return "TransactionalStateOperation{"
                + "operationType='" + operation + '\''
                + ", state=" + request
                + '}';
    }

    public enum OperationType {
        UPSERT,
        DELETE
    }
}
