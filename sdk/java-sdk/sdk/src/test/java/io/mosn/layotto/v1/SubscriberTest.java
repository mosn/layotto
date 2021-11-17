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

import com.google.protobuf.ByteString;
import com.google.protobuf.Empty;
import io.grpc.inprocess.InProcessChannelBuilder;
import io.grpc.inprocess.InProcessServerBuilder;
import io.grpc.testing.GrpcCleanupRule;
import io.mosn.layotto.v1.callback.GrpcAppCallbackImpl;
import io.mosn.layotto.v1.callback.component.pubsub.SubscriberRegistry;
import io.mosn.layotto.v1.callback.component.pubsub.SubscriberRegistryImpl;
import io.mosn.layotto.v1.mock.MySubscriber;
import org.junit.Before;
import org.junit.Rule;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.junit.runners.JUnit4;
import spec.proto.runtime.v1.AppCallbackGrpc;
import spec.proto.runtime.v1.AppCallbackProto;

import static org.junit.Assert.assertEquals;

@RunWith(JUnit4.class)
public class SubscriberTest {
    private final static String             pubsubName  = "redis";
    private final static String             topic       = "hello";

    @Rule
    public final GrpcCleanupRule            grpcCleanup = new GrpcCleanupRule();

    AppCallbackGrpc.AppCallbackBlockingStub blockingStub;

    @Before
    public void setUp() throws Exception {
        // Generate a unique in-process server name.
        String serverName = InProcessServerBuilder.generateName();
        SubscriberRegistry sr = new SubscriberRegistryImpl();
        sr.registerPubSubCallback(pubsubName, new MySubscriber(pubsubName, topic));

        // Create a server, add service, start, and register for automatic graceful shutdown.
        grpcCleanup.register(InProcessServerBuilder
            .forName(serverName).directExecutor().addService(new GrpcAppCallbackImpl(sr)).build().start());

        // Create a client channel and register for automatic graceful shutdown.
        blockingStub = AppCallbackGrpc.newBlockingStub(
            grpcCleanup.register(InProcessChannelBuilder.forName(serverName).directExecutor().build()));
    }

    @Test
    public void listTopicSubscriptions() throws Exception {
        AppCallbackProto.ListTopicSubscriptionsResponse subscriptionsResponse = blockingStub.listTopicSubscriptions(
            Empty.getDefaultInstance());
        assertEquals(1, subscriptionsResponse.getSubscriptionsCount());
        assertEquals("hello", subscriptionsResponse.getSubscriptions(0).getTopic());
    }

    @Test
    public void testOnEventSuccess() throws Exception {
        //{
        //    "contentType": "text/plain",
        //        "data": "d29ybGQ=",
        //        "id": "a31aa292-2703-4f29-be58-c1798e540619",
        //        "pubsubName": "redis",
        //        "source": "runtime",
        //        "specVersion": "1.0",
        //        "topic": "hello",
        //        "type": "com.runtime.event.sent"
        //}
        AppCallbackProto.TopicEventRequest req = AppCallbackProto.TopicEventRequest.newBuilder()
            .setData(ByteString.copyFrom("d29ybGQ=".getBytes()))
            .setId("a31aa292-2703-4f29-be58-c1798e540619")
            .setPubsubName("redis")
            .setSpecVersion("1.0")
            .setTopic("hello")
            .setType("com.runtime.event.sent")
            .build();
        AppCallbackProto.TopicEventResponse resp = blockingStub.onTopicEvent(req);
        assertEquals(resp.getStatusValue(), 0);
    }

    @Test
    public void whenOnEventPubsubNameWrong_thenDrop() throws Exception {
        //{
        //    "contentType": "text/plain",
        //        "data": "d29ybGQ=",
        //        "id": "a31aa292-2703-4f29-be58-c1798e540619",
        //        "pubsubName": "redis",
        //        "source": "runtime",
        //        "specVersion": "1.0",
        //        "topic": "hello1",
        //        "type": "com.runtime.event.sent"
        //}
        AppCallbackProto.TopicEventRequest req = AppCallbackProto.TopicEventRequest.newBuilder()
            .setData(ByteString.copyFrom("d29ybGQ=".getBytes()))
            .setId("a31aa292-2703-4f29-be58-c1798e540619")
            .setPubsubName("redis")
            .setSpecVersion("1.0")
            .setTopic("hello1")
            .setType("com.runtime.event.sent")
            .build();
        AppCallbackProto.TopicEventResponse resp = blockingStub.onTopicEvent(req);
        assertEquals(resp.getStatusValue(), 2);
    }

    @Test(expected = io.grpc.StatusRuntimeException.class)
    public void whenOnEventTopicWrong_thenError() throws Exception {
        //{
        //    "contentType": "text/plain",
        //        "data": "d29ybGQ=",
        //        "id": "a31aa292-2703-4f29-be58-c1798e540619",
        //        "pubsubName": "redis1",
        //        "source": "runtime",
        //        "specVersion": "1.0",
        //        "topic": "hello",
        //        "type": "com.runtime.event.sent"
        //}
        AppCallbackProto.TopicEventRequest req = AppCallbackProto.TopicEventRequest.newBuilder()
            .setData(ByteString.copyFrom("d29ybGQ=".getBytes()))
            .setId("a31aa292-2703-4f29-be58-c1798e540619")
            .setPubsubName("redis1")
            .setSpecVersion("1.0")
            .setTopic("hello")
            .setType("com.runtime.event.sent")
            .build();
        AppCallbackProto.TopicEventResponse resp = blockingStub.onTopicEvent(req);
        assertEquals(resp.getStatusValue(), 2);
    }
}
