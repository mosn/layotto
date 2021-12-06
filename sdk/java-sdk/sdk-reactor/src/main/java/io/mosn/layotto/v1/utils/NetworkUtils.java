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
package io.mosn.layotto.v1.utils;

import java.io.IOException;
import java.net.InetSocketAddress;
import java.net.Socket;

/**
 * Utility methods for network, internal to Dapr SDK.
 */
public final class NetworkUtils {

    private NetworkUtils() {
    }

    /**
     * Tries to connect to a socket, retrying every 1 second.
     *
     * @param host                  Host to connect to.
     * @param port                  Port to connect to.
     * @param timeoutInMilliseconds Timeout in milliseconds to give up trying.
     * @throws InterruptedException If retry is interrupted.
     */
    public static void waitForSocket(String host, int port, int timeoutInMilliseconds) throws InterruptedException {
        long started = System.currentTimeMillis();
        Retry.callWithRetry(() -> {
            try {
                try (Socket socket = new Socket()) {
                    // timeout cannot be negative.
                    // zero timeout means infinite, so 1 is the practical minimum.
                    int remainingTimeout = (int) Math.max(1, timeoutInMilliseconds - (System.currentTimeMillis() - started));
                    socket.connect(new InetSocketAddress(host, port), remainingTimeout);
                }
            } catch (IOException e) {
                throw new RuntimeException(e);
            }
        }, timeoutInMilliseconds);
    }
}

class Retry {

    private static final long RETRY_WAIT_MILLISECONDS = 1000;

    private Retry() {
    }

    static void callWithRetry(Runnable function, long retryTimeoutMilliseconds) throws InterruptedException {
        long started = System.currentTimeMillis();
        while (true) {
            Throwable exception;
            try {
                function.run();
                return;
            } catch (Exception e) {
                exception = e;
            } catch (AssertionError e) {
                exception = e;
            }

            long elapsed = System.currentTimeMillis() - started;
            if (elapsed >= retryTimeoutMilliseconds) {
                if (exception instanceof RuntimeException) {
                    throw (RuntimeException) exception;
                }

                throw new RuntimeException(exception);
            }

            long remaining = retryTimeoutMilliseconds - elapsed;
            Thread.sleep(Math.min(remaining, RETRY_WAIT_MILLISECONDS));
        }
    }
}
