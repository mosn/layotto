package spec.sdk.runtime.v1.client;

import spec.sdk.runtime.v1.domain.*;

public interface RuntimesClient extends
        HelloRuntimes,
        ConfigurationRuntimes,
        InvocationRuntimes,
        PubSubRuntimes,
        StateRuntimes {
}