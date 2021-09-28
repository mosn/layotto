package spec.sdk.runtime.v1.client;

import spec.sdk.runtime.v1.domain.*;

public interface RuntimeClient extends
        HelloRuntime,
        ConfigurationRuntime,
        InvocationRuntime,
        PubSubRuntime,
        StateRuntime,
        LockRuntime,
        SequencerRuntime {
}