# RPC API
## What is RPC API
Layotto's RPC API is based on [Mosn](https://mosn.io/en/)'s grpc handler, which provides multi-protocol, modular, intelligent, and secure **service calls**.

The interface of the RPC API are consistent with [Dapr](https://docs.dapr.io/developing-applications/building-blocks/service-invocation/service-invocation-overview/), you could see its details in [invoke.go](https://github.com/mosn/layotto/blob/3802c4591181fdbcfb7dd07cbbdbadeaaada650a/sdk/go-sdk/client/invoke.go).

Using Layotto RPC invocation, your application can reliably and securely communicate with other applications using the standard  HTTP or [X-Protocol](https://cloudnative.to/blog/x-protocol-common-address-solution/) protocols.

![sidecar](https://mosn.io/en/docs/concept/sidecar-pattern/sidecar-pattern.jpg)

Layotto uses a [sidecar](https://mosn.io/en/docs/concept/sidecar-pattern/) architecture. To invoke an application using Layotto, you use the invoke API on any Layotto instance. The sidecar programming model encourages each applications to talk to its own instance of Layotto. The Layotto instances discover and communicate with one another.


## When to use RPC API and what are the benefits?
In many environments with multiple services that need to communicate with each other, developers often ask themselves the following questions:

- How do I discover and invoke methods on different services? 
- How do I call other services securely with encryption and apply access control on the methods? 
- How do I handle retries and transient errors? 
- How do I use tracing to see a call graph with metrics to diagnose issues in production?

Layotto addresses these challenges by providing a service invocation API that acts as a combination of a reverse proxy with built-in service discovery, while leveraging built-in distributed tracing, metrics, error handling, encryption and more.

## How to use RPC API？
You can do RPC calls through grpc interface **InvokeService**. The API is defined in [api.go](https://github.com/mosn/layotto/blob/77e0a4b2af063ff9e365a933c4735655898de369/pkg/grpc/api.go).

The component needs to be configured before use. For detailed configuration instructions, see [configuration reference](https://mosn.io/layotto/#/en/configuration/overview), currently, Layotto uses MOSN layer 4 filter integrated with MOSN and runs on MOSN, so the configuration file used by Layotto is actually a MOSN configuration file. So you can also refer to the documentation of the [MOSN configuration file](https://mosn.io/en/docs/configuration/).

### Demo 1: Hello World -- Basic Golang HTTP Server
Quick start Document for this demo: [Hello World](https://mosn.io/layotto/#/en/start/rpc/helloworld)

The [echoserver](https://github.com/mosn/layotto/blob/77e0a4b2af063ff9e365a933c4735655898de369/demo/rpc/http/echoserver/echoserver.go) publish a simple interface in port 8889, the config file [example.json](https://github.com/mosn/layotto/blob/77e0a4b2af063ff9e365a933c4735655898de369/demo/rpc/http/example.json) use mosn's routing capabilities to forward the request whose id field in the http header is equal to **HelloService:1.0** to the local port 8889, and then the [echoclient](https://github.com/mosn/layotto/blob/b66b998f50901f8bd1cce035478579c1b47f986d/demo/rpc/http/echoclient/echoclient.go) do a RPC calls though grpc interface **InvokeService**. 

```golang
resp, err := cli.InvokeService(
		ctx,
		&runtimev1pb.InvokeServiceRequest{
			Id: "HelloService:1.0",
			Message: &runtimev1pb.CommonInvokeRequest{
				Method:      "/hello",
				ContentType: "",
				Data:        &anypb.Any{Value: []byte(*data)}}},
	)
```

### Demo 2: Dubbo JSON RPC
Quick start Document for this demo: [Dubbo JSON RPC Example](https://mosn.io/layotto/#/en/start/rpc/dubbo_json_rpc)

The server is [dubbo-go-samples](https://github.com/apache/dubbo-go-samples), the config file [example.json](https://github.com/mosn/layotto/blob/77e0a4b2af063ff9e365a933c4735655898de369/demo/rpc/dubbo_json_rpc/example.json) use the callback function [dubbo_json_rpc](https://github.com/mosn/layotto/blob/8db7a2297bd05d1b0c4452cc980d8f6412a82f3a/components/rpc/callback/dubbo_json_rpc.go) to generate a request header.
And then,the [client](https://github.com/mosn/layotto/blob/b66b998f50901f8bd1cce035478579c1b47f986d/demo/rpc/dubbo_json_rpc/dubbo_json_client/client.go) do a RPC calls though grpc interface **InvokeService**. 

```golang
resp, err := cli.InvokeService(
		ctx,
		&runtimev1pb.InvokeServiceRequest{
			Id: "com.ikurento.user.UserProvider",
			Message: &runtimev1pb.CommonInvokeRequest{
				Method:        "GetUser",
				ContentType:   "",
				Data:          &anypb.Any{Value: []byte(*data)},
				HttpExtension: &runtimev1pb.HTTPExtension{Verb: runtimev1pb.HTTPExtension_POST},
			},
		},
	)
```

## Implementation Principle
If you are interested in the implementation principle, or want to extend some functions, you can read [RPC design document](https://mosn.io/layotto/#/en/design/rpc/rpc-design-doc).