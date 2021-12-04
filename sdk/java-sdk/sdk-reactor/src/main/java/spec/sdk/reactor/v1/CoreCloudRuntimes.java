package spec.sdk.reactor.v1;

import spec.sdk.reactor.v1.domain.core.ConfigurationRuntimes;
import spec.sdk.reactor.v1.domain.core.InvocationRuntimes;
import spec.sdk.reactor.v1.domain.core.PubSubRuntimes;
import spec.sdk.reactor.v1.domain.core.StateRuntimes;

/**
 * Core Cloud Runtimes standard API defined.
 */
public interface CoreCloudRuntimes extends
        InvocationRuntimes,
        PubSubRuntimes,
        // BindingRuntimes,
        StateRuntimes,
        // SecretsRuntimes,
        ConfigurationRuntimes {
}
