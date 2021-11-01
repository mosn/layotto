package io.mosn.layotto.v1.grpc;

import io.grpc.ManagedChannel;
import io.mosn.layotto.v1.grpc.stub.StubManager;
import spec.proto.runtime.v1.RuntimeGrpc;
import spec.sdk.runtime.v1.client.RuntimeClient;

public interface GrpcRuntimeClient extends RuntimeClient {

    StubManager<RuntimeGrpc.RuntimeStub, RuntimeGrpc.RuntimeBlockingStub> getStubManager();
}