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
package io.mosn.layotto.v1.client.reactor;

import com.google.common.base.Strings;
import com.google.protobuf.Any;
import com.google.protobuf.ByteString;
import com.google.protobuf.Empty;
import io.grpc.CallOptions;
import io.grpc.Channel;
import io.grpc.ClientCall;
import io.grpc.ClientInterceptor;
import io.grpc.ForwardingClientCall;
import io.grpc.Metadata;
import io.grpc.MethodDescriptor;
import io.grpc.stub.StreamObserver;
import io.mosn.layotto.v1.config.Properties;
import io.mosn.layotto.v1.exceptions.LayottoException;
import io.mosn.layotto.v1.serializer.LayottoObjectSerializer;
import io.mosn.layotto.v1.utils.GrpcWrapper;
import io.mosn.layotto.v1.utils.NetworkUtils;
import io.mosn.layotto.v1.value.Headers;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;
import reactor.core.publisher.MonoSink;
import reactor.util.context.Context;
import spec.proto.runtime.v1.RuntimeGrpc;
import spec.proto.runtime.v1.RuntimeProto;
import spec.sdk.reactor.v1.domain.core.configuration.ConfigurationItem;
import spec.sdk.reactor.v1.domain.core.configuration.ConfigurationRequestItem;
import spec.sdk.reactor.v1.domain.core.configuration.SaveConfigurationRequest;
import spec.sdk.reactor.v1.domain.core.configuration.SubConfigurationResp;
import spec.sdk.reactor.v1.domain.core.invocation.HttpExtension;
import spec.sdk.reactor.v1.domain.core.invocation.InvokeMethodRequest;
import spec.sdk.reactor.v1.domain.core.pubsub.PublishEventRequest;
import spec.sdk.reactor.v1.domain.core.state.DeleteStateRequest;
import spec.sdk.reactor.v1.domain.core.state.ExecuteStateTransactionRequest;
import spec.sdk.reactor.v1.domain.core.state.GetBulkStateRequest;
import spec.sdk.reactor.v1.domain.core.state.GetStateRequest;
import spec.sdk.reactor.v1.domain.core.state.SaveStateRequest;
import spec.sdk.reactor.v1.domain.core.state.State;
import spec.sdk.reactor.v1.domain.core.state.StateOptions;
import spec.sdk.reactor.v1.domain.core.state.TransactionalStateOperation;
import spec.sdk.reactor.v1.utils.TypeRef;

import java.io.Closeable;
import java.io.IOException;
import java.util.List;
import java.util.Map;
import java.util.concurrent.ExecutionException;
import java.util.function.Consumer;
import java.util.stream.Collectors;

public class LayottoReactorClientGrpc extends AbstractLayottoReactorClient {

    /**
     * The GRPC managed channel to be used.
     */
    private final Closeable               channel;

    /**
     * The async gRPC stub.
     */
    private final RuntimeGrpc.RuntimeStub asyncStub;

    /**
     * Default access level constructor, in order to create an instance of this class.
     *
     * @param closeableChannel A closeable for a Managed GRPC channel
     * @param asyncStub        async gRPC stub
     */
    LayottoReactorClientGrpc(LayottoObjectSerializer objectSerializer,
                             LayottoObjectSerializer stateSerializer,
                             Closeable closeableChannel,
                             RuntimeGrpc.RuntimeStub asyncStub) {
        super(objectSerializer, stateSerializer);
        this.channel = closeableChannel;
        this.asyncStub = intercept(asyncStub);
    }

    @Override
    public <T> Mono<List<ConfigurationItem<T>>> getConfiguration(ConfigurationRequestItem configurationRequestItem,
                                                                 TypeRef<T> type) {
        // TODO: 2021/9/26
        return null;
    }

    @Override
    public Mono<Void> saveConfiguration(SaveConfigurationRequest saveConfigurationRequest) {
        // TODO: 2021/9/26
        return null;
    }

    @Override
    public Mono<Void> deleteConfiguration(ConfigurationRequestItem configurationRequestItem) {
        // TODO: 2021/9/26
        return null;
    }

    @Override
    public <T> Flux<SubConfigurationResp<T>> subscribeConfiguration(ConfigurationRequestItem configurationRequestItem,
                                                                    TypeRef<T> type) {
        // TODO: 2021/9/26
        return null;
    }

