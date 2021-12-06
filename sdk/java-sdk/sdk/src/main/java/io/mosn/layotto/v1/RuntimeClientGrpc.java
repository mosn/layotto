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

import com.google.common.base.Strings;
import com.google.protobuf.Any;
import com.google.protobuf.ByteString;
import com.google.protobuf.Empty;
import io.grpc.Metadata;
import io.grpc.StatusRuntimeException;
import io.grpc.stub.MetadataUtils;
import io.grpc.stub.StreamObserver;
import io.mosn.layotto.v1.config.RuntimeProperties;
import io.mosn.layotto.v1.exceptions.RuntimeClientException;
import io.mosn.layotto.v1.grpc.GrpcRuntimeClient;
import io.mosn.layotto.v1.grpc.stub.StubManager;
import io.mosn.layotto.v1.serializer.ObjectSerializer;
import org.slf4j.Logger;
import spec.proto.runtime.v1.RuntimeGrpc;
import spec.proto.runtime.v1.RuntimeProto;
import spec.sdk.runtime.v1.domain.file.DelFileRequest;
import spec.sdk.runtime.v1.domain.file.DelFileResponse;
import spec.sdk.runtime.v1.domain.file.FileInfo;
import spec.sdk.runtime.v1.domain.file.GetFileRequest;
import spec.sdk.runtime.v1.domain.file.GetFileResponse;
import spec.sdk.runtime.v1.domain.file.GetMetaRequest;
import spec.sdk.runtime.v1.domain.file.GetMeteResponse;
import spec.sdk.runtime.v1.domain.file.ListFileRequest;
import spec.sdk.runtime.v1.domain.file.ListFileResponse;
import spec.sdk.runtime.v1.domain.file.PutFileRequest;
import spec.sdk.runtime.v1.domain.file.PutFileResponse;
import spec.sdk.runtime.v1.domain.invocation.InvokeResponse;
import spec.sdk.runtime.v1.domain.state.DeleteStateRequest;
import spec.sdk.runtime.v1.domain.state.ExecuteStateTransactionRequest;
import spec.sdk.runtime.v1.domain.state.GetBulkStateRequest;
import spec.sdk.runtime.v1.domain.state.GetStateRequest;
import spec.sdk.runtime.v1.domain.state.SaveStateRequest;
import spec.sdk.runtime.v1.domain.state.State;
import spec.sdk.runtime.v1.domain.state.StateOptions;
import spec.sdk.runtime.v1.domain.state.TransactionalStateOperation;

import java.io.IOException;
import java.io.InputStream;
import java.io.PipedInputStream;
import java.io.PipedOutputStream;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.TimeUnit;

public class RuntimeClientGrpc extends AbstractRuntimeClient implements GrpcRuntimeClient {

    private static final String                                                         TIMEOUT_KEY = "timeout";
    private final StubManager<RuntimeGrpc.RuntimeStub, RuntimeGrpc.RuntimeBlockingStub> stubManager;

    RuntimeClientGrpc(Logger logger,
                      int timeoutMs,
                      ObjectSerializer stateSerializer,
                      StubManager<RuntimeGrpc.RuntimeStub, RuntimeGrpc.RuntimeBlockingStub> stubManager) {
        super(logger, timeoutMs, stateSerializer);
        this.stubManager = stubManager;
    }

    @Override
    public String sayHello(String name, int timeoutMillisecond) {
        try {
            // 1. prepare request
            RuntimeProto.SayHelloRequest req = RuntimeProto.SayHelloRequest.newBuilder()
                .setServiceName(name)
                .setName(name)
                .build();

            // 2. invoke
            RuntimeProto.SayHelloResponse response = stubManager.getBlockingStub()
                .withDeadlineAfter(timeoutMillisecond,
                    TimeUnit.MILLISECONDS)
                .sayHello(req);

            // 3. parse result
            return response.getHello();
        } catch (Exception e) {
            logger.error("sayHello error ", e);
            throw new RuntimeClientException(e);
        }
    }

