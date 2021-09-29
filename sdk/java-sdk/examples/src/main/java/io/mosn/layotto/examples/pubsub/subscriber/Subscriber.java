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