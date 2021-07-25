package io.mosn.layotto.v1.domain;

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

public class UnlockResponse {
    public UnlockResponseStatus status;
}