    @Override
    public <T> Mono<T> invokeMethod(InvokeMethodRequest invokeMethodRequest, TypeRef<T> type) {
        try {
            String appId = invokeMethodRequest.getAppId();
            String method = invokeMethodRequest.getMethod();
            Object body = invokeMethodRequest.getBody();
            HttpExtension httpExtension = invokeMethodRequest.getHttpExtension();
            RuntimeProto.InvokeServiceRequest envelope = this.buildInvokeServiceRequest(
                    httpExtension,
                    appId,
                    method,
                    body);
            // Regarding missing metadata in method invocation for gRPC:
            // gRPC to gRPC does not handle metadata in Layotto runtime proto.
            // gRPC to HTTP does not map correctly in Layotto runtime as per https://github.com/layotto/layotto/issues/2342

            return Mono.subscriberContext().flatMap(
                    context -> this.<RuntimeProto.InvokeResponse>createMono(
                            it -> intercept(context, asyncStub).invokeService(envelope, it)
                    )
            ).flatMap(
                    it -> {
                        try {
                            return Mono.justOrEmpty(objectSerializer.deserialize(it.getData().getValue().toByteArray(), type));
                        } catch (IOException e) {
                            throw LayottoException.propagate(e);
                        }
                    }
            );
        } catch (Exception ex) {
            return LayottoException.wrapMono(ex);
        }
    }

    @Override
    public Mono<Void> publishEvent(PublishEventRequest request) {
        try {
            String pubsubName = request.getPubsubName();
            String topic = request.getTopic();
            Object data = request.getData();
            RuntimeProto.PublishEventRequest.Builder envelopeBuilder = RuntimeProto.PublishEventRequest.newBuilder()
                    .setTopic(topic)
                    .setPubsubName(pubsubName)
                    .setData(ByteString.copyFrom(objectSerializer.serialize(data)));

            // Content-type can be overwritten on a per-request basis.
            // It allows CloudEvents to be handled differently, for example.
            String contentType = request.getContentType();
            if (contentType == null || contentType.isEmpty()) {
                contentType = objectSerializer.getContentType();
            }
            envelopeBuilder.setDataContentType(contentType);

            Map<String, String> metadata = request.getMetadata();
            if (metadata != null) {
                envelopeBuilder.putAllMetadata(metadata);
            }

            return Mono.subscriberContext().flatMap(
                    context ->
                            this.<Empty>createMono(
                                    it -> intercept(context, asyncStub).publishEvent(envelopeBuilder.build(), it)
                            )
            ).then();
        } catch (Exception ex) {
            return LayottoException.wrapMono(ex);
        }
    }

    @Override
    public <T> Mono<State<T>> getState(GetStateRequest request, TypeRef<T> type) {
        try {
            final String stateStoreName = request.getStoreName();
            final String key = request.getKey();
            final StateOptions options = request.getStateOptions();
            final Map<String, String> metadata = request.getMetadata();

            if ((stateStoreName == null) || (stateStoreName.trim().isEmpty())) {
                throw new IllegalArgumentException("State store name cannot be null or empty.");
            }
            if ((key == null) || (key.trim().isEmpty())) {
                throw new IllegalArgumentException("Key cannot be null or empty.");
            }
            RuntimeProto.GetStateRequest.Builder builder = RuntimeProto.GetStateRequest.newBuilder()
                    .setStoreName(stateStoreName)
                    .setKey(key);
            if (metadata != null) {
                builder.putAllMetadata(metadata);
            }
            if (options != null && options.getConsistency() != null) {
                builder.setConsistency(getGrpcStateConsistency(options));
            }

            RuntimeProto.GetStateRequest envelope = builder.build();

            return Mono.subscriberContext().flatMap(
                    context ->
                            this.<RuntimeProto.GetStateResponse>createMono(
                                    it -> intercept(context, asyncStub).getState(envelope, it)
                            )
            ).map(
                    it -> {
                        try {
                            return buildStateKeyValue(it, key, options, type);
                        } catch (IOException ex) {
                            throw LayottoException.propagate(ex);
                        }
                    }
            );
        } catch (Exception ex) {
            return LayottoException.wrapMono(ex);
        }
    }

