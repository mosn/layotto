package spec.sdk.runtime.v1;

import spec.sdk.runtime.v1.domain.core.*;

/**
 * Core Cloud Runtimes standard API defined.
 */
public interface CoreCloudRuntimes extends
        InvocationRuntimes,
        PubSubRuntimes,
        BindingRuntimes,
        StateRuntimes,
        //SecretsRuntimes,
        ConfigurationRuntimes {
}
