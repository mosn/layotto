/*
 * Copyright 2021 Layotto Authors
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package io.mosn.layotto.v1;

import io.grpc.Server;
import io.grpc.ServerBuilder;
import io.mosn.layotto.v1.callback.GrpcAppCallbackImpl;
import io.mosn.layotto.v1.callback.component.pubsub.Subscriber;
import io.mosn.layotto.v1.callback.component.pubsub.SubscriberRegistryImpl;
import io.mosn.layotto.v1.callback.component.pubsub.SubscriberRegistry;
import io.mosn.layotto.v1.grpc.ExceptionHandler;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicBoolean;

public class RuntimeServerGrpc {

    private static final Logger logger = LoggerFactory.getLogger(RuntimeServerGrpc.class.getName());

    private final    int           port;
    private volatile Server        server;
    private final    AtomicBoolean started = new AtomicBoolean(false);

    public RuntimeServerGrpc(int port) {
        this.port = port;
    }

    private final SubscriberRegistry subscriberRegistry = new SubscriberRegistryImpl();

    public void start() throws IOException {
        // 1 make sure at most once
        if (!started.compareAndSet(false, true)) {
            return;
        }
        // 2 start grpc server
        try {
            /* The port on which the server should run */
            server = ServerBuilder.forPort(port)
                    .addService(new GrpcAppCallbackImpl(subscriberRegistry))
                    .intercept(new ExceptionHandler())
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
                    RuntimeServerGrpc.this.stop();
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

    public SubscriberRegistry getPubSubRegistry() {
        return subscriberRegistry;
    }

    public void registerPubSubCallback(String pubsubName, Subscriber callback) {
        subscriberRegistry.registerPubSubCallback(pubsubName, callback);
    }
}