    @Override
    public <T> Mono<List<State<T>>> getBulkState(GetBulkStateRequest request, TypeRef<T> type) {
        try {
            final String stateStoreName = request.getStoreName();
            final List<String> keys = request.getKeys();
            final int parallelism = request.getParallelism();
            final Map<String, String> metadata = request.getMetadata();
            if ((stateStoreName == null) || (stateStoreName.trim().isEmpty())) {
                throw new IllegalArgumentException("State store name cannot be null or empty.");
            }
            if (keys == null || keys.isEmpty()) {
                throw new IllegalArgumentException("Key cannot be null or empty.");
            }

            if (parallelism < 0) {
                throw new IllegalArgumentException("Parallelism cannot be negative.");
            }
            RuntimeProto.GetBulkStateRequest.Builder builder = RuntimeProto.GetBulkStateRequest.newBuilder()
                    .setStoreName(stateStoreName)
                    .addAllKeys(keys)
                    .setParallelism(parallelism);
            if (metadata != null) {
                builder.putAllMetadata(metadata);
            }

            RuntimeProto.GetBulkStateRequest envelope = builder.build();

            return Mono.subscriberContext().flatMap(
                    context -> this.<RuntimeProto.GetBulkStateResponse>createMono(it -> intercept(context, asyncStub)
                            .getBulkState(envelope, it)
                    )
            ).map(
                    it ->
                            it
                                    .getItemsList()
                                    .stream()
                                    .map(b -> {
                                        try {
                                            return buildStateKeyValue(b, type);
                                        } catch (Exception e) {
                                            throw LayottoException.propagate(e);
                                        }
                                    })
                                    .collect(Collectors.toList())
            );
        } catch (Exception ex) {
            return LayottoException.wrapMono(ex);
        }
    }

    @Override
    public Mono<Void> executeStateTransaction(ExecuteStateTransactionRequest request) {
        try {
            final String stateStoreName = request.getStateStoreName();
            final List<TransactionalStateOperation<?>> operations = request.getOperations();
            final Map<String, String> metadata = request.getMetadata();
            if ((stateStoreName == null) || (stateStoreName.trim().isEmpty())) {
                throw new IllegalArgumentException("State store name cannot be null or empty.");
            }
            RuntimeProto.ExecuteStateTransactionRequest.Builder builder = RuntimeProto.ExecuteStateTransactionRequest
                    .newBuilder();
            builder.setStoreName(stateStoreName);
            if (metadata != null) {
                builder.putAllMetadata(metadata);
            }
            for (TransactionalStateOperation<?> operation : operations) {
                RuntimeProto.TransactionalStateOperation.Builder operationBuilder = RuntimeProto.TransactionalStateOperation
                        .newBuilder();
                operationBuilder.setOperationType(operation.getOperation().toString().toLowerCase());
                operationBuilder.setRequest(buildStateRequest(operation.getRequest()).build());
                builder.addOperations(operationBuilder.build());
            }
            RuntimeProto.ExecuteStateTransactionRequest req = builder.build();

            return Mono.subscriberContext().flatMap(
                    context -> this.<Empty>createMono(it -> intercept(context, asyncStub).executeStateTransaction(req, it))
            ).then();
        } catch (Exception e) {
            return LayottoException.wrapMono(e);
        }
    }

    @Override
    public Mono<Void> saveBulkState(SaveStateRequest request) {
        try {
            final String stateStoreName = request.getStoreName();
            final List<State<?>> states = request.getStates();
            if ((stateStoreName == null) || (stateStoreName.trim().isEmpty())) {
                throw new IllegalArgumentException("State store name cannot be null or empty.");
            }
            RuntimeProto.SaveStateRequest.Builder builder = RuntimeProto.SaveStateRequest.newBuilder();
            builder.setStoreName(stateStoreName);
            for (State<?> state : states) {
                builder.addStates(buildStateRequest(state).build());
            }
            RuntimeProto.SaveStateRequest req = builder.build();

            return Mono.subscriberContext().flatMap(
                    context -> this.<Empty>createMono(it -> intercept(context, asyncStub).saveState(req, it))
            ).then();
        } catch (Exception ex) {
            return LayottoException.wrapMono(ex);
        }
    }

