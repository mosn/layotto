package spec.sdk.runtime.v1.client;

/**
 * Cloud Runtimes Client Provider.
 */
public interface CloudRuntimesClientProvider {

    /**
     * Provide cloud runtimes client.
     *
     * @return the cloud runtimes client
     */
    CloudRuntimesClient provide();
}
