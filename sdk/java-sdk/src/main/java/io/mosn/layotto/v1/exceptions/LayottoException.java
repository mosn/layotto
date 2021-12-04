/*
 * Copyright (c) Microsoft Corporation and Layotto Contributors.
 * Licensed under the MIT License.
 */

package io.mosn.layotto.v1.exceptions;

import io.grpc.StatusRuntimeException;
import reactor.core.Exceptions;
import reactor.core.publisher.Mono;

import java.util.concurrent.Callable;

/**
 * A Layotto's specific exception.
 */
public class LayottoException extends RuntimeException {

    /**
     * Layotto's error code for this exception.
     */
    private final String errorCode;

    /**
     * New exception from a server-side generated error code and message.
     *
     * @param daprError Server-side error.
     */
    public LayottoException(LayottoError daprError) {
        this(daprError.getErrorCode(), daprError.getMessage());
    }

    /**
     * New exception from a server-side generated error code and message.
     *
     * @param daprError Client-side error.
     * @param cause     the cause (which is saved for later retrieval by the
     *                  {@link #getCause()} method).  (A {@code null} value is
     *                  permitted, and indicates that the cause is nonexistent or
     *                  unknown.)
     */
    public LayottoException(LayottoError daprError, Throwable cause) {
        this(daprError.getErrorCode(), daprError.getMessage(), cause);
    }

    /**
     * Wraps an exception into a LayottoException.
     *
     * @param exception the exception to be wrapped.
     */
    public LayottoException(Throwable exception) {
        this("UNKNOWN", exception.getMessage(), exception);
    }

    /**
     * New Exception from a client-side generated error code and message.
     *
     * @param errorCode Client-side error code.
     * @param message   Client-side error message.
     */
    public LayottoException(String errorCode, String message) {
        super(String.format("%s: %s", errorCode, message));
        this.errorCode = errorCode;
    }

    /**
     * New exception from a server-side generated error code and message.
     *
     * @param errorCode Client-side error code.
     * @param message   Client-side error message.
     * @param cause     the cause (which is saved for later retrieval by the
     *                  {@link #getCause()} method).  (A {@code null} value is
     *                  permitted, and indicates that the cause is nonexistent or
     *                  unknown.)
     */
    public LayottoException(String errorCode, String message, Throwable cause) {
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

    /**
     * Wraps an exception into LayottoException (if not already LayottoException).
     *
     * @param exception Exception to be wrapped.
     */
    public static void wrap(Throwable exception) {
        if (exception == null) {
            return;
        }

        throw propagate(exception);
    }

    /**
     * Wraps a callable with a try-catch to throw LayottoException.
     *
     * @param callable callable to be invoked.
     * @param <T>      type to be returned
     * @return object of type T.
     */
    public static <T> Callable<T> wrap(Callable<T> callable) {
        return () -> {
            try {
                return callable.call();
            } catch (Exception e) {
                wrap(e);
                return null;
            }
        };
    }

    /**
     * Wraps a runnable with a try-catch to throw LayottoException.
     *
     * @param runnable runnable to be invoked.
     * @return object of type T.
     */
    public static Runnable wrap(Runnable runnable) {
        return () -> {
            try {
                runnable.run();
            } catch (Exception e) {
                wrap(e);
            }
        };
    }

    /**
     * Wraps an exception into LayottoException (if not already LayottoException).
     *
     * @param exception Exception to be wrapped.
     * @param <T>       Mono's response type.
     * @return Mono containing LayottoException.
     */
    public static <T> Mono<T> wrapMono(Exception exception) {
        try {
            wrap(exception);
        } catch (Exception e) {
            return Mono.error(e);
        }

        return Mono.empty();
    }

    /**
     * Wraps an exception into LayottoException (if not already LayottoException).
     *
     * @param exception Exception to be wrapped.
     * @return wrapped RuntimeException
     */
    public static RuntimeException propagate(Throwable exception) {
        Exceptions.throwIfFatal(exception);

        if (exception instanceof LayottoException) {
            return (LayottoException) exception;
        }

        Throwable e = exception;
        while (e != null) {
            if (e instanceof StatusRuntimeException) {
                StatusRuntimeException statusRuntimeException = (StatusRuntimeException) e;
                return new LayottoException(
                        statusRuntimeException.getStatus().getCode().toString(),
                        statusRuntimeException.getStatus().getDescription(),
                        exception);
            }

            e = e.getCause();
        }

        if (exception instanceof IllegalArgumentException) {
            return (IllegalArgumentException) exception;
        }

        return new LayottoException(exception);
    }

    private static String emptyIfNull(String str) {
        if (str == null) {
            return "";
        }

        return str;
    }
}