    @Override
    public InvokeResponse<byte[]> invokeMethod(String appId, String methodName, byte[] data, Map<String, String> header, int timeoutMs) {
        try {
            // 1. prepare request
            final ByteString byteString = ByteString.copyFrom(data);
            Any anyData = Any.newBuilder()
                    .setValue(byteString)
                    .build();
            RuntimeProto.CommonInvokeRequest commonReq = RuntimeProto.CommonInvokeRequest.newBuilder()
                    .setMethod(methodName)
                    .setData(anyData)
                    .build();
            RuntimeProto.InvokeServiceRequest invokeReq = RuntimeProto.InvokeServiceRequest.newBuilder()
                    .setId(appId)
                    .setMessage(commonReq)
                    .build();
            // metadata
            Metadata metadata = new Metadata();
            for (Map.Entry<String, String> entry : header.entrySet()) {
                if (entry.getKey() != null && entry.getValue() != null) {
                    Metadata.Key<String> key = Metadata.Key.of(entry.getKey(), Metadata.ASCII_STRING_MARSHALLER);
                    metadata.put(key, entry.getValue());
                }
            }
            Metadata.Key<String> key = Metadata.Key.of(TIMEOUT_KEY, Metadata.ASCII_STRING_MARSHALLER);
            metadata.put(key, Integer.toString(timeoutMs));

            // 2. invoke
            RuntimeProto.InvokeResponse resp = this.stubManager.getBlockingStub()
                    .withDeadlineAfter(timeoutMs, TimeUnit.MILLISECONDS)
                    .withInterceptors(MetadataUtils.newAttachHeadersInterceptor(metadata))
                    .invokeService(invokeReq);

            // 3. parse result
            InvokeResponse<byte[]> result = new InvokeResponse<>();
            result.setContentType(resp.getContentType());
            byte[] bytes = new byte[]{};
            result.setData(bytes);
            if (resp.getData() == null) {
                return result;
            }
            if (resp.getData().getValue() == null) {
                return result;
            }
            if (resp.getData().getValue().toByteArray() == null) {
                return result;
            }
            result.setData(resp.getData().getValue().toByteArray());
            return result;
        } catch (Exception e) {
            logger.error("invokeMethod error ", e);
            throw new RuntimeClientException(e);
        }
    }

    @Override
    public void publishEvent(String pubsubName, String topicName, byte[] data, String contentType,
                             Map<String, String> metadata) {
        try {
            // 1. prepare data
            if (data == null) {
                data = new byte[] {};
            }
            final ByteString byteString = ByteString.copyFrom(data);
            // Content-type can be overwritten on a per-request basis.
            // It allows CloudEvents to be handled differently, for example.
            if (contentType == null || contentType.isEmpty()) {
                contentType = RuntimeProperties.DEFAULT_PUBSUB_CONTENT_TYPE;
            }

            // 2. prepare request
            RuntimeProto.PublishEventRequest.Builder envelopeBuilder = RuntimeProto.PublishEventRequest.newBuilder()
                .setTopic(topicName)
                .setPubsubName(pubsubName)
                .setData(byteString)
                .setDataContentType(contentType);
            // metadata
            if (metadata != null) {
                envelopeBuilder.putAllMetadata(metadata);
            }
            RuntimeProto.PublishEventRequest req = envelopeBuilder.build();

            // 3. invoke
            this.stubManager.getBlockingStub().publishEvent(req);
        } catch (Exception e) {
            logger.error("publishEvent error ", e);
            throw new RuntimeClientException(e);
        }
    }

    @Override
    public void saveBulkState(SaveStateRequest request, int timeoutMs) {
        final String stateStoreName = request.getStoreName();
        final List<State<?>> states = request.getStates();
        // 1. validate
        if ((stateStoreName == null) || (stateStoreName.trim().isEmpty())) {
            throw new IllegalArgumentException("State store name cannot be null or empty.");
        }
        try {
            // 2. prepare request
            RuntimeProto.SaveStateRequest.Builder builder = RuntimeProto.SaveStateRequest.newBuilder();
            builder.setStoreName(stateStoreName);
            for (State<?> state : states) {
                // convert request and do serialization
                RuntimeProto.StateItem stateItem = buildStateRequest(state)
                    .build();
                builder.addStates(stateItem);
            }
            RuntimeProto.SaveStateRequest req = builder.build();
            // 3. invoke
            this.stubManager.getBlockingStub()
                .withDeadlineAfter(timeoutMs, TimeUnit.MILLISECONDS)
                .saveState(req);
        } catch (Exception e) {
            logger.error("saveBulkState error ", e);
            throw new RuntimeClientException(e);
        }
    }

