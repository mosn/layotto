package io.mosn.layotto.v1;

import io.grpc.ManagedChannel;
import io.grpc.inprocess.InProcessChannelBuilder;
import io.grpc.inprocess.InProcessServerBuilder;
import io.grpc.testing.GrpcCleanupRule;
import io.mosn.layotto.v1.config.RuntimeProperties;
import io.mosn.layotto.v1.domain.ApiProtocol;
import io.mosn.layotto.v1.serializer.JSONSerializer;
import io.mosn.layotto.v1.serializer.ObjectSerializer;
import org.junit.Before;
import org.junit.Rule;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.junit.runners.JUnit4;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import spec.proto.runtime.v1.RuntimeGrpc;
import spec.proto.runtime.v1.RuntimeProto;
import spec.sdk.runtime.v1.client.RuntimeClient;

import static org.junit.Assert.assertEquals;
import static org.mockito.AdditionalAnswers.delegatesTo;
import static org.mockito.Mockito.mock;

@RunWith(JUnit4.class)
public class SayHelloTest {

    private static final Logger logger = LoggerFactory.getLogger(RuntimeClient.class.getName());

    private int timeoutMs = RuntimeProperties.DEFAULT_TIMEOUT_MS;

    private String ip = RuntimeProperties.DEFAULT_IP;

    private int port = RuntimeProperties.DEFAULT_PORT;

    private ApiProtocol protocol = RuntimeProperties.DEFAULT_API_PROTOCOL;

    private ObjectSerializer stateSerializer = new JSONSerializer();

    @Rule
    public final GrpcCleanupRule grpcCleanup = new GrpcCleanupRule();

    private final RuntimeGrpc.RuntimeImplBase serviceImpl =
            mock(RuntimeGrpc.RuntimeImplBase.class, delegatesTo(
                    new RuntimeGrpc.RuntimeImplBase() {
                        @Override
                        public void sayHello(RuntimeProto.SayHelloRequest request,
                                             io.grpc.stub.StreamObserver<spec.proto.runtime.v1.RuntimeProto.SayHelloResponse> responseObserver) {
                            responseObserver.onNext(
                                    RuntimeProto.SayHelloResponse.newBuilder().setHello("hi, " + request.getServiceName()).build());
                            responseObserver.onCompleted();
                        }
                    }));

    private RuntimeClient client;

    @Before
    public void setUp() throws Exception {
        // Generate a unique in-process server name.
        String serverName = InProcessServerBuilder.generateName();

        // Create a server, add service, start, and register for automatic graceful shutdown.
        grpcCleanup.register(InProcessServerBuilder
                .forName(serverName).directExecutor().addService(serviceImpl).build().start());

        // Create a client channel and register for automatic graceful shutdown.
        ManagedChannel channel = grpcCleanup.register(
                InProcessChannelBuilder.forName(serverName).directExecutor().build());

        // Create a HelloWorldClient using the in-process channel;

        client = new RuntimeClientBuilder()
                .buildGrpcWithExistingChannel(channel);
    }

    @Test
    public void sayHello() {
        String greet = client.sayHello("layotto");
        assertEquals("hi, layotto", greet);
    }

}
