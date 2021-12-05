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
package io.mosn.layotto.v1.exceptions;

import com.fasterxml.jackson.annotation.JsonAutoDetect;
import io.grpc.Status;

/**
 * Represents an error message from Layotto.
 */
@JsonAutoDetect(fieldVisibility = JsonAutoDetect.Visibility.ANY)
public class LayottoError {

    /**
     * Error code.
     */
    private String  errorCode;

    /**
     * Error Message.
     */
    private String  message;

    /**
     * Error code from gRPC.
     */
    private Integer code;

    /**
     * Gets the error code.
     *
     * @return Error code.
     */
    public String getErrorCode() {
        if ((errorCode == null) && (code != null)) {
            return Status.fromCodeValue(code).getCode().name();
        }
        return errorCode;
    }

    /**
     * Sets the error code.
     *
     * @param errorCode Error code.
     * @return This instance.
     */
    public LayottoError setErrorCode(String errorCode) {
        this.errorCode = errorCode;
        return this;
    }

    /**
     * Gets the error message.
     *
     * @return Error message.
     */
    public String getMessage() {
        return message;
    }

    /**
     * Sets the error message.
     *
     * @param message Error message.
     * @return This instance.
     */
    public LayottoError setMessage(String message) {
        this.message = message;
        return this;
    }
}
