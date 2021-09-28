package spec.proto.runtime.v1;

import static io.grpc.MethodDescriptor.generateFullMethodName;

/**
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.37.0)",
    comments = "Source: proto/runtime/v1/runtime.proto")
public final class RuntimeGrpc {

  private RuntimeGrpc() {}

  public static final String SERVICE_NAME = "spec.proto.runtime.v1.Runtime";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<spec.proto.runtime.v1.RuntimeProto.SayHelloRequest,
      spec.proto.runtime.v1.RuntimeProto.SayHelloResponse> getSayHelloMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "SayHello",
      requestType = spec.proto.runtime.v1.RuntimeProto.SayHelloRequest.class,
      responseType = spec.proto.runtime.v1.RuntimeProto.SayHelloResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<spec.proto.runtime.v1.RuntimeProto.SayHelloRequest,
      spec.proto.runtime.v1.RuntimeProto.SayHelloResponse> getSayHelloMethod() {
    io.grpc.MethodDescriptor<spec.proto.runtime.v1.RuntimeProto.SayHelloRequest, spec.proto.runtime.v1.RuntimeProto.SayHelloResponse> getSayHelloMethod;
    if ((getSayHelloMethod = RuntimeGrpc.getSayHelloMethod) == null) {
      synchronized (RuntimeGrpc.class) {
        if ((getSayHelloMethod = RuntimeGrpc.getSayHelloMethod) == null) {
          RuntimeGrpc.getSayHelloMethod = getSayHelloMethod =
              io.grpc.MethodDescriptor.<spec.proto.runtime.v1.RuntimeProto.SayHelloRequest, spec.proto.runtime.v1.RuntimeProto.SayHelloResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "SayHello"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  spec.proto.runtime.v1.RuntimeProto.SayHelloRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  spec.proto.runtime.v1.RuntimeProto.SayHelloResponse.getDefaultInstance()))
              .setSchemaDescriptor(new RuntimeMethodDescriptorSupplier("SayHello"))
              .build();
        }
      }
    }
    return getSayHelloMethod;
  }

  private static volatile io.grpc.MethodDescriptor<spec.proto.runtime.v1.RuntimeProto.InvokeServiceRequest,
      spec.proto.runtime.v1.RuntimeProto.InvokeResponse> getInvokeServiceMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "InvokeService",
      requestType = spec.proto.runtime.v1.RuntimeProto.InvokeServiceRequest.class,
      responseType = spec.proto.runtime.v1.RuntimeProto.InvokeResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<spec.proto.runtime.v1.RuntimeProto.InvokeServiceRequest,
      spec.proto.runtime.v1.RuntimeProto.InvokeResponse> getInvokeServiceMethod() {
    io.grpc.MethodDescriptor<spec.proto.runtime.v1.RuntimeProto.InvokeServiceRequest, spec.proto.runtime.v1.RuntimeProto.InvokeResponse> getInvokeServiceMethod;
    if ((getInvokeServiceMethod = RuntimeGrpc.getInvokeServiceMethod) == null) {
      synchronized (RuntimeGrpc.class) {
        if ((getInvokeServiceMethod = RuntimeGrpc.getInvokeServiceMethod) == null) {
          RuntimeGrpc.getInvokeServiceMethod = getInvokeServiceMethod =
              io.grpc.MethodDescriptor.<spec.proto.runtime.v1.RuntimeProto.InvokeServiceRequest, spec.proto.runtime.v1.RuntimeProto.InvokeResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "InvokeService"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  spec.proto.runtime.v1.RuntimeProto.InvokeServiceRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  spec.proto.runtime.v1.RuntimeProto.InvokeResponse.getDefaultInstance()))
              .setSchemaDescriptor(new RuntimeMethodDescriptorSupplier("InvokeService"))
              .build();
        }
      }
    }
    return getInvokeServiceMethod;
  }

  private static volatile io.grpc.MethodDescriptor<spec.proto.runtime.v1.RuntimeProto.GetConfigurationRequest,
      spec.proto.runtime.v1.RuntimeProto.GetConfigurationResponse> getGetConfigurationMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetConfiguration",
      requestType = spec.proto.runtime.v1.RuntimeProto.GetConfigurationRequest.class,
      responseType = spec.proto.runtime.v1.RuntimeProto.GetConfigurationResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<spec.proto.runtime.v1.RuntimeProto.GetConfigurationRequest,
      spec.proto.runtime.v1.RuntimeProto.GetConfigurationResponse> getGetConfigurationMethod() {
    io.grpc.MethodDescriptor<spec.proto.runtime.v1.RuntimeProto.GetConfigurationRequest, spec.proto.runtime.v1.RuntimeProto.GetConfigurationResponse> getGetConfigurationMethod;
    if ((getGetConfigurationMethod = RuntimeGrpc.getGetConfigurationMethod) == null) {
      synchronized (RuntimeGrpc.class) {
        if ((getGetConfigurationMethod = RuntimeGrpc.getGetConfigurationMethod) == null) {
          RuntimeGrpc.getGetConfigurationMethod = getGetConfigurationMethod =
              io.grpc.MethodDescriptor.<spec.proto.runtime.v1.RuntimeProto.GetConfigurationRequest, spec.proto.runtime.v1.RuntimeProto.GetConfigurationResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetConfiguration"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  spec.proto.runtime.v1.RuntimeProto.GetConfigurationRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  spec.proto.runtime.v1.RuntimeProto.GetConfigurationResponse.getDefaultInstance()))
              .setSchemaDescriptor(new RuntimeMethodDescriptorSupplier("GetConfiguration"))
              .build();
        }
      }
    }
    return getGetConfigurationMethod;
  }

  private static volatile io.grpc.MethodDescriptor<spec.proto.runtime.v1.RuntimeProto.SaveConfigurationRequest,
      com.google.protobuf.Empty> getSaveConfigurationMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "SaveConfiguration",
      requestType = spec.proto.runtime.v1.RuntimeProto.SaveConfigurationRequest.class,
      responseType = com.google.protobuf.Empty.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<spec.proto.runtime.v1.RuntimeProto.SaveConfigurationRequest,
      com.google.protobuf.Empty> getSaveConfigurationMethod() {
    io.grpc.MethodDescriptor<spec.proto.runtime.v1.RuntimeProto.SaveConfigurationRequest, com.google.protobuf.Empty> getSaveConfigurationMethod;
    if ((getSaveConfigurationMethod = RuntimeGrpc.getSaveConfigurationMethod) == null) {
      synchronized (RuntimeGrpc.class) {
        if ((getSaveConfigurationMethod = RuntimeGrpc.getSaveConfigurationMethod) == null) {
          RuntimeGrpc.getSaveConfigurationMethod = getSaveConfigurationMethod =
              io.grpc.MethodDescriptor.<spec.proto.runtime.v1.RuntimeProto.SaveConfigurationRequest, com.google.protobuf.Empty>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "SaveConfiguration"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  spec.proto.runtime.v1.RuntimeProto.SaveConfigurationRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.google.protobuf.Empty.getDefaultInstance()))
              .setSchemaDescriptor(new RuntimeMethodDescriptorSupplier("SaveConfiguration"))
              .build();
        }
      }
    }
    return getSaveConfigurationMethod;
  }

  private static volatile io.grpc.MethodDescriptor<spec.proto.runtime.v1.RuntimeProto.DeleteConfigurationRequest,
      com.google.protobuf.Empty> getDeleteConfigurationMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "DeleteConfiguration",
      requestType = spec.proto.runtime.v1.RuntimeProto.DeleteConfigurationRequest.class,
      responseType = com.google.protobuf.Empty.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<spec.proto.runtime.v1.RuntimeProto.DeleteConfigurationRequest,
      com.google.protobuf.Empty> getDeleteConfigurationMethod() {
    io.grpc.MethodDescriptor<spec.proto.runtime.v1.RuntimeProto.DeleteConfigurationRequest, com.google.protobuf.Empty> getDeleteConfigurationMethod;
    if ((getDeleteConfigurationMethod = RuntimeGrpc.getDeleteConfigurationMethod) == null) {
      synchronized (RuntimeGrpc.class) {
        if ((getDeleteConfigurationMethod = RuntimeGrpc.getDeleteConfigurationMethod) == null) {
          RuntimeGrpc.getDeleteConfigurationMethod = getDeleteConfigurationMethod =
              io.grpc.MethodDescriptor.<spec.proto.runtime.v1.RuntimeProto.DeleteConfigurationRequest, com.google.protobuf.Empty>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "DeleteConfiguration"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  spec.proto.runtime.v1.RuntimeProto.DeleteConfigurationRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.google.protobuf.Empty.getDefaultInstance()))
              .setSchemaDescriptor(new RuntimeMethodDescriptorSupplier("DeleteConfiguration"))
              .build();
        }
      }
    }
    return getDeleteConfigurationMethod;
  }

  private static volatile io.grpc.MethodDescriptor<spec.proto.runtime.v1.RuntimeProto.SubscribeConfigurationRequest,
      spec.proto.runtime.v1.RuntimeProto.SubscribeConfigurationResponse> getSubscribeConfigurationMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "SubscribeConfiguration",
      requestType = spec.proto.runtime.v1.RuntimeProto.SubscribeConfigurationRequest.class,
      responseType = spec.proto.runtime.v1.RuntimeProto.SubscribeConfigurationResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.BIDI_STREAMING)
  public static io.grpc.MethodDescriptor<spec.proto.runtime.v1.RuntimeProto.SubscribeConfigurationRequest,
      spec.proto.runtime.v1.RuntimeProto.SubscribeConfigurationResponse> getSubscribeConfigurationMethod() {
    io.grpc.MethodDescriptor<spec.proto.runtime.v1.RuntimeProto.SubscribeConfigurationRequest, spec.proto.runtime.v1.RuntimeProto.SubscribeConfigurationResponse> getSubscribeConfigurationMethod;
    if ((getSubscribeConfigurationMethod = RuntimeGrpc.getSubscribeConfigurationMethod) == null) {
      synchronized (RuntimeGrpc.class) {
        if ((getSubscribeConfigurationMethod = RuntimeGrpc.getSubscribeConfigurationMethod) == null) {
          RuntimeGrpc.getSubscribeConfigurationMethod = getSubscribeConfigurationMethod =
              io.grpc.MethodDescriptor.<spec.proto.runtime.v1.RuntimeProto.SubscribeConfigurationRequest, spec.proto.runtime.v1.RuntimeProto.SubscribeConfigurationResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.BIDI_STREAMING)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "SubscribeConfiguration"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  spec.proto.runtime.v1.RuntimeProto.SubscribeConfigurationRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  spec.proto.runtime.v1.RuntimeProto.SubscribeConfigurationResponse.getDefaultInstance()))
              .setSchemaDescriptor(new RuntimeMethodDescriptorSupplier("SubscribeConfiguration"))
              .build();
        }
      }
    }
    return getSubscribeConfigurationMethod;
  }

  private static volatile io.grpc.MethodDescriptor<spec.proto.runtime.v1.RuntimeProto.TryLockRequest,
      spec.proto.runtime.v1.RuntimeProto.TryLockResponse> getTryLockMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "TryLock",
      requestType = spec.proto.runtime.v1.RuntimeProto.TryLockRequest.class,
      responseType = spec.proto.runtime.v1.RuntimeProto.TryLockResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<spec.proto.runtime.v1.RuntimeProto.TryLockRequest,
      spec.proto.runtime.v1.RuntimeProto.TryLockResponse> getTryLockMethod() {
    io.grpc.MethodDescriptor<spec.proto.runtime.v1.RuntimeProto.TryLockRequest, spec.proto.runtime.v1.RuntimeProto.TryLockResponse> getTryLockMethod;
    if ((getTryLockMethod = RuntimeGrpc.getTryLockMethod) == null) {
      synchronized (RuntimeGrpc.class) {
        if ((getTryLockMethod = RuntimeGrpc.getTryLockMethod) == null) {
          RuntimeGrpc.getTryLockMethod = getTryLockMethod =
              io.grpc.MethodDescriptor.<spec.proto.runtime.v1.RuntimeProto.TryLockRequest, spec.proto.runtime.v1.RuntimeProto.TryLockResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "TryLock"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  spec.proto.runtime.v1.RuntimeProto.TryLockRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  spec.proto.runtime.v1.RuntimeProto.TryLockResponse.getDefaultInstance()))
              .setSchemaDescriptor(new RuntimeMethodDescriptorSupplier("TryLock"))
              .build();
        }
      }
    }
    return getTryLockMethod;
  }

  private static volatile io.grpc.MethodDescriptor<spec.proto.runtime.v1.RuntimeProto.UnlockRequest,
      spec.proto.runtime.v1.RuntimeProto.UnlockResponse> getUnlockMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Unlock",
      requestType = spec.proto.runtime.v1.RuntimeProto.UnlockRequest.class,
      responseType = spec.proto.runtime.v1.RuntimeProto.UnlockResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<spec.proto.runtime.v1.RuntimeProto.UnlockRequest,
      spec.proto.runtime.v1.RuntimeProto.UnlockResponse> getUnlockMethod() {
    io.grpc.MethodDescriptor<spec.proto.runtime.v1.RuntimeProto.UnlockRequest, spec.proto.runtime.v1.RuntimeProto.UnlockResponse> getUnlockMethod;
    if ((getUnlockMethod = RuntimeGrpc.getUnlockMethod) == null) {
      synchronized (RuntimeGrpc.class) {
        if ((getUnlockMethod = RuntimeGrpc.getUnlockMethod) == null) {
          RuntimeGrpc.getUnlockMethod = getUnlockMethod =
              io.grpc.MethodDescriptor.<spec.proto.runtime.v1.RuntimeProto.UnlockRequest, spec.proto.runtime.v1.RuntimeProto.UnlockResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "Unlock"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  spec.proto.runtime.v1.RuntimeProto.UnlockRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  spec.proto.runtime.v1.RuntimeProto.UnlockResponse.getDefaultInstance()))
              .setSchemaDescriptor(new RuntimeMethodDescriptorSupplier("Unlock"))
              .build();
        }
      }
    }
    return getUnlockMethod;
  }

  private static volatile io.grpc.MethodDescriptor<spec.proto.runtime.v1.RuntimeProto.GetNextIdRequest,
      spec.proto.runtime.v1.RuntimeProto.GetNextIdResponse> getGetNextIdMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetNextId",
      requestType = spec.proto.runtime.v1.RuntimeProto.GetNextIdRequest.class,
      responseType = spec.proto.runtime.v1.RuntimeProto.GetNextIdResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<spec.proto.runtime.v1.RuntimeProto.GetNextIdRequest,
      spec.proto.runtime.v1.RuntimeProto.GetNextIdResponse> getGetNextIdMethod() {
    io.grpc.MethodDescriptor<spec.proto.runtime.v1.RuntimeProto.GetNextIdRequest, spec.proto.runtime.v1.RuntimeProto.GetNextIdResponse> getGetNextIdMethod;
    if ((getGetNextIdMethod = RuntimeGrpc.getGetNextIdMethod) == null) {
      synchronized (RuntimeGrpc.class) {
        if ((getGetNextIdMethod = RuntimeGrpc.getGetNextIdMethod) == null) {
          RuntimeGrpc.getGetNextIdMethod = getGetNextIdMethod =
              io.grpc.MethodDescriptor.<spec.proto.runtime.v1.RuntimeProto.GetNextIdRequest, spec.proto.runtime.v1.RuntimeProto.GetNextIdResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetNextId"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  spec.proto.runtime.v1.RuntimeProto.GetNextIdRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  spec.proto.runtime.v1.RuntimeProto.GetNextIdResponse.getDefaultInstance()))
              .setSchemaDescriptor(new RuntimeMethodDescriptorSupplier("GetNextId"))
              .build();
        }
      }
    }
    return getGetNextIdMethod;
  }

  private static volatile io.grpc.MethodDescriptor<spec.proto.runtime.v1.RuntimeProto.GetStateRequest,
      spec.proto.runtime.v1.RuntimeProto.GetStateResponse> getGetStateMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetState",
      requestType = spec.proto.runtime.v1.RuntimeProto.GetStateRequest.class,
      responseType = spec.proto.runtime.v1.RuntimeProto.GetStateResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<spec.proto.runtime.v1.RuntimeProto.GetStateRequest,
      spec.proto.runtime.v1.RuntimeProto.GetStateResponse> getGetStateMethod() {
    io.grpc.MethodDescriptor<spec.proto.runtime.v1.RuntimeProto.GetStateRequest, spec.proto.runtime.v1.RuntimeProto.GetStateResponse> getGetStateMethod;
    if ((getGetStateMethod = RuntimeGrpc.getGetStateMethod) == null) {
      synchronized (RuntimeGrpc.class) {
        if ((getGetStateMethod = RuntimeGrpc.getGetStateMethod) == null) {
          RuntimeGrpc.getGetStateMethod = getGetStateMethod =
              io.grpc.MethodDescriptor.<spec.proto.runtime.v1.RuntimeProto.GetStateRequest, spec.proto.runtime.v1.RuntimeProto.GetStateResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetState"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  spec.proto.runtime.v1.RuntimeProto.GetStateRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  spec.proto.runtime.v1.RuntimeProto.GetStateResponse.getDefaultInstance()))
              .setSchemaDescriptor(new RuntimeMethodDescriptorSupplier("GetState"))
              .build();
        }
      }
    }
    return getGetStateMethod;
  }

  private static volatile io.grpc.MethodDescriptor<spec.proto.runtime.v1.RuntimeProto.GetBulkStateRequest,
      spec.proto.runtime.v1.RuntimeProto.GetBulkStateResponse> getGetBulkStateMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetBulkState",
      requestType = spec.proto.runtime.v1.RuntimeProto.GetBulkStateRequest.class,
      responseType = spec.proto.runtime.v1.RuntimeProto.GetBulkStateResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<spec.proto.runtime.v1.RuntimeProto.GetBulkStateRequest,
      spec.proto.runtime.v1.RuntimeProto.GetBulkStateResponse> getGetBulkStateMethod() {
    io.grpc.MethodDescriptor<spec.proto.runtime.v1.RuntimeProto.GetBulkStateRequest, spec.proto.runtime.v1.RuntimeProto.GetBulkStateResponse> getGetBulkStateMethod;
    if ((getGetBulkStateMethod = RuntimeGrpc.getGetBulkStateMethod) == null) {
      synchronized (RuntimeGrpc.class) {
        if ((getGetBulkStateMethod = RuntimeGrpc.getGetBulkStateMethod) == null) {
          RuntimeGrpc.getGetBulkStateMethod = getGetBulkStateMethod =
              io.grpc.MethodDescriptor.<spec.proto.runtime.v1.RuntimeProto.GetBulkStateRequest, spec.proto.runtime.v1.RuntimeProto.GetBulkStateResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetBulkState"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  spec.proto.runtime.v1.RuntimeProto.GetBulkStateRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  spec.proto.runtime.v1.RuntimeProto.GetBulkStateResponse.getDefaultInstance()))
              .setSchemaDescriptor(new RuntimeMethodDescriptorSupplier("GetBulkState"))
              .build();
        }
      }
    }
    return getGetBulkStateMethod;
  }

  private static volatile io.grpc.MethodDescriptor<spec.proto.runtime.v1.RuntimeProto.SaveStateRequest,
      com.google.protobuf.Empty> getSaveStateMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "SaveState",
      requestType = spec.proto.runtime.v1.RuntimeProto.SaveStateRequest.class,
      responseType = com.google.protobuf.Empty.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<spec.proto.runtime.v1.RuntimeProto.SaveStateRequest,
      com.google.protobuf.Empty> getSaveStateMethod() {
    io.grpc.MethodDescriptor<spec.proto.runtime.v1.RuntimeProto.SaveStateRequest, com.google.protobuf.Empty> getSaveStateMethod;
    if ((getSaveStateMethod = RuntimeGrpc.getSaveStateMethod) == null) {
      synchronized (RuntimeGrpc.class) {
        if ((getSaveStateMethod = RuntimeGrpc.getSaveStateMethod) == null) {
          RuntimeGrpc.getSaveStateMethod = getSaveStateMethod =
              io.grpc.MethodDescriptor.<spec.proto.runtime.v1.RuntimeProto.SaveStateRequest, com.google.protobuf.Empty>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "SaveState"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  spec.proto.runtime.v1.RuntimeProto.SaveStateRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.google.protobuf.Empty.getDefaultInstance()))
              .setSchemaDescriptor(new RuntimeMethodDescriptorSupplier("SaveState"))
              .build();
        }
      }
    }
    return getSaveStateMethod;
  }

  private static volatile io.grpc.MethodDescriptor<spec.proto.runtime.v1.RuntimeProto.DeleteStateRequest,
      com.google.protobuf.Empty> getDeleteStateMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "DeleteState",
      requestType = spec.proto.runtime.v1.RuntimeProto.DeleteStateRequest.class,
      responseType = com.google.protobuf.Empty.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<spec.proto.runtime.v1.RuntimeProto.DeleteStateRequest,
      com.google.protobuf.Empty> getDeleteStateMethod() {
    io.grpc.MethodDescriptor<spec.proto.runtime.v1.RuntimeProto.DeleteStateRequest, com.google.protobuf.Empty> getDeleteStateMethod;
    if ((getDeleteStateMethod = RuntimeGrpc.getDeleteStateMethod) == null) {
      synchronized (RuntimeGrpc.class) {
        if ((getDeleteStateMethod = RuntimeGrpc.getDeleteStateMethod) == null) {
          RuntimeGrpc.getDeleteStateMethod = getDeleteStateMethod =
              io.grpc.MethodDescriptor.<spec.proto.runtime.v1.RuntimeProto.DeleteStateRequest, com.google.protobuf.Empty>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "DeleteState"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  spec.proto.runtime.v1.RuntimeProto.DeleteStateRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.google.protobuf.Empty.getDefaultInstance()))
              .setSchemaDescriptor(new RuntimeMethodDescriptorSupplier("DeleteState"))
              .build();
        }
      }
    }
    return getDeleteStateMethod;
  }

  private static volatile io.grpc.MethodDescriptor<spec.proto.runtime.v1.RuntimeProto.DeleteBulkStateRequest,
      com.google.protobuf.Empty> getDeleteBulkStateMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "DeleteBulkState",
      requestType = spec.proto.runtime.v1.RuntimeProto.DeleteBulkStateRequest.class,
      responseType = com.google.protobuf.Empty.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<spec.proto.runtime.v1.RuntimeProto.DeleteBulkStateRequest,
      com.google.protobuf.Empty> getDeleteBulkStateMethod() {
    io.grpc.MethodDescriptor<spec.proto.runtime.v1.RuntimeProto.DeleteBulkStateRequest, com.google.protobuf.Empty> getDeleteBulkStateMethod;
    if ((getDeleteBulkStateMethod = RuntimeGrpc.getDeleteBulkStateMethod) == null) {
      synchronized (RuntimeGrpc.class) {
        if ((getDeleteBulkStateMethod = RuntimeGrpc.getDeleteBulkStateMethod) == null) {
          RuntimeGrpc.getDeleteBulkStateMethod = getDeleteBulkStateMethod =
              io.grpc.MethodDescriptor.<spec.proto.runtime.v1.RuntimeProto.DeleteBulkStateRequest, com.google.protobuf.Empty>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "DeleteBulkState"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  spec.proto.runtime.v1.RuntimeProto.DeleteBulkStateRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.google.protobuf.Empty.getDefaultInstance()))
              .setSchemaDescriptor(new RuntimeMethodDescriptorSupplier("DeleteBulkState"))
              .build();
        }
      }
    }
    return getDeleteBulkStateMethod;
  }

  private static volatile io.grpc.MethodDescriptor<spec.proto.runtime.v1.RuntimeProto.ExecuteStateTransactionRequest,
      com.google.protobuf.Empty> getExecuteStateTransactionMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "ExecuteStateTransaction",
      requestType = spec.proto.runtime.v1.RuntimeProto.ExecuteStateTransactionRequest.class,
      responseType = com.google.protobuf.Empty.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<spec.proto.runtime.v1.RuntimeProto.ExecuteStateTransactionRequest,
      com.google.protobuf.Empty> getExecuteStateTransactionMethod() {
    io.grpc.MethodDescriptor<spec.proto.runtime.v1.RuntimeProto.ExecuteStateTransactionRequest, com.google.protobuf.Empty> getExecuteStateTransactionMethod;
    if ((getExecuteStateTransactionMethod = RuntimeGrpc.getExecuteStateTransactionMethod) == null) {
      synchronized (RuntimeGrpc.class) {
        if ((getExecuteStateTransactionMethod = RuntimeGrpc.getExecuteStateTransactionMethod) == null) {
          RuntimeGrpc.getExecuteStateTransactionMethod = getExecuteStateTransactionMethod =
              io.grpc.MethodDescriptor.<spec.proto.runtime.v1.RuntimeProto.ExecuteStateTransactionRequest, com.google.protobuf.Empty>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "ExecuteStateTransaction"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  spec.proto.runtime.v1.RuntimeProto.ExecuteStateTransactionRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.google.protobuf.Empty.getDefaultInstance()))
              .setSchemaDescriptor(new RuntimeMethodDescriptorSupplier("ExecuteStateTransaction"))
              .build();
        }
      }
    }
    return getExecuteStateTransactionMethod;
  }

  private static volatile io.grpc.MethodDescriptor<spec.proto.runtime.v1.RuntimeProto.PublishEventRequest,
      com.google.protobuf.Empty> getPublishEventMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "PublishEvent",
      requestType = spec.proto.runtime.v1.RuntimeProto.PublishEventRequest.class,
      responseType = com.google.protobuf.Empty.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<spec.proto.runtime.v1.RuntimeProto.PublishEventRequest,
      com.google.protobuf.Empty> getPublishEventMethod() {
    io.grpc.MethodDescriptor<spec.proto.runtime.v1.RuntimeProto.PublishEventRequest, com.google.protobuf.Empty> getPublishEventMethod;
    if ((getPublishEventMethod = RuntimeGrpc.getPublishEventMethod) == null) {
      synchronized (RuntimeGrpc.class) {
        if ((getPublishEventMethod = RuntimeGrpc.getPublishEventMethod) == null) {
          RuntimeGrpc.getPublishEventMethod = getPublishEventMethod =
              io.grpc.MethodDescriptor.<spec.proto.runtime.v1.RuntimeProto.PublishEventRequest, com.google.protobuf.Empty>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "PublishEvent"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  spec.proto.runtime.v1.RuntimeProto.PublishEventRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.google.protobuf.Empty.getDefaultInstance()))
              .setSchemaDescriptor(new RuntimeMethodDescriptorSupplier("PublishEvent"))
              .build();
        }
      }
    }
    return getPublishEventMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static RuntimeStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<RuntimeStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<RuntimeStub>() {
        @java.lang.Override
        public RuntimeStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new RuntimeStub(channel, callOptions);
        }
      };
    return RuntimeStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static RuntimeBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<RuntimeBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<RuntimeBlockingStub>() {
        @java.lang.Override
        public RuntimeBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new RuntimeBlockingStub(channel, callOptions);
        }
      };
    return RuntimeBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static RuntimeFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<RuntimeFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<RuntimeFutureStub>() {
        @java.lang.Override
        public RuntimeFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new RuntimeFutureStub(channel, callOptions);
        }
      };
    return RuntimeFutureStub.newStub(factory, channel);
  }

  /**
   */
  public static abstract class RuntimeImplBase implements io.grpc.BindableService {

    /**
     * <pre>
     *SayHello used for test
     * </pre>
     */
    public void sayHello(spec.proto.runtime.v1.RuntimeProto.SayHelloRequest request,
        io.grpc.stub.StreamObserver<spec.proto.runtime.v1.RuntimeProto.SayHelloResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getSayHelloMethod(), responseObserver);
    }

    /**
     * <pre>
     * InvokeService do rpc calls
     * </pre>
     */
    public void invokeService(spec.proto.runtime.v1.RuntimeProto.InvokeServiceRequest request,
        io.grpc.stub.StreamObserver<spec.proto.runtime.v1.RuntimeProto.InvokeResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getInvokeServiceMethod(), responseObserver);
    }

    /**
     * <pre>
     * GetConfiguration gets configuration from configuration store.
     * </pre>
     */
    public void getConfiguration(spec.proto.runtime.v1.RuntimeProto.GetConfigurationRequest request,
        io.grpc.stub.StreamObserver<spec.proto.runtime.v1.RuntimeProto.GetConfigurationResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetConfigurationMethod(), responseObserver);
    }

    /**
     * <pre>
     * SaveConfiguration saves configuration into configuration store.
     * </pre>
     */
    public void saveConfiguration(spec.proto.runtime.v1.RuntimeProto.SaveConfigurationRequest request,
        io.grpc.stub.StreamObserver<com.google.protobuf.Empty> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getSaveConfigurationMethod(), responseObserver);
    }

    /**
     * <pre>
     * DeleteConfiguration deletes configuration from configuration store.
     * </pre>
     */
    public void deleteConfiguration(spec.proto.runtime.v1.RuntimeProto.DeleteConfigurationRequest request,
        io.grpc.stub.StreamObserver<com.google.protobuf.Empty> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getDeleteConfigurationMethod(), responseObserver);
    }

    /**
     * <pre>
     * SubscribeConfiguration gets configuration from configuration store and subscribe the updates.
     * </pre>
     */
    public io.grpc.stub.StreamObserver<spec.proto.runtime.v1.RuntimeProto.SubscribeConfigurationRequest> subscribeConfiguration(
        io.grpc.stub.StreamObserver<spec.proto.runtime.v1.RuntimeProto.SubscribeConfigurationResponse> responseObserver) {
      return io.grpc.stub.ServerCalls.asyncUnimplementedStreamingCall(getSubscribeConfigurationMethod(), responseObserver);
    }

    /**
     * <pre>
     * Distributed Lock API
     * A non-blocking method trying to get a lock with ttl.
     * </pre>
     */
    public void tryLock(spec.proto.runtime.v1.RuntimeProto.TryLockRequest request,
        io.grpc.stub.StreamObserver<spec.proto.runtime.v1.RuntimeProto.TryLockResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getTryLockMethod(), responseObserver);
    }

    /**
     */
    public void unlock(spec.proto.runtime.v1.RuntimeProto.UnlockRequest request,
        io.grpc.stub.StreamObserver<spec.proto.runtime.v1.RuntimeProto.UnlockResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getUnlockMethod(), responseObserver);
    }

    /**
     * <pre>
     * Sequencer API
     * Get next unique id with some auto-increment guarantee
     * </pre>
     */
    public void getNextId(spec.proto.runtime.v1.RuntimeProto.GetNextIdRequest request,
        io.grpc.stub.StreamObserver<spec.proto.runtime.v1.RuntimeProto.GetNextIdResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetNextIdMethod(), responseObserver);
    }

    /**
     * <pre>
     * Gets the state for a specific key.
     * </pre>
     */
    public void getState(spec.proto.runtime.v1.RuntimeProto.GetStateRequest request,
        io.grpc.stub.StreamObserver<spec.proto.runtime.v1.RuntimeProto.GetStateResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetStateMethod(), responseObserver);
    }

    /**
     * <pre>
     * Gets a bulk of state items for a list of keys
     * </pre>
     */
    public void getBulkState(spec.proto.runtime.v1.RuntimeProto.GetBulkStateRequest request,
        io.grpc.stub.StreamObserver<spec.proto.runtime.v1.RuntimeProto.GetBulkStateResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetBulkStateMethod(), responseObserver);
    }

    /**
     * <pre>
     * Saves an array of state objects
     * </pre>
     */
    public void saveState(spec.proto.runtime.v1.RuntimeProto.SaveStateRequest request,
        io.grpc.stub.StreamObserver<com.google.protobuf.Empty> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getSaveStateMethod(), responseObserver);
    }

    /**
     * <pre>
     * Deletes the state for a specific key.
     * </pre>
     */
    public void deleteState(spec.proto.runtime.v1.RuntimeProto.DeleteStateRequest request,
        io.grpc.stub.StreamObserver<com.google.protobuf.Empty> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getDeleteStateMethod(), responseObserver);
    }

    /**
     * <pre>
     * Deletes a bulk of state items for a list of keys
     * </pre>
     */
    public void deleteBulkState(spec.proto.runtime.v1.RuntimeProto.DeleteBulkStateRequest request,
        io.grpc.stub.StreamObserver<com.google.protobuf.Empty> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getDeleteBulkStateMethod(), responseObserver);
    }

    /**
     * <pre>
     * Executes transactions for a specified store
     * </pre>
     */
    public void executeStateTransaction(spec.proto.runtime.v1.RuntimeProto.ExecuteStateTransactionRequest request,
        io.grpc.stub.StreamObserver<com.google.protobuf.Empty> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getExecuteStateTransactionMethod(), responseObserver);
    }

    /**
     * <pre>
     * Publishes events to the specific topic.
     * </pre>
     */
    public void publishEvent(spec.proto.runtime.v1.RuntimeProto.PublishEventRequest request,
        io.grpc.stub.StreamObserver<com.google.protobuf.Empty> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getPublishEventMethod(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getSayHelloMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                spec.proto.runtime.v1.RuntimeProto.SayHelloRequest,
                spec.proto.runtime.v1.RuntimeProto.SayHelloResponse>(
                  this, METHODID_SAY_HELLO)))
          .addMethod(
            getInvokeServiceMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                spec.proto.runtime.v1.RuntimeProto.InvokeServiceRequest,
                spec.proto.runtime.v1.RuntimeProto.InvokeResponse>(
                  this, METHODID_INVOKE_SERVICE)))
          .addMethod(
            getGetConfigurationMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                spec.proto.runtime.v1.RuntimeProto.GetConfigurationRequest,
                spec.proto.runtime.v1.RuntimeProto.GetConfigurationResponse>(
                  this, METHODID_GET_CONFIGURATION)))
          .addMethod(
            getSaveConfigurationMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                spec.proto.runtime.v1.RuntimeProto.SaveConfigurationRequest,
                com.google.protobuf.Empty>(
                  this, METHODID_SAVE_CONFIGURATION)))
          .addMethod(
            getDeleteConfigurationMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                spec.proto.runtime.v1.RuntimeProto.DeleteConfigurationRequest,
                com.google.protobuf.Empty>(
                  this, METHODID_DELETE_CONFIGURATION)))
          .addMethod(
            getSubscribeConfigurationMethod(),
            io.grpc.stub.ServerCalls.asyncBidiStreamingCall(
              new MethodHandlers<
                spec.proto.runtime.v1.RuntimeProto.SubscribeConfigurationRequest,
                spec.proto.runtime.v1.RuntimeProto.SubscribeConfigurationResponse>(
                  this, METHODID_SUBSCRIBE_CONFIGURATION)))
          .addMethod(
            getTryLockMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                spec.proto.runtime.v1.RuntimeProto.TryLockRequest,
                spec.proto.runtime.v1.RuntimeProto.TryLockResponse>(
                  this, METHODID_TRY_LOCK)))
          .addMethod(
            getUnlockMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                spec.proto.runtime.v1.RuntimeProto.UnlockRequest,
                spec.proto.runtime.v1.RuntimeProto.UnlockResponse>(
                  this, METHODID_UNLOCK)))
          .addMethod(
            getGetNextIdMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                spec.proto.runtime.v1.RuntimeProto.GetNextIdRequest,
                spec.proto.runtime.v1.RuntimeProto.GetNextIdResponse>(
                  this, METHODID_GET_NEXT_ID)))
          .addMethod(
            getGetStateMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                spec.proto.runtime.v1.RuntimeProto.GetStateRequest,
                spec.proto.runtime.v1.RuntimeProto.GetStateResponse>(
                  this, METHODID_GET_STATE)))
          .addMethod(
            getGetBulkStateMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                spec.proto.runtime.v1.RuntimeProto.GetBulkStateRequest,
                spec.proto.runtime.v1.RuntimeProto.GetBulkStateResponse>(
                  this, METHODID_GET_BULK_STATE)))
          .addMethod(
            getSaveStateMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                spec.proto.runtime.v1.RuntimeProto.SaveStateRequest,
                com.google.protobuf.Empty>(
                  this, METHODID_SAVE_STATE)))
          .addMethod(
            getDeleteStateMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                spec.proto.runtime.v1.RuntimeProto.DeleteStateRequest,
                com.google.protobuf.Empty>(
                  this, METHODID_DELETE_STATE)))
          .addMethod(
            getDeleteBulkStateMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                spec.proto.runtime.v1.RuntimeProto.DeleteBulkStateRequest,
                com.google.protobuf.Empty>(
                  this, METHODID_DELETE_BULK_STATE)))
          .addMethod(
            getExecuteStateTransactionMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                spec.proto.runtime.v1.RuntimeProto.ExecuteStateTransactionRequest,
                com.google.protobuf.Empty>(
                  this, METHODID_EXECUTE_STATE_TRANSACTION)))
          .addMethod(
            getPublishEventMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
              new MethodHandlers<
                spec.proto.runtime.v1.RuntimeProto.PublishEventRequest,
                com.google.protobuf.Empty>(
                  this, METHODID_PUBLISH_EVENT)))
          .build();
    }
  }

  /**
   */
  public static final class RuntimeStub extends io.grpc.stub.AbstractAsyncStub<RuntimeStub> {
    private RuntimeStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected RuntimeStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new RuntimeStub(channel, callOptions);
    }

    /**
     * <pre>
     *SayHello used for test
     * </pre>
     */
    public void sayHello(spec.proto.runtime.v1.RuntimeProto.SayHelloRequest request,
        io.grpc.stub.StreamObserver<spec.proto.runtime.v1.RuntimeProto.SayHelloResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getSayHelloMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * InvokeService do rpc calls
     * </pre>
     */
    public void invokeService(spec.proto.runtime.v1.RuntimeProto.InvokeServiceRequest request,
        io.grpc.stub.StreamObserver<spec.proto.runtime.v1.RuntimeProto.InvokeResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getInvokeServiceMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * GetConfiguration gets configuration from configuration store.
     * </pre>
     */
    public void getConfiguration(spec.proto.runtime.v1.RuntimeProto.GetConfigurationRequest request,
        io.grpc.stub.StreamObserver<spec.proto.runtime.v1.RuntimeProto.GetConfigurationResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetConfigurationMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * SaveConfiguration saves configuration into configuration store.
     * </pre>
     */
    public void saveConfiguration(spec.proto.runtime.v1.RuntimeProto.SaveConfigurationRequest request,
        io.grpc.stub.StreamObserver<com.google.protobuf.Empty> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getSaveConfigurationMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * DeleteConfiguration deletes configuration from configuration store.
     * </pre>
     */
    public void deleteConfiguration(spec.proto.runtime.v1.RuntimeProto.DeleteConfigurationRequest request,
        io.grpc.stub.StreamObserver<com.google.protobuf.Empty> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getDeleteConfigurationMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * SubscribeConfiguration gets configuration from configuration store and subscribe the updates.
     * </pre>
     */
    public io.grpc.stub.StreamObserver<spec.proto.runtime.v1.RuntimeProto.SubscribeConfigurationRequest> subscribeConfiguration(
        io.grpc.stub.StreamObserver<spec.proto.runtime.v1.RuntimeProto.SubscribeConfigurationResponse> responseObserver) {
      return io.grpc.stub.ClientCalls.asyncBidiStreamingCall(
          getChannel().newCall(getSubscribeConfigurationMethod(), getCallOptions()), responseObserver);
    }

    /**
     * <pre>
     * Distributed Lock API
     * A non-blocking method trying to get a lock with ttl.
     * </pre>
     */
    public void tryLock(spec.proto.runtime.v1.RuntimeProto.TryLockRequest request,
        io.grpc.stub.StreamObserver<spec.proto.runtime.v1.RuntimeProto.TryLockResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getTryLockMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void unlock(spec.proto.runtime.v1.RuntimeProto.UnlockRequest request,
        io.grpc.stub.StreamObserver<spec.proto.runtime.v1.RuntimeProto.UnlockResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getUnlockMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * Sequencer API
     * Get next unique id with some auto-increment guarantee
     * </pre>
     */
    public void getNextId(spec.proto.runtime.v1.RuntimeProto.GetNextIdRequest request,
        io.grpc.stub.StreamObserver<spec.proto.runtime.v1.RuntimeProto.GetNextIdResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetNextIdMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * Gets the state for a specific key.
     * </pre>
     */
    public void getState(spec.proto.runtime.v1.RuntimeProto.GetStateRequest request,
        io.grpc.stub.StreamObserver<spec.proto.runtime.v1.RuntimeProto.GetStateResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetStateMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * Gets a bulk of state items for a list of keys
     * </pre>
     */
    public void getBulkState(spec.proto.runtime.v1.RuntimeProto.GetBulkStateRequest request,
        io.grpc.stub.StreamObserver<spec.proto.runtime.v1.RuntimeProto.GetBulkStateResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetBulkStateMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * Saves an array of state objects
     * </pre>
     */
    public void saveState(spec.proto.runtime.v1.RuntimeProto.SaveStateRequest request,
        io.grpc.stub.StreamObserver<com.google.protobuf.Empty> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getSaveStateMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * Deletes the state for a specific key.
     * </pre>
     */
    public void deleteState(spec.proto.runtime.v1.RuntimeProto.DeleteStateRequest request,
        io.grpc.stub.StreamObserver<com.google.protobuf.Empty> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getDeleteStateMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * Deletes a bulk of state items for a list of keys
     * </pre>
     */
    public void deleteBulkState(spec.proto.runtime.v1.RuntimeProto.DeleteBulkStateRequest request,
        io.grpc.stub.StreamObserver<com.google.protobuf.Empty> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getDeleteBulkStateMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * Executes transactions for a specified store
     * </pre>
     */
    public void executeStateTransaction(spec.proto.runtime.v1.RuntimeProto.ExecuteStateTransactionRequest request,
        io.grpc.stub.StreamObserver<com.google.protobuf.Empty> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getExecuteStateTransactionMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     * <pre>
     * Publishes events to the specific topic.
     * </pre>
     */
    public void publishEvent(spec.proto.runtime.v1.RuntimeProto.PublishEventRequest request,
        io.grpc.stub.StreamObserver<com.google.protobuf.Empty> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getPublishEventMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   */
  public static final class RuntimeBlockingStub extends io.grpc.stub.AbstractBlockingStub<RuntimeBlockingStub> {
    private RuntimeBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected RuntimeBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new RuntimeBlockingStub(channel, callOptions);
    }

    /**
     * <pre>
     *SayHello used for test
     * </pre>
     */
    public spec.proto.runtime.v1.RuntimeProto.SayHelloResponse sayHello(spec.proto.runtime.v1.RuntimeProto.SayHelloRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getSayHelloMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * InvokeService do rpc calls
     * </pre>
     */
    public spec.proto.runtime.v1.RuntimeProto.InvokeResponse invokeService(spec.proto.runtime.v1.RuntimeProto.InvokeServiceRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getInvokeServiceMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * GetConfiguration gets configuration from configuration store.
     * </pre>
     */
    public spec.proto.runtime.v1.RuntimeProto.GetConfigurationResponse getConfiguration(spec.proto.runtime.v1.RuntimeProto.GetConfigurationRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetConfigurationMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * SaveConfiguration saves configuration into configuration store.
     * </pre>
     */
    public com.google.protobuf.Empty saveConfiguration(spec.proto.runtime.v1.RuntimeProto.SaveConfigurationRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getSaveConfigurationMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * DeleteConfiguration deletes configuration from configuration store.
     * </pre>
     */
    public com.google.protobuf.Empty deleteConfiguration(spec.proto.runtime.v1.RuntimeProto.DeleteConfigurationRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getDeleteConfigurationMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * Distributed Lock API
     * A non-blocking method trying to get a lock with ttl.
     * </pre>
     */
    public spec.proto.runtime.v1.RuntimeProto.TryLockResponse tryLock(spec.proto.runtime.v1.RuntimeProto.TryLockRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getTryLockMethod(), getCallOptions(), request);
    }

    /**
     */
    public spec.proto.runtime.v1.RuntimeProto.UnlockResponse unlock(spec.proto.runtime.v1.RuntimeProto.UnlockRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getUnlockMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * Sequencer API
     * Get next unique id with some auto-increment guarantee
     * </pre>
     */
    public spec.proto.runtime.v1.RuntimeProto.GetNextIdResponse getNextId(spec.proto.runtime.v1.RuntimeProto.GetNextIdRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetNextIdMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * Gets the state for a specific key.
     * </pre>
     */
    public spec.proto.runtime.v1.RuntimeProto.GetStateResponse getState(spec.proto.runtime.v1.RuntimeProto.GetStateRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetStateMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * Gets a bulk of state items for a list of keys
     * </pre>
     */
    public spec.proto.runtime.v1.RuntimeProto.GetBulkStateResponse getBulkState(spec.proto.runtime.v1.RuntimeProto.GetBulkStateRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetBulkStateMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * Saves an array of state objects
     * </pre>
     */
    public com.google.protobuf.Empty saveState(spec.proto.runtime.v1.RuntimeProto.SaveStateRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getSaveStateMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * Deletes the state for a specific key.
     * </pre>
     */
    public com.google.protobuf.Empty deleteState(spec.proto.runtime.v1.RuntimeProto.DeleteStateRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getDeleteStateMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * Deletes a bulk of state items for a list of keys
     * </pre>
     */
    public com.google.protobuf.Empty deleteBulkState(spec.proto.runtime.v1.RuntimeProto.DeleteBulkStateRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getDeleteBulkStateMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * Executes transactions for a specified store
     * </pre>
     */
    public com.google.protobuf.Empty executeStateTransaction(spec.proto.runtime.v1.RuntimeProto.ExecuteStateTransactionRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getExecuteStateTransactionMethod(), getCallOptions(), request);
    }

    /**
     * <pre>
     * Publishes events to the specific topic.
     * </pre>
     */
    public com.google.protobuf.Empty publishEvent(spec.proto.runtime.v1.RuntimeProto.PublishEventRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getPublishEventMethod(), getCallOptions(), request);
    }
  }

  /**
   */
  public static final class RuntimeFutureStub extends io.grpc.stub.AbstractFutureStub<RuntimeFutureStub> {
    private RuntimeFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected RuntimeFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new RuntimeFutureStub(channel, callOptions);
    }

    /**
     * <pre>
     *SayHello used for test
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<spec.proto.runtime.v1.RuntimeProto.SayHelloResponse> sayHello(
        spec.proto.runtime.v1.RuntimeProto.SayHelloRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getSayHelloMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * InvokeService do rpc calls
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<spec.proto.runtime.v1.RuntimeProto.InvokeResponse> invokeService(
        spec.proto.runtime.v1.RuntimeProto.InvokeServiceRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getInvokeServiceMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * GetConfiguration gets configuration from configuration store.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<spec.proto.runtime.v1.RuntimeProto.GetConfigurationResponse> getConfiguration(
        spec.proto.runtime.v1.RuntimeProto.GetConfigurationRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetConfigurationMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * SaveConfiguration saves configuration into configuration store.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<com.google.protobuf.Empty> saveConfiguration(
        spec.proto.runtime.v1.RuntimeProto.SaveConfigurationRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getSaveConfigurationMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * DeleteConfiguration deletes configuration from configuration store.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<com.google.protobuf.Empty> deleteConfiguration(
        spec.proto.runtime.v1.RuntimeProto.DeleteConfigurationRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getDeleteConfigurationMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * Distributed Lock API
     * A non-blocking method trying to get a lock with ttl.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<spec.proto.runtime.v1.RuntimeProto.TryLockResponse> tryLock(
        spec.proto.runtime.v1.RuntimeProto.TryLockRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getTryLockMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<spec.proto.runtime.v1.RuntimeProto.UnlockResponse> unlock(
        spec.proto.runtime.v1.RuntimeProto.UnlockRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getUnlockMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * Sequencer API
     * Get next unique id with some auto-increment guarantee
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<spec.proto.runtime.v1.RuntimeProto.GetNextIdResponse> getNextId(
        spec.proto.runtime.v1.RuntimeProto.GetNextIdRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetNextIdMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * Gets the state for a specific key.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<spec.proto.runtime.v1.RuntimeProto.GetStateResponse> getState(
        spec.proto.runtime.v1.RuntimeProto.GetStateRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetStateMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * Gets a bulk of state items for a list of keys
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<spec.proto.runtime.v1.RuntimeProto.GetBulkStateResponse> getBulkState(
        spec.proto.runtime.v1.RuntimeProto.GetBulkStateRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetBulkStateMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * Saves an array of state objects
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<com.google.protobuf.Empty> saveState(
        spec.proto.runtime.v1.RuntimeProto.SaveStateRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getSaveStateMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * Deletes the state for a specific key.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<com.google.protobuf.Empty> deleteState(
        spec.proto.runtime.v1.RuntimeProto.DeleteStateRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getDeleteStateMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * Deletes a bulk of state items for a list of keys
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<com.google.protobuf.Empty> deleteBulkState(
        spec.proto.runtime.v1.RuntimeProto.DeleteBulkStateRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getDeleteBulkStateMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * Executes transactions for a specified store
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<com.google.protobuf.Empty> executeStateTransaction(
        spec.proto.runtime.v1.RuntimeProto.ExecuteStateTransactionRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getExecuteStateTransactionMethod(), getCallOptions()), request);
    }

    /**
     * <pre>
     * Publishes events to the specific topic.
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<com.google.protobuf.Empty> publishEvent(
        spec.proto.runtime.v1.RuntimeProto.PublishEventRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getPublishEventMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_SAY_HELLO = 0;
  private static final int METHODID_INVOKE_SERVICE = 1;
  private static final int METHODID_GET_CONFIGURATION = 2;
  private static final int METHODID_SAVE_CONFIGURATION = 3;
  private static final int METHODID_DELETE_CONFIGURATION = 4;
  private static final int METHODID_TRY_LOCK = 5;
  private static final int METHODID_UNLOCK = 6;
  private static final int METHODID_GET_NEXT_ID = 7;
  private static final int METHODID_GET_STATE = 8;
  private static final int METHODID_GET_BULK_STATE = 9;
  private static final int METHODID_SAVE_STATE = 10;
  private static final int METHODID_DELETE_STATE = 11;
  private static final int METHODID_DELETE_BULK_STATE = 12;
  private static final int METHODID_EXECUTE_STATE_TRANSACTION = 13;
  private static final int METHODID_PUBLISH_EVENT = 14;
  private static final int METHODID_SUBSCRIBE_CONFIGURATION = 15;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final RuntimeImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(RuntimeImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_SAY_HELLO:
          serviceImpl.sayHello((spec.proto.runtime.v1.RuntimeProto.SayHelloRequest) request,
              (io.grpc.stub.StreamObserver<spec.proto.runtime.v1.RuntimeProto.SayHelloResponse>) responseObserver);
          break;
        case METHODID_INVOKE_SERVICE:
          serviceImpl.invokeService((spec.proto.runtime.v1.RuntimeProto.InvokeServiceRequest) request,
              (io.grpc.stub.StreamObserver<spec.proto.runtime.v1.RuntimeProto.InvokeResponse>) responseObserver);
          break;
        case METHODID_GET_CONFIGURATION:
          serviceImpl.getConfiguration((spec.proto.runtime.v1.RuntimeProto.GetConfigurationRequest) request,
              (io.grpc.stub.StreamObserver<spec.proto.runtime.v1.RuntimeProto.GetConfigurationResponse>) responseObserver);
          break;
        case METHODID_SAVE_CONFIGURATION:
          serviceImpl.saveConfiguration((spec.proto.runtime.v1.RuntimeProto.SaveConfigurationRequest) request,
              (io.grpc.stub.StreamObserver<com.google.protobuf.Empty>) responseObserver);
          break;
        case METHODID_DELETE_CONFIGURATION:
          serviceImpl.deleteConfiguration((spec.proto.runtime.v1.RuntimeProto.DeleteConfigurationRequest) request,
              (io.grpc.stub.StreamObserver<com.google.protobuf.Empty>) responseObserver);
          break;
        case METHODID_TRY_LOCK:
          serviceImpl.tryLock((spec.proto.runtime.v1.RuntimeProto.TryLockRequest) request,
              (io.grpc.stub.StreamObserver<spec.proto.runtime.v1.RuntimeProto.TryLockResponse>) responseObserver);
          break;
        case METHODID_UNLOCK:
          serviceImpl.unlock((spec.proto.runtime.v1.RuntimeProto.UnlockRequest) request,
              (io.grpc.stub.StreamObserver<spec.proto.runtime.v1.RuntimeProto.UnlockResponse>) responseObserver);
          break;
        case METHODID_GET_NEXT_ID:
          serviceImpl.getNextId((spec.proto.runtime.v1.RuntimeProto.GetNextIdRequest) request,
              (io.grpc.stub.StreamObserver<spec.proto.runtime.v1.RuntimeProto.GetNextIdResponse>) responseObserver);
          break;
        case METHODID_GET_STATE:
          serviceImpl.getState((spec.proto.runtime.v1.RuntimeProto.GetStateRequest) request,
              (io.grpc.stub.StreamObserver<spec.proto.runtime.v1.RuntimeProto.GetStateResponse>) responseObserver);
          break;
        case METHODID_GET_BULK_STATE:
          serviceImpl.getBulkState((spec.proto.runtime.v1.RuntimeProto.GetBulkStateRequest) request,
              (io.grpc.stub.StreamObserver<spec.proto.runtime.v1.RuntimeProto.GetBulkStateResponse>) responseObserver);
          break;
        case METHODID_SAVE_STATE:
          serviceImpl.saveState((spec.proto.runtime.v1.RuntimeProto.SaveStateRequest) request,
              (io.grpc.stub.StreamObserver<com.google.protobuf.Empty>) responseObserver);
          break;
        case METHODID_DELETE_STATE:
          serviceImpl.deleteState((spec.proto.runtime.v1.RuntimeProto.DeleteStateRequest) request,
              (io.grpc.stub.StreamObserver<com.google.protobuf.Empty>) responseObserver);
          break;
        case METHODID_DELETE_BULK_STATE:
          serviceImpl.deleteBulkState((spec.proto.runtime.v1.RuntimeProto.DeleteBulkStateRequest) request,
              (io.grpc.stub.StreamObserver<com.google.protobuf.Empty>) responseObserver);
          break;
        case METHODID_EXECUTE_STATE_TRANSACTION:
          serviceImpl.executeStateTransaction((spec.proto.runtime.v1.RuntimeProto.ExecuteStateTransactionRequest) request,
              (io.grpc.stub.StreamObserver<com.google.protobuf.Empty>) responseObserver);
          break;
        case METHODID_PUBLISH_EVENT:
          serviceImpl.publishEvent((spec.proto.runtime.v1.RuntimeProto.PublishEventRequest) request,
              (io.grpc.stub.StreamObserver<com.google.protobuf.Empty>) responseObserver);
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
        case METHODID_SUBSCRIBE_CONFIGURATION:
          return (io.grpc.stub.StreamObserver<Req>) serviceImpl.subscribeConfiguration(
              (io.grpc.stub.StreamObserver<spec.proto.runtime.v1.RuntimeProto.SubscribeConfigurationResponse>) responseObserver);
        default:
          throw new AssertionError();
      }
    }
  }

  private static abstract class RuntimeBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    RuntimeBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return spec.proto.runtime.v1.RuntimeProto.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("Runtime");
    }
  }

  private static final class RuntimeFileDescriptorSupplier
      extends RuntimeBaseDescriptorSupplier {
    RuntimeFileDescriptorSupplier() {}
  }

  private static final class RuntimeMethodDescriptorSupplier
      extends RuntimeBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    RuntimeMethodDescriptorSupplier(String methodName) {
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
      synchronized (RuntimeGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new RuntimeFileDescriptorSupplier())
              .addMethod(getSayHelloMethod())
              .addMethod(getInvokeServiceMethod())
              .addMethod(getGetConfigurationMethod())
              .addMethod(getSaveConfigurationMethod())
              .addMethod(getDeleteConfigurationMethod())
              .addMethod(getSubscribeConfigurationMethod())
              .addMethod(getTryLockMethod())
              .addMethod(getUnlockMethod())
              .addMethod(getGetNextIdMethod())
              .addMethod(getGetStateMethod())
              .addMethod(getGetBulkStateMethod())
              .addMethod(getSaveStateMethod())
              .addMethod(getDeleteStateMethod())
              .addMethod(getDeleteBulkStateMethod())
              .addMethod(getExecuteStateTransactionMethod())
              .addMethod(getPublishEventMethod())
              .build();
        }
      }
    }
    return result;
  }
}
