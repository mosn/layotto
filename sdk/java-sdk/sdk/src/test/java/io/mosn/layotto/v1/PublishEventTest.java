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

import io.grpc.ManagedChannel;
import io.grpc.inprocess.InProcessChannelBuilder;
import io.grpc.inprocess.InProcessServerBuilder;
import io.grpc.testing.GrpcCleanupRule;
import io.mosn.layotto.v1.mock.MyPublishService;
import org.junit.Assert;
import org.junit.Before;
import org.junit.Rule;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.junit.runners.JUnit4;
import spec.proto.runtime.v1.RuntimeGrpc;
import spec.proto.runtime.v1.RuntimeProto;
import spec.sdk.runtime.v1.client.RuntimeClient;

import static org.mockito.AdditionalAnswers.delegatesTo;
import static org.mockito.Mockito.mock;

@RunWith(JUnit4.class)
public class PublishEventTest {
    @Rule
    public final GrpcCleanupRule              grpcCleanup = new GrpcCleanupRule();

    private final MyPublishService            mockService = new MyPublishService();

    private final RuntimeGrpc.RuntimeImplBase serviceImpl =
                                                                  mock(RuntimeGrpc.RuntimeImplBase.class,
                                                                      delegatesTo(mockService));

    private RuntimeClient                     client;

    @Before
    public void setUp() throws Exception {
        String serverName = InProcessServerBuilder.generateName();
        grpcCleanup.register(InProcessServerBuilder
            .forName(serverName).directExecutor()
            .addService(serviceImpl)
            .build().start());
        ManagedChannel channel = grpcCleanup.register(
            InProcessChannelBuilder.forName(serverName).directExecutor().build());
        client = new RuntimeClientBuilder()
            .buildGrpcWithExistingChannel(channel);
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
