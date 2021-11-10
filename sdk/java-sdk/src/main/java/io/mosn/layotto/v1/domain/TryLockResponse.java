package io.mosn.layotto.v1.domain;

public class  TryLockResponse {

    public boolean success;

    public boolean isSuccess() {
        return success;
    }

    public void setSuccess(boolean success) {
        this.success = success;
    }
}