    private <T> RuntimeProto.StateItem.Builder buildStateRequest(State<T> state) throws IOException {
        RuntimeProto.StateItem.Builder stateBuilder = RuntimeProto.StateItem.newBuilder();
        T value = state.getValue();
        // 1. serialize value
        byte[] bytes = null;
        if (value != null) {
            bytes = stateSerializer.serialize(value);
        }
        if (bytes != null) {
            stateBuilder.setValue(ByteString.copyFrom(bytes));
        }

        // 2. etag
        if (state.getEtag() != null) {
            RuntimeProto.Etag etag = RuntimeProto.Etag.newBuilder()
                .setValue(state.getEtag())
                .build();
            stateBuilder.setEtag(etag);
        }
        // 3. metadata
        if (state.getMetadata() != null) {
            stateBuilder.putAllMetadata(state.getMetadata());
        }
        // 4. key
        stateBuilder.setKey(state.getKey());
        // 5. StateOptions
        RuntimeProto.StateOptions.Builder optionBuilder = null;
        if (state.getOptions() != null) {
            StateOptions options = state.getOptions();
            optionBuilder = RuntimeProto.StateOptions.newBuilder();
            if (options.getConcurrency() != null) {
                RuntimeProto.StateOptions.StateConcurrency concurrency = getGrpcStateConcurrency(options);
                optionBuilder.setConcurrency(concurrency);
            }
            if (options.getConsistency() != null) {
                RuntimeProto.StateOptions.StateConsistency consistency = getGrpcStateConsistency(options);
                optionBuilder.setConsistency(consistency);
            }
        }
        if (optionBuilder != null) {
            RuntimeProto.StateOptions stateOptions = optionBuilder.build();
            stateBuilder.setOptions(stateOptions);
        }

        return stateBuilder;
    }

    @Override
    public void deleteState(DeleteStateRequest request, int timeoutMs) {
        final String stateStoreName = request.getStateStoreName();
        final String key = request.getKey();
        final StateOptions options = request.getStateOptions();
        final String etag = request.getEtag();
        final Map<String, String> metadata = request.getMetadata();

        // 1. validate
        if ((stateStoreName == null) || (stateStoreName.trim().isEmpty())) {
            throw new IllegalArgumentException("State store name cannot be null or empty.");
        }
        if ((key == null) || (key.trim().isEmpty())) {
            throw new IllegalArgumentException("Key cannot be null or empty.");
        }
        try {
            // 2. prepare request
            RuntimeProto.StateOptions.Builder optionBuilder = null;
            if (options != null) {
                optionBuilder = RuntimeProto.StateOptions.newBuilder();
                if (options.getConcurrency() != null) {
                    RuntimeProto.StateOptions.StateConcurrency concurrency = getGrpcStateConcurrency(options);
                    optionBuilder.setConcurrency(concurrency);
                }
                if (options.getConsistency() != null) {
                    RuntimeProto.StateOptions.StateConsistency consistency = getGrpcStateConsistency(options);
                    optionBuilder.setConsistency(consistency);
                }
            }
            RuntimeProto.DeleteStateRequest.Builder builder = RuntimeProto.DeleteStateRequest.newBuilder()
                .setStoreName(stateStoreName)
                .setKey(key);
            if (metadata != null) {
                builder.putAllMetadata(metadata);
            }
            if (etag != null) {
                RuntimeProto.Etag value = RuntimeProto.Etag.newBuilder()
                    .setValue(etag)
                    .build();
                builder.setEtag(value);
            }
            if (optionBuilder != null) {
                RuntimeProto.StateOptions stateOptions = optionBuilder.build();
                builder.setOptions(stateOptions);
            }
            RuntimeProto.DeleteStateRequest req = builder.build();

            // 3. invoke
            this.stubManager.getBlockingStub()
                .withDeadlineAfter(timeoutMs, TimeUnit.MILLISECONDS)
                .deleteState(req);
        } catch (Exception e) {
            logger.error("deleteState error ", e);
            throw new RuntimeClientException(e);
        }
    }

