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

import com.google.protobuf.Empty;
import io.grpc.inprocess.InProcessChannelBuilder;
import io.grpc.inprocess.InProcessServerBuilder;
import io.grpc.testing.GrpcCleanupRule;
import io.mosn.layotto.v1.callback.GrpcAppCallbackImpl;
import io.mosn.layotto.v1.callback.component.pubsub.PubSubRegistry;
import org.junit.Rule;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.junit.runners.JUnit4;
import org.mockito.Mockito;
import spec.proto.runtime.v1.AppCallbackGrpc;
import spec.proto.runtime.v1.AppCallbackProto;

import static org.junit.Assert.assertEquals;

@RunWith(JUnit4.class)
public class ListTopicSubscriptionsTest {

    @Rule
    public final GrpcCleanupRule grpcCleanup = new GrpcCleanupRule();

    @Test
    public void listTopicSubscriptions() throws Exception {
        // Generate a unique in-process server name.
        String serverName = InProcessServerBuilder.generateName();

        // Create a server, add service, start, and register for automatic graceful shutdown.
        grpcCleanup.register(InProcessServerBuilder
                .forName(serverName).directExecutor().addService(new GrpcAppCallbackImpl(
                        Mockito.mock(PubSubRegistry.class))).build().start());

        // Create a client channel and register for automatic graceful shutdown.
        AppCallbackGrpc.AppCallbackBlockingStub blockingStub = AppCallbackGrpc.newBlockingStub(
                grpcCleanup.register(InProcessChannelBuilder.forName(serverName).directExecutor().build()));

        AppCallbackProto.ListTopicSubscriptionsResponse subscriptionsResponse = blockingStub.listTopicSubscriptions(
                Empty.getDefaultInstance());
        assertEquals(0, subscriptionsResponse.getSubscriptionsCount());
        //assertEquals("hello", subscriptionsResponse.getSubscriptions(0).getTopic());
    }

}