    @Override
    public Mono<Void> deleteState(DeleteStateRequest request) {
        try {
            final String stateStoreName = request.getStateStoreName();
            final String key = request.getKey();
            final StateOptions options = request.getStateOptions();
            final String etag = request.getEtag();
            final Map<String, String> metadata = request.getMetadata();

            if ((stateStoreName == null) || (stateStoreName.trim().isEmpty())) {
                throw new IllegalArgumentException("State store name cannot be null or empty.");
            }
            if ((key == null) || (key.trim().isEmpty())) {
                throw new IllegalArgumentException("Key cannot be null or empty.");
            }

            RuntimeProto.StateOptions.Builder optionBuilder = null;
            if (options != null) {
                optionBuilder = RuntimeProto.StateOptions.newBuilder();
                if (options.getConcurrency() != null) {
                    optionBuilder.setConcurrency(getGrpcStateConcurrency(options));
                }
                if (options.getConsistency() != null) {
                    optionBuilder.setConsistency(getGrpcStateConsistency(options));
                }
            }
            RuntimeProto.DeleteStateRequest.Builder builder = RuntimeProto.DeleteStateRequest.newBuilder()
                    .setStoreName(stateStoreName)
                    .setKey(key);
            if (metadata != null) {
                builder.putAllMetadata(metadata);
            }
            if (etag != null) {
                builder.setEtag(RuntimeProto.Etag.newBuilder().setValue(etag).build());
            }

            if (optionBuilder != null) {
                builder.setOptions(optionBuilder.build());
            }

            RuntimeProto.DeleteStateRequest req = builder.build();

            return Mono.subscriberContext().flatMap(
                    context -> this.<Empty>createMono(it -> intercept(context, asyncStub).deleteState(req, it))
            ).then();
        } catch (Exception ex) {
            return LayottoException.wrapMono(ex);
        }
    }

    /**
     * Builds the object io.layotto.{@link RuntimeProto.InvokeServiceRequest} to be send based on the parameters.
     *
     * @param httpExtension Object for HttpExtension
     * @param appId         The application id to be invoked
     * @param method        The application method to be invoked
     * @param body          The body of the request to be send as part of the invocation
     * @param <K>           The Type of the Body
     * @return The object to be sent as part of the invocation.
     * @throws IOException If there's an issue serializing the request.
     */
    private <K> RuntimeProto.InvokeServiceRequest buildInvokeServiceRequest(
                                                                            HttpExtension httpExtension,
                                                                            String appId,
                                                                            String method,
                                                                            K body) throws IOException {
        if (httpExtension == null) {
            throw new IllegalArgumentException("HttpExtension cannot be null. Use HttpExtension.NONE instead.");
        }
        RuntimeProto.CommonInvokeRequest.Builder requestBuilder = RuntimeProto.CommonInvokeRequest.newBuilder();
        requestBuilder.setMethod(method);
        if (body != null) {
            byte[] byteRequest = objectSerializer.serialize(body);
            Any data = Any.newBuilder().setValue(ByteString.copyFrom(byteRequest)).build();
            requestBuilder.setData(data);
        } else {
            requestBuilder.setData(Any.newBuilder().build());
        }
        RuntimeProto.HTTPExtension.Builder httpExtensionBuilder = RuntimeProto.HTTPExtension.newBuilder();

        httpExtensionBuilder.setVerb(RuntimeProto.HTTPExtension.Verb.valueOf(httpExtension.getMethod().toString()))
            .setQuerystring(httpExtension.encodeQueryString());
        requestBuilder.setHttpExtension(httpExtensionBuilder.build());

        requestBuilder.setContentType(objectSerializer.getContentType());

        RuntimeProto.InvokeServiceRequest.Builder envelopeBuilder = RuntimeProto.InvokeServiceRequest.newBuilder()
            .setId(appId)
            .setMessage(requestBuilder.build());
        return envelopeBuilder.build();
    }

    private <T> State<T> buildStateKeyValue(
            RuntimeProto.BulkStateItem item,
            TypeRef<T> type) throws IOException {
        String key = item.getKey();
        String error = item.getError();
        if (!Strings.isNullOrEmpty(error)) {
            return new State<>(key, error);
        }

        ByteString payload = item.getData();
        byte[] data = payload == null ? null : payload.toByteArray();
        T value = stateSerializer.deserialize(data, type);
        String etag = item.getEtag();
        if (etag.equals("")) {
            etag = null;
        }
        return new State<>(key, value, etag, item.getMetadataMap(), null);
    }

    private <T> State<T> buildStateKeyValue(
            RuntimeProto.GetStateResponse response,
            String requestedKey,
            StateOptions stateOptions,
            TypeRef<T> type) throws IOException {
        ByteString payload = response.getData();
        byte[] data = payload == null ? null : payload.toByteArray();
        T value = stateSerializer.deserialize(data, type);
        String etag = response.getEtag();
        if (etag.equals("")) {
            etag = null;
        }
        return new State<>(requestedKey, value, etag, response.getMetadataMap(), stateOptions);
    }

