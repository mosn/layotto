/*
 * Copyright (c) Microsoft Corporation and Layotto Contributors.
 * Licensed under the MIT License.
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
    private String errorCode;

    /**
     * Error Message.
     */
    private String message;

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
