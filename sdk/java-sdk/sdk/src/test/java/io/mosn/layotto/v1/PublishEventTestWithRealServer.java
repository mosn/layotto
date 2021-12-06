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
import io.mosn.layotto.v1.grpc.ExceptionHandler;
import io.mosn.layotto.v1.grpc.GrpcRuntimeClient;
import io.mosn.layotto.v1.mock.MyPublishService;
import org.junit.After;
import org.junit.Assert;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.junit.runners.JUnit4;
import spec.proto.runtime.v1.RuntimeProto;

@RunWith(JUnit4.class)
public class PublishEventTestWithRealServer {

    private final MyPublishService mockService = new MyPublishService();

    private Server                 srv;
    private GrpcRuntimeClient      client;

    int                            port        = 9999;
    String                         ip          = "127.0.0.1";

    @Before
    public void setUp() throws Exception {
        // start grpc server
        /* The port on which the server should run */
        srv = ServerBuilder.forPort(port)
            .addService(mockService)
            .intercept(new ExceptionHandler())
            .build()
            .start();

        // build a client
        client = new RuntimeClientBuilder()
            .withIp(ip)
            .withPort(port)
            .withConnectionPoolSize(4)
            .withTimeout(2000)
            .buildGrpc();
    }

    @After
    public void shutdown() throws InterruptedException {
        client.shutdown();
        srv.shutdownNow();
    }

    @Test
    public void testPublishEvent() {
        client.publishEvent("redis", "hello", "word".getBytes());
        RuntimeProto.PublishEventRequest last = mockService.getLastReq();
        Assert.assertEquals(last.getPubsubName(), "redis");
        Assert.assertEquals(last.getTopic(), "hello");
        Assert.assertEquals(new String(last.getData().toByteArray()), "word");
    }
}
