package io.mosn.layotto.v1;

import io.grpc.Server;
import io.grpc.ServerBuilder;
import io.mosn.layotto.v1.callback.GrpcAppCallbackImpl;
import io.mosn.layotto.v1.callback.component.pubsub.PubSub;
import io.mosn.layotto.v1.callback.component.pubsub.PubSubClientRegistryImpl;
import io.mosn.layotto.v1.callback.component.pubsub.PubSubRegistry;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicBoolean;

public class RuntimeServer {

    private static final Logger logger = LoggerFactory.getLogger(RuntimeServer.class.getName());

    private final    int           port;
    private volatile Server        server;
    private final    AtomicBoolean started = new AtomicBoolean(false);

    public RuntimeServer(int port) {
        this.port = port;
    }

    private final PubSubRegistry pubSubRegistry = new PubSubClientRegistryImpl();

    public void start() throws IOException {
        // 1 make sure at most once
        if (!started.compareAndSet(false, true)) {
            return;
        }
        // 2 start grpc server
        try {
            /* The port on which the server should run */
            server = ServerBuilder.forPort(port)
                    .addService(new GrpcAppCallbackImpl(pubSubRegistry))
                    .build()
                    .start();
        } catch (Exception e) {
            // revert
            server = null;
            started.set(false);
            throw e;
        }
        logger.info("Server started, listening on {}", port);

        // 3 addShutdownHook
        Runtime.getRuntime().addShutdownHook(new Thread() {
            @Override
            public void run() {
                // Use stderr here since the logger may have been reset by its JVM shutdown hook.
                logger.info("*** shutting down gRPC server since JVM is shutting down");
                try {
                    RuntimeServer.this.stop();
                } catch (Exception e) {
                    logger.error("*** server shut down error", e);
                }
            }
        });
    }

    public void stop() throws InterruptedException {
        Server srv = this.server;
        if (srv != null) {
            srv.shutdown().awaitTermination(30, TimeUnit.SECONDS);
        }
    }

    /**
     * Await termination on the main thread since the grpc library uses daemon threads.
     */
    public void blockUntilShutdown() throws InterruptedException {
        Server srv = this.server;
        if (srv != null) {
            srv.awaitTermination();
        }
    }

    public PubSubRegistry getPubSubRegistry() {
        return pubSubRegistry;
    }

    public void registerPubSubCallback(String pubsubName, PubSub callback) {
        pubSubRegistry.registerPubSubCallback(pubsubName, callback);
    }
}