    private RuntimeProto.StateOptions.StateConcurrency getGrpcStateConcurrency(StateOptions options) {
        switch (options.getConcurrency()) {
            case FIRST_WRITE:
                return RuntimeProto.StateOptions.StateConcurrency.CONCURRENCY_FIRST_WRITE;
            case LAST_WRITE:
                return RuntimeProto.StateOptions.StateConcurrency.CONCURRENCY_LAST_WRITE;
            default:
                throw new IllegalArgumentException("Missing StateConcurrency mapping to gRPC Concurrency enum");
        }
    }

    private RuntimeProto.StateOptions.StateConsistency getGrpcStateConsistency(StateOptions options) {
        switch (options.getConsistency()) {
            case EVENTUAL:
                return RuntimeProto.StateOptions.StateConsistency.CONSISTENCY_EVENTUAL;
            case STRONG:
                return RuntimeProto.StateOptions.StateConsistency.CONSISTENCY_STRONG;
            default:
                throw new IllegalArgumentException("Missing Consistency mapping to gRPC Consistency enum");
        }
    }

    /**
     * Execute a transaction.
     *
     * @param request Request to execute transaction.
     */
    @Override
    public void executeStateTransaction(ExecuteStateTransactionRequest request) {
        final String stateStoreName = request.getStateStoreName();
        final List<TransactionalStateOperation<?>> operations = request.getOperations();
        final Map<String, String> metadata = request.getMetadata();

        // 1. validate
        assertTrue(stateStoreName != null && !stateStoreName.trim().isEmpty(),
            "stateStoreName cannot be null or empty.");
        assertTrue(operations != null && !operations.isEmpty(), "operations cannot be null or empty.");
        try {
            // 2. construct request object
            RuntimeProto.ExecuteStateTransactionRequest.Builder builder = RuntimeProto.ExecuteStateTransactionRequest
                .newBuilder();
            builder.setStoreName(stateStoreName);
            if (metadata != null) {
                builder.putAllMetadata(metadata);
            }
            for (TransactionalStateOperation<?> op : operations) {
                // validate each operation
                assertTrue(op.getOperation() != null, "operation cannot be null.");
                State<?> req = op.getRequest();
                assertTrue(req != null, "request cannot be null.");
                String k = req.getKey();
                assertTrue(k != null && !k.isEmpty(), "request cannot be null.");

                // build grpc request
                RuntimeProto.TransactionalStateOperation.Builder operationBuilder = RuntimeProto.TransactionalStateOperation
                    .newBuilder();
                String operationType = op.getOperation().toString().toLowerCase();
                operationBuilder.setOperationType(operationType);

                // convert request and do serialization
                RuntimeProto.StateItem stateItem = buildStateRequest(req)
                    .build();
                operationBuilder.setRequest(stateItem);

                builder.addOperations(operationBuilder.build());
            }
            RuntimeProto.ExecuteStateTransactionRequest req = builder.build();

            // 3. invoke grpc api
            this.stubManager.getBlockingStub().executeStateTransaction(req);
        } catch (IllegalArgumentException e) {
            logger.error("executeStateTransaction error ", e);
            throw e;
        } catch (Exception e) {
            logger.error("executeStateTransaction error ", e);
            throw new RuntimeClientException(e);
        }
    }

    private void assertTrue(boolean argumentAssertion, String errMsg) {
        if (!argumentAssertion) {
            throw new IllegalArgumentException(errMsg);
        }
    }

