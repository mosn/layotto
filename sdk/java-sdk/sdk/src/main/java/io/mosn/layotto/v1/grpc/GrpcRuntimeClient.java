package io.mosn.layotto.v1.grpc;

import io.grpc.ManagedChannel;
import spec.sdk.runtime.v1.client.RuntimeClient;

public interface GrpcRuntimeClient extends RuntimeClient {
    ManagedChannel[] getChannels();
}