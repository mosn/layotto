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

import java.util.Objects;

public class TransactionalStateOperation<T> {

    /**
     * The type of operation to be executed.
     */
    private final OperationType operation;

    /**
     * State values to be operated on.
     */
    private final State<T>      request;

    /**
     * Construct an immutable transactional state operation object.
     *
     * @param operationType The type of operation done.
     * @param state         The required state.
     */
    public TransactionalStateOperation(OperationType operationType, State<T> state) {
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