    @Override
    protected State<byte[]> doGetState(GetStateRequest request, int timeoutMs) {
        // 1. extract fields.
        final String stateStoreName = request.getStoreName();
        final String key = request.getKey();
        final StateOptions options = request.getStateOptions();
        final Map<String, String> metadata = request.getMetadata();
        try {
            // 2. construct request object
            RuntimeProto.GetStateRequest.Builder builder = RuntimeProto.GetStateRequest.newBuilder()
                    .setStoreName(stateStoreName)
                    .setKey(key);
            if (metadata != null) {
                builder.putAllMetadata(metadata);
            }
            if (options != null && options.getConsistency() != null) {
                RuntimeProto.StateOptions.StateConsistency consistency = getGrpcStateConsistency(options);
                builder.setConsistency(consistency);
            }
            RuntimeProto.GetStateRequest envelope = builder.build();

            // 3. invoke grpc api
            RuntimeGrpc.RuntimeBlockingStub stub = this.stubManager.getBlockingStub();
            if (timeoutMs > 0) {
                stub = stub.withDeadlineAfter(timeoutMs, TimeUnit.MILLISECONDS);
            }
            RuntimeProto.GetStateResponse getStateResponse = stub.getState(envelope);

            // 4. parse result
            // value
            final ByteString payload = getStateResponse.getData();
            byte[] value = payload == null ? null : payload.toByteArray();
            // etag
            String etag = getStateResponse.getEtag();
            if (etag != null && etag.isEmpty()) {
                etag = null;
            }
            return new State<>(key, value, etag, getStateResponse.getMetadataMap(), options);
        } catch (Exception e) {
            logger.error("getState error ", e);
            throw new RuntimeClientException(e);
        }
    }

    @Override
    protected List<State<byte[]>> doGetBulkState(GetBulkStateRequest request, int timeoutMs) {
        // 1. extract fields
        final String stateStoreName = request.getStoreName();
        final List<String> keys = request.getKeys();
        final int parallelism = request.getParallelism();
        final Map<String, String> metadata = request.getMetadata();

        try {
            // 2. construct request object
            RuntimeProto.GetBulkStateRequest.Builder builder = RuntimeProto.GetBulkStateRequest.newBuilder()
                    .setStoreName(stateStoreName)
                    .addAllKeys(keys)
                    .setParallelism(parallelism);
            if (metadata != null) {
                builder.putAllMetadata(metadata);
            }
            RuntimeProto.GetBulkStateRequest envelope = builder.build();

            // 3. invoke grpc API
            RuntimeGrpc.RuntimeBlockingStub stub = this.stubManager.getBlockingStub();
            if (timeoutMs > 0) {
                stub = stub.withDeadlineAfter(timeoutMs, TimeUnit.MILLISECONDS);
            }
            RuntimeProto.GetBulkStateResponse resp = stub.getBulkState(envelope);

            // 4. parse result
            List<RuntimeProto.BulkStateItem> itemsList = resp.getItemsList();
            List<State<byte[]>> result = new ArrayList<>(itemsList.size());
            for (RuntimeProto.BulkStateItem itm : itemsList) {
                State<byte[]> tState = parseGetStateResult(itm);
                result.add(tState);
            }
            return result;
        } catch (Exception e) {
            logger.error("getBulkState error ", e);
            throw new RuntimeClientException(e);
        }
    }

    private State<byte[]> parseGetStateResult(RuntimeProto.BulkStateItem bulkStateItem) throws IOException {
        final String key = bulkStateItem.getKey();
        // check error
        final String error = bulkStateItem.getError();
        if (!Strings.isNullOrEmpty(error)) {
            return new State<>(key, error);
        }
        // value
        final ByteString payload = bulkStateItem.getData();
        byte[] value = payload == null ? null : payload.toByteArray();
        // etag
        String etag = bulkStateItem.getEtag();
        if (etag != null && etag.isEmpty()) {
            etag = null;
        }
        return new State<>(key, value, etag, bulkStateItem.getMetadataMap(), null);
    }

    /**
     * Getter method for property <tt>stubManager</tt>.
     *
     * Do not use it !
     * This method is deprecated and might be refactored in the future.
     * We want this client to expose grpc Channels instead of grpc stubs.
     *
     * @return property value of stubManager
     */
    @Deprecated
    @Override
    public StubManager<RuntimeGrpc.RuntimeStub, RuntimeGrpc.RuntimeBlockingStub> getStubManager() {
        return stubManager;
    }

