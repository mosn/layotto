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

import spec.sdk.runtime.v1.domain.file.GetFileRequest;
import spec.sdk.runtime.v1.domain.file.PutFileRequest;
import spec.sdk.runtime.v1.domain.file.DelFileRequest;
import spec.sdk.runtime.v1.domain.file.ListFileResponse;
import spec.sdk.runtime.v1.domain.file.ListFileRequest;
import spec.sdk.runtime.v1.domain.file.GetMetaRequest;
import spec.sdk.runtime.v1.domain.file.GetMeteResponse;

public interface FileRuntime {

    /**
     * save or update file
     *
     * @param request
     * @param timeoutMs If the time is less than or equal to zero, the method will not wait at all.
     * @throws Exception Instance of RuntimeClientException Or StatusRuntimeException
     */
    void putFile(PutFileRequest request, int timeoutMs) throws Exception;

    /**
     * @param request
     * @param timeoutMs If the time is less than or equal to zero, the method will not wait at all.
     * @throws Exception Instance of RuntimeClientException Or StatusRuntimeException
     */
    void getFile(GetFileRequest request, int timeoutMs) throws Exception;

    /**
     * @param request
     * @param timeoutMs If the time is less than or equal to zero, the method will not wait at all.
     * @throws Exception Instance of RuntimeClientException Or StatusRuntimeException
     */
    ListFileResponse listFile(ListFileRequest request, int timeoutMs) throws Exception;

    /**
     * @param request
     * @param timeoutMs If the time is less than or equal to zero, the method will not wait at all.
     * @throws Exception Instance of RuntimeClientException Or StatusRuntimeException
     */
    void delFile(DelFileRequest request, int timeoutMs) throws Exception;

    /**
     * @param request
     * @param timeoutMs If the time is less than or equal to zero, the method will not wait at all.
     * @throws Exception Instance of RuntimeClientException Or StatusRuntimeException
     */
    GetMeteResponse getMeta(GetMetaRequest request, int timeoutMs) throws Exception;
}
