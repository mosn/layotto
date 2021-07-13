import {
  call_foreign_function,
  Context,
  FilterHeadersStatusValues,
  registerRootContext,
  RootContext,
  set_http_response_body,
  stream_context
} from "@nobodyiam/proxy-runtime";

export * from "@nobodyiam/proxy-runtime/proxy"; // this exports the required functions for the proxy to interact with us.

class BridgeHttpResponseRoot extends RootContext {
  createContext(context_id: u32): Context {
    return new BridgeHttpResponse(context_id, this);
  }
}

class BridgeHttpResponse extends Context {
  constructor(context_id: u32, root_context: BridgeHttpResponseRoot) {
    super(context_id, root_context);
  }

  onRequestHeaders(a: u32, end_of_stream: bool): FilterHeadersStatusValues {
    let name = stream_context.headers.request.get("Name");
    let result = call_foreign_function("SayHello", '{"service_name":"helloworld","name":"' + name + '"}');
    set_http_response_body(String.UTF8.decode(result));
    return FilterHeadersStatusValues.Continue;
  }
}

registerRootContext((context_id: u32) => {
  return new BridgeHttpResponseRoot(context_id);
}, "");

export function _start(): void {
}
