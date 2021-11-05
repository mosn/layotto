## trace

### Features

In [runtime_config.json](https://github.com/mosn/layotto/blob/main/configs/runtime_config.json), there is a paragraph about trace configuration as follows:

```json
[
  "tracing": {
    "enable": true,
    "driver": "SOFATracer",
    "config": {
      "generator": "mosntracing",
      "exporter": ["stdout"]
    }
  }
]
```
This configuration can turn on the trace capability of layotto. The user can specify the method of trace reporting and the generation method of spanId and traceId through configuration.

The corresponding caller code is in [client.go](https://github.com/mosn/layotto/blob/main/demo/flowcontrol/client.go), the trace of layotto is printed as follows:

![img.png](../../../img/trace/trace.png)


### Configuration parameter description

Trace configuration:

| Field name | Field type | Description |
| ---- | ---- | ---- |
| enable | boolean | Global switch, whether to enable trace|
| driver | String | driver represents the type of trace, mosn has SOFATracer and skywalking, users can expand |
| config | Object | Trace expansion configuration |

Trace expansion configuration:

| Field name | Field type | Description |
| ---- | ---- | ---- |
| generator | String | SpanId, traceId and other resource generation methods, users can expand by themselves |
| exporter | Array | The way users need to report by trace can be implemented and expanded by themselves |




### Trace design and expansion

#### Span structure:

```go
type Span struct {
    StartTime time.Time //The time when the request was received
    EndTime time.Time //Returned time
    traceId string //traceId
    spanId string //spanId
    parentSpanId string // parent spanId
    tags [xprotocol.TRACE_END]string //Expand the field, the component can store its own information in this field
    operationName string
}
```
The Span structure defines the data structure passed between layotto and its component, as shown in the following figure, component can pass its own information to layotto through tags, and layotto does
Unified trace report:

#### generator interface:

```go
type Generator interface {
    GetTraceId(ctx context.Context) string
    GetSpanId(ctx context.Context) string
    GenerateNewContext(ctx context.Context, span api.Span) context.Context
    GetParentSpanId(ctx context.Context) string
}
```

This interface corresponds to the generator configuration above. This interface is mainly used to generate traceId, spanId according to the received context, obtain the parent spanId and the function of the context passed to the component, the user
You can implement your own Generator, you can refer to the implementation of [OpenGenerator](../../../../diagnostics/genetator.go) in the code.

#### Exporter interface:

```go
type Exporter interface {
ExportSpan(s *Span)
}
```

The exporter interface defines how to report Span information to the remote end, corresponding to the exporter field in the configuration, which is an array and can be reported to multiple servers. Can
Refer to the implementation of [StdoutExporter](../../../../diagnostics/exporter_iml/stdout.go), which will print trace information to standard output.


#### Span context transfer:

##### Layotto test
```go
GenerateNewContext(ctx context.Context, span api.Span) context.Context
```

GenerateNewContext is used to generate a new context, and we can pass the context between contexts through mosnctx:

```go
ctx = mosnctx.WithValue(ctx, types.ContextKeyActiveSpan, span)
```
You can refer to the implementation of [OpenGenerator](../../../../diagnostics/genetator.go) in the code

##### Component test

In Component measurement, you can insert component information through [SetExtraComponentInfo](../../../../components/trace/utils.go),
For example, the following operations are performed on the interface [Hello](../../../../components/hello/helloworld/helloworld.go):

```go
trace.SetExtraComponentInfo(ctx, fmt.Sprintf("method: %+v", "hello"))
```

The results printed by trace are as follows:

![img.png](../../../img/trace/trace.png)