    private <T> RuntimeProto.StateItem.Builder buildStateRequest(State<T> state) throws IOException {
        byte[] bytes = stateSerializer.serialize(state.getValue());

        RuntimeProto.StateItem.Builder stateBuilder = RuntimeProto.StateItem.newBuilder();
        if (state.getEtag() != null) {
            stateBuilder.setEtag(RuntimeProto.Etag.newBuilder().setValue(state.getEtag()).build());
        }
        if (state.getMetadata() != null) {
            stateBuilder.putAllMetadata(state.getMetadata());
        }
        if (bytes != null) {
            stateBuilder.setValue(ByteString.copyFrom(bytes));
        }
        stateBuilder.setKey(state.getKey());
        RuntimeProto.StateOptions.Builder optionBuilder = null;
        if (state.getOptions() != null) {
            StateOptions options = state.getOptions();
            optionBuilder = RuntimeProto.StateOptions.newBuilder();
            if (options.getConcurrency() != null) {
                optionBuilder.setConcurrency(getGrpcStateConcurrency(options));
            }
            if (options.getConsistency() != null) {
                optionBuilder.setConsistency(getGrpcStateConsistency(options));
            }
        }
        if (optionBuilder != null) {
            stateBuilder.setOptions(optionBuilder.build());
        }
        return stateBuilder;
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

    // -- Lifecycle Functions

    @Override
    public Mono<Void> waitForSidecar(int timeoutInMilliseconds) {
        return Mono.fromRunnable(() -> {
            try {
                NetworkUtils.waitForSocket(Properties.SIDECAR_IP.get(), Properties.GRPC_PORT.get(), timeoutInMilliseconds);
            } catch (InterruptedException e) {
                throw new RuntimeException(e);
            }
        });
    }

    @Override
    public Mono<Void> shutdown() {
        return Mono.subscriberContext()
            // FIXME: 2021/9/26 Refer to Dapr
            // .flatMap(context ->
            //     this.<Empty>createMono(it ->
            //         intercept(context, asyncStub)
            //                 .shutdown(Empty.getDefaultInstance(), it)))
            .then();
    }

    private <T> Mono<T> createMono(Consumer<StreamObserver<T>> consumer) {
        return Mono.create(sink ->
                LayottoException
                        .wrap(() -> consumer.accept(createStreamObserver(sink)))
                        .run());
    }

    private <T> StreamObserver<T> createStreamObserver(MonoSink<T> sink) {
        return new StreamObserver<T>() {
            @Override
            public void onNext(T value) {
                sink.success(value);
            }

            @Override
            public void onError(Throwable t) {
                sink.error(LayottoException.propagate(new ExecutionException(t)));
            }

            @Override
            public void onCompleted() {
                sink.success();
            }
        };
    }

    /**
     * Closes the ManagedChannel for GRPC.
     *
     * @throws IOException on exception.
     * @see io.grpc.ManagedChannel#shutdown()
     */
    @Override
    public void close() throws Exception {
        if (channel != null) {
            LayottoException
                    .wrap(() -> {
                        channel.close();
                        return true;
                    })
                    .call();
        }
    }

    /**
     * Populates GRPC client with interceptors for telemetry.
     *
     * @param context Reactor's context.
     * @param client  GRPC client for Layotto.
     * @return Client after adding interceptors.
     */
    private static RuntimeGrpc.RuntimeStub intercept(Context context, RuntimeGrpc.RuntimeStub client) {
        return GrpcWrapper.intercept(context, client);
    }

    /**
     * Populates GRPC client with interceptors.
     *
     * @param client GRPC client for Layotto.
     * @return Client after adding interceptors.
     */
    private static RuntimeGrpc.RuntimeStub intercept(RuntimeGrpc.RuntimeStub client) {
        ClientInterceptor interceptor = new ClientInterceptor() {
            @Override
            public <ReqT, RespT> ClientCall<ReqT, RespT> interceptCall(MethodDescriptor<ReqT, RespT> methodDescriptor,
                                                                       CallOptions callOptions,
                                                                       Channel channel) {
                ClientCall<ReqT, RespT> clientCall = channel.newCall(methodDescriptor, callOptions);
                return new ForwardingClientCall.SimpleForwardingClientCall<ReqT, RespT>(clientCall) {
                    @Override
                    public void start(final Listener<RespT> responseListener, final Metadata metadata) {
                        String layottoApiToken = Properties.API_TOKEN.get();
                        if (layottoApiToken != null) {
                            metadata.put(Metadata.Key.of(Headers.DAPR_API_TOKEN, Metadata.ASCII_STRING_MARSHALLER),
                                layottoApiToken);
                        }
                        super.start(responseListener, metadata);
                    }
                };
            }
        };
        return client.withInterceptors(interceptor);
    }
}
