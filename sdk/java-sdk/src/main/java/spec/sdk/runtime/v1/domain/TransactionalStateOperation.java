package spec.sdk.runtime.v1.domain;

public class TransactionalStateOperation<T> {

    /**
     * The type of operation to be executed.
     */
    private final OperationType operation;

    /**
     * State values to be operated on.
     */
    private final State<T> request;

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

    public enum OperationType {
        UPSERT,
        DELETE
    }
}
