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
 *
 */
package io.mosn.layotto.v1.exceptions;

/**
 * A Runtime's specific exception.
 */
public class RuntimeClientException extends RuntimeException {

    /**
     * Runtime's error code for this exception.
     */
    private final String errorCode;

    /**
     * New exception from a server-side generated error code and message.
     *
     * @param runtimeError Server-side error.
     */
    public RuntimeClientException(RuntimeError runtimeError) {
        this(runtimeError.getErrorCode(), runtimeError.getMessage());
    }

    /**
     * New exception from a server-side generated error code and message.
     *
     * @param runtimeError Client-side error.
     * @param cause        the cause (which is saved for later retrieval by the {@link #getCause()} method).  (A {@code null} value is
     *                     permitted, and indicates that the cause is nonexistent or unknown.)
     */
    public RuntimeClientException(RuntimeError runtimeError, Throwable cause) {
        this(runtimeError.getErrorCode(), runtimeError.getMessage(), cause);
    }

    /**
     * Wraps an exception into a RuntimeException.
     *
     * @param exception the exception to be wrapped.
     */
    public RuntimeClientException(Throwable exception) {
        this("UNKNOWN", exception.getMessage(), exception);
    }

    /**
     * New Exception from a client-side generated error code and message.
     *
     * @param errorCode Client-side error code.
     * @param message   Client-side error message.
     */
    public RuntimeClientException(String errorCode, String message) {
        super(String.format("%s: %s", errorCode, message));
        this.errorCode = errorCode;
    }

    /**
     * New exception from a server-side generated error code and message.
     *
     * @param errorCode Client-side error code.
     * @param message   Client-side error message.
     * @param cause     the cause (which is saved for later retrieval by the {@link #getCause()} method).  (A {@code null} value is
     *                  permitted, and indicates that the cause is nonexistent or unknown.)
     */
    public RuntimeClientException(String errorCode, String message, Throwable cause) {
        super(String.format("%s: %s", errorCode, emptyIfNull(message)), cause);
        this.errorCode = errorCode;
    }

    /**
     * Returns the exception's error code.
     *
     * @return Error code.
     */
    public String getErrorCode() {
        return this.errorCode;
    }

    private static String emptyIfNull(String str) {
        if (str == null) {
            return "";
        }
        return str;
    }
}
