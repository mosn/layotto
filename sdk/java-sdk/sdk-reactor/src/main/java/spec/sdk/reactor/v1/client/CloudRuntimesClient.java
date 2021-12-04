package spec.sdk.reactor.v1.client;

import spec.sdk.reactor.v1.CoreCloudRuntimes;
import spec.sdk.reactor.v1.EnhancedCloudRuntimes;
import reactor.core.publisher.Mono;

/**
 * Cloud Runtimes JVM Client.
 */
public interface CloudRuntimesClient extends AutoCloseable, CoreCloudRuntimes, EnhancedCloudRuntimes {

    /**
     * Waits for the sidecar, giving up after timeout.
     *
     * @param timeoutInMilliseconds Timeout in milliseconds to wait for sidecar.
     * @return a Mono plan of type Void.
     */
    Mono<Void> waitForSidecar(int timeoutInMilliseconds);

    /**
     * Gracefully shutdown the cloud-runtimes runtime.
     *
     * @return a Mono plan of type Void.
     */
    Mono<Void> shutdown();
}
