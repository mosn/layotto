package io.mosn.layotto.v1.domain;

public class  TryLockRequest {
    public String storeName;
    public String resourceId;
    public String lockOwner;
    public int expire;
}