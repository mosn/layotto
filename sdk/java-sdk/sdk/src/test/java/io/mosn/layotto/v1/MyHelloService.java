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
import io.grpc.stub.StreamObserver;
import spec.proto.runtime.v1.RuntimeGrpc;
import spec.proto.runtime.v1.RuntimeProto;

import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class MyHelloService extends RuntimeGrpc.RuntimeImplBase {

    @Override
    public void sayHello(RuntimeProto.SayHelloRequest request, StreamObserver<RuntimeProto.SayHelloResponse> responseObserver) {

        RuntimeProto.SayHelloResponse resp = RuntimeProto.SayHelloResponse.newBuilder()
                .setHello("hi, " + request.getServiceName()).build();
        responseObserver.onNext(resp);
        responseObserver.onCompleted();
    }
}