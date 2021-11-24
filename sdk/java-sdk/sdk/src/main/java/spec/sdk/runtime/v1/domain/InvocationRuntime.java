/*
 * Copyright 2021 Layotto Authors
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
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
    InvokeResponse<byte[]> invokeMethod(String appId, String methodName, byte[] data, Map<String, String> header,
                                        int timeoutMs);
}