    @Override
    public void shutdown() {
        stubManager.destroy();
    }

    @Override
    public PutFileResponse putFile(PutFileRequest request, int timeoutMs) throws Exception {

        checkParamOfPutFile(request);

        PutFileFuture putFuture = new PutFileFuture(request.getFileName());
        StreamObserver<RuntimeProto.PutFileRequest> observer = createPutFileObserver(putFuture, timeoutMs);

        observer.onNext(buildPutFileMetaDataRequest(request.getStoreName(), request.getFileName(),
            request.getMetaData()));

        byte[] buf = new byte[4096];
        for (int size = request.getIn().read(buf); size > 0; size = request.getIn().read(buf)) {
            observer.onNext(buildPutFileDataRequest(buf, size));
        }

        observer.onCompleted();

        putFuture.awaitDone(timeoutMs);

        return new PutFileResponse();
    }

    @Override
    public GetFileResponse getFile(GetFileRequest request, int timeoutMs) throws Exception {

        checkParamOfGetFile(request);

        GetFilePipe pipe = new GetFilePipe(request.getFileName());

        stubManager.
            getAsyncStub().
            getFile(
                buildGetFileRequest(
                    request.getStoreName(),
                    request.getFileName(),
                    request.getMetaData()),
                pipe);

        return new GetFileResponse(pipe.getReader());
    }

    @Override
    public ListFileResponse listFile(ListFileRequest request, int timeoutMs) throws Exception {

        checkParamOfListFile(request);

        RuntimeProto.ListFileResp response = stubManager.
            getBlockingStub().
            withDeadlineAfter(timeoutMs, TimeUnit.MILLISECONDS).
            listFile(
                buildListFileRequest(
                    request.getStoreName(),
                    request.getName(),
                    request.getMarker(),
                    request.getPageSize(),
                    request.getMetaData()));

        return buildListFileResponse(response);
    }

    @Override
    public DelFileResponse delFile(DelFileRequest request, int timeoutMs) throws Exception {

        checkParamOfDeleteFile(request);

        stubManager.
            getBlockingStub().
            withDeadlineAfter(timeoutMs, TimeUnit.MILLISECONDS).
            delFile(
                buildDelFileRequest(request.getStoreName(), request.getFileName(), request.getMetaData()));

        return new DelFileResponse();
    }

    @Override
    public GetMeteResponse getFileMeta(GetMetaRequest request, int timeoutMs) throws Exception {

        checkParamOfGetFileMeta(request);

        RuntimeProto.GetFileMetaResponse resp = stubManager.
            getBlockingStub().
            withDeadlineAfter(timeoutMs, TimeUnit.MILLISECONDS).
            getFileMeta(
                buildGetFileMetaRequest(request.getStoreName(), request.getFileName(), request.getMetaData()));

        return buildGetFileMetaResponse(resp);
    }

    private void checkParamOfGetFile(GetFileRequest request) {

        // check request
        if (request == null) {
            throw new IllegalArgumentException("miss request");
        }

        // check store name
        if (request.getStoreName() == null) {
            throw new IllegalArgumentException("miss store name");
        }

        // check file name
        if (request.getFileName() == null) {
            throw new IllegalArgumentException("miss file name");
        }
    }

    private void checkParamOfPutFile(PutFileRequest request) {

        // check request
        if (request == null) {
            throw new IllegalArgumentException("miss request");
        }

        // check store name
        if (request.getStoreName() == null) {
            throw new IllegalArgumentException("miss store name");
        }

        // check file name
        if (request.getFileName() == null) {
            throw new IllegalArgumentException("miss file name");
        }

        // check input stream
        if (request.getIn() == null) {
            throw new IllegalArgumentException("miss file stream");
        }
    }

    private void checkParamOfListFile(ListFileRequest request) {

        // check request
        if (request == null) {
            throw new IllegalArgumentException("miss request");
        }

        // check store name
        if (request.getStoreName() == null) {
            throw new IllegalArgumentException("miss store name");
        }
    }

