package io.mosn.layotto.v1.exceptions;

/**
 * A Runtime's specific exception.
 */
public class RuntimeClientException extends RuntimeException {

    /**
     * Runtime's error code for this exception.
     */
    private String errorCode;

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
