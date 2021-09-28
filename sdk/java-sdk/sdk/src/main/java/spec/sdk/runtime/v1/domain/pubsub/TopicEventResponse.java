package spec.sdk.runtime.v1.domain.pubsub;

public class TopicEventResponse {
    private TopicEventResponseStatus status;

    /**
     * Getter method for property <tt>status</tt>.
     *
     * @return property value of status
     */
    public TopicEventResponseStatus getStatus() {
        return status;
    }

    /**
     * Setter method for property <tt>status</tt>.
     *
     * @param status value to be assigned to property status
     */
    public void setStatus(TopicEventResponseStatus status) {
        this.status = status;
    }
}
