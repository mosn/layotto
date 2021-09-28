package spec.sdk.runtime.v1.domain.state;

import java.util.List;
import java.util.Map;

public class ExecuteStateTransactionRequest {

    /**
     * Name of the state store.
     */
    private final String stateStoreName;

    /**
     * Transactional operations list.
     */
    private List<TransactionalStateOperation<?>> operations;

    /**
     * Metadata used for transactional operations.
     */
    private Map<String, String> metadata;

    public ExecuteStateTransactionRequest(String stateStoreName) {
        this.stateStoreName = stateStoreName;
    }

    public String getStateStoreName() {
        return stateStoreName;
    }

    /**
     * Getter method for property <tt>operations</tt>.
     *
     * @return property value of operations
     */
    public List<TransactionalStateOperation<?>> getOperations() {
        return operations;
    }

    /**
     * Setter method for property <tt>operations</tt>.
     *
     * @param operations value to be assigned to property operations
     */
    public void setOperations(List<TransactionalStateOperation<?>> operations) {
        this.operations = operations;
    }

    /**
     * Getter method for property <tt>metadata</tt>.
     *
     * @return property value of metadata
     */
    public Map<String, String> getMetadata() {
        return metadata;
    }

    /**
     * Setter method for property <tt>metadata</tt>.
     *
     * @param metadata value to be assigned to property metadata
     */
    public void setMetadata(Map<String, String> metadata) {
        this.metadata = metadata;
    }
}
