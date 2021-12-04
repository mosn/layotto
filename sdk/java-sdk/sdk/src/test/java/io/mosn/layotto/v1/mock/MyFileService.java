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

import com.google.protobuf.ByteString;
import com.google.protobuf.Empty;
import io.grpc.stub.StreamObserver;
import spec.proto.runtime.v1.RuntimeGrpc;
import spec.proto.runtime.v1.RuntimeProto;

public class MyFileService extends RuntimeGrpc.RuntimeImplBase {

    @Override
    public StreamObserver<RuntimeProto.PutFileRequest> putFile(StreamObserver<Empty> responseObserver) {

        return new StreamObserver<RuntimeProto.PutFileRequest>() {

            @Override
            public void onNext(RuntimeProto.PutFileRequest putFileRequest) {

                String log = String.format("put file store name %s, meta %d, file name %s, data size %d",
                    putFileRequest.getStoreName(),
                    putFileRequest.getMetadataCount(),
                    putFileRequest.getName(),
                    putFileRequest.getData().size());

                System.out.println(log);
            }

            @Override
            public void onError(Throwable throwable) {
                if (throwable != null) {
                    System.err.println("put file err: " + throwable);
                }
            }

            @Override
            public void onCompleted() {
                System.out.println("finished put file");
                responseObserver.onNext(Empty.newBuilder().build());
                responseObserver.onCompleted();
            }
        };
    }

    @Override
    public void getFile(RuntimeProto.GetFileRequest request,
                        StreamObserver<RuntimeProto.GetFileResponse> responseObserver) {

        String echo = String.format("get file store name %s, meta %d, file name %s",
            request.getStoreName(),
            request.getMetadataCount(),
            request.getName());

        responseObserver.onNext(
            RuntimeProto.GetFileResponse.newBuilder().
                setData(
                    ByteString.copyFrom(echo.getBytes())).
                build());

        responseObserver.onCompleted();
    }

    @Override
    public void delFile(RuntimeProto.DelFileRequest request, StreamObserver<Empty> responseObserver) {

        String log = String.format("del file store name %s, meta %d, file name %s",
            request.getRequest().getStoreName(),
            request.getRequest().getMetadataCount(),
            request.getRequest().getName());

        System.out.println(log);

        responseObserver.onNext(Empty.newBuilder().build());
        responseObserver.onCompleted();
    }

    @Override
    public void listFile(RuntimeProto.ListFileRequest request,
                         StreamObserver<RuntimeProto.ListFileResp> responseObserver) {

        String echo = String.format("put file store name %s, meta %d",
            request.getRequest().getStoreName(),
            request.getRequest().getMetadataCount());

        responseObserver.onNext(
            RuntimeProto.ListFileResp.newBuilder().
                addFiles(
                    RuntimeProto.FileInfo.newBuilder().
                        setFileName(echo).
                        setSize(100).
                        setLastModified("2021-11-23 10:24:11").
                        putMetadata("k1", "v1").
                        build()).
                setMarker("marker").
                setIsTruncated(true).
                build()
            );

        responseObserver.onCompleted();
    }

    @Override
    public void getFileMeta(RuntimeProto.GetFileMetaRequest request,
                            StreamObserver<RuntimeProto.GetFileMetaResponse> responseObserver) {

        String log = String.format("get file meta store name %s, meta %d, file name %s",
            request.getRequest().getStoreName(),
            request.getRequest().getMetadataCount(),
            request.getRequest().getName());

        System.out.println(log);

        responseObserver.onNext(
            RuntimeProto.GetFileMetaResponse.newBuilder().
                setSize(100).
                setLastModified("2021-11-22 10:24:11").
                setResponse(
                    RuntimeProto.FileMeta.newBuilder().
                        putMetadata(
                            "k1",
                            RuntimeProto.FileMetaValue.newBuilder().
                                addValue("v1").
                                addValue("v2").
                                build()).
                        build()).
                build());

        responseObserver.onCompleted();
    }
}
