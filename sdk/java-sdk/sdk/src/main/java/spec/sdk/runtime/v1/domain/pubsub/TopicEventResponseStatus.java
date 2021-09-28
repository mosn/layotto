package spec.sdk.runtime.v1.domain.pubsub;

// TopicEventResponseStatus allows apps to have finer control over handling of the message.
public enum TopicEventResponseStatus {
    // SUCCESS is the default behavior: message is acknowledged and not retried or logged.
    SUCCESS(0),
    // RETRY status signals runtime to retry the message as part of an expected scenario (no warning is logged).
    RETRY(1),
    // DROP status signals runtime to drop the message as part of an unexpected scenario (warning is logged).
    DROP(2);

    int idx;

    TopicEventResponseStatus(int idx) {
        this.idx = idx;
    }

    /**
     * Getter method for property <tt>idx</tt>.
     *
     * @return property value of idx
     */
    public int getIdx() {
        return idx;
    }
}
