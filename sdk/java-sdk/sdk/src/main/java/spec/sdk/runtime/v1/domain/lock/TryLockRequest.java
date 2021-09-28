package spec.sdk.runtime.v1.domain.lock;

public class  TryLockRequest {
    public String storeName;
    public String resourceId;
    public String lockOwner;
    public int expire;
}