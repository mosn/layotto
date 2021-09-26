package io.mosn.layotto.v1.domain;

public class UnlockResponse {

    public UnlockResponseStatus status;

    public UnlockResponseStatus getStatus() {
        return status;
    }

    public void setStatus(UnlockResponseStatus status) {
        this.status = status;
    }
}

enum UnlockResponseStatus {

    SUCCESS(0),
    LOCK_UNEXIST(1),
    LOCK_BELONG_TO_OTHERS(2),
    INTERNAL_ERROR(3);

    private final int value;

    UnlockResponseStatus(int value) {
        this.value = value;
    }
}
