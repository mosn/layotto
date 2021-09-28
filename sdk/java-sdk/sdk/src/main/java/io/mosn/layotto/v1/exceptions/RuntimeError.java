package io.mosn.layotto.v1.exceptions;

public enum RuntimeError {
    ;

    /**
     * Error code.
     */
    private String errorCode;

    /**
     * Error Message.
     */
    private String message;

    /**
     * Getter method for property <tt>errorCode</tt>.
     *
     * @return property value of errorCode
     */
    public String getErrorCode() {
        return errorCode;
    }

    /**
     * Getter method for property <tt>message</tt>.
     *
     * @return property value of message
     */
    public String getMessage() {
        return message;
    }
}
