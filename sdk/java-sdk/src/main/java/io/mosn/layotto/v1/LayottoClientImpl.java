package io.mosn.layotto.v1;

import com.google.protobuf.Empty;
import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import spec.proto.runtime.v1.RuntimeGrpc;
import spec.proto.runtime.v1.RuntimeProto;

public class LayottoClientImpl implements LayottoClient {

    ManagedChannel channel;

    public LayottoClientImpl(String name, int port) {
        channel = ManagedChannelBuilder.forAddress(name, port)
                .usePlaintext()
                .build();
    }

    @Override
    public RuntimeProto.SayHelloResponse SayHello(RuntimeProto.SayHelloRequest req) {
        return RuntimeGrpc.newBlockingStub(channel).sayHello(req);
    }

    @Override
    public RuntimeProto.InvokeResponse InvokeService(RuntimeProto.InvokeServiceRequest req) {
        return RuntimeGrpc.newBlockingStub(channel).invokeService(req);
    }


    @Override
    public RuntimeProto.GetConfigurationResponse GetConfiguration(RuntimeProto.GetConfigurationRequest req) {
        return RuntimeGrpc.newBlockingStub(channel).getConfiguration(req);
    }

    @Override
    public Empty SaveConfiguration(RuntimeProto.SaveConfigurationRequest req) {
        return RuntimeGrpc.newBlockingStub(channel).saveConfiguration(req);
    }

    @Override
    public Empty DeleteConfiguration(RuntimeProto.DeleteConfigurationRequest req) {
        return RuntimeGrpc.newBlockingStub(channel).deleteConfiguration(req);
    }

    @Override
    public RuntimeProto.SubscribeConfigurationResponse SubscribeConfiguration(RuntimeProto.SubscribeConfigurationRequest req) {
        return null;
    }

    @Override
    public RuntimeProto.TryLockResponse TryLock(RuntimeProto.TryLockRequest req) {

        return RuntimeGrpc.newBlockingStub(channel).tryLock(req);
    }

    @Override
    public RuntimeProto.UnlockResponse Unlock(RuntimeProto.UnlockRequest req) {
        return null;
    }

    @Override
    public RuntimeProto.GetNextIdResponse GetNextId(RuntimeProto.GetNextIdRequest req) {
        return RuntimeGrpc.newBlockingStub(channel).getNextId(req);
    }

    @Override
    public RuntimeProto.GetStateResponse GetState(RuntimeProto.GetStateRequest req) {
        return null;
    }

    @Override
    public RuntimeProto.GetBulkStateResponse GetBulkState(RuntimeProto.GetBulkStateRequest req) {
        return null;
    }

    @Override
    public Empty SaveState(RuntimeProto.SaveStateRequest req) {
        return null;
    }

    @Override
    public Empty DeleteState(RuntimeProto.DeleteStateRequest req) {
        return null;
    }

    @Override
    public Empty DeleteBulkState(RuntimeProto.DeleteBulkStateRequest req) {
        return null;
    }

    @Override
    public Empty ExecuteStateTransaction(RuntimeProto.ExecuteStateTransactionRequest req) {
        return null;
    }

    @Override
    public Empty PublishEvent(RuntimeProto.PublishEventRequest req) {
        return null;
    }

    public static void main(String[] args) {
        LayottoClient client = new LayottoClientImpl("127.0.0.1", 34904);

        for (int i = 0; i < 10; i++) {
            RuntimeProto.GetNextIdResponse getNextIdResponse = client.GetNextId(RuntimeProto.GetNextIdRequest.newBuilder().setKey("key_xxx").setStoreName("redis").build());

            System.out.println(getNextIdResponse.getNextId());
        }

    }
}