    private void checkParamOfDeleteFile(DelFileRequest request) {

        // check request
        if (request == null) {
            throw new IllegalArgumentException("miss request");
        }

        // check file name
        if (request.getFileName() == null) {
            throw new IllegalArgumentException("miss file name");
        }

        // check store name
        if (request.getStoreName() == null) {
            throw new IllegalArgumentException("miss store name");
        }
    }

    private void checkParamOfGetFileMeta(GetMetaRequest request) {

        // check request
        if (request == null) {
            throw new IllegalArgumentException("miss request");
        }

        // check store name
        if (request.getStoreName() == null) {
            throw new IllegalArgumentException("miss store name");
        }

        // check file name
        if (request.getFileName() == null) {
            throw new IllegalArgumentException("miss file name");
        }
    }

    private class PutFileFuture implements StreamObserver<Empty> {

        private final String         fileName;
        private final CountDownLatch latch;

        private volatile Throwable   t;

        PutFileFuture(String fileName) {
            this.fileName = fileName;
            this.latch = new CountDownLatch(1);
        }

        @Override
        public void onNext(Empty value) {
            logger.info(String.format("put File %s successfully", this.fileName));
        }

        @Override
        public void onError(Throwable t) {
            logger.error(String.format("put File error, file=%s", this.fileName), t);
            this.t = t;
            this.latch.countDown();
        }

        @Override
        public void onCompleted() {
            logger.info(String.format("put File %s complete", this.fileName));
            latch.countDown();
        }

        public void awaitDone(int timeoutMs) throws Exception {

            boolean finished = latch.await(timeoutMs, TimeUnit.MILLISECONDS);
            if (!finished) {
                String tip = String.format("put file timeout, file=%s", fileName);
                throw new RuntimeClientException("PUT_FILE", tip);
            }

            // do not wrap for grpc Exception
            if (t instanceof StatusRuntimeException) {
                throw (StatusRuntimeException) t;
            }

            // wrap exception for non grpc Exception
            if (t != null) {
                throw new RuntimeClientException(t);
            }
        }
    }

    private class PipeFileInputStream extends PipedInputStream {

        private volatile Throwable cause;

        PipeFileInputStream(PipedOutputStream out) throws IOException {
            super(out);
        }

        @Override
        public synchronized int read() throws IOException {

            checkCause();

            return super.read();
        }

        @Override
        public int read(byte[] b) throws IOException {

            checkCause();

            return super.read(b);
        }

        @Override
        public synchronized int read(byte[] b, int off, int len) throws IOException {

            checkCause();

            return super.read(b, off, len);
        }

        private void checkCause() throws IOException {
            if (this.cause != null) {
                this.close();
                throw new IOException(this.cause);
            }
        }

        public void setCause(Throwable cause) {
            this.cause = cause;
        }
    }

    private class GetFilePipe implements StreamObserver<RuntimeProto.GetFileResponse> {

        private final String              fileName;
        private final PipeFileInputStream reader;
        private final PipedOutputStream   dataSource;

        GetFilePipe(String fileName) throws IOException {
            this.fileName = fileName;
            this.dataSource = new PipedOutputStream();
            this.reader = new PipeFileInputStream(this.dataSource);
        }

        @Override
        public void onNext(RuntimeProto.GetFileResponse value) {

            logger.info(String.format("get File %s successfully", this.fileName));

            pipe(value.getData().toByteArray());
        }

        @Override
        public void onError(Throwable t) {

            logger.error(String.format("get File error, file=%s", this.fileName), t);

            reader.setCause(t);

            pipe(t.getMessage().getBytes());
        }

        @Override
        public void onCompleted() {

            logger.info(String.format("get File %s complete", this.fileName));

            close();
        }

        private void pipe(byte[] data) {
            try {
                dataSource.write(data);
            } catch (IOException e) {
                logger.error("get File transform err", e);
            }
        }

        private void close() {
            try {
                dataSource.close();
            } catch (IOException e) {
                logger.error("get File close stream err", e);
            }
        }

