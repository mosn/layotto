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
package io.mosn.layotto.v1.mock;

import com.google.protobuf.Empty;
import io.grpc.stub.StreamObserver;
import spec.proto.runtime.v1.RuntimeGrpc;
import spec.proto.runtime.v1.RuntimeProto;

public class MyPublishService extends RuntimeGrpc.RuntimeImplBase {
    RuntimeProto.PublishEventRequest lastReq = null;

    @Override
    public void publishEvent(RuntimeProto.PublishEventRequest request, StreamObserver<Empty> responseObserver) {
        lastReq = request;
        responseObserver.onNext(null);
        responseObserver.onCompleted();
    }

    /**
     * Getter method for property <tt>lastReq</tt>.
     *
     * @return property value of lastReq
     */
    public RuntimeProto.PublishEventRequest getLastReq() {
        return lastReq;
    }
}