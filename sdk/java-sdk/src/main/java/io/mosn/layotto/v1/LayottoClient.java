package io.mosn.layotto.v1;

import com.google.protobuf.Empty;
import spec.proto.runtime.v1.RuntimeProto;

public interface LayottoClient {
    RuntimeProto.SayHelloResponse SayHello(RuntimeProto.SayHelloRequest req);

    RuntimeProto.InvokeResponse InvokeService(RuntimeProto.InvokeServiceRequest req);

    RuntimeProto.GetConfigurationResponse GetConfiguration(RuntimeProto.GetConfigurationRequest req);

    Empty SaveConfiguration(RuntimeProto.SaveConfigurationRequest req);

    Empty DeleteConfiguration(RuntimeProto.DeleteConfigurationRequest req);

    RuntimeProto.SubscribeConfigurationResponse SubscribeConfiguration(RuntimeProto.SubscribeConfigurationRequest req);

    RuntimeProto.TryLockResponse TryLock(RuntimeProto.TryLockRequest req);

    RuntimeProto.UnlockResponse Unlock(RuntimeProto.UnlockRequest req);

    RuntimeProto.GetNextIdResponse GetNextId(RuntimeProto.GetNextIdRequest req);

    RuntimeProto.GetStateResponse GetState(RuntimeProto.GetStateRequest req);

    RuntimeProto.GetBulkStateResponse GetBulkState(RuntimeProto.GetBulkStateRequest req);

    Empty SaveState(RuntimeProto.SaveStateRequest req);

    Empty DeleteState(RuntimeProto.DeleteStateRequest req);

    Empty DeleteBulkState(RuntimeProto.DeleteBulkStateRequest req);

    Empty ExecuteStateTransaction(RuntimeProto.ExecuteStateTransactionRequest req);

    Empty PublishEvent(RuntimeProto.PublishEventRequest req);
}
