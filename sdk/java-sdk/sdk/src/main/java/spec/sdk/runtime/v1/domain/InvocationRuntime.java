package spec.sdk.runtime.v1.domain;

import spec.sdk.runtime.v1.domain.invocation.InvokeResponse;

import java.util.Map;

public interface InvocationRuntime {

    InvokeResponse<byte[]> invokeMethod(String appId, String methodName, byte[] data, Map<String, String> header);

    /**
     * Invoke a service method.
     *
     * @param appId
     * @param methodName
     * @param data
     * @param header
     * @param timeoutMs  can be customized every time a service method is called, since different services provide different SLA.
     * @return
     */
    InvokeResponse<byte[]> invokeMethod(String appId, String methodName, byte[] data, Map<String, String> header, int timeoutMs);
}
