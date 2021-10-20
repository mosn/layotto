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
 *
 */
package io.mosn.layotto.examples.pubsub.subscriber;

import com.alibaba.fastjson.JSON;
import io.mosn.layotto.v1.RuntimeServer;

import java.util.concurrent.Semaphore;

public class Subscriber {

    /**
     * This is the entry point for this example app, which subscribes to a topic.
     *
     * @throws Exception An Exception on startup.
     */
    public static void main(String[] args) throws Exception {
        RuntimeServer srv = new RuntimeServer(9999);
        RawPubSub pubsub = new RawPubSub("redis");
        pubsub.subscribe("hello", request -> {
            System.out.println(JSON.toJSONString(request));
        });
        pubsub.subscribe("topic1", request -> {
            System.out.println(JSON.toJSONString(request));
        });
        srv.registerPubSubCallback(pubsub.getComponentName(), pubsub);
        Semaphore sm = new Semaphore(0);
        srv.start();
        sm.acquire();
    }
}