        public InputStream getReader() {
            return this.reader;
        }
    }

    private StreamObserver<RuntimeProto.PutFileRequest> createPutFileObserver(
                                                                              StreamObserver<Empty> callBackObserver,
                                                                              int timeoutMs) {

        return stubManager.
            getAsyncStub().
            withDeadlineAfter(timeoutMs, TimeUnit.MILLISECONDS).
            putFile(callBackObserver);
    }

    private RuntimeProto.PutFileRequest buildPutFileMetaDataRequest(String storeName,
                                                                    String fileName,
                                                                    Map<String, String> meta) {
        return RuntimeProto.PutFileRequest.
            newBuilder().
            setStoreName(storeName).
            setName(fileName).
            putAllMetadata(meta).
            build();
    }

    private RuntimeProto.PutFileRequest buildPutFileDataRequest(byte[] bytes, int size) {

        return RuntimeProto.PutFileRequest.
            newBuilder().
            setData(ByteString.copyFrom(bytes, 0, size)).
            build();
    }

    private RuntimeProto.GetFileRequest buildGetFileRequest(String storeName,
                                                            String fileName,
                                                            Map<String, String> meta) {

        return RuntimeProto.GetFileRequest.
            newBuilder().
            setStoreName(storeName).
            setName(fileName).
            putAllMetadata(meta).
            build();
    }

    private RuntimeProto.ListFileRequest buildListFileRequest(String storeName, String name, String marker,
                                                              int pageSize, Map<String, String> meta) {

        RuntimeProto.FileRequest fileRequest = RuntimeProto.FileRequest.
            newBuilder().
            setStoreName(storeName).
            setName(name).
            putAllMetadata(meta).
            build();

        return RuntimeProto.ListFileRequest.
            newBuilder().
            setRequest(fileRequest).
            setMarker(marker).
            setPageSize(pageSize).
            build();
    }

    private RuntimeProto.DelFileRequest buildDelFileRequest(String storeName,
                                                            String fileName,
                                                            Map<String, String> meta) {

        RuntimeProto.FileRequest fileRequest = RuntimeProto.FileRequest.
            newBuilder().
            setStoreName(storeName).
            setName(fileName).
            putAllMetadata(meta).
            build();

        return RuntimeProto.DelFileRequest.
            newBuilder().
            setRequest(fileRequest).
            build();
    }

    private RuntimeProto.GetFileMetaRequest buildGetFileMetaRequest(String storeName,
                                                                    String fileName,
                                                                    Map<String, String> meta) {

        RuntimeProto.FileRequest fileRequest = RuntimeProto.FileRequest.
            newBuilder().
            setStoreName(storeName).
            setName(fileName).
            putAllMetadata(meta).
            build();

        return RuntimeProto.GetFileMetaRequest.
            newBuilder().
            setRequest(fileRequest).
            build();
    }

    private GetMeteResponse buildGetFileMetaResponse(RuntimeProto.GetFileMetaResponse resp) {

        Map<String, String[]> metas = new HashMap<>();
        resp.getResponse().
                getMetadataMap().
                forEach(
                        (s, fileMetaValue) ->
                                metas.put(s, fileMetaValue.getValueList().toArray(new String[0])));

        GetMeteResponse result = new GetMeteResponse();
        result.setSize(resp.getSize());
        result.setLastModified(resp.getLastModified());
        result.setMeta(metas);

        return result;
    }

    private ListFileResponse buildListFileResponse(RuntimeProto.ListFileResp resp) {

        FileInfo[] files = resp.getFilesList().
                stream().
                map(
                        fileInfo ->
                                new FileInfo(
                                        fileInfo.getFileName(),
                                        fileInfo.getSize(),
                                        fileInfo.getLastModified(),
                                        fileInfo.getMetadataMap())).
                toArray(FileInfo[]::new);

        ListFileResponse result = new ListFileResponse();
        result.setFiles(files);
        result.setTruncated(resp.getIsTruncated());
        result.setMarker(resp.getMarker());

        return result;
    }
}
