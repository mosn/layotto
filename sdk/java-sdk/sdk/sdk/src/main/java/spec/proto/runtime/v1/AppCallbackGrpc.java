package spec.proto.runtime.v1;

import static io.grpc.MethodDescriptor.generateFullMethodName;

/**
 * <pre>
 * AppCallback V1 allows user application to interact with runtime.
 * User application needs to implement AppCallback service if it needs to
 * receive message from runtime.
 * </pre>
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.37.0)",
    comments = "Source: proto/runtime/v1/appcallback.proto")
public final class AppCallbackGrpc {

  private AppCallbackGrpc() {}

  public static final String SERVICE_NAME = "spec.proto.runtime.v1.AppCallback";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<com.google.protobuf.Empty,
      spec.proto.runtime.v1.AppCallbackProto.ListTopicSubscriptionsResponse> getListTopicSubscriptionsMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "ListTopicSubscriptions",
      requestType = com.google.protobuf.Empty.class,
      responseType = spec.proto.runtime.v1.AppCallbackProto.ListTopicSubscriptionsResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<com.google.protobuf.Empty,
      spec.proto.runtime.v1.AppCallbackProto.ListTopicSubscriptionsResponse> getListTopicSubscriptionsMethod() {
    io.grpc.MethodDescriptor<com.google.protobuf.Empty, spec.proto.runtime.v1.AppCallbackProto.ListTopicSubscriptionsResponse> getListTopicSubscriptionsMethod;
    if ((getListTopicSubscriptionsMethod = AppCallbackGrpc.getListTopicSubscriptionsMethod) == null) {
      synchronized (AppCallbackGrpc.class) {
        if ((getListTopicSubscriptionsMethod = AppCallbackGrpc.getListTopicSubscriptionsMethod) == null) {
          AppCallbackGrpc.getListTopicSubscriptionsMethod = getListTopicSubscriptionsMethod =
              io.grpc.MethodDescriptor.<com.google.protobuf.Empty, spec.proto.runtime.v1.AppCallbackProto.ListTopicSubscriptionsResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "ListTopicSubscriptions"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.google.protobuf.Empty.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  spec.proto.runtime.v1.AppCallbackProto.ListTopicSubscriptionsResponse.getDefaultInstance()))
              .setSchemaDescriptor(new AppCallbackMethodDescriptorSupplier("ListTopicSubscriptions"))
              .build();
        }
      }
    }
    return getListTopicSubscriptionsMethod;
  }

  private static volatile io.grpc.MethodDescriptor<spec.proto.runtime.v1.AppCallbackProto.TopicEventRequest,
      spec.proto.runtime.v1.AppCallbackProto.TopicEventResponse> getOnTopicEventMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "OnTopicEvent",
      requestType = spec.proto.runtime.v1.AppCallbackProto.TopicEventRequest.class,
      responseType = spec.proto.runtime.v1.AppCallbackProto.TopicEventResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<spec.proto.runtime.v1.AppCallbackProto.TopicEventRequest,
      spec.proto.runtime.v1.AppCallbackProto.TopicEventResponse> getOnTopicEventMethod() {
    io.grpc.MethodDescriptor<spec.proto.runtime.v1.AppCallbackProto.TopicEventRequest, spec.proto.runtime.v1.AppCallbackProto.TopicEventResponse> getOnTopicEventMethod;
    if ((getOnTopicEventMethod = AppCallbackGrpc.getOnTopicEventMethod) == null) {
      synchronized (AppCallbackGrpc.class) {
        if ((getOnTopicEventMethod = AppCallbackGrpc.getOnTopicEventMethod) == null) {
          AppCallbackGrpc.getOnTopicEventMethod = getOnTopicEventMethod =
              io.grpc.MethodDescriptor.<spec.proto.runtime.v1.AppCallbackProto.TopicEventRequest, spec.proto.runtime.v1.AppCallbackProto.TopicEventResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "OnTopicEvent"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  spec.proto.runtime.v1.AppCallbackProto.TopicEventRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  spec.proto.runtime.v1.AppCallbackProto.TopicEventResponse.getDefaultInstance()))
              .setSchemaDescriptor(new AppCallbackMethodDescriptorSupplier("OnTopicEvent"))
              .build();
        }
      }
    }
    return getOnTopicEventMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static AppCallbackStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<AppCallbackStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<AppCallbackStub>() {
        @java.lang.Override
        public AppCallbackStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new AppCallbackStub(channel, callOptions);
        }
      };
    return AppCallbackStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static AppCallbackBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<AppCallbackBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<AppCallbackBlockingStub>() {
        @java.lang.Override
        public AppCallbackBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new AppCallbackBlockingStub(channel, callOptions);
        }
      };
    return AppCallbackBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static AppCallbackFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<AppCallbackFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<AppCallbackFutureStub>() {
        @java.lang.Override
        public AppCallbackFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new AppCallbackFutureStub(channel, callOptions);
        }
      };
    return AppCallbackFutureStub.newStub(factory, channel);
  }

  /**
   * <pre>
   * AppCallback V1 allows user application to interact with runtime.
   * User application needs to implement AppCallback service if it needs to
   * receive message from runtime.
   * </pre>
   */
  public static abstract class AppCallbackImplBase implements io.grpc.BindableService {

    /**
     * <pre>
     * Lists all topics subscribed by this app.
     * </pre>
     */
    public void listTopicSubscriptions(com.google.protobuf.Empty request,
        io.grpc.stub.StreamObserver<spec.proto.runtime.v1.AppCallbackProto.ListTopicSubscriptionsResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getListTopicSubscriptionsMethod(), responseObserver);
    }

    /**
     * <pre>
     * Subscribes events from Pubsub
     * </pre>
     */
    public void onTopicEvent(spec.proto.runtime.v1.AppCallbackProto.TopicEventRequest request,
        io.grpc.stub.StreamObserver<spec.proto.runtime.v1.AppCallbackProto.TopicEventResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getOnTopicEventMethod(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getListTopicSubscriptionsMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                com.google.protobuf.Empty,
                spec.proto.runtime.v1.AppCallbackProto.ListTopicSubscriptionsResponse>(
                  this, METHODID_LIST_TOPIC_SUBSCRIPTIONS)))
          .addMethod(
            getOnTopicEventMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                spec.proto.runtime.v1.AppCallbackProto.TopicEventRequest,
                spec.proto.runtime.v1.AppCallbackProto.TopicEventResponse>(
                  this, METHODID_ON_TOPIC_EVENT)))
          .build();
    }
  }

  /**
   * <pre>
   * AppCallback V1 allows user application to interact with runtime.
   * User application needs to implement AppCallback service if it needs to
   * receive message from runtime.
   * </pre>
   */
  public static final class AppCallbackStub extends io.grpc.stub.AbstractAsyncStub<AppCallbackStub> {
    private AppCallbackStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected AppCallbackStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new AppCallbackStub(channel, callOptions);
    }

    /**
     * <pre>
     * Lists all topics subscribed by this app.
     * </pre>
     */
    public void listTopicSubscriptions(com.google.protobuf.Empty request,
        io.grpc.stub.StreamObserver<spec.proto.runtime.v1.AppCallbackProto.ListTopicSubscriptionsResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getListTopicSubscriptionsMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * Subscribes events from Pubsub
     * </pre>
     */
    public void onTopicEvent(spec.proto.runtime.v1.AppCallbackProto.TopicEventRequest request,
        io.grpc.stub.StreamObserver<spec.proto.runtime.v1.AppCallbackProto.TopicEventResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getOnTopicEventMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   * <pre>
   * AppCallback V1 allows user application to interact with runtime.
   * User application needs to implement AppCallback service if it needs to
   * receive message from runtime.
   * </pre>
   */
  public static final class AppCallbackBlockingStub extends io.grpc.stub.AbstractBlockingStub<AppCallbackBlockingStub> {
    private AppCallbackBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected AppCallbackBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new AppCallbackBlockingStub(channel, callOptions);
    }

    /**
     * <pre>
     * Lists all topics subscribed by this app.
     * </pre>
     */
    public spec.proto.runtime.v1.AppCallbackProto.ListTopicSubscriptionsResponse listTopicSubscriptions(com.google.protobuf.Empty request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getListTopicSubscriptionsMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * Subscribes events from Pubsub
     * </pre>
     */
    public spec.proto.runtime.v1.AppCallbackProto.TopicEventResponse onTopicEvent(spec.proto.runtime.v1.AppCallbackProto.TopicEventRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getOnTopicEventMethod(), getCallOptions(), request);
    }
  }

  /**
   * <pre>
   * AppCallback V1 allows user application to interact with runtime.
   * User application needs to implement AppCallback service if it needs to
   * receive message from runtime.
   * </pre>
   */
  public static final class AppCallbackFutureStub extends io.grpc.stub.AbstractFutureStub<AppCallbackFutureStub> {
    private AppCallbackFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected AppCallbackFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new AppCallbackFutureStub(channel, callOptions);
    }

    /**
     * <pre>
     * Lists all topics subscribed by this app.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<spec.proto.runtime.v1.AppCallbackProto.ListTopicSubscriptionsResponse> listTopicSubscriptions(
        com.google.protobuf.Empty request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getListTopicSubscriptionsMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * Subscribes events from Pubsub
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<spec.proto.runtime.v1.AppCallbackProto.TopicEventResponse> onTopicEvent(
        spec.proto.runtime.v1.AppCallbackProto.TopicEventRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getOnTopicEventMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_LIST_TOPIC_SUBSCRIPTIONS = 0;
  private static final int METHODID_ON_TOPIC_EVENT = 1;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final AppCallbackImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(AppCallbackImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_LIST_TOPIC_SUBSCRIPTIONS:
          serviceImpl.listTopicSubscriptions((com.google.protobuf.Empty) request,
              (io.grpc.stub.StreamObserver<spec.proto.runtime.v1.AppCallbackProto.ListTopicSubscriptionsResponse>) responseObserver);
          break;
        case METHODID_ON_TOPIC_EVENT:
          serviceImpl.onTopicEvent((spec.proto.runtime.v1.AppCallbackProto.TopicEventRequest) request,
              (io.grpc.stub.StreamObserver<spec.proto.runtime.v1.AppCallbackProto.TopicEventResponse>) responseObserver);
          break;
        default:
          throw new AssertionError();
      }
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public io.grpc.stub.StreamObserver<Req> invoke(
        io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        default:
          throw new AssertionError();
      }
    }
  }

  private static abstract class AppCallbackBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    AppCallbackBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return spec.proto.runtime.v1.AppCallbackProto.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("AppCallback");
    }
  }

  private static final class AppCallbackFileDescriptorSupplier
      extends AppCallbackBaseDescriptorSupplier {
    AppCallbackFileDescriptorSupplier() {}
  }

  private static final class AppCallbackMethodDescriptorSupplier
      extends AppCallbackBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    AppCallbackMethodDescriptorSupplier(String methodName) {
      this.methodName = methodName;
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.MethodDescriptor getMethodDescriptor() {
      return getServiceDescriptor().findMethodByName(methodName);
    }
  }

  private static volatile io.grpc.ServiceDescriptor serviceDescriptor;

  public static io.grpc.ServiceDescriptor getServiceDescriptor() {
    io.grpc.ServiceDescriptor result = serviceDescriptor;
    if (result == null) {
      synchronized (AppCallbackGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new AppCallbackFileDescriptorSupplier())
              .addMethod(getListTopicSubscriptionsMethod())
              .addMethod(getOnTopicEventMethod())
              .build();
        }
      }
    }
    return result;
  }
}
